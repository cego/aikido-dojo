package docsfetch

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
	return FullDigest(t.Context(), srv.Client(), srv.URL+"/llms-full.txt")
}

// The changelog is the last page; a post there must not trigger a sync of the docs.
func TestFullDigestIgnoresTheChangelog(t *testing.T) {
	body := "# A\n\n# Changelog\n\nPost 1\n"
	srv := rawServer(t, func(string) string { return body })
	got, err := digestOf(t, srv)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256([]byte("# A\n"))
	if want := hex.EncodeToString(sum[:]); got != want {
		t.Errorf("digest = %s, want %s", got, want)
	}
	body = "# A\n\n# Changelog\n\nPost 2\n"
	if again, err := digestOf(t, srv); err != nil || again != got {
		t.Errorf("digest = %s, %v; want a changelog post to leave it at %s", again, err, got)
	}
	body = "# A edited\n\n# Changelog\n\nPost 2\n"
	if edited, err := digestOf(t, srv); err != nil || edited == got {
		t.Errorf("digest = %s, %v; want a docs edit to move it", edited, err)
	}
}

func TestFullDigestRefusesWhatIsNotLlmsFull(t *testing.T) {
	srv := rawServer(t, func(string) string { return "<html>Just a moment...</html>" })
	if _, err := digestOf(t, srv); err == nil || !strings.Contains(err.Error(), "not the help centre's llms-full.txt") {
		t.Errorf("err = %v, want a challenge page refused", err)
	}
}
