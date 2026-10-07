package cli

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/zalando/go-keyring"

	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/config"
)

// testJWT is shaped like Aikido's access tokens: a JWT whose scope claim
// lists the API client's scopes (live, 2026-09-29).
func testJWT(scope string) string {
	enc := base64.RawURLEncoding
	return enc.EncodeToString([]byte(`{"alg":"RS256","typ":"JWT"}`)) + "." +
		enc.EncodeToString([]byte(`{"client_id":"id","scope":"`+scope+`"}`)) + ".signature"
}

func answer(secret string) func(context.Context, string) (string, error) {
	return func(context.Context, string) (string, error) { return secret, nil }
}

func writeConfig(t *testing.T, vars map[string]string, content string) {
	t.Helper()
	if err := os.WriteFile(vars[config.EnvConfig], []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func readConfig(t *testing.T, vars map[string]string) config.File {
	t.Helper()
	f, err := config.Load(vars[config.EnvConfig])
	if err != nil {
		t.Fatal(err)
	}
	return f
}

// noPair leaves the config file and the keychain to say who calls run as.
func noPair(vars map[string]string) {
	delete(vars, config.EnvClientID)
	delete(vars, config.EnvClientSecret)
}

func TestLoginStoresAProfileOnceAikidoAcceptsIt(t *testing.T) {
	f, env, vars := newFake(t, respond(`{}`))
	noPair(vars)
	var prompts []string
	env.ReadSecret = func(_ context.Context, prompt string) (string, error) {
		prompts = append(prompts, prompt)
		return "typed-secret", nil
	}
	if out := mustRun(t, env, "--profile", "cego", "auth", "login", "--client-id", "id-1", "--region", "us"); out != `{"profile":"cego","client_id":"id-1","region":"us"}`+"\n" {
		t.Errorf("stdout = %q", out)
	}
	if len(prompts) != 1 || !strings.Contains(prompts[0], `"cego"`) {
		t.Errorf("prompts = %q, want one naming the profile", prompts)
	}
	if got := f.seenLogins(); !slices.Equal(got, []string{"id-1:typed-secret"}) {
		t.Errorf("token requests = %q, want one with the typed pair", got)
	}
	if hosts := f.rewrite.seenHosts(); !slices.Equal(hosts, []string{"app.us.aikido.dev"}) {
		t.Errorf("hosts = %q, want only the us token endpoint", hosts)
	}
	want := config.File{DefaultProfile: "cego", Profiles: map[string]config.Profile{"cego": {ClientID: "id-1", Region: "us"}}}
	if got := readConfig(t, vars); !reflect.DeepEqual(got, want) {
		t.Errorf("config = %+v, want %+v", got, want)
	}
	if s, err := keyring.Get("aikido-dojo", "cego/client_secret"); err != nil || s != "typed-secret" {
		t.Errorf("keychain secret = %q, %v", s, err)
	}
	// The next process runs as the new default profile, on the token login cached.
	mustRun(t, env, "workspace", "get")
	if n := len(f.seenLogins()); n != 1 {
		t.Errorf("token requests = %d, want still 1", n)
	}
}

func TestLoginTakesTheSecretFromTheEnvironment(t *testing.T) {
	f, env, vars := newFake(t, respond(`{}`))
	delete(vars, config.EnvClientID)
	mustRun(t, env, "auth", "login", "--client-id", "id-1")
	if got := f.seenLogins(); !slices.Equal(got, []string{"id-1:secret"}) {
		t.Errorf("token requests = %q, want the environment's secret and no prompt", got)
	}
	want := config.File{DefaultProfile: "default", Profiles: map[string]config.Profile{"default": {ClientID: "id-1", Region: "eu"}}}
	if got := readConfig(t, vars); !reflect.DeepEqual(got, want) {
		t.Errorf("config = %+v, want %+v", got, want)
	}
}

func TestLoginKeepsTheOtherProfiles(t *testing.T) {
	_, env, vars := newFake(t, respond(`{}`))
	delete(vars, config.EnvClientID)
	writeConfig(t, vars, `{"default_profile":"cego","profiles":{"cego":{"client_id":"id-1","region":"eu"}}}`)
	mustRun(t, env, "--profile", "ci", "auth", "login", "--client-id", "id-2", "--region", "au")
	want := config.File{DefaultProfile: "cego", Profiles: map[string]config.Profile{
		"cego": {ClientID: "id-1", Region: "eu"}, "ci": {ClientID: "id-2", Region: "au"}}}
	if got := readConfig(t, vars); !reflect.DeepEqual(got, want) {
		t.Errorf("config = %+v, want %+v", got, want)
	}
}

// A shared config in a directory the user can't write still allows a login
// that changes nothing in it.
func TestLoginLeavesAnUnchangedConfigAlone(t *testing.T) {
	_, env, vars := newFake(t, respond(`{}`))
	delete(vars, config.EnvClientID)
	writeConfig(t, vars, `{"default_profile":"cego","profiles":{"cego":{"client_id":"id-1","region":"eu"}}}`)
	dir := filepath.Dir(vars[config.EnvConfig])
	if err := os.Chmod(dir, 0o500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })
	mustRun(t, env, "auth", "login")
	if s, err := keyring.Get("aikido-dojo", "cego/client_secret"); err != nil || s != "secret" {
		t.Errorf("keychain secret = %q, %v", s, err)
	}
}

func TestLoginRefusedStoresNothing(t *testing.T) {
	f, env, vars := newFake(t, respond(`{}`))
	noPair(vars)
	f.issueTokens(func(string, string) (int, string) { return http.StatusUnauthorized, "" })
	env.ReadSecret = answer("wrong")
	_, stderr, code := run(t, env, "--profile", "cego", "auth", "login", "--client-id", "id-1")
	if e := errorOf(t, stderr); code != clierr.ExitAuth || e.Code != "auth_failed" {
		t.Errorf("exit %d, error %+v; want 3 auth_failed", code, e)
	}
	if _, err := os.Stat(vars[config.EnvConfig]); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("config file: %v, want none written", err)
	}
	for _, item := range []string{"client_secret", "access_token"} {
		if _, err := keyring.Get("aikido-dojo", "cego/"+item); !errors.Is(err, keyring.ErrNotFound) {
			t.Errorf("keychain %s: %v, want nothing stored", item, err)
		}
	}
}

