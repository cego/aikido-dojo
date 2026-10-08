// Package cli builds the aikido-dojo command tree from the generated catalog
// and runs one command line.
package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cego/aikido-dojo/internal/clierr"
)

// Env is everything a run touches outside the process, so tests can supply their own.
type Env struct {
	Stdin      io.Reader
	Stdout     io.Writer
	Stderr     io.Writer
	Getenv     func(string) string
	CacheDir   func() (string, error)
	Transport  http.RoundTripper                                        // nil in production; tests route requests to a fake server
	ReadSecret func(ctx context.Context, prompt string) (string, error) // the hidden prompt auth login asks with
	StdoutTTY  bool                                                     // stdout is a terminal, where a person reads indented JSON
}

// app holds a run's environment and the global flags' values.
type app struct {
	env     Env
	out     *output
	config  string
	profile string
	debug   bool
}

// Run executes one command line and returns the process exit code.
func Run(ctx context.Context, args []string, env Env) int {
	root, err := buildRoot(env) //nolint:contextcheck // commands get ctx from root.ExecuteContext below
	if err != nil {
		return clierr.Report(env.Stderr, err)
	}
	root.SetArgs(args)
	if err := root.ExecuteContext(ctx); err != nil {
		return clierr.Report(env.Stderr, err)
	}
	return clierr.ExitOK
}

func newRoot(env Env) (*cobra.Command, *app) {
	a := &app{env: env, out: &output{w: env.Stdout, pretty: env.StdoutTTY}}
	root := &cobra.Command{
		Use:   "aikido-dojo",
		Short: "An unofficial command-line client for the Aikido Security public API",
		Long: "aikido-dojo calls the Aikido Security public API. It is unofficial and not affiliated with Aikido Security.\n\n" +
			"Commands are <resource> <verb>, for example: aikido-dojo repo list",
		Args:          unknownCommand,
		RunE:          showHelp,
		SilenceErrors: true,
		SilenceUsage:  true,
		// cobra applies its default of 2 only on its own unknown-command path, which unknownCommand replaces.
		SuggestionsMinimumDistance: 2,
	}
	root.SetIn(env.Stdin)
	root.SetOut(env.Stdout)
	root.SetErr(env.Stderr)
	root.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		return &clierr.Error{Code: "usage", Message: err.Error(), Hint: "run " + cmd.CommandPath() + " --help", Exit: clierr.ExitUsage}
	})
	// Only the root's help gets the footer: it is where an agent starts.
	root.SetUsageTemplate(root.UsageTemplate() + `{{if not .HasParent}}
Find the commands for a task:  aikido-dojo search "<what you want to do>"
See what a command takes:      aikido-dojo schema <resource> <verb>
{{end}}`)
	pf := root.PersistentFlags()
	pf.StringVar(&a.config, "config", "", "the config file (default ~/.config/aikido-dojo/config.json; also AIKIDO_DOJO_CONFIG)")
	pf.StringVar(&a.profile, "profile", "", "the profile to run as (also AIKIDO_DOJO_PROFILE)")
	pf.BoolVar(&a.debug, "debug", false, "log each request's method, URL, status and timing to stderr")
	return root, a
}

// unknownCommand checks the arguments of every command that only groups
// others. Without it, a mistyped verb prints help and exits 0, and a mistyped
// resource exits 1; an agent needs a usage error.
func unknownCommand(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	}
	hint := "run " + cmd.CommandPath() + " --help"
	if s := cmd.SuggestionsFor(args[0]); len(s) > 0 {
		hint = "did you mean " + strings.Join(s, " or ") + "?"
	}
	return &clierr.Error{Code: "unknown_command", Message: fmt.Sprintf("unknown command %q for %s", args[0], cmd.CommandPath()),
		Hint: hint, Exit: clierr.ExitUsage}
}

func showHelp(cmd *cobra.Command, _ []string) error {
	if err := cmd.Help(); err != nil {
		return fmt.Errorf("print help: %w", err)
	}
	return nil
}
