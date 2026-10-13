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

// fullServer serves llms-full.txt as parts; part i links to part i+1 unless it is the last.
func fullServer(t *testing.T, parts []string, nextHost string) *httptest.Server {
	t.Helper()
	var srv *httptest.Server
	srv = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		i := 0
		if r.URL.Path != "/llms-full.txt" {
			if _, err := fmt.Sscanf(r.URL.Path, "/llms-full.txt/%d", &i); err != nil || i >= len(parts) {
				http.NotFound(w, r)
				return
			}
		}
		host := srv.URL
		if nextHost != "" {
			host = nextHost
		}
		fmt.Fprintf(w, "> This is part %d of %d.\n", i+1, len(parts))
		if i+1 < len(parts) {
			fmt.Fprintf(w, "> Next part: %s/llms-full.txt/%d\n", host, i+1)
		}
		fmt.Fprint(w, "\n"+parts[i])
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestFullDigestCoversEveryPart(t *testing.T) {
	srv := fullServer(t, []string{"# A\n", "# B\n"}, "")
	got, err := FullDigest(t.Context(), srv.Client(), srv.URL+"/llms-full.txt", mustURL(t, srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	h := sha256.New()
	fmt.Fprintf(h, "> This is part 1 of 2.\n> Next part: %s/llms-full.txt/1\n\n# A\n", srv.URL)
	fmt.Fprint(h, "> This is part 2 of 2.\n\n# B\n")
	if want := hex.EncodeToString(h.Sum(nil)); got != want {
		t.Errorf("digest = %s, want %s", got, want)
	}
	other := fullServer(t, []string{"# A\n", "# B changed\n"}, "")
	changed, err := FullDigest(t.Context(), other.Client(), other.URL+"/llms-full.txt", mustURL(t, other.URL))
	if err != nil {
		t.Fatal(err)
	}
	if changed == got {
		t.Error("a change in the last part left the digest unchanged")
	}
}

func TestFullDigestRefusesAPartOffOrigin(t *testing.T) {
	srv := fullServer(t, []string{"# A\n", "# B\n"}, "https://evil.test")
	_, err := FullDigest(t.Context(), srv.Client(), srv.URL+"/llms-full.txt", mustURL(t, srv.URL))
	if err == nil || !strings.Contains(err.Error(), "is not on") {
		t.Errorf("err = %v, want the off-origin part refused", err)
	}
}
