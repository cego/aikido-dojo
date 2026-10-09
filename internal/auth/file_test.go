package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/zalando/go-keyring"

	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/config"
)

var errNoKeychain = errors.New("the name org.freedesktop.secrets was not provided by any .service files")

// fileProfile is a profile kept in a credentials file whose directory doesn't exist yet.
func (ts *tokenServer) fileProfile(profile string) config.Resolved {
	r := ts.resolved(profile)
	r.Storage = config.StorageFile
	r.Credentials = filepath.Join(ts.dir, "new", "credentials.json")
	return r
}

func readCredentials(t *testing.T, path string) map[string]map[string]string {
	t.Helper()
	data, err := os.ReadFile(path) //nolint:gosec // the test's own temp file
	if err != nil {
		t.Fatal(err)
	}
	var f struct {
		Profiles map[string]map[string]string `json:"profiles"`
	}
	if err := json.Unmarshal(data, &f); err != nil {
		t.Fatalf("%s: %v", data, err)
	}
	return f.Profiles
}

func writeCredentials(t *testing.T, path, content string, mode os.FileMode) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, mode); err != nil {
		t.Fatal(err)
	}
}

func TestLoginKeepsTheSecretInAFileWithNoKeychain(t *testing.T) {
	ts := newTokenServer(t, 0, "")
	keyring.MockInitWithError(errNoKeychain)
	r := ts.fileProfile("cego")
	r.Secret = testSecret
	if err := Login(t.Context(), ts.Client(), r, false, noRecord); err != nil {
		t.Fatal(err)
	}
	for p, want := range map[string]os.FileMode{r.Credentials: 0o600, filepath.Dir(r.Credentials): 0o700} {
		if info, err := os.Stat(p); err != nil || info.Mode().Perm() != want {
			t.Errorf("%s: %v, %v; want mode %v", p, info.Mode().Perm(), err, want)
		}
	}
	items := readCredentials(t, r.Credentials)["cego"]
	var cached cachedToken
	if err := json.Unmarshal([]byte(items["access_token"]), &cached); err != nil || items["client_secret"] != testSecret || cached.Token != "tok-1" {
		t.Errorf("stored %q, %v; want the secret and tok-1", items, err)
	}
	// The next process reads both from the file, with no second token call.
	src, err := NewSource(ts.Client(), ts.fileProfile("cego"))
	if err != nil {
		t.Fatal(err)
	}
	if tok, err := src.Token(t.Context()); err != nil || tok != "tok-1" || ts.calls.Load() != 1 {
		t.Errorf("token = %q, %v after %d calls; want tok-1 after 1", tok, err, ts.calls.Load())
	}
	if err := src.Invalidate(); err != nil {
		t.Fatal(err)
	}
	if _, ok := readCredentials(t, r.Credentials)["cego"]["access_token"]; ok {
		t.Error("Invalidate left the cached token in the file")
	}
}

func TestNewSourceNeedsAStoredSecretInTheFile(t *testing.T) {
	keyring.MockInitWithError(errNoKeychain)
	ts := newTokenServer(t, 0, "")
	_, err := NewSource(ts.Client(), ts.fileProfile("cego"))
	if e := wantCode(t, err, "no_secret", clierr.ExitAuth); !strings.Contains(e.Hint, "auth login --profile cego") {
		t.Errorf("hint = %q", e.Hint)
	}
}

