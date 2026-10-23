package cli

import (
	"encoding/json"
	"maps"
	"net/http"
	"slices"
	"sync"
	"testing"

	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/helpdocs"
	"github.com/cego/aikido-dojo/internal/helpdocs/helpdocstest"
)

// Each question's page must rank in the top 3. testEnv has no network, so this also shows docs search needs none.
func TestDocsSearchFindsThePageForAQuestion(t *testing.T) {
	tests := map[string]string{
		"change the app language":            "getting-started/setting-up-your-account/change-language",
		"limit access to specific IPs":       "getting-started/setting-up-your-account/limit-aikido-access-to-specific-ips",
		"invite users without a git account": "getting-started/automated-user-management/invite-users-to-aikido-without-a-git-account",
		"connect github enterprise":          "code-scanning/connect-your-source-code/connect-github-enterprise",
		"jfrog artifactory private registry": "miscellaneous-integrations/private-registries/jfrog-artifactory",
	}
	env, _ := testEnv(t)
	for query, want := range tests {
		var hits []helpdocs.Hit
		if err := json.Unmarshal([]byte(mustRun(t, env, "docs", "search", query)), &hits); err != nil {
			t.Fatal(err)
		}
		var top []string
		for _, h := range hits[:min(3, len(hits))] {
			top = append(top, h.Path)
		}
		if !slices.Contains(top, want) {
			t.Errorf("docs search %q: top 3 %q lacks %s", query, top, want)
		}
	}
}

func TestDocsGetPrintsThePage(t *testing.T) {
	env, _ := testEnv(t)
	for _, ref := range []string{
		"getting-started/setting-up-your-account/change-language",
		"https://help.aikido.dev/getting-started/setting-up-your-account/change-language",
	} {
		if out := mustRun(t, env, "docs", "get", ref, "--jq", ".title"); out != "Change Language\n" {
			t.Errorf("docs get %s --jq .title = %q", ref, out)
		}
	}
}

func TestDocsSearchWithoutAMatchPrintsAnEmptyArray(t *testing.T) {
	env, _ := testEnv(t)
	if out := mustRun(t, env, "docs", "search", "zzqqxx"); out != "[]\n" {
		t.Errorf("stdout = %q, want []", out)
	}
}

func TestDocsMistakes(t *testing.T) {
	for _, tt := range []struct {
		args []string
		exit int
		code string
	}{
		{[]string{"docs", "search"}, clierr.ExitUsage, "usage"},
		{[]string{"docs", "search", "the of a"}, clierr.ExitUsage, "invalid_input"},
		{[]string{"docs", "get"}, clierr.ExitUsage, "usage"},
		{[]string{"docs", "get", "a", "b"}, clierr.ExitUsage, "usage"},
		{[]string{"docs", "get", "no/such/page"}, clierr.ExitNotFound, "not_found"},
		{[]string{"docs", "get", "../../go.mod"}, clierr.ExitNotFound, "not_found"},
		{[]string{"docs", "get", "https://evil.example/readme"}, clierr.ExitNotFound, "not_found"},
		{[]string{"docs", "fetch"}, clierr.ExitUsage, "unknown_command"},
	} {
		env, _ := testEnv(t)
		stdout, stderr, code := run(t, env, tt.args...)
		if e := errorOf(t, stderr); code != tt.exit || e.Code != tt.code || stdout != "" {
			t.Errorf("%q: exit %d, error %+v, stdout %q; want exit %d %s", tt.args, code, e, stdout, tt.exit, tt.code)
		}
	}
}

const languagePage = "getting-started/setting-up-your-account/change-language"

func TestDocsSyncFetchesWhatAikidoChanged(t *testing.T) {
	edited := "# Change Language\n\nPick a language under zebracorn settings.\n"
	help := helpdocstest.Handler("../helpdocs", map[string]string{languagePage: edited})
	var (
		mu     sync.Mutex
		agents = map[string]bool{}
	)
	env, _ := testEnv(t)
	env.Transport = helpdocstest.Transport(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		agents[r.UserAgent()] = true
		mu.Unlock()
		help.ServeHTTP(w, r)
	}))
	var r helpdocs.SyncReport
	// It changes nothing in Aikido, so read-only mode allows it.
	if err := json.Unmarshal([]byte(mustRun(t, env, "--read-only", "docs", "sync")), &r); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(r.Changed, []string{languagePage}) || r.Snapshot != "2099-01-01" || r.Requests == 0 {
		t.Errorf("report = %+v, want %s changed at 2099-01-01", r, languagePage)
	}
	if !maps.Equal(agents, map[string]bool{userAgent: true}) {
		t.Errorf("user agents = %v, want only %s", agents, userAgent)
	}
	env.Transport = noNetwork{}
	if out := mustRun(t, env, "docs", "get", languagePage, "--jq", ".markdown"); out != edited+"\n" {
		t.Errorf("docs get = %q, want the synced page", out)
	}
	if out := mustRun(t, env, "docs", "search", "zebracorn", "--jq", ".[0].path"); out != languagePage+"\n" {
		t.Errorf("docs search zebracorn = %q, want %s", out, languagePage)
	}
}

func TestDocsSyncFailureKeepsTheCopy(t *testing.T) {
	env, _ := testEnv(t)
	stdout, stderr, code := run(t, env, "docs", "sync")
	if e := errorOf(t, stderr); code != clierr.ExitUnexpected || e.Code != "docs_sync_failed" || e.Hint == "" || stdout != "" {
		t.Errorf("exit %d, error %+v, stdout %q; want exit 1 docs_sync_failed with a hint", code, e, stdout)
	}
	if out := mustRun(t, env, "docs", "get", languagePage, "--jq", ".title"); out != "Change Language\n" {
		t.Errorf("docs get after a failed sync = %q", out)
	}
}
