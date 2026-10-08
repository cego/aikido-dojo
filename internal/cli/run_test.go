package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/cego/aikido-dojo/internal/catalog"
	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/config"
)

// fakeAPI stands in for Aikido: it issues tokens and hands API requests to
// handler, recording each one.
type fakeAPI struct {
	mu       sync.Mutex
	requests []seen
	logins   []string // id:secret of each token request
	token    func(id, secret string) (int, string)
	rewrite  *rewrite
}

type seen struct {
	Method string
	URI    string
	Body   string
	Query  url.Values
	Auth   string
}

func newFake(t *testing.T, handler http.HandlerFunc) (*fakeAPI, Env, map[string]string) {
	t.Helper()
	f := &fakeAPI{}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		id, secret, _ := r.BasicAuth()
		f.mu.Lock()
		f.logins = append(f.logins, id+":"+secret)
		issue := f.token
		f.mu.Unlock()
		status, tok := http.StatusOK, "tok"
		if issue != nil {
			status, tok = issue(id, secret)
		}
		if status != http.StatusOK {
			w.WriteHeader(status)
			fmt.Fprint(w, `{"error":"invalid_client","error_description":"The provided credentials are invalid"}`)
			return
		}
		fmt.Fprintf(w, `{"access_token":%q,"expires_in":3600}`, tok)
	})
	mux.HandleFunc("/api/public/v1/", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		f.mu.Lock()
		f.requests = append(f.requests, seen{Method: r.Method, URI: r.RequestURI, Body: string(body), Query: r.URL.Query(), Auth: r.Header.Get("Authorization")})
		f.mu.Unlock()
		handler(w, r)
	})
	srv := httptest.NewTLSServer(mux)
	t.Cleanup(srv.Close)
	env, vars := testEnv(t)
	vars[config.EnvClientID], vars[config.EnvClientSecret] = "id", "secret"
	f.rewrite = &rewrite{to: strings.TrimPrefix(srv.URL, "https://"), next: srv.Client().Transport}
	env.Transport = f.rewrite
	return f, env, vars
}

func (f *fakeAPI) seen() []seen {
	f.mu.Lock()
	defer f.mu.Unlock()
	return slices.Clone(f.requests)
}

// issueTokens decides each token request's status and token, in place of
// "tok" for everyone.
func (f *fakeAPI) issueTokens(issue func(id, secret string) (int, string)) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.token = issue
}

func (f *fakeAPI) seenLogins() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return slices.Clone(f.logins)
}

// first is the first API request; the test fails if there was none.
func (f *fakeAPI) first(t *testing.T) seen {
	t.Helper()
	s := f.seen()
	if len(s) == 0 {
		t.Fatal("no API request was made")
	}
	return s[0]
}

// respond answers with body as JSON, with the Content-Type Aikido sends (live, 2026-10-08).
func respond(body string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, body)
	}
}

// pages serves pages[n] for ?page=n and a 404 past the end.
func pages(bodies ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var n int
		if _, err := fmt.Sscan(r.URL.Query().Get("page"), &n); err != nil || n >= len(bodies) {
			http.Error(w, `{"error":"no such page"}`, http.StatusNotFound)
			return
		}
		fmt.Fprint(w, bodies[n])
	}
}

func mustRun(t *testing.T, env Env, args ...string) string {
	t.Helper()
	stdout, stderr, code := run(t, env, args...)
	if code != clierr.ExitOK {
		t.Fatalf("%q: exit %d: %s", args, code, stderr)
	}
	return stdout
}

func TestPathArguments(t *testing.T) {
	f, env, _ := newFake(t, respond(`{}`))
	mustRun(t, env, "cve", "get", "a/b?c")
	if got := f.first(t).URI; got != "/api/public/v1/cve/a%2Fb%3Fc" {
		t.Errorf("URI = %q, want the argument escaped into one path segment", got)
	}
	for _, args := range [][]string{{"cve", "get", "."}, {"cve", "get", ".."}, {"repo", "get", "abc"}} {
		_, stderr, code := run(t, env, args...)
		if e := errorOf(t, stderr); code != clierr.ExitUsage || e.Code != "invalid_input" {
			t.Errorf("%q: exit %d, error %+v; want invalid_input", args, code, e)
		}
	}
	if _, stderr, _ := run(t, env, "repo", "get", "abc"); !strings.Contains(errorOf(t, stderr).Message, `<code_repo_id>: "abc" is not an integer`) {
		t.Errorf("error = %s, want the argument named", stderr)
	}
	if n := len(f.seen()); n != 1 {
		t.Errorf("API calls = %d, want only the first command's", n)
	}
	// The spec bounds user_id to 0..1; the overlay drops that wrong bound.
	mustRun(t, env, "user", "get", "5")
	if got := f.seen()[1].URI; got != "/api/public/v1/users/5" {
		t.Errorf("URI = %q, want user 5", got)
	}
}

