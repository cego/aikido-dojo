package cli

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/search"
)

func searchFor(t *testing.T, env Env, words ...string) []search.Hit {
	t.Helper()
	var hits []search.Hit
	if err := json.Unmarshal([]byte(mustRun(t, env, append([]string{"search"}, words...)...)), &hits); err != nil {
		t.Fatal(err)
	}
	return hits
}

// Each task's command must rank in the top 3. testEnv has no network, so
// this also shows search needs none.
func TestSearchFindsTheCommandForATask(t *testing.T) {
	tests := map[string]string{
		"list repositories":            "repo list",
		"ignore a finding":             "issue ignore",
		"export issues":                "issue export",
		"rotate the firewall token":    "firewall-app-token rotate",
		"add a note to an issue group": "issue-group-note create",
		"cloud assets":                 "cloud-asset list",
		"add a user to a team":         "team-member create",
		"snooze an issue":              "issue snooze",
		"list containers":              "container list",
		"delete a webhook":             "webhook delete",
	}
	env, _ := testEnv(t)
	for query, want := range tests {
		var top []string
		for _, h := range searchFor(t, env, query) {
			top = append(top, h.Command)
		}
		if !slices.Contains(top[:min(3, len(top))], want) {
			t.Errorf("search %q: top 3 of %q lacks %s", query, top, want)
		}
	}
}

func TestSearchJoinsItsArguments(t *testing.T) {
	env, _ := testEnv(t)
	quoted, words := searchFor(t, env, "list repositories"), searchFor(t, env, "list", "repositories")
	if !slices.Equal(quoted, words) || len(quoted) != searchHits {
		t.Errorf("quoted %v, words %v; want the same %d hits", quoted, words, searchHits)
	}
}

func TestSearchWithoutAMatchPrintsAnEmptyArray(t *testing.T) {
	env, _ := testEnv(t)
	if out := mustRun(t, env, "search", "zzqqxx"); out != "[]\n" {
		t.Errorf("stdout = %q, want []", out)
	}
}

func TestSearchMistakes(t *testing.T) {
	for _, tt := range []struct {
		args []string
		code string
	}{
		{[]string{"search"}, "usage"},
		{[]string{"search", "the of a"}, "invalid_input"},
	} {
		env, _ := testEnv(t)
		_, stderr, code := run(t, env, tt.args...)
		if e := errorOf(t, stderr); code != clierr.ExitUsage || e.Code != tt.code {
			t.Errorf("%q: exit %d, error %+v; want exit 2 %s", tt.args, code, e, tt.code)
		}
	}
}
