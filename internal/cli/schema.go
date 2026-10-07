package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/cego/aikido-dojo/internal/catalog"
	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/ops"
)

type schemaDoc struct {
	Command       string         `json:"command"`
	Usage         string         `json:"usage"`
	Method        string         `json:"method"`
	Path          string         `json:"path"`
	Scope         string         `json:"scope,omitempty"`
	Args          map[string]any `json:"args"`
	Flags         map[string]any `json:"flags"`
	Body          map[string]any `json:"body,omitempty"`
	Response      map[string]any `json:"response,omitempty"`
	ResponseTypes []string       `json:"response_types,omitempty"`
}

func (a *app) schemaCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "schema <resource> <verb>",
		Short: "Print what a command takes and prints, as JSON Schema, without calling the API",
		Long: "schema prints a generated command's inputs and output as JSON Schema: args (its positional arguments), " +
			"flags (keyed by flag name), body (what --body and --body-file take) and response (what the command prints). " +
			"A list command prints one array of every item, so its response is that array, not one page.\n\n" +
			"It reads only what is built into aikido-dojo and calls nothing.",
		Example: "  aikido-dojo schema issue-group ignore",
		Args:    someArgs("a command, such as: repo list"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.schema(cmd.Root(), strings.Fields(strings.Join(args, " ")))
		},
	}
}

func (a *app) schema(root *cobra.Command, words []string) error {
	name := strings.Join(words, " ")
	target, rest, err := root.Find(words)
	if err != nil || target == root || len(rest) > 0 {
		return &clierr.Error{Code: "unknown_command", Message: fmt.Sprintf("no command %q", name),
			Hint: `find one with aikido-dojo search "<what you want to do>"`, Exit: clierr.ExitUsage}
	}
	id := target.Annotations["operation"]
	switch {
	case id == "" && target.HasSubCommands():
		names := subcommandNames(target)
		return &clierr.Error{Code: "usage", Message: name + " groups commands: " + strings.Join(names, ", "),
			Hint: "name one, for example: aikido-dojo schema " + name + " " + names[0], Exit: clierr.ExitUsage}
	case id == "":
		return &clierr.Error{Code: "no_schema", Message: name + " is built into aikido-dojo, not generated from the API spec, so it has no schema",
			Hint: "run aikido-dojo " + name + " --help", Exit: clierr.ExitUsage}
	}
	sc, err := catalog.Schemas(id)
	if err != nil {
		return err
	}
	op, ok := catalog.Command(sc.Command)
	if !ok {
		return fmt.Errorf("the catalog has no command %s", sc.Command)
	}
	return a.printJSON(schemaDoc{
		Command: sc.Command, Usage: target.UseLine(), Method: sc.Method, Path: sc.Path, Scope: sc.Scope,
		Args: sc.Args, Flags: flagSchema(target, op, sc), Body: sc.Body, Response: sc.Response, ResponseTypes: sc.ResponseTypes,
	})
}

// flagSchema describes every flag of cmd that sets a value, keyed by flag
// name: the query parameters, the body fields that got a flag, and --limit
// on lists. --body and --body-file take the body schema.
func flagSchema(cmd *cobra.Command, op ops.Op, sc ops.SchemaSet) map[string]any {
	query, body := properties(sc.Flags), properties(sc.Body)
	props := map[string]any{}
	var required []any
	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		switch {
		case f.Name == "help" || f.Name == "body" || f.Name == "body-file":
			return
		case f.Name == "limit" && op.Paging != nil:
			props[f.Name] = map[string]any{"type": "integer", "minimum": 0, "description": f.Usage}
			return
		}
		for _, p := range op.Flags {
			if ops.FlagName(p.Name) == f.Name {
				props[f.Name] = query[p.Name]
				if p.Required {
					required = append(required, f.Name)
				}
				return
			}
		}
		if op.Body != nil {
			for _, p := range op.Body.Fields {
				if ops.FlagName(p.Name) == f.Name {
					props[f.Name] = body[p.Name]
					return
				}
			}
		}
	})
	s := map[string]any{"type": "object", "properties": props}
	if len(required) > 0 {
		s["required"] = required
	}
	return s
}

func properties(s map[string]any) map[string]any {
	p, _ := s["properties"].(map[string]any)
	return p
}