func TestQueryHoldsOnlyTheFlagsSet(t *testing.T) {
	f, env, _ := newFake(t, respond(`[]`))
	if out := mustRun(t, env, "repo", "list", "--filter-name", "web", "--include-inactive=false"); out != "[]\n" {
		t.Errorf("stdout = %q, want []", out)
	}
	want := url.Values{"filter_name": {"web"}, "include_inactive": {"false"}, "page": {"0"}, "per_page": {"200"}}
	if got := f.first(t).Query; !reflect.DeepEqual(got, want) {
		t.Errorf("query = %v, want %v", got, want)
	}
}

func TestListFlagsAreCommaSeparated(t *testing.T) {
	f, env, _ := newFake(t, respond(`{"assets":[],"hasMore":false}`))
	mustRun(t, env, "cloud-asset", "list", "--cloud-id", "1,2", "--provider", "aws", "--provider", "gcp")
	q := f.first(t).Query
	if q.Get("cloud_id") != "1,2" || q.Get("provider") != "aws,gcp" || q.Get("limit") != "100" {
		t.Errorf("query = %v", q)
	}
}

func TestFlagMistakes(t *testing.T) {
	tests := []struct {
		args       []string
		code, text string
	}{
		{[]string{"cloud-asset", "list", "--cloud-id", "1,x"}, "invalid_input", `--cloud-id: "x" is not an integer`},
		{[]string{"code-quality-pr-finding", "list", "--code-repo-id", "1"}, "invalid_input", "--pr-number: required"},
		{[]string{"repo", "list", "--limit", "-1"}, "invalid_input", "--limit"},
		{[]string{"repo", "list", "--limit", "x"}, "usage", "--limit"},
		{[]string{"code-scanning-token", "update", "--access-token", "x"}, "usage", "unknown flag: --access-token"},
	}
	for _, tt := range tests {
		f, env, _ := newFake(t, respond(`[]`))
		_, stderr, code := run(t, env, tt.args...)
		if e := errorOf(t, stderr); code != clierr.ExitUsage || e.Code != tt.code || !strings.Contains(e.Message, tt.text) {
			t.Errorf("%q: exit %d, error %+v; want %s containing %q", tt.args, code, e, tt.code, tt.text)
		}
		if n := len(f.rewrite.seenHosts()); n != 0 {
			t.Errorf("%q: %d requests, want none before the input is valid", tt.args, n)
		}
	}
}

func TestUnknownEnumValueIsSentWithAWarning(t *testing.T) {
	f, env, _ := newFake(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Has-Next-Page", "false")
		fmt.Fprint(w, `[]`)
	})
	_, stderr, code := run(t, env, "issue-group", "list", "--filter-status", "auto-ignored")
	if code != clierr.ExitOK || !strings.Contains(stderr, `{"warning":{"message":"--filter-status: \"auto-ignored\" is not one of`) {
		t.Errorf("exit %d, stderr %q; want 0 and a warning", code, stderr)
	}
	if got := f.first(t).Query.Get("filter_status"); got != "auto-ignored" {
		t.Errorf("filter_status = %q, want it sent", got)
	}
}