func TestLoginMistakesCostNoCall(t *testing.T) {
	for _, tt := range []struct {
		args   []string
		secret string
		code   string
	}{
		{[]string{"auth", "login"}, "s", "no_client_id"},
		{[]string{"auth", "login", "--client-id", "x", "--region", "mars"}, "s", "unknown_region"},
		{[]string{"auth", "login", "--client-id", "x"}, " ", "empty_secret"},
		{[]string{"auth", "login", "extra"}, "s", "usage"},
	} {
		f, env, vars := newFake(t, respond(`{}`))
		noPair(vars)
		env.ReadSecret = answer(tt.secret)
		_, stderr, code := run(t, env, tt.args...)
		if e := errorOf(t, stderr); code != clierr.ExitUsage || e.Code != tt.code || e.Hint == "" {
			t.Errorf("%q: exit %d, error %+v; want exit 2 %s with a hint", tt.args, code, e, tt.code)
		}
		if n := len(f.seenLogins()); n != 0 {
			t.Errorf("%q: %d token requests, want none", tt.args, n)
		}
	}
}

func TestStatusReportsTheTokensScopes(t *testing.T) {
	f, env, _ := newFake(t, respond(`{}`))
	f.issueTokens(func(string, string) (int, string) { return http.StatusOK, testJWT("issues:read basics:read extra:read") })
	var got authStatus
	if err := json.Unmarshal([]byte(mustRun(t, env, "auth", "status")), &got); err != nil {
		t.Fatal(err)
	}
	if got.Profile != "" || got.Source != "environment" || got.Method != "client_credentials" || got.ClientID != "id" || got.Region != "eu" {
		t.Errorf("status = %+v", got)
	}
	s := got.Scopes
	if !slices.Equal(s.Granted, []string{"basics:read", "extra:read", "issues:read"}) || len(s.Unknown) != 0 ||
		!slices.Contains(s.Denied, "issues:write") || slices.Contains(s.Denied, "issues:read") {
		t.Errorf("scopes = %+v", s)
	}
	if n := len(f.seen()); n != 0 {
		t.Errorf("API calls = %d, want none: the token answers on its own", n)
	}
}

