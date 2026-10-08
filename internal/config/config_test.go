package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/cego/aikido-dojo/internal/clierr"
)

func env(vars map[string]string) func(string) string {
	return func(k string) string { return vars[k] }
}

func wantCode(t *testing.T, err error, code string, exit int) clierr.Error {
	t.Helper()
	var e *clierr.Error
	if !errors.As(err, &e) {
		t.Fatalf("err = %v, want a *clierr.Error with code %s", err, code)
	}
	if e.Code != code || e.Exit != exit {
		t.Fatalf("code, exit = %s, %d, want %s, %d (%v)", e.Code, e.Exit, code, exit, err)
	}
	return *e
}

func TestPath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	tests := []struct {
		name  string
		flags Flags
		vars  map[string]string
		want  string
	}{
		{name: "flag wins", flags: Flags{Config: "/f.json"}, vars: map[string]string{EnvConfig: "/e.json"}, want: "/f.json"},
		{name: "then the environment", vars: map[string]string{EnvConfig: "/e.json"}, want: "/e.json"},
		{name: "then the home directory", want: filepath.Join(home, ".config", "aikido-dojo", "config.json")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Path(tt.flags, env(tt.vars))
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	write := func(t *testing.T, content string) string {
		t.Helper()
		p := filepath.Join(t.TempDir(), "config.json")
		if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
		return p
	}

	t.Run("a missing file is an empty config", func(t *testing.T) {
		f, err := Load(filepath.Join(t.TempDir(), "none.json"))
		if err != nil || f.DefaultProfile != "" || len(f.Profiles) != 0 {
			t.Fatalf("got %+v, %v", f, err)
		}
	})

	t.Run("reads profiles", func(t *testing.T) {
		f, err := Load(write(t, `{"default_profile":"cego","profiles":{"cego":{"client_id":"AIK_CLIENT_x","region":"us"}}}`))
		if err != nil {
			t.Fatal(err)
		}
		if f.DefaultProfile != "cego" || f.Profiles["cego"] != (Profile{ClientID: "AIK_CLIENT_x", Region: "us"}) {
			t.Errorf("got %+v", f)
		}
	})

	for _, tc := range []struct{ name, content, key string }{
		{"a top-level secret", `{"client_secret":"x"}`, "client_secret"},
		{"a secret nested in a profile", `{"profiles":{"cego":{"client_id":"a","client_secret":"x"}}}`, "client_secret"},
		{"a token", `{"profiles":{"cego":{"client_id":"a","access_token":"x"}}}`, "access_token"},
		{"a password", `{"Password":"x"}`, "Password"},
	} {
		t.Run("refuses "+tc.name, func(t *testing.T) {
			p := write(t, tc.content)
			_, err := Load(p)
			e := wantCode(t, err, "secret_in_config", clierr.ExitUsage)
			if !strings.Contains(e.Message, tc.key) || !strings.Contains(e.Message, p) {
				t.Errorf("message %q should name %q and %s", e.Message, tc.key, p)
			}
			if !strings.Contains(e.Hint, "keychain") {
				t.Errorf("hint %q should say where the secret belongs", e.Hint)
			}
		})
	}

	t.Run("accepts profile names that sound secret", func(t *testing.T) {
		f, err := Load(write(t, `{"profiles":{"github-token":{"client_id":"a"},"secrets-team":{"client_id":"b"}}}`))
		if err != nil {
			t.Fatal(err)
		}
		if len(f.Profiles) != 2 {
			t.Errorf("profiles = %v", f.Profiles)
		}
	})

	for _, tc := range []struct{ name, content string }{
		{"a second JSON value", `{"profiles":{}} {"client_secret":"x"}`},
		{"trailing garbage", `{"profiles":{}} garbage`},
	} {
		t.Run("rejects "+tc.name, func(t *testing.T) {
			_, err := Load(write(t, tc.content))
			wantCode(t, err, "bad_config", clierr.ExitUsage)
		})
	}

	t.Run("reports an unreadable file", func(t *testing.T) {
		_, err := Load(t.TempDir())
		if err == nil || !strings.Contains(err.Error(), "read config") {
			t.Fatalf("err = %v, want a read error", err)
		}
	})

	t.Run("rejects an unknown field", func(t *testing.T) {
		_, err := Load(write(t, `{"profiles":{"cego":{"clientid":"a"}}}`))
		wantCode(t, err, "bad_config", clierr.ExitUsage)
	})

	t.Run("rejects a profile that is not an object", func(t *testing.T) {
		_, err := Load(write(t, `{"profiles":{"cego":1}}`))
		wantCode(t, err, "bad_config", clierr.ExitUsage)
	})

	t.Run("rejects invalid JSON", func(t *testing.T) {
		_, err := Load(write(t, `{`))
		wantCode(t, err, "bad_config", clierr.ExitUsage)
	})
}

func TestResolve(t *testing.T) {
	file := File{
		DefaultProfile: "cego",
		Profiles: map[string]Profile{
			"cego": {ClientID: "AIK_CLIENT_cego", Region: "us"},
			"ci":   {ClientID: "AIK_CLIENT_ci"},
			"bare": {},
		},
	}
	pair := map[string]string{EnvClientID: "AIK_CLIENT_env", EnvClientSecret: "s"}
	tests := []struct {
		name  string
		file  File
		flags Flags
		vars  map[string]string
		want  Resolved
	}{
		{
			name: "--profile beats the environment pair and the default", file: file,
			flags: Flags{Profile: "ci"}, vars: pair,
			want: Resolved{Profile: "ci", ClientID: "AIK_CLIENT_ci", Region: "eu", Host: "app.aikido.dev"},
		},
		{
			name: "AIKIDO_DOJO_PROFILE beats the environment pair", file: file,
			vars: map[string]string{EnvProfile: "ci", EnvClientID: "AIK_CLIENT_env", EnvClientSecret: "s"},
			want: Resolved{Profile: "ci", ClientID: "AIK_CLIENT_ci", Region: "eu", Host: "app.aikido.dev"},
		},
		{
			name: "the environment pair beats default_profile and stores nothing", file: file, vars: pair,
			want: Resolved{ClientID: "AIK_CLIENT_env", Secret: "s", Region: "eu", Host: "app.aikido.dev"},
		},
		{
			name: "default_profile when nothing else is set", file: file,
			want: Resolved{Profile: "cego", ClientID: "AIK_CLIENT_cego", Region: "us", Host: "app.us.aikido.dev"},
		},
		{
			name: "AIKIDO_DOJO_REGION beats the profile", file: file,
			vars: map[string]string{EnvRegion: "me"},
			want: Resolved{Profile: "cego", ClientID: "AIK_CLIENT_cego", Region: "me", Host: "app.me.aikido.dev"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Resolve(tt.file, tt.flags, env(tt.vars))
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("got  %+v\nwant %+v", got, tt.want)
			}
		})
	}
}