func TestBody(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		stdin string
		want  string
	}{
		{"field flags", []string{"issue-group-note", "create", "5", "--note", "hi"}, "", `{"note":"hi"}`},
		{"field flags win over --body", []string{"issue-group-note", "create", "5", "--body", `{"note":"a","cve_id":"CVE-1"}`, "--note", "b"}, "",
			`{"cve_id":"CVE-1","note":"b"}`},
		{"--body-file - reads stdin", []string{"issue-group-note", "create", "5", "--body-file", "-"}, `{"note":"x"}`, `{"note":"x"}`},
		{"a required object body defaults to {}", []string{"issue-group", "ignore", "5"}, "", `{}`},
		{"an optional body is left out", []string{"issue", "unsnooze", "5"}, "", ``},
		{"an integer field", []string{"team-member", "create", "3", "--user-id", "7"}, "", `{"user_id":7}`},
		{"a boolean field", []string{"repo-dev-dep-scan", "update", "1", "--dev-dep-scanning-enabled"}, "", `{"dev_dep_scanning_enabled":true}`},
		{"a credential through --body-file", []string{"code-scanning-token", "update", "--body-file", "-"}, `{"access_token":"x"}`, `{"access_token":"x"}`},
		{"HTML characters stay as typed", []string{"issue-group-note", "create", "5", "--note", "<b> & co"}, "", `{"note":"<b> & co"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, env, _ := newFake(t, respond(`{}`))
			env.Stdin = strings.NewReader(tt.stdin)
			mustRun(t, env, tt.args...)
			if got := f.first(t).Body; got != tt.want {
				t.Errorf("body = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBodyFile(t *testing.T) {
	f, env, _ := newFake(t, respond(`{}`))
	path := filepath.Join(t.TempDir(), "body.json")
	if err := os.WriteFile(path, []byte(`{"note":"from a file"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	mustRun(t, env, "issue-group-note", "create", "5", "--body-file", path)
	if got := f.first(t).Body; got != `{"note":"from a file"}` {
		t.Errorf("body = %q", got)
	}
}

func TestBodyMistakes(t *testing.T) {
	tests := []struct {
		args []string
		text string
	}{
		{[]string{"issue-group-severity", "update", "5", "--adjusted-severity", "high"}, "body.reason: required"},
		{[]string{"issue-group-note", "create", "5", "--body", "{}", "--body-file", "-"}, "--body or --body-file, not both"},
		{[]string{"issue-group-note", "create", "5", "--body", "{"}, "invalid JSON"},
		{[]string{"issue-group-note", "create", "5", "--body", "{}{}"}, "unexpected data after the JSON value"},
		{[]string{"issue-group-note", "create", "5", "--body", "[1]", "--note", "x"}, "must be a JSON object"},
		{[]string{"issue-group-note", "create", "5", "--body", `{"note":1}`}, "body.note: want string, got number"},
		{[]string{"issue-group-note", "create", "5", "--body-file", "/nonexistent/body.json"}, "--body-file"},
		{[]string{"firewall-bot-list", "update", "1"}, "needs a body"},
		{[]string{"code-scanning-token", "update", "--body", `{"access_token":"x"}`}, "access_token is a credential"},
	}
	for _, tt := range tests {
		f, env, _ := newFake(t, respond(`{}`))
		_, stderr, code := run(t, env, tt.args...)
		if e := errorOf(t, stderr); code != clierr.ExitUsage || e.Code != "invalid_input" || !strings.Contains(e.Message, tt.text) {
			t.Errorf("%q: exit %d, error %+v; want invalid_input containing %q", tt.args, code, e, tt.text)
		}
		if n := len(f.rewrite.seenHosts()); n != 0 {
			t.Errorf("%q: %d requests, want none", tt.args, n)
		}
	}
}

func TestListPagesThroughAndStopsAtTheLimit(t *testing.T) {
	f, env, _ := newFake(t, pages(`[1,2]`, `[3]`, `[]`))
	if out := mustRun(t, env, "repo", "list"); out != "[1,2,3]\n" {
		t.Errorf("stdout = %q", out)
	}
	f2, env2, _ := newFake(t, pages(`[1,2]`, `[3]`, `[]`))
	if out := mustRun(t, env2, "repo", "list", "--limit", "2"); out != "[1,2]\n" {
		t.Errorf("stdout = %q", out)
	}
	if len(f.seen()) != 3 || len(f2.seen()) != 1 {
		t.Errorf("calls = %d and %d, want 3 and 1", len(f.seen()), len(f2.seen()))
	}
}

func TestListFailureLeavesStdoutEmpty(t *testing.T) {
	_, env, _ := newFake(t, pages(`[1,2]`))
	stdout, stderr, code := run(t, env, "repo", "list")
	if stdout != "" || code != clierr.ExitNotFound || errorOf(t, stderr).Code != "not_found" {
		t.Errorf("stdout %q, exit %d, stderr %q; want nothing printed and the page's error", stdout, code, stderr)
	}
}

func TestRawOutputIsCopiedUnchanged(t *testing.T) {
	_, env, _ := newFake(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		fmt.Fprint(w, "a,b\n1,2\n")
	})
	if out := mustRun(t, env, "repo-license", "export", "1", "--format", "csv"); out != "a,b\n1,2\n" {
		t.Errorf("stdout = %q", out)
	}
}

func TestNoContentPrintsNothing(t *testing.T) {
	_, env, _ := newFake(t, func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })
	if out := mustRun(t, env, "repo", "scan", "1"); out != "" {
		t.Errorf("stdout = %q, want nothing", out)
	}
}

func TestRegionPicksTheHost(t *testing.T) {
	f, env, vars := newFake(t, respond(`{}`))
	vars[config.EnvRegion] = "us"
	mustRun(t, env, "workspace", "get")
	hosts := f.rewrite.seenHosts()
	if len(hosts) != 2 || hosts[0] != "app.us.aikido.dev" || hosts[1] != "app.us.aikido.dev" {
		t.Errorf("hosts = %v, want the token and API calls on app.us.aikido.dev", hosts)
	}
}

