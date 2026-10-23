// Package helpdocstest serves a help centre for tests: a copy of its pages, as GitBook serves help.aikido.dev, with the edits a test makes.
package helpdocstest

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"maps"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
)

const origin = "https://help.aikido.dev"

// Edited is the lastmod of every page a test edits, later than any real one.
const Edited = "2099-01-01T00:00:00.000Z"

// Handler serves the copy in dir, its pages/ and snapshot.json, with edits replacing or adding pages by path.
func Handler(dir string, edits map[string]string) http.Handler {
	site := sync.OnceValues(func() (*help, error) { return read(dir, edits) })
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h, err := site()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		p := strings.TrimPrefix(r.URL.Path, "/")
		switch {
		case p == "llms-full.txt":
			fmt.Fprint(w, "> This is part 1 of 1\n")
			for _, path := range slices.Sorted(maps.Keys(edits)) {
				fmt.Fprintf(w, "\n%s\n%s", path, edits[path])
			}
		case p == "llms.txt":
			fmt.Fprint(w, "# Aikido\n\n## Docs\n\n")
			for _, path := range slices.Sorted(maps.Keys(h.pages)) {
				fmt.Fprintf(w, "- [%s](%s/%s.md)\n", path, origin, path)
			}
		case p == "sitemap-pages.xml":
			fmt.Fprint(w, "<urlset>")
			for _, path := range slices.Sorted(maps.Keys(h.mods)) {
				loc := origin + "/" + path
				if path == "readme" {
					loc = origin + "/"
				}
				fmt.Fprintf(w, "<url><loc>%s</loc><lastmod>%s</lastmod></url>", loc, h.mods[path])
			}
			fmt.Fprint(w, "</urlset>")
		case strings.HasSuffix(p, ".md") && h.pages[strings.TrimSuffix(p, ".md")] != "":
			fmt.Fprint(w, h.pages[strings.TrimSuffix(p, ".md")])
		default:
			http.NotFound(w, r)
		}
	})
}

type help struct {
	pages map[string]string
	mods  map[string]string
}

func read(dir string, edits map[string]string) (*help, error) {
	data, err := os.ReadFile(filepath.Join(dir, "snapshot.json")) //nolint:gosec // the test's own copy
	if err != nil {
		return nil, fmt.Errorf("read the snapshot: %w", err)
	}
	var snap struct {
		LastMod map[string]string `json:"lastmod"`
	}
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("decode the snapshot: %w", err)
	}
	h := &help{pages: map[string]string{}, mods: snap.LastMod}
	pages := os.DirFS(filepath.Join(dir, "pages"))
	err = fs.WalkDir(pages, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		md, err := fs.ReadFile(pages, p)
		if err != nil {
			return fmt.Errorf("read %s: %w", p, err)
		}
		h.pages[strings.TrimSuffix(p, ".md")] = string(md)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("read the pages: %w", err)
	}
	for p, md := range edits {
		h.pages[p], h.mods[p] = md, Edited
	}
	return h, nil
}

// Transport sends every request to a TLS server running h, whatever host the request names.
func Transport(t testing.TB, h http.Handler) http.RoundTripper {
	t.Helper()
	srv := httptest.NewTLSServer(h)
	t.Cleanup(srv.Close)
	tr, ok := srv.Client().Transport.(*http.Transport)
	if !ok {
		t.Fatal("httptest's client has no *http.Transport")
	}
	tr = tr.Clone()
	tr.DialContext = func(ctx context.Context, network, _ string) (net.Conn, error) {
		var d net.Dialer
		return d.DialContext(ctx, network, srv.Listener.Addr().String())
	}
	// httptest's certificate names example.com, whatever host the request meant.
	tr.TLSClientConfig.ServerName = "example.com"
	return tr
}
