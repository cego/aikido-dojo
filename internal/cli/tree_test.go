package cli

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/cego/aikido-dojo/internal/api"
	"github.com/cego/aikido-dojo/internal/catalog"
	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/ops"
)

var update = flag.Bool("update", false, "rewrite testdata/tree.golden")

func nopRun(*cobra.Command, command, []string) error { return nil }

// treeWith builds the root with the given operations under it.
func treeWith(t *testing.T, all []ops.Op, run runFunc) *cobra.Command {
	t.Helper()
	env, _ := testEnv(t)
	root, _ := newRoot(env)
	if err := addGenerated(root, all, run); err != nil {
		t.Fatal(err)
	}
	return root
}

// execute runs args against root and returns the error and everything printed.
func execute(root *cobra.Command, args ...string) (string, error) {
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs(args)
	err := root.Execute()
	return out.String(), err
}

func TestTreeGolden(t *testing.T) {
	got := renderTree(treeWith(t, catalog.All, nopRun))
	golden := filepath.Join("testdata", "tree.golden")
	if *update {
		if err := os.WriteFile(golden, []byte(got), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	if got != string(want) {
		t.Errorf("the command tree changed; review it and run go test ./internal/cli -run TestTreeGolden -update\n%s", firstDiff(string(want), got))
	}
}

var endNames = map[api.End]string{api.EndEmpty: "empty", api.EndHeader: "header", api.EndField: "field"}

// renderTree lists every generated command with what it calls, its scope,
// paging and flags, so a spec or overlay change is a reviewable diff.
func renderTree(root *cobra.Command) string {
	byID := map[string]ops.Op{}
	for _, op := range catalog.All {
		byID[op.ID] = op
	}
	var b strings.Builder
	var walk func(*cobra.Command)
	walk = func(c *cobra.Command) {
		for _, sub := range c.Commands() {
			if op, ok := byID[sub.Annotations["operation"]]; ok {
				renderCommand(&b, sub, op)
			}
			walk(sub)
		}
	}
	walk(root)
	return b.String()
}

func renderCommand(b *strings.Builder, cmd *cobra.Command, op ops.Op) {
	fmt.Fprintf(b, "%s\n\t%s %s", strings.TrimPrefix(cmd.CommandPath(), "aikido-dojo ")+strings.TrimPrefix(cmd.Use, cmd.Name()), op.Method, op.Path)
	if op.Scope != "" {
		fmt.Fprintf(b, " scope=%s", op.Scope)
	}
	if p := op.Paging; p != nil {
		fmt.Fprintf(b, " paged=%s:%d:%s", p.SizeParam, p.Size, endNames[p.End])
		if p.Items != "" {
			fmt.Fprintf(b, ":items=%s", p.Items)
		}
		if p.More != "" {
			fmt.Fprintf(b, ":more=%s", p.More)
		}
	}
	if op.Destructive {
		b.WriteString(" destructive")
	}
	b.WriteString("\n")
	required, body := map[string]bool{}, map[string]bool{}
	for _, p := range op.Flags {
		required[ops.FlagName(p.Name)] = p.Required
	}
	if op.Body != nil {
		for _, p := range op.Body.Fields {
			required[ops.FlagName(p.Name)], body[ops.FlagName(p.Name)] = p.Required, true
		}
	}
	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		fmt.Fprintf(b, "\t--%s %s", f.Name, f.Value.Type())
		if required[f.Name] {
			b.WriteString(" required")
		}
		if body[f.Name] {
			b.WriteString(" body")
		}
		b.WriteString("\n")
	})
}

func firstDiff(want, got string) string {
	w, g := strings.Split(want, "\n"), strings.Split(got, "\n")
	for i := 0; i < len(w) || i < len(g); i++ {
		var wl, gl string
		if i < len(w) {
			wl = w[i]
		}
		if i < len(g) {
			gl = g[i]
		}
		if wl != gl {
			return fmt.Sprintf("line %d:\n  want %q\n  got  %q", i+1, wl, gl)
		}
	}
	return ""
}