func TestBadRequestHintFromTheOverlay(t *testing.T) {
	_, env, _ := newFake(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"status_code":400,"reason_phrase":"This feature is not enabled on your workspace. Please contact support to enable it."}`)
	})
	var hint string
	for _, op := range catalog.All {
		if op.Command == "issue list" {
			hint = op.BadRequestHint
		}
	}
	_, stderr, code := run(t, env, "issue", "list", "--issue-ids", "1")
	e := errorOf(t, stderr)
	if code != clierr.ExitUsage || hint == "" || e.Hint != hint || !strings.Contains(e.Message, "This feature is not enabled") {
		t.Errorf("exit %d, error %+v; want the overlay's hint %q", code, e, hint)
	}
}

func TestCredentialsAreNeededOnlyForCalls(t *testing.T) {
	_, env, vars := newFake(t, respond(`{}`))
	delete(vars, config.EnvClientID)
	delete(vars, config.EnvClientSecret)
	if _, _, code := run(t, env, "repo", "get", "abc"); code != clierr.ExitUsage {
		t.Errorf("a usage mistake without credentials: exit %d, want 2", code)
	}
	if _, stderr, code := run(t, env, "workspace", "get"); code != clierr.ExitAuth || errorOf(t, stderr).Code != "no_credentials" {
		t.Errorf("a call without credentials: exit %d, stderr %q; want 3 no_credentials", code, stderr)
	}
}

func TestDebugLogsTheCall(t *testing.T) {
	_, env, _ := newFake(t, respond(`{}`))
	_, stderr, code := run(t, env, "--debug", "workspace", "get")
	if code != clierr.ExitOK || !strings.Contains(stderr, "debug: GET https://app.aikido.dev/api/public/v1/workspace -> 200") {
		t.Errorf("exit %d, stderr %q", code, stderr)
	}
}

func TestIntegersAreBase10(t *testing.T) {
	f, env, _ := newFake(t, respond(`{}`))
	mustRun(t, env, "team-member", "create", "3", "--user-id", "010")
	if got := f.first(t).Body; got != `{"user_id":10}` {
		t.Errorf("body = %q, want 010 sent as 10", got)
	}
	_, stderr, code := run(t, env, "team-member", "create", "3", "--user-id", "0x10")
	if code != clierr.ExitUsage || !strings.Contains(errorOf(t, stderr).Message, `"0x10" is not a base-10 integer`) {
		t.Errorf("exit %d, stderr %q; want 0x10 refused", code, stderr)
	}
	_, env2, _ := newFake(t, pages(`[1,2,3,4,5,6,7,8,9,10,11,12]`, `[]`))
	if out := mustRun(t, env2, "repo", "list", "--limit", "010"); out != "[1,2,3,4,5,6,7,8,9,10]\n" {
		t.Errorf("stdout = %q, want ten items", out)
	}
}

// An empty argument would drop its path segment and call another endpoint.
func TestEmptyArgumentIsRefused(t *testing.T) {
	f, env, _ := newFake(t, respond(`{}`))
	_, stderr, code := run(t, env, "cve", "get", "")
	if e := errorOf(t, stderr); code != clierr.ExitUsage || e.Code != "invalid_input" || len(f.seen()) != 0 {
		t.Errorf("exit %d, error %+v; want invalid_input before any call", code, e)
	}
}

// Ctrl-C must end a wait on stdin that never closes, as a terminal or a stuck pipe is.
func TestReadBodyFileStopsOnCancel(t *testing.T) {
	r, w := io.Pipe()
	defer w.Close()
	ctx, cancel := context.WithCancel(t.Context())
	time.AfterFunc(20*time.Millisecond, cancel)
	if _, err := readBodyFile(ctx, "-", r); !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want context.Canceled", err)
	}
}

// Ctrl-C while a body is read from stdin ends the run like Ctrl-C during a
// call: exit 1, not a usage error.
func TestCtrlCOnStdinIsNotAUsageError(t *testing.T) {
	for _, args := range [][]string{
		{"team", "create", "--body-file", "-"},
		{"api", "POST", "/teams", "--input", "-"},
	} {
		_, env, _ := newFake(t, respond(`{}`))
		r, w := io.Pipe()
		defer w.Close()
		env.Stdin = r
		var stderr bytes.Buffer
		env.Stdout, env.Stderr = io.Discard, &stderr
		ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
		code := Run(ctx, args, env)
		cancel()
		if code != clierr.ExitUnexpected || strings.Contains(stderr.String(), "invalid_input") {
			t.Errorf("%q: exit %d, stderr %s; want exit 1", args, code, stderr.String())
		}
	}
}