func TestResolveFailures(t *testing.T) {
	file := File{Profiles: map[string]Profile{"bare": {}}}
	tests := []struct {
		name     string
		file     File
		flags    Flags
		vars     map[string]string
		code     string
		exit     int
		wantText string
	}{
		{name: "nothing configured", code: "no_credentials", exit: clierr.ExitAuth, wantText: EnvClientSecret},
		{
			name: "only the client ID is set", vars: map[string]string{EnvClientID: "AIK_CLIENT_env"},
			code: "incomplete_credentials", exit: clierr.ExitUsage, wantText: EnvClientSecret,
		},
		{
			name: "only the secret is set", vars: map[string]string{EnvClientSecret: "s"},
			code: "incomplete_credentials", exit: clierr.ExitUsage, wantText: EnvClientID,
		},
		{name: "unknown profile", file: file, flags: Flags{Profile: "nope"}, code: "unknown_profile", exit: clierr.ExitUsage, wantText: "nope"},
		{name: "profile without a client ID", file: file, flags: Flags{Profile: "bare"}, code: "bad_config", exit: clierr.ExitUsage, wantText: "client_id"},
		{
			name: "default_profile that does not exist", file: File{DefaultProfile: "gone"},
			code: "unknown_profile", exit: clierr.ExitUsage, wantText: "gone",
		},
		{
			name: "unknown region", vars: map[string]string{EnvClientID: "a", EnvClientSecret: "s", EnvRegion: "mars"},
			code: "unknown_region", exit: clierr.ExitUsage, wantText: "au, eu, me, us",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Resolve(tt.file, tt.flags, env(tt.vars))
			e := wantCode(t, err, tt.code, tt.exit)
			if !strings.Contains(e.Message+" "+e.Hint, tt.wantText) {
				t.Errorf("message/hint %q / %q should mention %q", e.Message, e.Hint, tt.wantText)
			}
		})
	}
}

