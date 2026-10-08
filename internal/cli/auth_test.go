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
	readOnly(t, vars)
	mustRun(t, env, "auth", "login")
	if s, err := keyring.Get("aikido-dojo", "cego/client_secret"); err != nil || s != "secret" {
		t.Errorf("keychain secret = %q, %v", s, err)
	}
}

// readOnly makes the config's directory unwritable, as a shared one is.
func readOnly(t *testing.T, vars map[string]string) {
	t.Helper()
	dir := filepath.Dir(vars[config.EnvConfig])
	if err := os.Chmod(dir, 0o500); err != nil { //nolint:gosec // a directory, which needs its execute bit
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o700) }) //nolint:gosec // a directory, which needs its execute bit
}

// An entry that omits its region, in a file with no default profile, is
// already what login would record.
func TestLoginLeavesASharedConfigAlone(t *testing.T) {
	_, env, vars := newFake(t, respond(`{}`))
	delete(vars, config.EnvClientID)
	shared := `{"profiles":{"github":{"client_id":"id-1"},"gitlab":{"client_id":"id-2"}}}`
	writeConfig(t, vars, shared)
	readOnly(t, vars)
	if out := mustRun(t, env, "--profile", "gitlab", "auth", "login"); out != `{"profile":"gitlab","client_id":"id-2","region":"eu"}`+"\n" {
		t.Errorf("stdout = %q", out)
	}
	if data, err := os.ReadFile(vars[config.EnvConfig]); err != nil || string(data) != shared {
		t.Errorf("config = %s, %v; want it unchanged", data, err)
	}
}

// A login whose profile can't be saved keeps the working one: the secret
// stays the one the config's client ID was stored with.
func TestLoginThatCantSaveKeepsTheWorkingProfile(t *testing.T) {
	f, env, vars := newFake(t, respond(`{}`))
	delete(vars, config.EnvClientID)
	writeConfig(t, vars, `{"default_profile":"cego","profiles":{"cego":{"client_id":"id-1","region":"eu"}}}`)
	if err := keyring.Set("aikido-dojo", "cego/client_secret", "old"); err != nil {
		t.Fatal(err)
	}
	readOnly(t, vars)
	_, stderr, code := run(t, env, "auth", "login", "--client-id", "id-2")
	if e := errorOf(t, stderr); code == clierr.ExitOK || e.Code != "config_not_saved" || !strings.Contains(e.Hint, "--config") {
		t.Errorf("exit %d, error %+v; want config_not_saved with a hint naming --config", code, e)
	}
	if s, err := keyring.Get("aikido-dojo", "cego/client_secret"); err != nil || s != "old" {
		t.Errorf("keychain secret = %q, %v; want the old one kept", s, err)
	}
	if _, err := keyring.Get("aikido-dojo", "cego/access_token"); !errors.Is(err, keyring.ErrNotFound) {
		t.Errorf("cached token: %v, want none for the unsaved client", err)
	}
	if got := f.seenLogins(); !slices.Equal(got, []string{"id-2:secret"}) {
		t.Errorf("token requests = %q, want the one check of id-2", got)
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
	f.issueTokens(func(string, string) (int, string) {
		return http.StatusOK, testJWT("issues:read basics:read extra:read")
	})
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
		{"--profile", "cego", "--debug", "auth", "login"},
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

// The auth commands print one object; --ndjson is refused before they act.
func TestAuthRefusesNDJSONBeforeActing(t *testing.T) {
	for _, args := range [][]string{
		{"auth", "login", "--client-id", "cid", "--ndjson"},
		{"--profile", "cego", "auth", "logout", "--ndjson"},
		{"auth", "status", "--ndjson"},
	} {
		f, env, vars := newFake(t, respond(`{}`))
		if err := keyring.Set("aikido-dojo", "cego/client_secret", "kept"); err != nil {
			t.Fatal(err)
		}
		_, stderr, code := run(t, env, args...)
		if code != clierr.ExitUsage || len(f.seenLogins()) != 0 {
			t.Errorf("%q: exit %d, %d token requests, stderr %s; want a usage error before acting", args, code, len(f.seenLogins()), stderr)
		}
		if s, err := keyring.Get("aikido-dojo", "cego/client_secret"); err != nil || s != "kept" {
			t.Errorf("%q: secret %q, %v; want it untouched", args, s, err)
		}
		if _, err := os.Stat(vars[config.EnvConfig]); !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("%q: config written: %v", args, err)
		}
	}
}

func TestLogoutWithAFailingFilterSaysItTookEffect(t *testing.T) {
	env, _ := testEnv(t)
	if err := keyring.Set("aikido-dojo", "cego/client_secret", "x"); err != nil {
		t.Fatal(err)
	}
	_, stderr, code := run(t, env, "--profile", "cego", "auth", "logout", "--jq", ".[0]")
	if e := errorOf(t, stderr); code != clierr.ExitUnexpected || e.Code != "output_failed" || !strings.Contains(e.Message, "auth logout succeeded") {
		t.Errorf("exit %d, error %+v; want output_failed saying logout took effect", code, e)
	}
}

func TestLoginWithAFailingFilterSaysItTookEffect(t *testing.T) {
	_, env, vars := newFake(t, respond(`{}`))
	delete(vars, config.EnvClientID)
	_, stderr, code := run(t, env, "auth", "login", "--client-id", "id-1", "--jq", ".[0]")
	if e := errorOf(t, stderr); code != clierr.ExitUnexpected || e.Code != "output_failed" || !strings.Contains(e.Message, "auth login succeeded") {
		t.Errorf("exit %d, error %+v; want output_failed saying login took effect", code, e)
	}
	if s, err := keyring.Get("aikido-dojo", "default/client_secret"); err != nil || s != "secret" {
		t.Errorf("secret %q, %v; want the login stored", s, err)
	}
}
