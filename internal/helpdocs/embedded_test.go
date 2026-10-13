package helpdocs

import (
	"io/fs"
	"testing"

	"github.com/cego/aikido-dojo/internal/docsfetch"
)

func TestOriginMatchesTheFetcher(t *testing.T) {
	if Origin != docsfetch.Origin {
		t.Errorf("Origin = %s, docsfetch fetches from %s", Origin, docsfetch.Origin)
	}
}

// The real copy is whole, dated, and every page has its "# " title.
func TestEmbeddedCopy(t *testing.T) {
	l, err := embedded()
	if err != nil {
		t.Fatal(err)
	}
	if l.snapshot == "" {
		t.Error("the copy has no snapshot date")
	}
	pages := map[string]bool{}
	for _, s := range l.secs {
		pages[s.page.path] = true
		if s.page.title == s.page.path {
			t.Errorf("%s has no # title", s.page.path)
		}
	}
	if len(pages) < 500 {
		t.Errorf("the copy has %d pages, want the whole help centre", len(pages))
	}
	for p := range pages {
		md, err := fs.ReadFile(l.fsys, p+".md")
		if err != nil {
			t.Fatal(err)
		}
		if docsfetch.HasAgentInstructions(md) {
			t.Errorf("%s carries GitBook's instructions for agents", p)
		}
	}
}

func BenchmarkSearchFromCold(b *testing.B) {
	for b.Loop() {
		pages, err := fs.Sub(pagesFS, "pages")
		if err != nil {
			b.Fatal(err)
		}
		l, err := load(pages, snapshotJSON)
		if err != nil {
			b.Fatal(err)
		}
		l.search("limit access to specific IPs", 10)
	}
}
