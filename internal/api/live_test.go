//go:build live

package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/cego/aikido-dojo/internal/auth"
	"github.com/cego/aikido-dojo/internal/clierr"
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

	// Each style is read far enough to cross page boundaries; a repeated id is
	// what 1-indexed or overlapping pages would produce.
	t.Run("array paging: code repos", func(t *testing.T) {
		n := pageThrough(t, c, Request{Method: http.MethodGet, Path: "/repositories/code", Scope: "repositories:read"},
			Paging{Style: PageArray, SizeParam: "per_page", Size: 10}, 60)
		t.Logf("read %d repos", n)
	})

	t.Run("header paging: open issue groups", func(t *testing.T) {
		n := pageThrough(t, c, Request{Method: http.MethodGet, Path: "/open-issue-groups", Scope: "issues:read"},
			Paging{Style: PageHeader, SizeParam: "per_page", Size: 10}, 35)
		t.Logf("read %d issue groups", n)
	})

	t.Run("envelope paging: cloud assets", func(t *testing.T) {
		n := pageThrough(t, c, Request{Method: http.MethodGet, Path: "/clouds/assets", Scope: "clouds:read"},
			Paging{Style: PageEnvelope, SizeParam: "limit", Size: 10, Items: "assets", More: "hasMore"}, 25)
		t.Logf("read %d cloud assets", n)
	})
}

// pageThrough reads up to limit items and fails on an item without an id or
// with one seen before. A client without the scope skips the subtest.
func pageThrough(t *testing.T, c *Client, req Request, p Paging, limit int) int {
	t.Helper()
	seen := map[string]bool{}
	for item, err := range c.Items(t.Context(), req, p) {
		if err != nil {
			var e *clierr.Error
			if errors.As(err, &e) && e.Code == "missing_scope" {
				t.Skipf("the API client lacks %s", req.Scope)
			}
			t.Fatal(err)
		}
		var v struct {
			ID json.RawMessage `json:"id"`
		}
		if err := json.Unmarshal(item, &v); err != nil || len(v.ID) == 0 {
			t.Fatalf("item without an id: %.200s", item)
		}
		if seen[string(v.ID)] {
			t.Fatalf("id %s appeared twice: pages overlap", v.ID)
		}
		seen[string(v.ID)] = true
		if len(seen) == limit {
			break
		}
	}
	return len(seen)
}
