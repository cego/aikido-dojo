package auth

import (
	"errors"
	"strings"
	"testing"

	"github.com/zalando/go-keyring"

	"github.com/cego/aikido-dojo/internal/clierr"
)

func noRecord() error { return nil }

func TestLoginStoresTheSecretAndTheToken(t *testing.T) {
	keyring.MockInit()
	ts := newTokenServer(t, 0, "")
	r := ts.resolved("cego")
	r.Secret = testSecret
	if err := Login(t.Context(), ts.Client(), r, noRecord); err != nil {
		t.Fatal(err)
	}
	if got, err := keyring.Get(keychainService, account("cego", "client_secret")); err != nil || got != testSecret {
		t.Errorf("stored secret = %q, %v", got, err)
	}
	// The next process uses the token Login got, with no second token call.
	src, err := NewSource(ts.Client(), ts.resolved("cego"))
	if err != nil {
		t.Fatal(err)
	}
	if tok, err := src.Token(t.Context()); err != nil || tok != "tok-1" || ts.calls.Load() != 1 {
		t.Errorf("token = %q, %v after %d calls; want tok-1 after 1", tok, err, ts.calls.Load())
	}
}

func TestLoginReplacesTheCachedToken(t *testing.T) {
	keyring.MockInit()
	ts := newTokenServer(t, 0, "")
	r := ts.resolved("cego")
	r.Secret = testSecret
	for range 2 {
		if err := Login(t.Context(), ts.Client(), r, noRecord); err != nil {
			t.Fatal(err)
		}
	}
	src, err := NewSource(ts.Client(), ts.resolved("cego"))
	if err != nil {
		t.Fatal(err)
	}
	if tok, err := src.Token(t.Context()); err != nil || tok != "tok-2" {
		t.Errorf("token = %q, %v; want the second login's tok-2", tok, err)
	}
}

func TestLoginStoresNothingWhenAikidoRefuses(t *testing.T) {
	keyring.MockInit()
	ts := newTokenServer(t, 0, "")
	r := ts.resolved("cego")
	r.Secret = "wrong"
	e := wantCode(t, Login(t.Context(), ts.Client(), r, noRecord), "auth_failed", clierr.ExitAuth)
	if !strings.Contains(e.Hint, "Aikido's workspace settings") || strings.Contains(e.Hint, "auth login") {
		t.Errorf("hint = %q, want it to point at the API client, not back at auth login", e.Hint)
	}
	for _, item := range []string{"client_secret", "access_token"} {
		if _, err := keyring.Get(keychainService, account("cego", item)); !errors.Is(err, keyring.ErrNotFound) {
			t.Errorf("%s: %v, want nothing stored", item, err)
		}
	}
}

func TestLoginStoresNothingWhenRecordFails(t *testing.T) {
	keyring.MockInit()
	ts := newTokenServer(t, 0, "")
	r := ts.resolved("cego")
	r.Secret = testSecret
	if err := Login(t.Context(), ts.Client(), r, func() error { return errors.New("disk full") }); err == nil || !strings.Contains(err.Error(), "disk full") {
		t.Errorf("err = %v, want record's error", err)
	}
	for _, item := range []string{"client_secret", "access_token"} {
		if _, err := keyring.Get(keychainService, account("cego", item)); !errors.Is(err, keyring.ErrNotFound) {
			t.Errorf("%s: %v, want nothing stored", item, err)
		}
	}
}

func TestLoginReportsAKeychainFailure(t *testing.T) {
	ts := newTokenServer(t, 0, "")
	keyring.MockInitWithError(errors.New("keychain locked"))
	r := ts.resolved("cego")
	r.Secret = testSecret
	if err := Login(t.Context(), ts.Client(), r, noRecord); err == nil || !strings.Contains(err.Error(), "keychain locked") {
		t.Errorf("err = %v", err)
	}
}

func TestForget(t *testing.T) {
	keyring.MockInit()
	storeSecret(t, "cego")
	if err := keyring.Set(keychainService, account("cego", "access_token"), "{}"); err != nil {
		t.Fatal(err)
	}
	if removed, err := Forget("cego"); err != nil || !removed {
		t.Fatalf("Forget = %v, %v; want true", removed, err)
	}
	for _, item := range []string{"client_secret", "access_token"} {
		if _, err := keyring.Get(keychainService, account("cego", item)); !errors.Is(err, keyring.ErrNotFound) {
			t.Errorf("%s: %v, want it deleted", item, err)
		}
	}
	if removed, err := Forget("cego"); err != nil || removed {
		t.Errorf("a second Forget = %v, %v; want false, nil", removed, err)
	}
}

func TestForgetReportsAKeychainFailure(t *testing.T) {
	keyring.MockInitWithError(errors.New("keychain locked"))
	if _, err := Forget("cego"); err == nil || !strings.Contains(err.Error(), "keychain locked") {
		t.Errorf("err = %v", err)
	}
}
