package api

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/zalando/go-keyring"

	"github.com/cego/aikido-dojo/internal/auth"
	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/config"
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
		{name: "Aikido's reason_phrase is the message", req: Request{Method: http.MethodGet, Path: "/issues/detail/bulk"}, status: 400,
			body: `{"status_code":400,"reason_phrase":"This feature is not enabled on your workspace."}`,
			code: "bad_request", exit: clierr.ExitUsage, wantText: "GET /issues/detail/bulk: This feature is not enabled on your workspace.", wantAttempts: 1},
		{name: "a failed write is never retried", req: Request{Method: http.MethodPost, Path: "/teams", Body: []byte(`{}`)}, status: 500,
			code: "server_error", exit: clierr.ExitUnexpected, wantText: "POST /teams", wantAttempts: 1},
		{name: "a long plain-text body becomes the status text", req: Request{Method: http.MethodGet, Path: "/x"}, status: 404,
			body: strings.Repeat("x", 300), code: "not_found", exit: clierr.ExitNotFound, wantText: "GET /x: Not Found", wantAttempts: 1},
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

func TestDoReportsATokenFailure(t *testing.T) {
	ta := newTestAPI(t, func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, `{}`) }, nil)
	src, err := auth.NewSource(ta.client.http, config.Resolved{ClientID: "id", Secret: "wrong", Region: "eu", Host: ta.host})
	if err != nil {
		t.Fatal(err)
	}
	ta.client.tokens = src
	err = doErr(t, ta.client, Request{Method: http.MethodGet, Path: "/workspace"})
	wantCode(t, err, "auth_failed", clierr.ExitAuth)
}

func TestDebugLogsATransportError(t *testing.T) {
	srv := httptest.NewTLSServer(http.NotFoundHandler())
	transport, u := srv.Client().Transport, srv.URL+"/x"
	srv.Close()
	var log bytes.Buffer
	hc := &http.Client{Transport: Debug(transport, &log)}
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, u, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := hc.Do(req)
	if err == nil {
		defer resp.Body.Close()
		t.Fatal("request to a closed server succeeded")
	}
	if want := "debug: GET " + u + " -> "; !strings.Contains(log.String(), want) {
		t.Errorf("debug log = %q, want a line starting %q", log.String(), want)
	}
}

func TestDoReportsATransportFailure(t *testing.T) {
	ta := newTestAPI(t, func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, `{}`) }, nil)
	resp, err := ta.client.Do(t.Context(), Request{Method: http.MethodGet, Path: "/workspace"})
	if err != nil {
		t.Fatal(err)
	}
	readBody(t, resp)
	ta.close()
	err = doErr(t, ta.client, Request{Method: http.MethodGet, Path: "/workspace"})
	if !strings.Contains(err.Error(), "GET /workspace: ") {
		t.Errorf("err = %v, want the failed request named", err)
	}
}

func TestDoReportsAFailedTokenInvalidation(t *testing.T) {
	keyring.MockInit()
	if err := keyring.Set("aikido-dojo", "cego/client_secret", testSecret); err != nil {
		t.Fatal(err)
	}
	ta := newTestAPI(t, func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusUnauthorized) }, nil)
	src, err := auth.NewSource(ta.client.http, config.Resolved{Profile: "cego", ClientID: "id", Region: "eu", Host: ta.host})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := src.Token(t.Context()); err != nil {
		t.Fatal(err)
	}
	ta.client.tokens = src
	keyring.MockInitWithError(errors.New("keychain locked"))
	err = doErr(t, ta.client, Request{Method: http.MethodGet, Path: "/workspace"})
	if !strings.Contains(err.Error(), "drop the cached access token: keychain locked") {
		t.Errorf("err = %v", err)
	}
}