func TestResolvedNeverEncodesTheSecret(t *testing.T) {
	b, err := json.Marshal(Resolved{ClientID: "id", Secret: "s3cr3t"})
	if err != nil || strings.Contains(string(b), "s3cr3t") {
		t.Errorf("Resolved encodes as %s, %v; want no secret", b, err)
	}
}

func TestSave(t *testing.T) {
	path := filepath.Join(t.TempDir(), "new", "config.json")
	f := File{DefaultProfile: "cego", Profiles: map[string]Profile{"cego": {ClientID: "id", Region: "us"}, "ci": {ClientID: "id2"}}}
	if err := Save(path, f); err != nil {
		t.Fatal(err)
	}
	got, err := Load(path)
	if err != nil || !reflect.DeepEqual(got, f) {
		t.Errorf("Load after Save = %+v, %v; want %+v", got, err, f)
	}
	for p, want := range map[string]os.FileMode{path: 0o600, filepath.Dir(path): 0o700} {
		info, err := os.Stat(p)
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm() != want {
			t.Errorf("%s: mode %v, want %v", p, info.Mode().Perm(), want)
		}
	}
	if entries, _ := os.ReadDir(filepath.Dir(path)); len(entries) != 1 {
		t.Errorf("the directory holds %d files, want only the config", len(entries))
	}
}

func TestSaveFailureLeavesNoTempFile(t *testing.T) {
	dir := t.TempDir()
	// A directory where the file should go makes the final rename fail.
	path := filepath.Join(dir, "config.json")
	if err := os.Mkdir(path, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := Save(path, File{}); err == nil || !strings.Contains(err.Error(), path) {
		t.Errorf("err = %v, want one naming %s", err, path)
	}
	if entries, _ := os.ReadDir(dir); len(entries) != 1 {
		t.Errorf("the directory holds %d entries, want no temp file left", len(entries))
	}
}

func TestRegionHost(t *testing.T) {
	if h, err := RegionHost("us"); err != nil || h != "app.us.aikido.dev" {
		t.Errorf("RegionHost(us) = %q, %v", h, err)
	}
	_, err := RegionHost("mars")
	wantCode(t, err, "unknown_region", clierr.ExitUsage)
}

// A config kept elsewhere and linked in, as dotfiles often are, stays linked,
// and an existing file keeps its mode.
func TestSaveWritesThroughASymlink(t *testing.T) {
	target := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(target, []byte("{}\n"), 0o644); err != nil { //nolint:gosec // the mode under test
		t.Fatal(err)
	}
	link := filepath.Join(t.TempDir(), "config.json")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	f := File{DefaultProfile: "cego", Profiles: map[string]Profile{"cego": {ClientID: "id"}}}
	if err := Save(link, f); err != nil {
		t.Fatal(err)
	}
	if info, err := os.Lstat(link); err != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Errorf("%s is no longer a symlink: %v", link, err)
	}
	if got, err := Load(target); err != nil || !reflect.DeepEqual(got, f) {
		t.Errorf("target = %+v, %v; want the saved file", got, err)
	}
	if info, err := os.Stat(target); err != nil || info.Mode().Perm() != 0o644 {
		t.Errorf("target mode = %v, %v; want 0644 kept", info.Mode().Perm(), err)
	}
}

// In a directory others can write, such as /tmp, a link may have been planted
// to aim the write at another file, so the link itself is replaced.
func TestSaveDoesNotFollowALinkInASharedDirectory(t *testing.T) {
	target := filepath.Join(t.TempDir(), "victim")
	if err := os.WriteFile(target, []byte("keep\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	shared := t.TempDir()
	if err := os.Chmod(shared, 0o777); err != nil { //nolint:gosec // a directory others can write is the case under test
		t.Fatal(err)
	}
	link := filepath.Join(shared, "config.json")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	if err := Save(link, File{}); err != nil {
		t.Fatal(err)
	}
	if data, err := os.ReadFile(target); err != nil || string(data) != "keep\n" { //nolint:gosec // the test's own temp file
		t.Errorf("target = %q, %v; want it untouched", data, err)
	}
	if info, err := os.Lstat(link); err != nil || info.Mode()&os.ModeSymlink != 0 {
		t.Errorf("%s: %v; want the link replaced by the config", link, err)
	}
}