func TestEveryOperationIsACommand(t *testing.T) {
	n := 0
	var walk func(*cobra.Command)
	walk = func(c *cobra.Command) {
		for _, sub := range c.Commands() {
			if sub.Annotations["operation"] != "" {
				n++
			}
			walk(sub)
		}
	}
	walk(treeWith(t, catalog.All, nopRun))
	if n != len(catalog.All) {
		t.Errorf("commands = %d, want %d", n, len(catalog.All))
	}
}

func TestResourceShortListsItsVerbs(t *testing.T) {
	cmd, _, err := treeWith(t, catalog.All, nopRun).Find([]string{"repo"})
	if err != nil || cmd.Short != "activate, clone, deactivate, delete, get, import, list, scan" {
		t.Errorf("repo = %q, %v", cmd.Short, err)
	}
}

func TestHelpNamesTheScopeAndTheOverlayNote(t *testing.T) {
	out, err := execute(treeWith(t, catalog.All, nopRun), "issue-group", "ignore", "--help")
	if err != nil || !strings.Contains(out, "Needs the issues:write scope.") || !strings.Contains(out, "Aikido discards --reason") {
		t.Errorf("help = %q, %v", out, err)
	}
}

func TestArgumentMistakes(t *testing.T) {
	tests := []struct {
		args       []string
		code, text string
	}{
		{[]string{"repo", "get"}, "usage", "takes <code_repo_id>, got 0"},
		{[]string{"workspace", "get", "extra"}, "usage", "takes no arguments, got 1"},
		{[]string{"repo", "lis"}, "unknown_command", "did you mean list?"},
		{[]string{"repoo"}, "unknown_command", "did you mean repo"},
	}
	for _, tt := range tests {
		_, err := execute(treeWith(t, catalog.All, nopRun), tt.args...)
		var e *clierr.Error
		if !errors.As(err, &e) || e.Code != tt.code || e.Exit != clierr.ExitUsage || !strings.Contains(e.Message+" "+e.Hint, tt.text) {
			t.Errorf("%q: err = %v, want %s containing %q", tt.args, err, tt.code, tt.text)
		}
	}
}

func TestBodyFieldClashingWithAGlobalFlagIsBodyOnly(t *testing.T) {
	op := ops.Op{ID: "x", Command: "thing update", Method: "PUT", Path: "/things", Body: &ops.Body{Object: true,
		Fields: []ops.Param{{Name: "name", Kind: ops.String}, {Name: "profile", Kind: ops.String}}}}
	var got command
	root := treeWith(t, []ops.Op{op}, func(_ *cobra.Command, c command, _ []string) error { got = c; return nil })
	cmd, _, err := root.Find([]string{"thing", "update"})
	if err != nil {
		t.Fatal(err)
	}
	if cmd.LocalFlags().Lookup("profile") != nil || cmd.LocalFlags().Lookup("name") == nil || !strings.Contains(cmd.Long, "profile share a name with a built-in flag") {
		t.Errorf("flags or help wrong: %q", cmd.Long)
	}
	if _, err := execute(root, "thing", "update", "--name", "n"); err != nil {
		t.Fatal(err)
	}
	if len(got.bodyFlags) != 1 || got.bodyFlags[0].Name != "name" {
		t.Errorf("body flags = %+v, want only name", got.bodyFlags)
	}
}

func TestQueryFlagClashingWithAGlobalFlagFailsTheBuild(t *testing.T) {
	env, _ := testEnv(t)
	root, _ := newRoot(env)
	op := ops.Op{ID: "x", Command: "thing list", Method: "GET", Path: "/things", Flags: []ops.Param{{Name: "profile", Kind: ops.String}}}
	if err := addGenerated(root, []ops.Op{op}, nopRun); err == nil || !strings.Contains(err.Error(), "--profile clashes") {
		t.Errorf("err = %v, want the clash reported", err)
	}
}
