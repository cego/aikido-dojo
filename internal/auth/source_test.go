package auth

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/zalando/go-keyring"

	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/config"
)

const testSecret = "s3cr3t-value"

type tokenServer struct {
	*httptest.Server
	calls atomic.Int32
}

// newTokenServer issues tok-1, tok-2, … to client "id" with testSecret, or answers every call with status and body when status is non-zero.
func newTokenServer(t *testing.T, status int, body string) *tokenServer {
	t.Helper()
	ts := &tokenServer{}
	ts.Server = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/oauth/token" || r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		n := ts.calls.Add(1)
		if status != 0 {
			w.WriteHeader(status)
			fmt.Fprint(w, body)
			return
		}
		if id, secret, ok := r.BasicAuth(); !ok || id != "id" || secret != testSecret {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{"error":"invalid_client","error_description":"The provided credentials are invalid"}`)
			return
		}
		if r.FormValue("grant_type") != "client_credentials" {
			t.Errorf("grant_type = %q", r.FormValue("grant_type"))
		}
		fmt.Fprintf(w, `{"access_token":"tok-%d","expires_in":3600,"token_type":"bearer"}`, n)
	}))
	t.Cleanup(ts.Close)
	return ts
}

func (ts *tokenServer) resolved(profile string) config.Resolved {
	r := config.Resolved{Profile: profile, ClientID: "id", Region: "eu", Host: strings.TrimPrefix(ts.URL, "https://")}
	if profile == "" {
		r.Secret = testSecret
	}
	return r
}

func storeSecret(t *testing.T, profile string) {
	t.Helper()
	if err := keyring.Set(keychainService, account(profile, "client_secret"), testSecret); err != nil {
		t.Fatal(err)
	}
}

func wantCode(t *testing.T, err error, code string, exit int) clierr.Error {
	t.Helper()
	var e *clierr.Error
	if !errors.As(err, &e) || e.Code != code || e.Exit != exit {
		t.Fatalf("err = %v, want code %s exit %d", err, code, exit)
	}
	return *e
}

func TestNewSourceNeedsAStoredSecret(t *testing.T) {
	keyring.MockInit()
	ts := newTokenServer(t, 0, "")
	_, err := NewSource(ts.Client(), ts.resolved("cego"))
	e := wantCode(t, err, "no_secret", clierr.ExitAuth)
	if !strings.Contains(e.Hint, "auth login --profile cego") {
		t.Errorf("hint = %q", e.Hint)
	}
}

func TestNewSourceReportsAKeychainFailure(t *testing.T) {
	keyring.MockInitWithError(errors.New("keychain locked"))
	ts := newTokenServer(t, 0, "")
	_, err := NewSource(ts.Client(), ts.resolved("cego"))
	if err == nil || !strings.Contains(err.Error(), "keychain locked") {
		t.Fatalf("err = %v, want the keychain failure", err)
	}
}

func TestTokenIsCachedInTheKeychainAcrossProcesses(t *testing.T) {
	keyring.MockInit()
	storeSecret(t, "cego")
	ts := newTokenServer(t, 0, "")
	for range 2 { // a fresh Source per iteration stands in for a new process
		s, err := NewSource(ts.Client(), ts.resolved("cego"))
		if err != nil {
			t.Fatal(err)
		}
		tok, err := s.Token(t.Context())
		if err != nil {
			t.Fatal(err)
		}
		if tok != "tok-1" {
			t.Errorf("token = %s, want the cached tok-1", tok)
		}
	}
	if n := ts.calls.Load(); n != 1 {
		t.Errorf("token endpoint called %d times, want 1", n)
	}
}

func TestTokenIsRefetchedNearExpiry(t *testing.T) {
	keyring.MockInit()
	storeSecret(t, "cego")
	ts := newTokenServer(t, 0, "")
	s, err := NewSource(ts.Client(), ts.resolved("cego"))
	if err != nil {
		t.Fatal(err)
	}
	start := time.Unix(1_800_000_000, 0)
	s.now = func() time.Time { return start }
	if _, err := s.Token(t.Context()); err != nil {
		t.Fatal(err)
	}
	s.now = func() time.Time { return start.Add(time.Hour - 30*time.Second) }
	tok, err := s.Token(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if tok != "tok-2" {
		t.Errorf("token = %s, want a fresh tok-2 thirty seconds before expiry", tok)
	}
}

func TestEnvironmentProfileStoresNothing(t *testing.T) {
	keyring.MockInit()
	ts := newTokenServer(t, 0, "")
	s, err := NewSource(ts.Client(), ts.resolved(""))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Token(t.Context()); err != nil {
		t.Fatal(err)
	}
	if _, err := keyring.Get(keychainService, account("", "access_token")); !errors.Is(err, keyring.ErrNotFound) {
		t.Errorf("keychain lookup err = %v, want nothing stored", err)
	}
}

func TestInvalidateForcesANewToken(t *testing.T) {
	keyring.MockInit()
	storeSecret(t, "cego")
	ts := newTokenServer(t, 0, "")
	s, err := NewSource(ts.Client(), ts.resolved("cego"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Token(t.Context()); err != nil {
		t.Fatal(err)
	}
	if err := s.Invalidate(); err != nil {
		t.Fatal(err)
	}
	if tok, err := s.Token(t.Context()); err != nil || tok != "tok-2" {
		t.Errorf("token = %s, %v, want tok-2", tok, err)
	}
}

func TestADamagedCacheEntryIsReplaced(t *testing.T) {
	keyring.MockInit()
	storeSecret(t, "cego")
	if err := keyring.Set(keychainService, account("cego", "access_token"), "garbage"); err != nil {
		t.Fatal(err)
	}
	ts := newTokenServer(t, 0, "")
	s, err := NewSource(ts.Client(), ts.resolved("cego"))
	if err != nil {
		t.Fatal(err)
	}
	if tok, err := s.Token(t.Context()); err != nil || tok != "tok-1" {
		t.Errorf("token = %s, %v, want a fetched tok-1", tok, err)
	}
}

func TestTokenFailures(t *testing.T) {
	tests := []struct {
		name     string
		profile  string
		status   int
		body     string
		code     string
		exit     int
		wantHint string
	}{
		{
			name: "a rejected secret for a stored profile", profile: "cego", status: 401,
			body: `{"error":"invalid_client","error_description":"The provided credentials are invalid"}`,
			code: "auth_failed", exit: clierr.ExitAuth, wantHint: "auth login --profile cego",
		},
		{
			name: "a rejected secret from the environment", status: 401,
			body: `{"error":"invalid_client","error_description":"The provided credentials are invalid"}`,
			code: "auth_failed", exit: clierr.ExitAuth, wantHint: config.EnvClientSecret,
		},
		{
			name: "a Public app", status: 401,
			body: `{"error":"invalid_client","error_description":"Your app is public, only meant for 3-legged oauth flow"}`,
			code: "auth_failed", exit: clierr.ExitAuth, wantHint: "Public app",
		},
		{name: "rate limited", status: 429, body: `{"error":"too many"}`, code: "rate_limited", exit: clierr.ExitRateLimited, wantHint: "20 calls per minute"},
		{name: "server failure", status: 502, body: `<html>bad gateway</html>`, code: "server_error", exit: clierr.ExitUnexpected, wantHint: "retry"},
		{name: "200 without a token", status: 200, body: `{}`, code: "auth_failed", exit: clierr.ExitAuth, wantHint: "API client"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyring.MockInit()
			if tt.profile != "" {
				storeSecret(t, tt.profile)
			}
			ts := newTokenServer(t, tt.status, tt.body)
			s, err := NewSource(ts.Client(), ts.resolved(tt.profile))
			if err != nil {
				t.Fatal(err)
			}
			_, err = s.Token(t.Context())
			e := wantCode(t, err, tt.code, tt.exit)
			if !strings.Contains(e.Hint, tt.wantHint) {
				t.Errorf("hint = %q, want it to mention %q", e.Hint, tt.wantHint)
			}
			if strings.Contains(err.Error(), testSecret) {
				t.Errorf("error leaks the secret: %v", err)
			}
		})
	}
}
