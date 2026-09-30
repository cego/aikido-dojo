package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/cego/aikido-dojo/internal/auth"
	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/config"
	"github.com/cego/aikido-dojo/internal/ratelimit"
)

const testSecret = "s3cr3t-value"

type testAPI struct {
	client     *Client
	slept      []time.Duration
	tokenCalls atomic.Int32
	host       string
	close      func()
}

// newTestAPI serves the token endpoint and routes /api/public/v1/ to handler.
// Retries record their waits in slept instead of sleeping. debug, if non-nil,
// receives the debug log.
func newTestAPI(t *testing.T, handler http.HandlerFunc, debug io.Writer) *testAPI {
	t.Helper()
	ta := &testAPI{}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		if id, secret, ok := r.BasicAuth(); !ok || id != "id" || secret != testSecret {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		fmt.Fprintf(w, `{"access_token":"tok-%d","expires_in":3600,"token_type":"bearer"}`, ta.tokenCalls.Add(1))
	})
	mux.Handle("/api/public/v1/", handler)
	srv := httptest.NewTLSServer(mux)
	t.Cleanup(srv.Close)
	ta.close = srv.Close
	hc := srv.Client()
	if debug != nil {
		hc = &http.Client{Transport: Debug(hc.Transport, debug)}
	}
	ta.host = strings.TrimPrefix(srv.URL, "https://")
	src, err := auth.NewSource(hc, config.Resolved{ClientID: "id", Secret: testSecret, Region: "eu", Host: ta.host})
	if err != nil {
		t.Fatal(err)
	}
	ta.client = New(hc, ta.host, src, ratelimit.New(t.TempDir(), "id"))
	ta.client.sleep = func(_ context.Context, d time.Duration) error {
		ta.slept = append(ta.slept, d)
		return nil
	}
	return ta
}

func wantCode(t *testing.T, err error, code string, exit int) clierr.Error {
	t.Helper()
	var e *clierr.Error
	if !errors.As(err, &e) || e.Code != code || e.Exit != exit {
		t.Fatalf("err = %v, want code %s exit %d", err, code, exit)
	}
	return *e
}

// doErr runs a request that must fail and returns its error.
func doErr(t *testing.T, c *Client, req Request) error {
	t.Helper()
	resp, err := c.Do(t.Context(), req)
	if err == nil {
		defer resp.Body.Close()
		t.Fatalf("%s %s succeeded, want an error", req.Method, req.Path)
	}
	return err
}

func readBody(t *testing.T, resp *http.Response) string {
	t.Helper()
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