// Like ssh with StrictModes, a credentials file others could read or replace is refused, on reads and on writes.
func TestAnInsecureCredentialsFileIsRefused(t *testing.T) {
	for _, tc := range []struct {
		name  string
		setup func(t *testing.T, path string)
		hint  func(path string) string
	}{
		{
			name:  "readable by others",
			setup: func(t *testing.T, path string) { writeCredentials(t, path, `{"profiles":{}}`, 0o644) },
			hint:  func(path string) string { return "chmod 600 " + path },
		},
		{
			name: "a symlink",
			setup: func(t *testing.T, path string) {
				target := filepath.Join(t.TempDir(), "elsewhere.json")
				writeCredentials(t, target, `{"profiles":{}}`, 0o600)
				if err := os.Symlink(target, path); err != nil {
					t.Fatal(err)
				}
			},
			hint: func(path string) string { return "remove " + path },
		},
		{
			name: "in a group-writable directory",
			setup: func(t *testing.T, path string) {
				writeCredentials(t, path, `{"profiles":{}}`, 0o600)
				if err := os.Chmod(filepath.Dir(path), 0o770); err != nil { //nolint:gosec // the mode under test
					t.Fatal(err)
				}
			},
			hint: func(path string) string { return "chmod go-w " + filepath.Dir(path) },
		},
		{
			name: "in a world-writable directory",
			setup: func(t *testing.T, path string) {
				writeCredentials(t, path, `{"profiles":{}}`, 0o600)
				if err := os.Chmod(filepath.Dir(path), 0o703); err != nil { //nolint:gosec // the mode under test
					t.Fatal(err)
				}
			},
			hint: func(path string) string { return "chmod go-w " + filepath.Dir(path) },
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			keyring.MockInitWithError(errNoKeychain)
			ts := newTokenServer(t, 0, "")
			r := ts.fileProfile("cego")
			r.Credentials = filepath.Join(t.TempDir(), "credentials.json")
			tc.setup(t, r.Credentials)
			_, err := NewSource(ts.Client(), r)
			e := wantCode(t, err, "insecure_credentials_file", clierr.ExitUsage)
			if !strings.Contains(e.Message, r.Credentials) && !strings.Contains(e.Message, filepath.Dir(r.Credentials)) {
				t.Errorf("message = %q, want the path named", e.Message)
			}
			if !strings.Contains(e.Hint, tc.hint(r.Credentials)) {
				t.Errorf("hint = %q, want %q", e.Hint, tc.hint(r.Credentials))
			}
			r.Secret = testSecret
			wantCode(t, Login(t.Context(), ts.Client(), r, false, noRecord), "insecure_credentials_file", clierr.ExitUsage)
		})
	}
}

// A test can't chown, so the file is checked against another user's uid instead.
func TestACredentialsFileOfAnotherUserIsRefused(t *testing.T) {
	c := credentialsFile{path: filepath.Join(t.TempDir(), "credentials.json"), profile: "cego"}
	writeCredentials(t, c.path, `{"profiles":{}}`, 0o600)
	if err := c.check(os.Getuid()); err != nil {
		t.Fatalf("own file: %v", err)
	}
	e := wantCode(t, c.check(os.Getuid()+1), "insecure_credentials_file", clierr.ExitUsage)
	if !strings.Contains(e.Message, "another user") || !strings.Contains(e.Hint, "remove "+c.path) {
		t.Errorf("message %q, hint %q", e.Message, e.Hint)
	}
}

// Each writer opens its own lock, so they exclude each other as separate
// processes would: a token cached for one profile can't drop another's login.
func TestConcurrentWritersKeepEveryProfile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "credentials.json")
	var wg sync.WaitGroup
	for i := range 40 {
		wg.Go(func() {
			if err := (credentialsFile{path: path, profile: fmt.Sprint("p", i)}).set(clientSecret, testSecret); err != nil {
				t.Error(err)
			}
		})
	}
	wg.Wait()
	if n := len(readCredentials(t, path)); n != 40 {
		t.Errorf("the file holds %d profiles, want all 40", n)
	}
	if entries, _ := os.ReadDir(filepath.Dir(path)); len(entries) != 2 {
		t.Errorf("the directory holds %d files, want the credentials and their lock", len(entries))
	}
}

func TestADamagedCredentialsFileIsReportedNotReplaced(t *testing.T) {
	keyring.MockInitWithError(errNoKeychain)
	ts := newTokenServer(t, 0, "")
	r := ts.fileProfile("cego")
	writeCredentials(t, r.Credentials, `{"profiles":`, 0o600)
	_, err := NewSource(ts.Client(), r)
	wantCode(t, err, "bad_credentials_file", clierr.ExitUsage)
	r.Secret = testSecret
	wantCode(t, Login(t.Context(), ts.Client(), r, false, noRecord), "bad_credentials_file", clierr.ExitUsage)
	if data, err := os.ReadFile(r.Credentials); err != nil || string(data) != `{"profiles":` {
		t.Errorf("file = %q, %v; want it untouched", data, err)
	}
}

