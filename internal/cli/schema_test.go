package cli

import (
	"encoding/json"
	"maps"
	"net/http"
	"slices"
	"testing"

	"github.com/cego/aikido-dojo/internal/clierr"
)

func schemaOf(t *testing.T, args ...string) schemaDoc {
	t.Helper()
	env, _ := testEnv(t)
	var doc schemaDoc
	if err := json.Unmarshal([]byte(mustRun(t, env, append([]string{"schema"}, args...)...)), &doc); err != nil {
		t.Fatal(err)
	}
	return doc
}

func TestSchemaKeysFlagsByFlagName(t *testing.T) {
	doc := schemaOf(t, "issue-group", "list")
	flags := properties(doc.Flags)
	if flags["filter-code-repo-id"] == nil || flags["filter_code_repo_id"] != nil || flags["limit"] == nil {
		t.Errorf("flags = %q, want --filter-code-repo-id and --limit by flag name", slices.Sorted(maps.Keys(flags)))
	}
	if doc.Command != "issue-group list" || doc.Usage != "aikido-dojo issue-group list [flags]" || doc.Scope != "issues:read" || doc.Method != http.MethodGet {
		t.Errorf("doc = %+v", doc)
	}
	// A list prints one array of every item, not a page.
	if doc.Response["type"] != "array" {
		t.Errorf("response type = %v, want array", doc.Response["type"])
	}
}

func TestSchemaOfAnEnvelopeListIsItsItems(t *testing.T) {
	doc := schemaOf(t, "cloud-asset", "list")
	if doc.Response["type"] != "array" || properties(doc.Flags)["region"] == nil {
		t.Errorf("response = %v, flags = %v; want the assets array and the region filter", doc.Response["type"], doc.Flags)
	}
}

func TestSchemaHasBodyFieldFlagsAndPositionalArgs(t *testing.T) {
	doc := schemaOf(t, "issue-group-note", "create")
	if properties(doc.Flags)["note"] == nil || properties(doc.Body)["note"] == nil || properties(doc.Args)["issue_group_id"] == nil {
		t.Errorf("flags %v, body %v, args %v; want note in flags and body, issue_group_id in args", doc.Flags, doc.Body, doc.Args)
	}
	if doc.Usage != "aikido-dojo issue-group-note create <issue_group_id> [flags]" {
		t.Errorf("usage = %q", doc.Usage)
	}
}

func TestSchemaLeavesCredentialsToTheBody(t *testing.T) {
	doc := schemaOf(t, "code-scanning-token update")
	if properties(doc.Flags)["access-token"] != nil || properties(doc.Body)["access_token"] == nil {
		t.Errorf("flags %v, body %v; want access_token only in the body", doc.Flags, doc.Body)
	}
}

func TestSchemaMistakes(t *testing.T) {
	for _, tt := range []struct {
		args []string
		code string
	}{
		{nil, "usage"},
		{[]string{"nosuch", "thing"}, "unknown_command"},
		{[]string{"repo", "nosuch"}, "unknown_command"},
		{[]string{"repo"}, "usage"},
		{[]string{"version"}, "no_schema"},
	} {
		env, _ := testEnv(t)
		_, stderr, code := run(t, env, append([]string{"schema"}, tt.args...)...)
		if e := errorOf(t, stderr); code != clierr.ExitUsage || e.Code != tt.code || e.Hint == "" {
			t.Errorf("schema %q: exit %d, error %+v; want exit 2 %s with a hint", tt.args, code, e, tt.code)
		}
	}
}
