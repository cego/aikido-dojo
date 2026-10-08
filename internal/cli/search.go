package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cego/aikido-dojo/internal/catalog"
	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/search"
)

// searchHits is enough to show a good match's neighbours without burying it.
const searchHits = 10

func (a *app) searchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search <words...>",
		Short: "Find the commands for a task, from words describing it",
		Long: "search ranks every generated command against your words, using an index of the API's summaries, " +
			"descriptions and parameter names built into aikido-dojo. It runs locally and sends nothing anywhere.\n\n" +
			"It prints up to 10 commands as a JSON array of {command, summary, score}, best first. " +
			"aikido-dojo schema <resource> <verb> then shows what a command takes and prints.",
		Example: `  aikido-dojo search "ignore a finding"`,
		Args:    someArgs("words describing the task"),
		RunE:    func(cmd *cobra.Command, args []string) error { return a.search(cmd.Context(), strings.Join(args, " ")) },
	}
}

func (a *app) search(ctx context.Context, query string) error {
	if len(search.Tokens(query)) == 0 {
		return &clierr.Error{Code: "invalid_input", Message: fmt.Sprintf("%q has no words to search for", query),
			Hint: `describe the task, for example: aikido-dojo search "list repositories"`, Exit: clierr.ExitUsage}
	}
	idx, err := catalog.SearchIndex()
	if err != nil {
		return err
	}
	return a.printJSON(ctx, idx.Rank(query, searchHits))
}
