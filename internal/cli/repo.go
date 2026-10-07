package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/url"
	"os/exec"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cego/aikido-dojo/internal/api"
	"github.com/cego/aikido-dojo/internal/catalog"
	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/config"
	"github.com/cego/aikido-dojo/internal/ops"
)

// repoMatch is an Aikido code repo with the working directory's origin.
// Profile is empty for the AIKIDO_DOJO_CLIENT_ID pair.
type repoMatch struct {
	Profile string          `json:"profile,omitempty"`
	Repo    json.RawMessage `json:"repo"`
}

// candidate is a profile repo current asks, or the reason it can't be asked.
type candidate struct {
	label string
	r     config.Resolved
	err   error
}

// remote is a repo URL reduced to what its SSH and HTTPS forms share.
type remote struct {
	key  string // lowercased host/path
	name string // the last path segment as written, which Aikido stores as the repo's name
}

// addRepoCurrent puts the hand-written repo current beside the generated repo verbs.
func (a *app) addRepoCurrent(root *cobra.Command) error {
	repo, _, err := root.Find([]string{"repo"})
	list, ok := catalog.Command("repo list")
	if err != nil || repo == root || !ok || list.Paging == nil {
		return errors.New("the catalog has no paged repo list for repo current")
	}
	repo.AddCommand(a.repoCurrentCmd(list))
	repo.Short = strings.Join(subcommandNames(repo), ", ")
	return nil
}

func (a *app) repoCurrentCmd(list ops.Op) *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "Find the Aikido code repo of this directory's git origin, in every profile",
		Long: "current reads git remote get-url origin in the working directory and finds the Aikido code repo with " +
			"that URL. It asks the profile --profile or AIKIDO_DOJO_PROFILE selects; otherwise the AIKIDO_DOJO_CLIENT_ID " +
			"pair, if set, and every profile in the config file. A profile that fails is skipped with a warning.\n\n" +
			"It prints a JSON array of {profile, repo}, where repo is the entry repo list prints; pass its id, with " +
			"that --profile, to other commands. SSH and HTTPS remotes of one repo match. Inactive repos aren't searched.\n\n" +
			"Needs the " + list.Scope + " scope.",
		Args: exactArgs(nil),
		RunE: func(cmd *cobra.Command, _ []string) error { return a.repoCurrent(cmd.Context(), list) },
	}
}

