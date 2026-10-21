package docsfetch

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

// fullSite is llms-full.txt as parts a test can change between requests.
type fullSite struct {
	mu    sync.Mutex
	parts []string
}

func (s *fullSite) set(i int, text string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.parts[i] = text
}

// fullServer serves site's parts; part i links to part i+1 unless it is the last.
func fullServer(t *testing.T, site *fullSite, nextHost string) *httptest.Server {
	t.Helper()
	var srv *httptest.Server
	srv = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		site.mu.Lock()
		defer site.mu.Unlock()
		i := 0
		if r.URL.Path != "/llms-full.txt" {
			if _, err := fmt.Sscanf(r.URL.Path, "/llms-full.txt/%d", &i); err != nil || i >= len(site.parts) {
				http.NotFound(w, r)
				return
			}
		}
		host := srv.URL
		if nextHost != "" {
			host = nextHost
		}
		fmt.Fprintf(w, "> This is part %d of %d.\n", i+1, len(site.parts))
		if i+1 < len(site.parts) {
			fmt.Fprintf(w, "> Next part: %s/llms-full.txt/%d\n", host, i+1)
		}
		fmt.Fprint(w, "\n"+site.parts[i])
	}))
	t.Cleanup(srv.Close)
	return srv
}

// rawServer answers every path with body, so a test can serve what GitBook shouldn't.
func rawServer(t *testing.T, body func(path string) string) *httptest.Server {
	t.Helper()
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body(r.URL.Path))
	}))
	t.Cleanup(srv.Close)
	return srv
}

func digestOf(t *testing.T, srv *httptest.Server) (string, error) {
	t.Helper()
	return FullDigest(t.Context(), srv.Client(), srv.URL+"/llms-full.txt", mustURL(t, srv.URL))
}

func TestFullDigestCoversEveryPart(t *testing.T) {
	site := &fullSite{parts: []string{"# A\n", "# B\n"}}
	srv := fullServer(t, site, "")
	got, err := digestOf(t, srv)
	if err != nil {
		t.Fatal(err)
	}
	h := sha256.New()
	fmt.Fprintf(h, "> This is part 1 of 2.\n> Next part: %s/llms-full.txt/1\n\n# A\n", srv.URL)
	fmt.Fprint(h, "> This is part 2 of 2.\n\n# B\n")
	if want := hex.EncodeToString(h.Sum(nil)); got != want {
		t.Errorf("digest = %s, want %s", got, want)
	}
	site.set(1, "# B changed\n")
	if changed, err := digestOf(t, srv); err != nil || changed == got {
		t.Errorf("digest = %s, %v; want a change in the last part to move it", changed, err)
	}
}

// The changelog is the last page; a post there must not trigger a sync of the docs.
func TestFullDigestIgnoresTheChangelog(t *testing.T) {
	site := &fullSite{parts: []string{"# A\n", "# B\n\n# Changelog\n\nPost 1\n"}}
	srv := fullServer(t, site, "")
	before, err := digestOf(t, srv)
	if err != nil {
		t.Fatal(err)
	}
	site.set(1, "# B\n\n# Changelog\n\nPost 2\n")
	if got, err := digestOf(t, srv); err != nil || got != before {
		t.Errorf("digest = %s, %v; want a changelog post to leave it at %s", got, err, before)
	}
	site.set(1, "# B edited\n\n# Changelog\n\nPost 2\n")
	if got, err := digestOf(t, srv); err != nil || got == before {
		t.Errorf("digest = %s, %v; want a docs edit before the changelog to move it", got, err)
	}
}

func TestFullDigestRefusesAPartOffOrigin(t *testing.T) {
	srv := fullServer(t, &fullSite{parts: []string{"# A\n", "# B\n"}}, "https://evil.test")
	if _, err := digestOf(t, srv); err == nil || !strings.Contains(err.Error(), "is not on") {
		t.Errorf("err = %v, want the off-origin part refused", err)
	}
}

func TestFullDigestRefusesWhatIsNotLlmsFull(t *testing.T) {
	srv := rawServer(t, func(string) string { return "<html>Just a moment...</html>" })
	if _, err := digestOf(t, srv); err == nil || !strings.Contains(err.Error(), "not part 1 of GitBook's llms-full.txt") {
		t.Errorf("err = %v, want a challenge page refused", err)
	}
}

func TestFullDigestNeedsEveryPart(t *testing.T) {
	var srv *httptest.Server
	srv = rawServer(t, func(path string) string {
		if path == "/llms-full.txt" {
			return "> This is part 1 of 3.\n> Next part: " + srv.URL + "/llms-full.txt/1\n\n# A\n"
		}
		return "> This is part 2 of 3.\n\n# B\n"
	})
	if _, err := digestOf(t, srv); err == nil || !strings.Contains(err.Error(), "ends at part 2 of 3") {
		t.Errorf("err = %v, want a chain that stops early refused", err)
	}
}
