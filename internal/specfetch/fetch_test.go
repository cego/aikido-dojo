package specfetch

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// docsServer serves an llms.txt listing pages in order, plus one handler per page path.
func docsServer(t *testing.T, pages []string, handlers map[string]http.HandlerFunc) *httptest.Server {
	t.Helper()
	var srv *httptest.Server
	srv = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/llms.txt" {
			var b strings.Builder
			b.WriteString("## API Reference: Issues\n")
			for _, p := range pages {
				fmt.Fprintf(&b, "- [%s](%s%s): page\n", p, srv.URL, p)
			}
			fmt.Fprint(w, b.String())
			return
		}
		h, ok := handlers[r.URL.Path]
		if !ok {
			http.NotFound(w, r)
			return
		}
		h(w, r)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func servePage(updated, doc string) http.HandlerFunc {
	body := "---\nupdatedAt: " + updated + "\n---\n\n# Title\n\n# OpenAPI definition\n\n```json\n" + doc + "\n```\n"
	return func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, body) }
}

func opDoc(path string) string {
	return `{"openapi":"3.1.0","info":{"title":"T","version":"1"},"paths":{"` + path + `":{"get":{"operationId":"x"}}}}`
}

func TestAssembleMergesEveryListedPage(t *testing.T) {
	srv := docsServer(t, []string{"/reference/a.md", "/reference/b.md"}, map[string]http.HandlerFunc{
		"/reference/a.md": servePage("2026-05-27T10:04:31.000Z", opDoc("/a")),
		"/reference/b.md": servePage("2026-06-01T08:00:00.000Z", opDoc("/b")),
	})
	spec, updated, err := Assemble(t.Context(), srv.Client(), srv.URL+"/llms.txt")
	if err != nil {
		t.Fatal(err)
	}
	paths, err := object(spec["paths"], "paths")
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 2 || paths["/a"] == nil || paths["/b"] == nil {
		t.Errorf("paths = %v, want /a and /b", paths)
	}
	if want := time.Date(2026, 6, 1, 8, 0, 0, 0, time.UTC); !updated.Equal(want) {
		t.Errorf("updated = %s, want %s", updated, want)
	}
}

func TestAssembleFailures(t *testing.T) {
	blockUntilCancelled := func(_ http.ResponseWriter, r *http.Request) { <-r.Context().Done() }
	tests := []struct {
		name     string
		pages    []string
		handlers map[string]http.HandlerFunc
		wantErr  string
	}{
		{
			name:     "names a missing page",
			pages:    []string{"/reference/a.md", "/reference/b.md"},
			handlers: map[string]http.HandlerFunc{"/reference/a.md": servePage("2026-05-27T10:04:31Z", opDoc("/a"))},
			wantErr:  "/reference/b.md: HTTP 404",
		},
		{
			name:  "fails fast while other pages hang",
			pages: []string{"/reference/slow1.md", "/reference/slow2.md", "/reference/bad.md"},
			handlers: map[string]http.HandlerFunc{
				"/reference/slow1.md": blockUntilCancelled,
				"/reference/slow2.md": blockUntilCancelled,
				"/reference/bad.md":   func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusInternalServerError) },
			},
			wantErr: "/reference/bad.md: HTTP 500",
		},
		{
			name:  "refuses to follow a redirect",
			pages: []string{"/reference/a.md"},
			handlers: map[string]http.HandlerFunc{"/reference/a.md": func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "https://elsewhere.test/evil.md", http.StatusFound)
			}},
			wantErr: "refusing redirect to https://elsewhere.test/evil.md",
		},
		{
			name:  "rejects an oversized page",
			pages: []string{"/reference/big.md"},
			handlers: map[string]http.HandlerFunc{"/reference/big.md": func(w http.ResponseWriter, _ *http.Request) {
				fmt.Fprint(w, strings.Repeat("x", maxBodyBytes+1))
			}},
			wantErr: "exceeds",
		},
		{
			name:  "names a page without a spec",
			pages: []string{"/reference/prose.md"},
			handlers: map[string]http.HandlerFunc{"/reference/prose.md": func(w http.ResponseWriter, _ *http.Request) {
				fmt.Fprint(w, "---\nupdatedAt: 2026-05-27T10:04:31Z\n---\n# Guide\n")
			}},
			wantErr: `/reference/prose.md: no "# OpenAPI definition" heading`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := docsServer(t, tt.pages, tt.handlers)
			// Assemble must return long before this deadline; reaching it means it waited on the hanging pages.
			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancel()
			_, _, err := Assemble(ctx, srv.Client(), srv.URL+"/llms.txt")
			if ctx.Err() != nil {
				t.Fatal("Assemble returned only at the test deadline")
			}
			wantErr(t, err, tt.wantErr)
		})
	}
}