func TestStatusWithoutAScopeClaim(t *testing.T) {
	_, env, _ := newFake(t, respond(`{}`))
	stdout, stderr, code := run(t, env, "auth", "status")
	var got authStatus
	if err := json.Unmarshal([]byte(stdout), &got); err != nil || code != clierr.ExitOK {
		t.Fatalf("exit %d, %v: %s", code, err, stderr)
	}
	if len(got.Scopes.Granted) != 0 || len(got.Scopes.Denied) != 0 || !slices.Contains(got.Scopes.Unknown, "issues:read") {
		t.Errorf("scopes = %+v, want every needed scope unknown", got.Scopes)
	}
	if !strings.Contains(stderr, `{"warning":`) {
		t.Errorf("stderr = %q, want a warning", stderr)
	}
}

func TestStatusOfAStoredProfile(t *testing.T) {
	_, env, vars := newFake(t, respond(`{}`))
	noPair(vars)
	writeConfig(t, vars, `{"default_profile":"cego","profiles":{"cego":{"client_id":"id-1","region":"au"}}}`)
	if err := keyring.Set("aikido-dojo", "cego/client_secret", "stored"); err != nil {
		t.Fatal(err)
	}
	var got authStatus
	if err := json.Unmarshal([]byte(mustRun(t, env, "auth", "status")), &got); err != nil {
		t.Fatal(err)
	}
	if got.Profile != "cego" || got.Source != "keychain" || got.ClientID != "id-1" || got.Region != "au" {
		t.Errorf("status = %+v", got)
	}
}

func TestAuthNeverPrintsASecretOrAToken(t *testing.T) {
	f, env, vars := newFake(t, respond(`{}`))
	vars[config.EnvClientSecret] = "env-secret-value"
	tok := testJWT("basics:read")
	f.issueTokens(func(string, string) (int, string) { return http.StatusOK, tok })
	for _, args := range [][]string{
		{"--profile", "cego", "auth", "login"},
		{"auth", "status"},
		{"--debug", "auth", "status"},
		{"--profile", "cego", "--debug", "auth", "status"},
		{"--profile", "cego", "auth", "logout"},
	} {
		stdout, stderr, code := run(t, env, args...)
		if code != clierr.ExitOK {
			t.Errorf("%q: exit %d: %s", args, code, stderr)
		}
		for _, secret := range []string{"env-secret-value", tok} {
			if strings.Contains(stdout+stderr, secret) {
				t.Errorf("%q printed a secret or token", args)
			}
		}
	}
}

func TestLogout(t *testing.T) {
	_, env, vars := newFake(t, respond(`{}`))
	noPair(vars)
	writeConfig(t, vars, `{"default_profile":"cego","profiles":{"cego":{"client_id":"id-1"}}}`)
	for _, item := range []string{"client_secret", "access_token"} {
		if err := keyring.Set("aikido-dojo", "cego/"+item, "x"); err != nil {
			t.Fatal(err)
		}
	}
	if out := mustRun(t, env, "auth", "logout"); out != `{"profile":"cego","removed":true}`+"\n" {
		t.Errorf("stdout = %q", out)
	}
	for _, item := range []string{"client_secret", "access_token"} {
		if _, err := keyring.Get("aikido-dojo", "cego/"+item); !errors.Is(err, keyring.ErrNotFound) {
			t.Errorf("keychain %s: %v, want it deleted", item, err)
		}
	}
	if _, ok := readConfig(t, vars).Profiles["cego"]; !ok {
		t.Error("logout removed the profile from the config file")
	}
	if out := mustRun(t, env, "auth", "logout"); out != `{"profile":"cego","removed":false}`+"\n" {
		t.Errorf("a second logout: stdout = %q", out)
	}
}

func TestLogoutNeedsAProfile(t *testing.T) {
	env, _ := testEnv(t)
	_, stderr, code := run(t, env, "auth", "logout")
	if e := errorOf(t, stderr); code != clierr.ExitUsage || e.Code != "usage" || e.Hint == "" {
		t.Errorf("exit %d, error %+v; want a usage error", code, e)
	}
}
