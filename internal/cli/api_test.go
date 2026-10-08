package cli

import (
	"net/http"
	"strings"
	"testing"

	"github.com/cego/aikido-dojo/internal/clierr"
)

func TestAPISendsTheRequest(t *testing.T) {
	f, env, _ := newFake(t, respond(`{"ok":true}`))
	out := mustRun(t, env, "api", "get", "repositories/code?per_page=5", "-f", "filter_name=a,b", "-f", "x=1", "-f", "x=2")
	s := f.first(t)
	if out != `{"ok":true}` || s.Method != http.MethodGet || s.Body != "" || !strings.HasPrefix(s.URI, "/api/public/v1/repositories/code?") {
		t.Errorf("stdout %q, request %+v", out, s)
	}
	if s.Query.Get("per_page") != "5" || s.Query.Get("filter_name") != "a,b" || strings.Join(s.Query["x"], " ") != "1 2" {
		t.Errorf("query = %v, want the path's query, the comma kept, and both x values", s.Query)
	}
}

// An escaped slash inside one segment is an ID, not a way out of the prefix.
func TestAPIKeepsAnEscapedSlashInAnID(t *testing.T) {
	f, env, _ := newFake(t, respond(`{}`))
	mustRun(t, env, "api", "GET", "/containers/org%2Fimage")
	if got := f.first(t).URI; got != "/api/public/v1/containers/org%2Fimage" {
		t.Errorf("URI = %q, want the escaped slash kept", got)
	}
}

func TestAPIFieldsBecomeABodyForWrites(t *testing.T) {
	f, env, _ := newFake(t, respond(`{}`))
	mustRun(t, env, "api", "POST", "/issues/groups/12/notes", "-f", "note=<b>&ok", "-f", "cve_id=CVE-1")
	if got := f.first(t).Body; got != `{"cve_id":"CVE-1","note":"<b>&ok"}` {
		t.Errorf("body = %s", got)
	}
}

func TestAPIInputIsTheBody(t *testing.T) {
	f, env, _ := newFake(t, respond(`{}`))
	env.Stdin = strings.NewReader(`{"reason": "x"}`)
	mustRun(t, env, "api", "PUT", "/issues/groups/12/ignore", "--input", "-", "-f", "dry=1")
	if s := f.first(t); s.Body != `{"reason": "x"}` || s.Query.Get("dry") != "1" {
		t.Errorf("body %q, query %v; want the input as sent and -f in the query", s.Body, s.Query)
	}
}

func TestAPINamesTheScopeOfAMatchingCommand(t *testing.T) {
	_, env, _ := newFake(t, func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusForbidden) })
	_, stderr, code := run(t, env, "api", "PUT", "/issues/groups/12/ignore")
	if e := errorOf(t, stderr); code != clierr.ExitForbidden || e.Code != "missing_scope" || !strings.Contains(e.Hint, "issues:write") {
		t.Errorf("exit %d, error %+v; want missing_scope naming issues:write", code, e)
	}
	if _, stderr, _ := run(t, env, "api", "GET", "/nosuch"); errorOf(t, stderr).Code != "forbidden" {
		t.Errorf("an unknown path: %s, want forbidden", stderr)
	}
}

func TestAPIMistakesCostNoCall(t *testing.T) {
	f, env, _ := newFake(t, respond(`{}`))
	for _, args := range [][]string{
		{"api", "GET"},
		{"api", "PATCH", "/x"},
		{"api", "GET", "https://app.aikido.dev/api/public/v1/x"},
		{"api", "GET", "/../oauth/token"},
		{"api", "GET", "/issues/%2e%2e/x"},
		{"api", "GET", "/issues%2f..%2f..%2foauth%2ftoken"},
		{"api", "GET", "/%2e%2e%2f%2e%2e%2foauth"},
		{"api", "GET", "/..%5c..%5coauth"},
		{"api", "GET", `/x\..\oauth`},
		{"api", "GET", "/x\ny"},
		{"api", "GET", "/issues/%zz"},
		{"api", "GET", "/x?a=%zz"},
		{"api", "POST", "/x", "-f", "novalue"},
		{"api", "POST", "/x", "-f", "=v"},
		{"api", "POST", "/x", "-f", "a=1", "-f", "a=2"},
		{"api", "POST", "/x", "--input", "-"},
	} {
		_, stderr, code := run(t, env, args...)
		if code != clierr.ExitUsage || errorOf(t, stderr).Hint == "" {
			t.Errorf("%q: exit %d, stderr %s; want a usage error with a hint", args, code, stderr)
		}
	}
	if n, logins := len(f.seen()), len(f.seenLogins()); n != 0 || logins != 0 {
		t.Errorf("API calls = %d and token requests = %d, want none", n, logins)
	}
}

// Only the path is checked for a URL of its own; a query value may hold one.
func TestAPIAllowsAURLInTheQuery(t *testing.T) {
	f, env, _ := newFake(t, respond(`{}`))
	mustRun(t, env, "api", "GET", "/repositories/code?next=https://example.com/a")
	if got := f.first(t).Query.Get("next"); got != "https://example.com/a" {
		t.Errorf("next = %q", got)
	}
}
