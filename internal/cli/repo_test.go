package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zalando/go-keyring"

	"github.com/cego/aikido-dojo/internal/clierr"
)

// inRepo runs the rest of the test in a new directory: a git checkout whose
// origin is origin, or, with origin "", a directory outside any checkout.
// git reads no user or system config, whose insteadOf rules would rewrite
// the remote.
func inRepo(t *testing.T, origin string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	dir := t.TempDir()
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_CEILING_DIRECTORIES", filepath.Dir(dir))
	if origin != "" {
		for _, args := range [][]string{{"init", "-q"}, {"remote", "add", "origin", origin}} {
			cmd := exec.CommandContext(t.Context(), "git", args...) //nolint:gosec // git with this helper's own arguments
			cmd.Dir = dir
			// Only the subcommand is named: the origin may be a real tenant's URL.
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("git %s: %v: %s", args[0], err, out)
			}
		}
	}
	t.Chdir(dir)
}

// repos serves the repo list: the items whose name equals filter_name when
// it is set, as Aikido filters (live, 2026-10-02), else every item; then an
// empty page.
func repos(items ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") != "0" {
			fmt.Fprint(w, `[]`)
			return
		}
		filter := r.URL.Query().Get("filter_name")
		var page []string
		for _, item := range items {
			var repo struct {
				Name string `json:"name"`
			}
			if json.Unmarshal([]byte(item), &repo) == nil && (filter == "" || filter == repo.Name) {
				page = append(page, item)
			}
		}
		fmt.Fprint(w, "["+strings.Join(page, ",")+"]")
	}
}

func TestParseRemote(t *testing.T) {
	tests := []struct{ in, key, name string }{
		{"git@gitlab.example.com:Group/Sub/Repo.git", "gitlab.example.com/group/sub/repo", "Repo"},
		{"ssh://git@gitlab.example.com:2222/Group/Sub/Repo.git", "gitlab.example.com/group/sub/repo", "Repo"},
		{"https://user:hunter2@GitLab.example.com/Group/Sub/Repo.git/", "gitlab.example.com/group/sub/repo", "Repo"},
		{"https://gitlab.example.com/group/sub/repo", "gitlab.example.com/group/sub/repo", "repo"},
		{"gitlab.example.com:group/repo.git", "gitlab.example.com/group/repo", "repo"},
		{"https://api.github.com/repos/aikidemo/compression-service", "github.com/aikidemo/compression-service", "compression-service"},
		{"git@github.com:aikidemo/compression-service.git", "github.com/aikidemo/compression-service", "compression-service"},
	}
	for _, tt := range tests {
		got, ok := parseRemote(tt.in)
		if !ok || got.key != tt.key || got.name != tt.name {
			t.Errorf("parseRemote(%q) = %+v, %v; want %s and %s", tt.in, got, ok, tt.key, tt.name)
		}
	}
	// The last two are malformed remotes git can't use, whose credential
	// would otherwise land in the key that not_found quotes.
	for _, in := range []string{"", "/srv/git/repo.git", "../repo", "file:///srv/repo.git", "https://host.example.com/",
		"user:hunter2@gitlab.example.com:g/r.git", "https://user:12/hunter2@gitlab.example.com/g/r.git"} {
		if got, ok := parseRemote(in); ok {
			t.Errorf("parseRemote(%q) = %+v, want no remote", in, got)
		}
	}
}

func TestRepoCurrentFindsTheRepoInOneCall(t *testing.T) {
	inRepo(t, "git@gitlab.example.com:Group/Sub/Repo.git")
	f, env, _ := newFake(t, repos(
		`{"id":8,"name":"Repo","url":"https://gitlab.example.com/other/repo.git"}`,
		`{"id":7,"name":"Repo","url":"https://gitlab.example.com/group/sub/repo.git"}`,
	))
	want := `[{"repo":{"id":7,"name":"Repo","url":"https://gitlab.example.com/group/sub/repo.git"}}]` + "\n"
	if out := mustRun(t, env, "repo", "current"); out != want {
		t.Errorf("stdout = %s, want %s", out, want)
	}
	if s := f.seen(); len(s) != 1 || s[0].Query.Get("filter_name") != "Repo" || s[0].Query.Get("per_page") != "200" {
		t.Errorf("requests = %+v, want one filtered page", s)
	}
}

func TestRepoCurrentFallsBackToEveryRepo(t *testing.T) {
	inRepo(t, "https://gitlab.example.com/Group/Sub/Repo.git")
	// A repo renamed in Aikido: the filter misses it, the full list has it.
	f, env, _ := newFake(t, repos(`{"id":9,"name":"Renamed","url":"https://gitlab.example.com/group/sub/repo"}`))
	if out := mustRun(t, env, "repo", "current"); !strings.Contains(out, `"id":9`) {
		t.Errorf("stdout = %s, want repo 9", out)
	}
	if s := f.seen(); len(s) != 2 || s[1].Query.Has("filter_name") {
		t.Errorf("requests = %+v, want the filtered page, then the unfiltered one", s)
	}
}