func TestLoginMovesAProfileIntoTheFile(t *testing.T) {
	keyring.MockInit()
	ts := newTokenServer(t, 0, "")
	storeSecret(t, "cego")
	if err := keyring.Set(keychainService, account("cego", "access_token"), "{}"); err != nil {
		t.Fatal(err)
	}
	r := ts.fileProfile("cego")
	r.Secret = testSecret
	if err := Login(t.Context(), ts.Client(), r, true, noRecord); err != nil {
		t.Fatal(err)
	}
	for _, item := range []string{"client_secret", "access_token"} {
		if _, err := keyring.Get(keychainService, account("cego", item)); !errors.Is(err, keyring.ErrNotFound) {
			t.Errorf("keychain %s: %v, want it deleted", item, err)
		}
	}
	if readCredentials(t, r.Credentials)["cego"]["client_secret"] != testSecret {
		t.Error("the file holds no secret")
	}
}

// A machine moved to the file has no keychain to delete the old copy from, and that is no failure.
func TestLoginMovingIntoTheFileIgnoresAMissingKeychain(t *testing.T) {
	keyring.MockInitWithError(errNoKeychain)
	ts := newTokenServer(t, 0, "")
	r := ts.fileProfile("cego")
	r.Secret = testSecret
	if err := Login(t.Context(), ts.Client(), r, true, noRecord); err != nil {
		t.Fatalf("err = %v, want the login to succeed", err)
	}
}

func TestLoginMovesAProfileBackToTheKeychain(t *testing.T) {
	keyring.MockInit()
	ts := newTokenServer(t, 0, "")
	file := ts.fileProfile("cego")
	file.Secret = testSecret
	other := ts.fileProfile("ci")
	other.Secret = testSecret
	for _, r := range []config.Resolved{file, other} {
		if err := Login(t.Context(), ts.Client(), r, false, noRecord); err != nil {
			t.Fatal(err)
		}
	}
	back := file
	back.Storage = ""
	if err := Login(t.Context(), ts.Client(), back, true, noRecord); err != nil {
		t.Fatal(err)
	}
	if got, err := keyring.Get(keychainService, account("cego", "client_secret")); err != nil || got != testSecret {
		t.Errorf("keychain secret = %q, %v", got, err)
	}
	if profiles := readCredentials(t, file.Credentials); len(profiles) != 1 || profiles["ci"] == nil {
		t.Errorf("file holds %v, want only ci", profiles)
	}
	ci := other
	ci.Storage = ""
	if err := Login(t.Context(), ts.Client(), ci, true, noRecord); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(file.Credentials); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("credentials file: %v, want it removed with its last profile", err)
	}
}

func TestForgetAFileProfileWithNoKeychain(t *testing.T) {
	keyring.MockInitWithError(errNoKeychain)
	ts := newTokenServer(t, 0, "")
	r := ts.fileProfile("cego")
	r.Secret = testSecret
	if err := Login(t.Context(), ts.Client(), r, false, noRecord); err != nil {
		t.Fatal(err)
	}
	if removed, err := Forget(r); err != nil || !removed {
		t.Fatalf("Forget = %v, %v; want true", removed, err)
	}
	if _, err := os.Lstat(r.Credentials); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("credentials file: %v, want it removed with its last profile", err)
	}
	if removed, err := Forget(r); err != nil || removed {
		t.Errorf("a second Forget = %v, %v; want false, nil", removed, err)
	}
}

// A keychain profile's logout also clears a copy in the file, which a failed move may have left.
func TestForgetClearsBothStores(t *testing.T) {
	keyring.MockInit()
	r := storedProfile(t, "cego")
	writeCredentials(t, r.Credentials, `{"profiles":{"cego":{"client_secret":"s","access_token":"{}"}}}`, 0o600)
	if removed, err := Forget(r); err != nil || !removed {
		t.Fatalf("Forget = %v, %v; want true, for the secret in the file", removed, err)
	}
	if _, err := os.Lstat(r.Credentials); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("credentials file: %v, want it removed", err)
	}
}
