package cli

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/helpdocs"
	"github.com/cego/aikido-dojo/internal/ops"
	"github.com/cego/aikido-dojo/internal/search"
)

// docsHits matches search: enough to show a good match's neighbours.
const docsHits = 10

// docsTimeout bounds each request of docs sync, as internal/tools/docsfetch does.
const docsTimeout = 30 * time.Second

func (a *app) docsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "docs",
		Short:                      "Search and read Aikido's help centre, from a copy built into aikido-dojo or fetched with docs sync",
		Args:                       unknownCommand,
		RunE:                       showHelp,
		SuggestionsMinimumDistance: 2,
	}
	cmd.AddCommand(a.docsSearchCmd(), a.docsGetCmd(), a.docsSyncCmd())
	return cmd
}

func (a *app) docsSearchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search <words...>",
		Short: "Find the help centre pages that answer a question",
		Long: "docs search ranks every section of Aikido's help centre (help.aikido.dev) against your words, using the copy " +
			"built into aikido-dojo, or a newer one docs sync fetched. It runs locally and sends nothing anywhere.\n\n" +
			"It prints up to 10 pages as a JSON array of {path, title, section, url, snippet, score}, best first, each at " +
			"the section that matched best. aikido-dojo docs get <path> then prints a page.",
		Example: `  aikido-dojo docs search "limit access to specific IPs"`,
		Args:    someArgs("words describing the question"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.docsSearch(cmd.Context(), strings.Join(args, " "))
		},
	}
}

func (a *app) docsSearch(ctx context.Context, query string) error {
	if len(search.Tokens(query)) == 0 {
		return &clierr.Error{Code: "invalid_input", Message: fmt.Sprintf("%q has no words to search for", query),
			Hint: `describe the question, for example: aikido-dojo docs search "connect a GitLab repository"`, Exit: clierr.ExitUsage}
	}
	hits, err := a.docsStore().Search(query, docsHits)
	if err != nil {
		return err
	}
	return a.printJSON(ctx, hits)
}

func (a *app) docsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <path>",
		Short: "Print a help centre page, from aikido-dojo's copy",
		Long: "docs get prints one page of Aikido's help centre as JSON: {path, title, url, snapshot, markdown}. " +
			"<path> is a path docs search printed, or the page's help.aikido.dev URL. snapshot dates the copy; url is " +
			"the live page.\n\nPrint only the markdown with --jq .markdown",
		Example: "  aikido-dojo docs get getting-started/setting-up-your-account/change-language --jq .markdown",
		Args:    exactArgs([]ops.Param{{Name: "path"}}),
		RunE:    func(cmd *cobra.Command, args []string) error { return a.docsGet(cmd.Context(), args[0]) },
	}
}

func (a *app) docsGet(ctx context.Context, ref string) error {
	page, err := a.docsStore().Get(ref)
	if errors.Is(err, helpdocs.ErrNotFound) {
		return &clierr.Error{Code: "not_found", Message: fmt.Sprintf("no help centre page %q in aikido-dojo's copy", ref),
			Hint: `find one with: aikido-dojo docs search "<words>"`, Exit: clierr.ExitNotFound}
	}
	if err != nil {
		return err
	}
	return a.printJSON(ctx, page)
}

func (a *app) docsSyncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "Fetch the help centre pages Aikido changed since aikido-dojo's copy",
		Long: "docs sync updates aikido-dojo's copy of Aikido's help centre from help.aikido.dev, so docs search and docs get " +
			"read the pages Aikido changed since this release. It fetches only the pages that changed; when none did, it " +
			"fetches only the site's llms-full.txt. The copy goes in your cache directory. Updating aikido-dojo discards it, " +
			"and the docs commands read the copy built into the new release until the next docs sync.\n\n" +
			"It prints what it did as JSON: changed and removed list page paths, requests counts the requests it sent, and " +
			"snapshot dates the copy.",
		Example: "  aikido-dojo docs sync",
		Args:    exactArgs(nil),
		RunE:    func(cmd *cobra.Command, _ []string) error { return a.docsSync(cmd.Context()) },
	}
}

func (a *app) docsSync(ctx context.Context) error {
	cache, err := a.env.CacheDir()
	if err != nil {
		return fmt.Errorf("find the cache directory for the docs: %w", err)
	}
	hc := newHTTPClient(a.env, strings.TrimPrefix(helpdocs.Origin, "https://"), a.debug)
	hc.Timeout = docsTimeout
	r, err := docsStoreIn(cache).Sync(ctx, hc)
	if err != nil {
		return &clierr.Error{Code: "docs_sync_failed", Message: "docs sync failed", Err: err,
			Hint: "the docs commands still read the copy they read before; retry later", Exit: clierr.ExitUnexpected}
	}
	return a.printJSON(ctx, r)
}

// docsStore is the copy the docs commands read; without a cache directory, that is the built-in one.
func (a *app) docsStore() helpdocs.Store {
	cache, err := a.env.CacheDir()
	if err != nil {
		return helpdocs.Store{}
	}
	return docsStoreIn(cache)
}

func docsStoreIn(cache string) helpdocs.Store {
	return helpdocs.Store{Dir: filepath.Join(cache, "aikido-dojo", "docs"), Build: build}
}
