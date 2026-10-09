package cli

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/cego/aikido-dojo/internal/catalog"
	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/config"
	"github.com/cego/aikido-dojo/internal/ops"
)

// complete asks cobra's __complete as a shell does, and returns the candidates without descriptions, and the directive.
func complete(t *testing.T, env Env, args ...string) ([]string, cobra.ShellCompDirective) {
	t.Helper()
	stdout, stderr, code := run(t, env, append([]string{cobra.ShellCompRequestCmd}, args...)...)
	lines := strings.Split(strings.TrimSuffix(stdout, "\n"), "\n")
	last, ok := strings.CutPrefix(lines[len(lines)-1], ":")
	d, err := strconv.Atoi(last)
	if code != clierr.ExitOK || !ok || err != nil {
		t.Fatalf("%q: exit %d, stdout %q, stderr %q; want candidates and a directive", args, code, stdout, stderr)
	}
	var got []string
	for _, line := range lines[:len(lines)-1] {
		value, _, _ := strings.Cut(line, "\t")
		got = append(got, value)
	}
	return got, cobra.ShellCompDirective(d)
}

func TestCompleteProfilesFromTheConfigFile(t *testing.T) {
	env, vars := testEnv(t)
	if got, d := complete(t, env, "--profile", ""); len(got) != 0 || d != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("no config file: %q, directive %d; want nothing and no file names", got, d)
	}
	writeConfig(t, vars, `{"default_profile":"prod","profiles":{"prod":{"client_id":"a"},"dev":{"client_id":"b"}}}`)
	for _, args := range [][]string{{"--profile", ""}, {"repo", "list", "--profile", ""}} {
		if got, d := complete(t, env, args...); !slices.Equal(got, []string{"dev", "prod"}) || d != cobra.ShellCompDirectiveNoFileComp {
			t.Errorf("%q: %q, directive %d; want dev and prod, no file names", args, got, d)
		}
	}
	other := filepath.Join(t.TempDir(), "other.json")
	if err := os.WriteFile(other, []byte(`{"profiles":{"ci":{"client_id":"c"}}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if got, _ := complete(t, env, "--config", other, "--profile", ""); !slices.Equal(got, []string{"ci"}) {
		t.Errorf("--config %s: %q, want the profiles of that file", other, got)
	}
	for _, broken := range []string{`{`, `{"profiles":{"x":{"client_id":"a","client_secret":"s"}}}`} {
		writeConfig(t, vars, broken)
		if got, d := complete(t, env, "--profile", ""); len(got) != 0 || d != cobra.ShellCompDirectiveNoFileComp {
			t.Errorf("config %s: %q, directive %d; want nothing and no file names", broken, got, d)
		}
	}
}

func TestCompleteRegionsFromTheRegionTable(t *testing.T) {
	env, _ := testEnv(t)
	got, d := complete(t, env, "auth", "login", "--region", "")
	if !slices.Equal(got, []string{"au", "eu", "me", "us"}) || d != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("%q, directive %d; want au, eu, me and us, no file names", got, d)
	}
	for _, region := range got {
		if _, err := config.RegionHost(region); err != nil {
			t.Errorf("%s: %v", region, err)
		}
	}
}

func TestCompleteEnumValues(t *testing.T) {
	env, _ := testEnv(t)
	for _, tt := range []struct{ args, want []string }{
		{[]string{"repo-license", "export", "1", "--format", ""}, []string{"csv", "sbom", "sbom_spdx"}},
		{[]string{"cloud-asset", "list", "--provider", ""}, []string{"aws", "azure", "gcp", "alibaba", "oci", "supabase", "render"}},
		{[]string{"issue-group-severity", "update", "1", "--adjusted-severity", ""}, []string{"critical", "high", "medium", "low"}},
		{[]string{"pr-check-config", "update", "--post-code-quality-inline-comments-min-severity", ""}, []string{"low", "medium", "high", "critical"}},
		{[]string{"endpoint-exception", "list", ""}, []string{"npm", "pypi", "vscode", "open_vsx", "maven", "nuget", "chrome", "golang", "skills_sh"}},
	} {
		if got, d := complete(t, env, tt.args...); !slices.Equal(got, tt.want) || d != cobra.ShellCompDirectiveNoFileComp {
			t.Errorf("%q: %q, directive %d; want %q, no file names", tt.args, got, d, tt.want)
		}
	}
}

func TestEveryEnumCompletesTheSpecsValues(t *testing.T) {
	env, _ := testEnv(t)
	root, err := buildRoot(env, nil)
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	for _, op := range catalog.All {
		sc, err := catalog.Schemas(op.ID)
		if err != nil {
			t.Fatal(err)
		}
		cmd, _, err := root.Find(strings.Fields(op.Command))
		if err != nil {
			t.Fatal(err)
		}
		check := func(args []string, schema any) {
			want := specEnum(schema)
			if want == nil {
				return
			}
			n++
			if got, d := complete(t, env, args...); !slices.Equal(got, want) || d != cobra.ShellCompDirectiveNoFileComp {
				t.Errorf("%q: %q, directive %d; want %q, no file names", args, got, d, want)
			}
		}
		words := strings.Fields(op.Command)
		cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
			check(slices.Concat(words, []string{"--" + f.Name, ""}), flagProperty(op, sc, f.Name))
		})
		for i, p := range op.Args {
			check(slices.Concat(words, slices.Repeat([]string{"1"}, i), []string{""}), properties(sc.Args)[p.Name])
		}
	}
	if n < 60 {
		t.Errorf("checked %d enums, want the spec's 60 or more", n)
	}
}

func TestFileFlagsCompleteFileNames(t *testing.T) {
	env, _ := testEnv(t)
	for _, args := range [][]string{
		{"--config", ""},
		{"issue-group-note", "create", "12", "--body-file", ""},
		{"api", "POST", "/issues/groups/12/notes", "--input", ""},
		{"repo-exclude-path", "create", "1", "--path", ""},
	} {
		if got, d := complete(t, env, args...); len(got) != 0 || d != cobra.ShellCompDirectiveDefault {
			t.Errorf("%q: %q, directive %d; want the shell's file names", args, got, d)
		}
	}
}

func TestOtherFlagsCompleteNothing(t *testing.T) {
	env, _ := testEnv(t)
	root, err := buildRoot(env, nil)
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	var walk func(*cobra.Command)
	walk = func(cmd *cobra.Command) {
		words := strings.Fields(cmd.CommandPath())[1:]
		op, generated := catalog.Command(strings.Join(words, " "))
		var sc ops.SchemaSet
		if generated {
			if sc, err = catalog.Schemas(op.ID); err != nil {
				t.Fatal(err)
			}
		}
		cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
			if f.NoOptDefVal != "" || fileFlag(cmd, f.Name) || specEnum(flagProperty(op, sc, f.Name)) != nil ||
				cmd == root && f.Name == "profile" || cmd.CommandPath() == "aikido-dojo auth login" && f.Name == "region" {
				return
			}
			n++
			args := slices.Concat(words, []string{"--" + f.Name, ""})
			if got, d := complete(t, env, args...); len(got) != 0 || d != cobra.ShellCompDirectiveNoFileComp {
				t.Errorf("%q: %q, directive %d; want nothing and no file names", args, got, d)
			}
		})
		for _, sub := range cmd.Commands() {
			walk(sub)
		}
	}
	walk(root)
	if n < 250 {
		t.Errorf("checked %d flags, want the tree's 250 or more", n)
	}
}

func TestArgumentsCompleteNothing(t *testing.T) {
	env, _ := testEnv(t)
	for _, args := range [][]string{
		{"repo", "get", ""},
		{"issue-group-note", "create", ""},
		{"search", ""},
		{"search", "ignore", ""},
		{"version", ""},
		{"repo", "current", ""},
		{"auth", "status", ""},
	} {
		if got, d := complete(t, env, args...); len(got) != 0 || d != cobra.ShellCompDirectiveNoFileComp {
			t.Errorf("%q: %q, directive %d; want nothing and no file names", args, got, d)
		}
	}
}

func TestCompleteSchemaCommands(t *testing.T) {
	env, _ := testEnv(t)
	resources, d := complete(t, env, "schema", "")
	for _, r := range []string{"repo", "issue", "issue-group", "cloud-asset"} {
		if !slices.Contains(resources, r) {
			t.Errorf("schema: %q lacks %s", resources, r)
		}
	}
	for _, builtin := range []string{"auth", "version", "search", "schema", "api"} {
		if slices.Contains(resources, builtin) {
			t.Errorf("schema: %q offers %s, which has no schema", resources, builtin)
		}
	}
	if d != cobra.ShellCompDirectiveNoFileComp || !slices.IsSorted(resources) || len(slices.Compact(slices.Clone(resources))) != len(resources) {
		t.Errorf("schema: %q, directive %d; want sorted, unique, no file names", resources, d)
	}
	var want []string
	for _, op := range catalog.All {
		if verb, ok := strings.CutPrefix(op.Command, "issue-group "); ok {
			want = append(want, verb)
		}
	}
	slices.Sort(want)
	if got, d := complete(t, env, "schema", "issue-group", ""); !slices.Equal(got, want) || d != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("schema issue-group: %q, directive %d; want %q", got, d, want)
	}
	if got, _ := complete(t, env, "schema", "repo", ""); !slices.Contains(got, "list") || slices.Contains(got, "current") {
		t.Errorf("schema repo: %q; want the generated verbs, not current", got)
	}
	if got, d := complete(t, env, "schema", "issue-group", "list", ""); len(got) != 0 || d != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("schema issue-group list: %q, directive %d; want nothing", got, d)
	}
}

func TestCompleteAPIMethodsAndPaths(t *testing.T) {
	env, _ := testEnv(t)
	if got, d := complete(t, env, "api", ""); !slices.Equal(got, slices.Sorted(slices.Values(apiMethods))) || d != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("api: %q, directive %d; want %q", got, d, apiMethods)
	}
	gets := map[string]bool{}
	for _, op := range catalog.All {
		if op.Method == http.MethodGet {
			gets[op.Path] = true
		}
	}
	got, d := complete(t, env, "api", "get", "")
	if want := slices.Sorted(maps.Keys(gets)); !slices.Equal(got, want) || d != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("api get: %q, directive %d; want every GET path once, sorted: %q", got, d, want)
	}
	if !slices.Contains(got, "/repositories/code") || !slices.Contains(got, "/repositories/code/{code_repo_id}") || slices.Contains(got, "/repositories/code/{code_repo_id}/scan") {
		t.Errorf("api get: %q; want /repositories/code and its template, not the POST-only scan", got)
	}
	if got, d := complete(t, env, "api", "GET", "/repositories/code", ""); len(got) != 0 || d != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("api GET /repositories/code: %q, directive %d; want nothing", got, d)
	}
}

func TestCompletionCallsNothing(t *testing.T) {
	env, vars := testEnv(t)
	writeConfig(t, vars, `{"default_profile":"prod","profiles":{"prod":{"client_id":"a"}}}`)
	env.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		t.Errorf("completion sent %s %s", r.Method, r.URL.Redacted())
		return nil, errors.New("completion calls nothing")
	})
	env.CacheDir = func() (string, error) {
		t.Error("completion looked for the cache directory")
		return "", errors.New("completion caches nothing")
	}
	env.ReadSecret = func(context.Context, string) (string, error) {
		t.Error("completion asked for a secret")
		return "", errors.New("completion needs no secret")
	}
	for _, args := range [][]string{
		{""},
		{"repo", ""},
		{"repo", "list", "--"},
		{"--profile", ""},
		{"auth", "login", "--region", ""},
		{"issue", "export", "--format", ""},
		{"issue-group-note", "create", "12", "--body-file", ""},
		{"repo", "get", ""},
		{"endpoint-exception", "list", ""},
		{"schema", "issue-group", ""},
		{"api", "GET", ""},
		{"repo", "current", ""},
		{"auth", "status", ""},
	} {
		complete(t, env, args...)
	}
}

// fileFlag reports whether cmd's --name reads a file, or names one in the checkout.
func fileFlag(cmd *cobra.Command, name string) bool {
	switch name {
	case "config":
		return !cmd.HasParent()
	case "input":
		return cmd.Name() == "api"
	case "body-file", "path":
		return cmd.Annotations["operation"] != ""
	}
	return false
}

// flagProperty is a generated flag's schema: its query parameter's, else its body field's.
func flagProperty(op ops.Op, sc ops.SchemaSet, flag string) any {
	for _, p := range op.Flags {
		if ops.FlagName(p.Name) == flag {
			return properties(sc.Flags)[p.Name]
		}
	}
	if op.Body != nil {
		for _, p := range op.Body.Fields {
			if ops.FlagName(p.Name) == flag {
				return properties(sc.Body)[p.Name]
			}
		}
	}
	return nil
}

// specEnum is what schema's enum, or a list's items' enum, allows a shell to type: null is left out.
func specEnum(schema any) []string {
	s, _ := schema.(map[string]any)
	if s["type"] == "array" {
		s, _ = s["items"].(map[string]any)
	}
	enum, _ := s["enum"].([]any)
	var values []string
	for _, v := range enum {
		if v != nil {
			values = append(values, fmt.Sprint(v))
		}
	}
	return values
}
