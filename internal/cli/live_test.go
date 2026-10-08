//go:build live

package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/config"
)

// Generated commands against a real workspace. Run with a read-only API client:
//
//	AIKIDO_DOJO_CLIENT_ID=… AIKIDO_DOJO_CLIENT_SECRET=… go test -tags live ./internal/cli/
func TestLiveCommands(t *testing.T) {
	env := liveEnv(t)
	tests := []struct {
		args  []string
		limit int // > 0: the output is an array of at most this many items
	}{
		{args: []string{"workspace", "get"}},
		{args: []string{"repo", "list", "--limit", "5"}, limit: 5},
		{args: []string{"container", "list", "--limit", "5"}, limit: 5},
		{args: []string{"pr-check-config", "list", "--limit", "5"}, limit: 5},
		{args: []string{"issue-group", "list", "--limit", "150"}, limit: 150},
		{args: []string{"cloud-asset", "list", "--limit", "5"}, limit: 5},
	}
	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			stdout, stderr, code := run(t, env, tt.args...)
			if code != 0 {
				t.Fatalf("exit %d: %s", code, stderr)
			}
			if tt.limit == 0 {
				if !json.Valid([]byte(stdout)) {
					t.Fatalf("stdout is not JSON: %.200s", stdout)
				}
				return
			}
			var items []json.RawMessage
			if err := json.Unmarshal([]byte(stdout), &items); err != nil || len(items) > tt.limit {
				t.Fatalf("stdout = %.200s (%v); want an array of at most %d", stdout, err, tt.limit)
			}
			t.Logf("%d items", len(items))
		})
	}
}

// liveEnv runs as the AIKIDO_DOJO_CLIENT_ID pair, with no config file, the
// real network, and one rate-limit window for the test's calls.
func liveEnv(t *testing.T) Env {
	t.Helper()
	if os.Getenv(config.EnvClientID) == "" || os.Getenv(config.EnvClientSecret) == "" {
		t.Skip("set " + config.EnvClientID + " and " + config.EnvClientSecret)
	}
	noConfig := filepath.Join(t.TempDir(), "config.json")
	cache := t.TempDir()
	return Env{
		Stdin: strings.NewReader(""),
		Getenv: func(k string) string {
			if k == config.EnvConfig {
				return noConfig
			}
			return os.Getenv(k)
		},
		CacheDir: func() (string, error) { return cache, nil },
	}
}

// Hand-written commands against a real workspace. Run with a read-only API client:
//
//	AIKIDO_DOJO_CLIENT_ID=… AIKIDO_DOJO_CLIENT_SECRET=… go test -tags live -run TestLiveHandWritten ./internal/cli/
func TestLiveHandWritten(t *testing.T) {
	env := liveEnv(t)
	t.Run("auth status", func(t *testing.T) {
		var st authStatus
		if err := json.Unmarshal([]byte(mustRun(t, env, "auth", "status")), &st); err != nil {
			t.Fatal(err)
		}
		if st.Source != "environment" || len(st.Scopes.Granted) == 0 || len(st.Scopes.Unknown) != 0 {
			t.Errorf("status = %+v, want the token's claimed scopes", st)
		}
		for _, s := range st.Scopes.Granted {
			if strings.HasSuffix(s, ":write") {
				t.Errorf("a READ pair holds %s", s)
			}
		}
		t.Logf("%d scopes granted, %d denied", len(st.Scopes.Granted), len(st.Scopes.Denied))
	})
	t.Run("api", func(t *testing.T) {
		var ws map[string]json.RawMessage
		if err := json.Unmarshal([]byte(mustRun(t, env, "api", "GET", "/workspace")), &ws); err != nil || len(ws) == 0 {
			t.Errorf("api GET /workspace = %v, %v; want a JSON object", ws, err)
		}
	})
	t.Run("repo current", func(t *testing.T) {
		var listed []struct {
			ID  json.RawMessage `json:"id"`
			URL string          `json:"url"`
		}
		if err := json.Unmarshal([]byte(mustRun(t, env, "repo", "list", "--limit", "1")), &listed); err != nil {
			t.Fatal(err)
		}
		if len(listed) == 0 {
			t.Skip("the workspace has no code repos")
		}
		inRepo(t, listed[0].URL)
		var found []repoMatch
		if err := json.Unmarshal([]byte(mustRun(t, env, "repo", "current")), &found); err != nil {
			t.Fatal(err)
		}
		var repo struct {
			ID json.RawMessage `json:"id"`
		}
		if len(found) != 1 || json.Unmarshal(found[0].Repo, &repo) != nil || string(repo.ID) != string(listed[0].ID) {
			t.Errorf("repo current found %d repos, want the one repo list gave", len(found))
		}
	})
}

// Output and safety flags against a real workspace, with GETs only.
func TestLiveOutput(t *testing.T) {
	env := liveEnv(t)
	t.Run("jq", func(t *testing.T) {
		out := strings.TrimSpace(mustRun(t, env, "repo", "list", "--limit", "2", "--jq", "length"))
		if n, err := strconv.Atoi(out); err != nil || n > 2 {
			t.Errorf("--jq length = %q, %v; want a number up to 2", out, err)
		}
	})
	t.Run("ndjson", func(t *testing.T) {
		out := mustRun(t, env, "issue-group", "list", "--limit", "3", "--ndjson")
		lines := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
		if out == "" {
			lines = nil
		}
		if len(lines) > 3 {
			t.Errorf("%d lines, want at most 3", len(lines))
		}
		for _, l := range lines {
			if !strings.HasPrefix(l, "{") || !json.Valid([]byte(l)) {
				t.Errorf("line %.80s is not a JSON object", l)
			}
		}
		t.Logf("%d lines", len(lines))
	})
	t.Run("read-only", func(t *testing.T) {
		if _, stderr, code := run(t, env, "--read-only", "team", "create", "--name", "x"); code != clierr.ExitRefused {
			t.Errorf("exit %d: %s; want 7", code, stderr)
		}
	})
}
