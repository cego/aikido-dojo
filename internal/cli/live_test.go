//go:build live

package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cego/aikido-dojo/internal/config"
)

// Generated commands against a real workspace. Run with a read-only API client:
//
//	AIKIDO_DOJO_CLIENT_ID=… AIKIDO_DOJO_CLIENT_SECRET=… go test -tags live ./internal/cli/
func TestLiveCommands(t *testing.T) {
	if os.Getenv(config.EnvClientID) == "" || os.Getenv(config.EnvClientSecret) == "" {
		t.Skip("set " + config.EnvClientID + " and " + config.EnvClientSecret)
	}
	noConfig := filepath.Join(t.TempDir(), "config.json")
	cache := t.TempDir() // shared, so the calls below are paced by one window
	env := Env{
		Stdin: strings.NewReader(""),
		Getenv: func(k string) string {
			if k == config.EnvConfig {
				return noConfig
			}
			return os.Getenv(k)
		},
		CacheDir: func() (string, error) { return cache, nil },
	}
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
