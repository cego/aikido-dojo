package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cego/aikido-dojo/internal/clierr"
)

func get(t *testing.T, c *http.Client, url string) (string, error) {
	t.Helper()
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("GET %s: %w", url, err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	return string(b), err
}

func TestRedirectsStayOnTheAPIHost(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/same", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://app.aikido.dev/target", http.StatusFound)
	})
	mux.HandleFunc("/target", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "ok") })
	mux.HandleFunc("/away", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://elsewhere.example/x", http.StatusFound)
	})
	mux.HandleFunc("/plain", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://app.aikido.dev/target", http.StatusFound)
	})
	srv := httptest.NewTLSServer(mux)
	t.Cleanup(srv.Close)
	c := newHTTPClient(Env{Transport: &rewrite{to: strings.TrimPrefix(srv.URL, "https://"), next: srv.Client().Transport}}, "app.aikido.dev", false)

	if body, err := get(t, c, "https://app.aikido.dev/same"); err != nil || body != "ok" {
		t.Errorf("same-host redirect = %q, %v; want it followed", body, err)
	}
	for _, path := range []string{"/away", "/plain"} {
		_, err := get(t, c, "https://app.aikido.dev"+path)
		var e *clierr.Error
		if !errors.As(err, &e) || e.Code != "redirect_refused" {
			t.Errorf("%s: err = %v, want redirect_refused", path, err)
		}
	}
}

func TestHTTPClientBoundsOnlyTheWaitForHeaders(t *testing.T) {
	c := newHTTPClient(Env{}, "app.aikido.dev", false)
	tr, ok := c.Transport.(*http.Transport)
	if c.Timeout != 0 || !ok || tr.ResponseHeaderTimeout != headerTimeout {
		t.Errorf("timeout %v, transport %T; want no overall timeout and a %v header timeout", c.Timeout, c.Transport, headerTimeout)
	}
}

func TestDebugLogsRequestsToStderr(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "ok") }))
	t.Cleanup(srv.Close)
	var log bytes.Buffer
	env := Env{Stderr: &log, Transport: &rewrite{to: strings.TrimPrefix(srv.URL, "https://"), next: srv.Client().Transport}}
	if _, err := get(t, newHTTPClient(env, "app.aikido.dev", true), "https://app.aikido.dev/x"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(log.String(), "debug: GET https://app.aikido.dev/x -> 200") {
		t.Errorf("debug log = %q", log.String())
	}
}