func TestRepoCurrentAcrossProfiles(t *testing.T) {
	inRepo(t, "git@gitlab.example.com:group/repo.git")
	match := `{"id":3,"name":"repo","url":"https://gitlab.example.com/group/repo.git"}`
	f, env, vars := newFake(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Header.Get("Authorization") {
		case "Bearer tok-a":
			repos(match)(w, r)
		case "Bearer tok-b":
			w.WriteHeader(http.StatusForbidden)
		default:
			repos()(w, r)
		}
	})
	f.issueTokens(func(id, _ string) (int, string) { return http.StatusOK, "tok-" + id })
	writeConfig(t, vars, `{"profiles":{"a":{"client_id":"a"},"b":{"client_id":"b"}}}`)
	for _, p := range []string{"a", "b"} {
		if err := keyring.Set("aikido-dojo", p+"/client_secret", "s"); err != nil {
			t.Fatal(err)
		}
	}
	stdout, stderr, code := run(t, env, "repo", "current")
	if want := `[{"profile":"a","repo":` + match + `}]` + "\n"; code != clierr.ExitOK || stdout != want {
		t.Errorf("exit %d, stdout %s; want %s", code, stdout, want)
	}
	if !strings.Contains(stderr, `profile \"b\" skipped`) || !strings.Contains(stderr, "repositories:read") {
		t.Errorf("stderr = %s, want a warning that b lacks repositories:read", stderr)
	}
	if n := len(f.seenLogins()); n != 3 {
		t.Errorf("token requests = %d, want one each for the environment pair, a and b", n)
	}
}

// A profile that can't be resolved is skipped like one that fails a call.
func TestRepoCurrentSkipsAMisconfiguredProfile(t *testing.T) {
	inRepo(t, "git@gitlab.example.com:g/r.git")
	f, env, vars := newFake(t, repos(`{"id":1,"name":"r","url":"https://gitlab.example.com/g/r.git"}`))
	noPair(vars)
	writeConfig(t, vars, `{"profiles":{"a":{"client_id":"a"},"bare":{},"far":{"client_id":"f","region":"mars"}}}`)
	if err := keyring.Set("aikido-dojo", "a/client_secret", "s"); err != nil {
		t.Fatal(err)
	}
	stdout, stderr, code := run(t, env, "repo", "current")
	if code != clierr.ExitOK || !strings.HasPrefix(stdout, `[{"profile":"a",`) {
		t.Errorf("exit %d, stdout %s, stderr %s; want profile a's match", code, stdout, stderr)
	}
	for _, want := range []string{`profile \"bare\" skipped`, `profile \"far\" skipped`} {
		if !strings.Contains(stderr, want) {
			t.Errorf("stderr = %s, want %s", stderr, want)
		}
	}
	if n := len(f.seenLogins()); n != 1 {
		t.Errorf("token requests = %d, want only profile a's", n)
	}
}

func TestRepoCurrentAsksOnlyTheSelectedProfile(t *testing.T) {
	inRepo(t, "git@gitlab.example.com:g/r.git")
	f, env, vars := newFake(t, repos(`{"id":1,"name":"r","url":"https://gitlab.example.com/g/r.git"}`))
	writeConfig(t, vars, `{"profiles":{"a":{"client_id":"a"},"b":{"client_id":"b"}}}`)
	if err := keyring.Set("aikido-dojo", "b/client_secret", "s"); err != nil {
		t.Fatal(err)
	}
	out := mustRun(t, env, "--profile", "b", "repo", "current")
	if logins := f.seenLogins(); !strings.HasPrefix(out, `[{"profile":"b",`) || len(logins) != 1 || !strings.HasPrefix(logins[0], "b:") {
		t.Errorf("stdout %s, token requests %q; want only profile b", out, logins)
	}
}

func TestRepoCurrentNotFoundHidesTheCredential(t *testing.T) {
	inRepo(t, "https://user:hunter2@gitlab.example.com/g/r.git")
	_, env, _ := newFake(t, repos(`{"id":1,"name":"r","url":"https://gitlab.example.com/other/r.git"}`))
	stdout, stderr, code := run(t, env, "repo", "current")
	if e := errorOf(t, stderr); code != clierr.ExitNotFound || e.Code != "not_found" || !strings.Contains(e.Message, "gitlab.example.com/g/r") {
		t.Errorf("exit %d, error %+v; want not_found naming the URL", code, e)
	}
	if strings.Contains(stdout+stderr, "hunter2") {
		t.Errorf("output holds the remote's credential: %s %s", stdout, stderr)
	}
}

func TestRepoCurrentReportsTheErrorWhenEveryProfileFails(t *testing.T) {
	inRepo(t, "git@gitlab.example.com:g/r.git")
	_, env, _ := newFake(t, func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusForbidden) })
	_, stderr, code := run(t, env, "repo", "current")
	if e := errorOf(t, stderr); code != clierr.ExitForbidden || e.Code != "missing_scope" {
		t.Errorf("exit %d, error %+v; want 4 missing_scope", code, e)
	}
	// With one profile, the error says it all; a warning would repeat it.
	if strings.Contains(stderr, `"warning"`) {
		t.Errorf("stderr = %s, want the error alone", stderr)
	}
}

func TestRepoCurrentNeedsAnOrigin(t *testing.T) {
	for _, tt := range []struct{ name, origin, code string }{
		{"outside a checkout", "", "no_git_remote"},
		{"a local path", "/srv/git/repo.git", "unsupported_remote"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			inRepo(t, tt.origin)
			f, env, _ := newFake(t, repos())
			_, stderr, code := run(t, env, "repo", "current")
			if e := errorOf(t, stderr); code != clierr.ExitUsage || e.Code != tt.code || e.Hint == "" {
				t.Errorf("exit %d, error %+v; want exit 2 %s with a hint", code, e, tt.code)
			}
			if n := len(f.seen()); n != 0 {
				t.Errorf("API calls = %d, want none", n)
			}
		})
	}
}

func TestRepoCurrentIsAmongTheRepoVerbs(t *testing.T) {
	env, _ := testEnv(t)
	root, err := buildRoot(env)
	if err != nil {
		t.Fatal(err)
	}
	repo, _, err := root.Find([]string{"repo"})
	if err != nil || repo.Short != "activate, clone, current, deactivate, delete, get, import, list, scan" {
		t.Errorf("repo = %q, %v", repo.Short, err)
	}
}
