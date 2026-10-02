package cli

import (
	"fmt"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/ops"
)

// command is one generated command: its descriptor and the body fields that got a flag.
type command struct {
	op        ops.Op
	bodyFlags []ops.Param
}

type runFunc func(cmd *cobra.Command, c command, args []string) error

// addGenerated adds a command per resource and, under it, one per operation.
func addGenerated(root *cobra.Command, all []ops.Op, run runFunc) error {
	resources := map[string]*cobra.Command{}
	verbs := map[string][]string{}
	for _, op := range all {
		resource, verb, _ := strings.Cut(op.Command, " ")
		parent, ok := resources[resource]
		if !ok {
			parent = &cobra.Command{Use: resource, Args: unknownCommand, RunE: showHelp, SuggestionsMinimumDistance: 2}
			resources[resource] = parent
			root.AddCommand(parent)
		}
		cmd, err := opCommand(root, verb, op, run)
		if err != nil {
			return err
		}
		parent.AddCommand(cmd)
		verbs[resource] = append(verbs[resource], verb)
	}
	for resource, parent := range resources {
		slices.Sort(verbs[resource])
		parent.Short = strings.Join(verbs[resource], ", ")
	}
	return nil
}

func opCommand(root *cobra.Command, verb string, op ops.Op, run runFunc) (*cobra.Command, error) {
	use := verb
	for _, a := range op.Args {
		use += " <" + a.Name + ">"
	}
	cmd := &cobra.Command{Use: use, Short: op.Summary, Args: exactArgs(op.Args), Annotations: map[string]string{"operation": op.ID}}
	fs := cmd.Flags()
	taken := func(name string) bool {
		return name == "help" || fs.Lookup(name) != nil || root.PersistentFlags().Lookup(name) != nil
	}
	if op.Paging != nil {
		fs.Int("limit", 0, "stop after this many items; 0 prints every item")
	}
	if op.Body != nil {
		fs.String("body", "", "the request body as JSON")
		fs.String("body-file", "", "a file holding the request body as JSON, or - for stdin")
	}
	for _, p := range op.Flags {
		if taken(ops.FlagName(p.Name)) {
			return nil, fmt.Errorf("%s: --%s clashes with a built-in flag", op.Command, ops.FlagName(p.Name))
		}
		addFlag(fs, p)
	}
	c := command{op: op}
	var clashes []string
	if op.Body != nil {
		for _, f := range op.Body.Fields {
			if taken(ops.FlagName(f.Name)) {
				clashes = append(clashes, f.Name)
				continue
			}
			addFlag(fs, f)
			c.bodyFlags = append(c.bodyFlags, f)
		}
	}
	cmd.Long = longHelp(op, clashes)
	cmd.RunE = func(cmd *cobra.Command, args []string) error { return run(cmd, c, args) }
	return cmd, nil
}

// addFlag registers p with a zero default: the runner sends only the flags
// the user set, so Aikido's own defaults apply to the rest.
func addFlag(fs *pflag.FlagSet, p ops.Param) {
	name := ops.FlagName(p.Name)
	switch p.Kind {
	case ops.Boolean:
		fs.Bool(name, false, p.Usage)
	case ops.Integer:
		fs.Int64(name, 0, p.Usage)
	case ops.StringList, ops.IntegerList:
		fs.StringSlice(name, nil, p.Usage)
	default:
		fs.String(name, "", p.Usage)
	}
}

func exactArgs(params []ops.Param) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == len(params) {
			return nil
		}
		want := "no arguments"
		if len(params) > 0 {
			names := make([]string, len(params))
			for i, p := range params {
				names[i] = "<" + p.Name + ">"
			}
			want = strings.Join(names, " ")
		}
		return &clierr.Error{Code: "usage", Message: fmt.Sprintf("%s takes %s, got %d", cmd.CommandPath(), want, len(args)),
			Hint: "run " + cmd.CommandPath() + " --help", Exit: clierr.ExitUsage}
	}
}

func longHelp(op ops.Op, clashes []string) string {
	var parts, notes []string
	for _, s := range []string{op.Summary, op.Description, op.Help} {
		if s != "" && !slices.Contains(parts, s) {
			parts = append(parts, s)
		}
	}
	if op.Scope != "" {
		notes = append(notes, "Needs the "+op.Scope+" scope.")
	}
	if op.Body != nil {
		notes = append(notes, "Fields without a flag, nested ones included, go in --body or --body-file.")
		if len(op.Body.Secret) > 0 {
			notes = append(notes, "Credentials ("+strings.Join(op.Body.Secret, ", ")+") are accepted only there, which keeps them out of shell history.")
		}
		if len(clashes) > 0 {
			notes = append(notes, strings.Join(clashes, ", ")+" share a name with a built-in flag, so they have no flag of their own.")
		}
	}
	if len(notes) > 0 {
		parts = append(parts, strings.Join(notes, " "))
	}
	return strings.Join(parts, "\n\n")
}
