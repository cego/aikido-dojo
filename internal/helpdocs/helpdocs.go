// Package helpdocs searches and reads the copy of Aikido's help centre built into aikido-dojo, or a newer one docs sync fetched; internal/tools/docsfetch refreshes the built-in copy.
package helpdocs

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"sync"
)

//go:embed all:pages
var pagesFS embed.FS

//go:embed snapshot.json
var snapshotJSON []byte

// Origin is the help centre the copy comes from.
const Origin = "https://help.aikido.dev"

// ErrNotFound means a ref names no page of the copy.
var ErrNotFound = errors.New("no such page")

// Hit is a page search found, at the section that matched best.
type Hit struct {
	Path    string  `json:"path"`
	Title   string  `json:"title"`
	Section string  `json:"section,omitempty"`
	URL     string  `json:"url"`
	Snippet string  `json:"snippet"`
	Score   float64 `json:"score"`
}

// Page is one page of the copy; Snapshot is the date of its newest page.
type Page struct {
	Path     string `json:"path"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	Snapshot string `json:"snapshot"`
	Markdown string `json:"markdown"`
}

// embedded indexes the copy on first use, so only docs commands pay for it.
var embedded = sync.OnceValues(func() (*library, error) {
	pages, err := fs.Sub(pagesFS, "pages")
	if err != nil {
		return nil, fmt.Errorf("open the built-in docs: %w", err)
	}
	return load(pages, snapshotJSON)
})
