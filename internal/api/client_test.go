package api

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/cego/aikido-dojo/internal/clierr"
)

func TestDoSendsTheTokenAndTheBody(t *testing.T) {
	ta := newTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if got := r.Header.Get("Authorization"); got != "Bearer tok-1" {
			t.Errorf("Authorization = %q", got)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q", got)
		}
		if r.URL.Path != "/api/public/v1/issues/groups/7/notes" || r.URL.Query().Get("x") != "1" || string(body) != `{"note":"n"}` {
			t.Errorf("got %s %s?%s %s", r.Method, r.URL.Path, r.URL.RawQuery, body)
		}
		fmt.Fprint(w, `{"ok":true}`)
	}, nil)
	resp, err := ta.client.Do(t.Context(), Request{
		Method: http.MethodPost, Path: "/issues/groups/7/notes",
		Query: map[string][]string{"x": {"1"}}, Body: []byte(`{"note":"n"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := readBody(t, resp); got != `{"ok":true}` {
		t.Errorf("body = %s", got)
	}
}

func TestDoRefreshesARejectedTokenOnce(t *testing.T) {
	ta := newTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "Bearer tok-1" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		fmt.Fprint(w, `{}`)
	}, nil)
	resp, err := ta.client.Do(t.Context(), Request{Method: http.MethodGet, Path: "/workspace"})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if n := ta.tokenCalls.Load(); n != 2 {
		t.Errorf("token calls = %d, want 2", n)
	}
}

func TestDoFailures(t *testing.T) {
	tests := []struct {
		name         string
		req          Request
		status       int
		contentType  string
		body         string
		code         string
		exit         int
		wantText     string
		wantAttempts int32
	}{
		{name: "401 after a fresh token", req: Request{Method: http.MethodGet, Path: "/workspace"}, status: 401, body: `{"error":"Invalid request."}`,
			code: "auth_failed", exit: clierr.ExitAuth, wantText: "Invalid request.", wantAttempts: 2},
		{name: "403 names the missing scope", req: Request{Method: http.MethodGet, Path: "/issues/counts", Scope: "issues:read"}, status: 403,
			code: "missing_scope", exit: clierr.ExitForbidden, wantText: "issues:read scope on the API client in Aikido's workspace settings", wantAttempts: 1},
		{name: "403 without a known scope", req: Request{Method: http.MethodGet, Path: "/x"}, status: 403,
			code: "forbidden", exit: clierr.ExitForbidden, wantText: "workspace settings", wantAttempts: 1},
		{name: "404 keeps Aikido's message", req: Request{Method: http.MethodGet, Path: "/issues/9"}, status: 404, body: `{"message":"Issue not found"}`,
			code: "not_found", exit: clierr.ExitNotFound, wantText: "GET /issues/9: Issue not found", wantAttempts: 1},
		{name: "400 is a usage error", req: Request{Method: http.MethodGet, Path: "/issues/export"}, status: 400, body: `{"error":"Request too big"}`,
			code: "bad_request", exit: clierr.ExitUsage, wantText: "Request too big", wantAttempts: 1},
		{name: "a failed write is never retried", req: Request{Method: http.MethodPost, Path: "/teams", Body: []byte(`{}`)}, status: 500,
			code: "server_error", exit: clierr.ExitUnexpected, wantText: "POST /teams", wantAttempts: 1},
		{name: "an HTML error page becomes its status text", req: Request{Method: http.MethodPost, Path: "/teams"}, status: 503,
			contentType: "text/html", body: "<html><h1>Service Unavailable</h1></html>",
			code: "server_error", exit: clierr.ExitUnexpected, wantText: "POST /teams: Service Unavailable", wantAttempts: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var attempts atomic.Int32
			ta := newTestAPI(t, func(w http.ResponseWriter, _ *http.Request) {
				attempts.Add(1)
				if tt.contentType != "" {
					w.Header().Set("Content-Type", tt.contentType)
				}
				w.WriteHeader(tt.status)
				fmt.Fprint(w, tt.body)
			}, nil)
			err := doErr(t, ta.client, tt.req)
			e := wantCode(t, err, tt.code, tt.exit)
			if e.HTTPStatus != tt.status {
				t.Errorf("http_status = %d, want %d", e.HTTPStatus, tt.status)
			}
			if !strings.Contains(e.Message+" "+e.Hint, tt.wantText) {
				t.Errorf("message/hint = %q / %q, want %q", e.Message, e.Hint, tt.wantText)
			}
			if got := attempts.Load(); got != tt.wantAttempts {
				t.Errorf("attempts = %d, want %d", got, tt.wantAttempts)
			}
		})
	}
}

func TestDebugLogNeverShowsCredentials(t *testing.T) {
	var log bytes.Buffer
	ta := newTestAPI(t, func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, `{}`) }, &log)
	resp, err := ta.client.Do(t.Context(), Request{Method: http.MethodGet, Path: "/workspace"})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	out := log.String()
	for _, want := range []string{
		"debug: POST https://" + ta.host + "/api/oauth/token -> 200 (",
		"debug: GET https://" + ta.host + "/api/public/v1/workspace -> 200 (",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("debug log lacks %q:\n%s", want, out)
		}
	}
	for _, secret := range []string{testSecret, "tok-1", "Basic", "Bearer"} {
		if strings.Contains(out, secret) {
			t.Errorf("debug log leaks %q:\n%s", secret, out)
		}
	}
}

func TestDoAsksForGzipAndDecodesIt(t *testing.T) {
	ta := newTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			t.Errorf("Accept-Encoding = %q, want gzip", r.Header.Get("Accept-Encoding"))
		}
		w.Header().Set("Content-Encoding", "gzip")
		zw := gzip.NewWriter(w)
		fmt.Fprint(zw, `[{"id":1}]`)
		if err := zw.Close(); err != nil {
			t.Error(err)
		}
	}, nil)
	resp, err := ta.client.Do(t.Context(), Request{Method: http.MethodGet, Path: "/issues/export"})
	if err != nil {
		t.Fatal(err)
	}
	if got := readBody(t, resp); got != `[{"id":1}]` {
		t.Errorf("body = %q, want the decompressed JSON", got)
	}
}
