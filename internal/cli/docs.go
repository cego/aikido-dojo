package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/helpdocs"
	"github.com/cego/aikido-dojo/internal/ops"
	"github.com/cego/aikido-dojo/internal/search"
)

// docsHits matches search: enough to show a good match's neighbours.
const docsHits = 10

func (a *app) docsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "docs",
		Short:                      "Search and read Aikido's help centre, from a copy built into aikido-dojo",
		Args:                       unknownCommand,
		RunE:                       showHelp,
		SuggestionsMinimumDistance: 2,
	}
	cmd.AddCommand(a.docsSearchCmd(), a.docsGetCmd())
	return cmd
}

func (a *app) docsSearchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search <words...>",
		Short: "Find the help centre pages that answer a question",
		Long: "docs search ranks every section of Aikido's help centre (help.aikido.dev) against your words, using a copy " +
			"built into aikido-dojo. It runs locally and sends nothing anywhere.\n\n" +
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
	hits, err := helpdocs.Search(query, docsHits)
	if err != nil {
		return err
	}
	return a.printJSON(ctx, hits)
}

func (a *app) docsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <path>",
		Short: "Print a help centre page, from the copy built into aikido-dojo",
		Long: "docs get prints one page of Aikido's help centre as JSON: {path, title, url, snapshot, markdown}. " +
			"<path> is a path docs search printed, or the page's help.aikido.dev URL. snapshot dates the copy; url is " +
			"the live page.\n\nPrint only the markdown with --jq .markdown",
		Example: "  aikido-dojo docs get getting-started/setting-up-your-account/change-language --jq .markdown",
		Args:    exactArgs([]ops.Param{{Name: "path"}}),
		RunE:    func(cmd *cobra.Command, args []string) error { return a.docsGet(cmd.Context(), args[0]) },
	}
}

func (a *app) docsGet(ctx context.Context, ref string) error {
	page, err := helpdocs.Get(ref)
	if errors.Is(err, helpdocs.ErrNotFound) {
		return &clierr.Error{Code: "not_found", Message: fmt.Sprintf("no help centre page %q in aikido-dojo's copy", ref),
			Hint: `find one with: aikido-dojo docs search "<words>"`, Exit: clierr.ExitNotFound}
	}
	if err != nil {
		return err
	}
	return a.printJSON(ctx, page)
}
