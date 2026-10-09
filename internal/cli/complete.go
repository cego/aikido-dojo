package cli

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cego/aikido-dojo/internal/catalog"
	"github.com/cego/aikido-dojo/internal/config"
	"github.com/cego/aikido-dojo/internal/ops"
)

// completing reports whether args are a request from the scripts aikido-dojo completion prints.
func completing(args []string) bool {
	return len(args) > 0 && (args[0] == cobra.ShellCompRequestCmd || args[0] == cobra.ShellCompNoDescRequestCmd)
}

// addCompletion registers what flags and generated arguments complete to; cobra keeps flag completions for the whole process, so only a completion request pays for them.
func (a *app) addCompletion(root *cobra.Command) error {
	root.CompletionOptions.SetDefaultShellCompDirective(cobra.ShellCompDirectiveNoFileComp)
	for _, c := range []struct {
		path []string
		flag string
		fn   cobra.CompletionFunc
	}{
		{nil, "config", files},
		{nil, "profile", a.profileNames},
		{[]string{"auth", "login"}, "region", cobra.FixedCompletions(config.Regions(), cobra.ShellCompDirectiveNoFileComp)},
		{[]string{"api"}, "input", files},
	} {
		cmd, _, err := root.Find(c.path)
		if err != nil {
			return fmt.Errorf("find %s: %w", strings.Join(c.path, " "), err)
		}
		if err := register(cmd, c.flag, c.fn); err != nil {
			return err
		}
	}
	for _, op := range catalog.All {
		cmd, _, err := root.Find(strings.Fields(op.Command))
		if err != nil {
			return fmt.Errorf("find %s: %w", op.Command, err)
		}
		if err := completeOp(cmd, op); err != nil {
			return err
		}
	}
	return nil
}

// completeOp offers file names for a generated command's --body-file and --path, and the spec's values where a flag or argument has an enum.
func completeOp(cmd *cobra.Command, op ops.Op) error {
	sc, err := catalog.Schemas(op.ID)
	if err != nil {
		return err
	}
	for _, name := range []string{"body-file", "path"} {
		if cmd.LocalFlags().Lookup(name) != nil {
			if err := register(cmd, name, files); err != nil {
				return err
			}
		}
	}
	for name, schema := range properties(flagSchema(cmd, op, sc)) {
		if values := enumValues(schema); values != nil {
			if err := register(cmd, name, cobra.FixedCompletions(values, cobra.ShellCompDirectiveNoFileComp)); err != nil {
				return err
			}
		}
	}
	args := make([][]cobra.Completion, len(op.Args))
	for i, p := range op.Args {
		args[i] = enumValues(properties(sc.Args)[p.Name])
	}
	cmd.ValidArgsFunction = func(_ *cobra.Command, typed []string, _ string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if len(typed) < len(args) {
			return args[len(typed)], cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return nil
}

// enumValues lists what schema's enum, or a list's items' enum, allows; null is left out, as no flag can send it.
func enumValues(schema any) []cobra.Completion {
	s, _ := schema.(map[string]any)
	if items, ok := s["items"].(map[string]any); ok {
		s = items
	}
	enum, _ := s["enum"].([]any)
	var values []cobra.Completion
	for _, v := range enum {
		if v != nil {
			values = append(values, fmt.Sprint(v))
		}
	}
	return values
}

func files(*cobra.Command, []string, string) ([]cobra.Completion, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveDefault
}

func register(cmd *cobra.Command, flag string, fn cobra.CompletionFunc) error {
	if err := cmd.RegisterFlagCompletionFunc(flag, fn); err != nil {
		return fmt.Errorf("complete %s --%s: %w", cmd.CommandPath(), flag, err)
	}
	return nil
}

// profileNames completes --profile from the config file; a missing or broken one offers none, as a shell has nowhere to show the error.
func (a *app) profileNames(*cobra.Command, []string, string) ([]cobra.Completion, cobra.ShellCompDirective) {
	_, f, err := a.configFile()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return slices.Sorted(maps.Keys(f.Profiles)), cobra.ShellCompDirectiveNoFileComp
}
