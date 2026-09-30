//go:build live

package api

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/cego/aikido-dojo/internal/auth"
	"github.com/cego/aikido-dojo/internal/config"
	"github.com/cego/aikido-dojo/internal/ratelimit"
)

// Read-only calls against a real workspace. Run with a read-only API client:
//
//	AIKIDO_DOJO_CLIENT_ID=… AIKIDO_DOJO_CLIENT_SECRET=… go test -tags live ./internal/api/
func TestLive(t *testing.T) {
	if os.Getenv(config.EnvClientID) == "" || os.Getenv(config.EnvClientSecret) == "" {
		t.Skip("set " + config.EnvClientID + " and " + config.EnvClientSecret)
	}
	r, err := config.Resolve(config.File{}, config.Flags{}, os.Getenv)
	if err != nil {
		t.Fatal(err)
	}
	hc := &http.Client{Timeout: 30 * time.Second}
	src, err := auth.NewSource(hc, r)
	if err != nil {
		t.Fatal(err)
	}
	c := New(hc, r.Host, src, ratelimit.New(t.TempDir(), r.ClientID))

	t.Run("workspace", func(t *testing.T) {
		resp, err := c.Do(t.Context(), Request{Method: http.MethodGet, Path: "/workspace", Scope: "basics:read"})
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		var ws map[string]json.RawMessage
		if err := json.NewDecoder(resp.Body).Decode(&ws); err != nil || len(ws) == 0 {
			t.Fatalf("workspace = %v, %v; want a JSON object", ws, err)
		}
	})

	t.Run("code repos page through", func(t *testing.T) {
		n := 0
		paging := Paging{Style: PageArray, SizeParam: "per_page", Size: 10}
		for _, err := range c.Items(t.Context(), Request{Method: http.MethodGet, Path: "/repositories/code", Scope: "repositories:read"}, paging) {
			if err != nil {
				t.Fatal(err)
			}
			if n++; n == 25 {
				break
			}
		}
		t.Logf("read %d repos", n)
	})
}