func (a *app) repoCurrent(ctx context.Context, list ops.Op) error {
	origin, err := gitOrigin(ctx)
	if err != nil {
		return err
	}
	want, ok := parseRemote(origin)
	if !ok {
		// The remote may hold a credential, so it is described rather than quoted.
		return &clierr.Error{Code: "unsupported_remote", Message: "origin is not a URL with a host and a path",
			Hint: "repo current needs an origin such as git@host:group/repo.git or https://host/group/repo.git", Exit: clierr.ExitUsage}
	}
	candidates, err := a.everyProfile()
	if err != nil {
		return err
	}
	matches := []repoMatch{}
	var firstErr error
	searched := 0
	for _, c := range candidates {
		err := c.err
		var repo json.RawMessage
		if err == nil {
			repo, err = a.findRepo(ctx, c.r, list, want)
		}
		if ctx.Err() != nil {
			return fmt.Errorf("repo current: %w", ctx.Err())
		}
		if err != nil {
			clierr.Warn(a.env.Stderr, c.label+" skipped: "+explain(err))
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		searched++
		if repo != nil {
			matches = append(matches, repoMatch{Profile: c.r.Profile, Repo: repo})
		}
	}
	if searched == 0 {
		return firstErr
	}
	if len(matches) == 0 {
		return &clierr.Error{Code: "not_found", Message: "no Aikido code repo has the URL " + want.key,
			Hint: "check that the repo is connected in Aikido, and that a profile here sees its workspace: aikido-dojo auth status", Exit: clierr.ExitNotFound}
	}
	return a.printJSON(matches)
}

// everyProfile is who repo current asks: the profile --profile or
// AIKIDO_DOJO_PROFILE selects, else the environment pair, if set, and every
// stored profile. Of those, one that can't be resolved comes with its error,
// so the others are still asked.
func (a *app) everyProfile() ([]candidate, error) {
	_, f, err := a.configFile()
	if err != nil {
		return nil, err
	}
	flags := a.flags()
	if flags.Profile != "" || a.env.Getenv(config.EnvProfile) != "" || len(f.Profiles) == 0 {
		r, err := config.Resolve(f, flags, a.env.Getenv)
		if err != nil {
			return nil, err
		}
		return []candidate{{label: profileLabel(r.Profile), r: r}}, nil
	}
	var out []candidate
	if a.env.Getenv(config.EnvClientID) != "" || a.env.Getenv(config.EnvClientSecret) != "" {
		// An empty file, so Resolve takes the pair rather than the default profile.
		r, err := config.Resolve(config.File{}, flags, a.env.Getenv)
		out = append(out, candidate{label: profileLabel(""), r: r, err: err})
	}
	for _, name := range slices.Sorted(maps.Keys(f.Profiles)) {
		flags.Profile = name
		r, err := config.Resolve(f, flags, a.env.Getenv)
		out = append(out, candidate{label: profileLabel(name), r: r, err: err})
	}
	return out, nil
}

// findRepo looks in one profile's workspace for the repo with want's URL.
// filter_name matches a repo's name exactly, and the name is the last segment
// of the repo's path for 744 of 751 repos (live, 2026-10-02), so one filtered
// call usually finds it; a miss pages through every repo.
func (a *app) findRepo(ctx context.Context, r config.Resolved, list ops.Op, want remote) (json.RawMessage, error) {
	c, err := a.clientFor(r)
	if err != nil {
		return nil, err
	}
	repo, err := repoWithURL(ctx, c, list, url.Values{"filter_name": {want.name}}, want.key)
	if err != nil || repo != nil {
		return repo, err
	}
	return repoWithURL(ctx, c, list, nil, want.key)
}

// repoWithURL pages through the repos query selects and returns the first
// with the URL key, or nil. A workspace holds each URL once (751 URLs for
// 751 repos, live 2026-10-02), so the first is the only one.
func repoWithURL(ctx context.Context, c *api.Client, list ops.Op, query url.Values, key string) (json.RawMessage, error) {
	req := api.Request{Method: list.Method, Path: list.Path, Query: query, Scope: list.Scope}
	for item, err := range c.Items(ctx, req, *list.Paging) {
		if err != nil {
			return nil, err
		}
		var repo struct {
			URL string `json:"url"`
		}
		// An item without a url string can't be the repo, so it is passed over.
		if json.Unmarshal(item, &repo) == nil {
			if got, ok := parseRemote(repo.URL); ok && got.key == key {
				return item, nil
			}
		}
	}
	return nil, nil
}

// gitOrigin runs git remote get-url origin, which applies the user's
// insteadOf rules, so the URL is the one git really fetches from.
func gitOrigin(ctx context.Context) (string, error) {
	var stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, "git", "remote", "get-url", "origin")
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if errors.Is(err, exec.ErrNotFound) {
		return "", &clierr.Error{Code: "no_git", Message: "git is not installed or not on PATH",
			Hint: "install git, or find the repo with aikido-dojo repo list --filter-name <name>", Exit: clierr.ExitUsage}
	}
	if err != nil {
		return "", &clierr.Error{Code: "no_git_remote", Message: "read the git origin: " + strings.TrimSpace(stderr.String()), Err: err,
			Hint: "run it inside a git checkout that has an origin remote", Exit: clierr.ExitUsage}
	}
	return strings.TrimSpace(string(out)), nil
}

// parseRemote reduces a git remote, or the url Aikido stores for a repo, to
// host/path: git@host:group/repo.git, ssh://git@host:22/group/repo and
// https://user@host/group/repo.git all become host/group/repo. A local path
// isn't a remote.
func parseRemote(raw string) (remote, bool) {
	s := strings.TrimSpace(raw)
	var host, path string
	if strings.Contains(s, "://") {
		u, err := url.Parse(s)
		if err != nil {
			return remote{}, false
		}
		host, path = u.Hostname(), u.Path
	} else {
		// git's scp-like form, [user@]host:path, has its colon before any slash.
		before, after, ok := strings.Cut(s, ":")
		if !ok || strings.Contains(before, "/") {
			return remote{}, false
		}
		if _, h, found := strings.Cut(before, "@"); found {
			before = h
		}
		host, path = before, after
	}
	host = strings.ToLower(host)
	path = strings.TrimSuffix(strings.Trim(path, "/"), ".git")
	// Aikido stores a GitHub repo as its API URL (the listCodeRepos example in the spec).
	if rest, ok := strings.CutPrefix(path, "repos/"); ok && host == "api.github.com" {
		host, path = "github.com", rest
	}
	if host == "" || path == "" {
		return remote{}, false
	}
	return remote{key: host + "/" + strings.ToLower(path), name: path[strings.LastIndex(path, "/")+1:]}, true
}

func profileLabel(name string) string {
	if name == "" {
		return "the " + config.EnvClientID + " pair"
	}
	return "profile " + strconv.Quote(name)
}

// explain is err's message with its hint, for a warning that stands in for the error.
func explain(err error) string {
	var e *clierr.Error
	if errors.As(err, &e) && e.Hint != "" {
		return err.Error() + "; " + e.Hint
	}
	return err.Error()
}
