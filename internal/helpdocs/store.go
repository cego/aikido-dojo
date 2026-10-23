package helpdocs

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/cego/aikido-dojo/internal/docsfetch"
)

// buildFile names the aikido-dojo build that synced a copy.
const buildFile = "build"

// Store is where docs sync keeps a copy newer than the built-in one; the zero Store reads only the built-in copy.
type Store struct {
	Dir   string // holds the synced copy: pages/, snapshot.json and build
	Build string // this aikido-dojo; a copy another build synced is ignored, so an update brings back the built-in copy
}

// SyncReport says what a sync fetched; Snapshot dates the copy docs commands now read.
type SyncReport struct {
	docsfetch.Report
	Snapshot string `json:"snapshot"`
}

// Search returns up to n pages for query, best first.
func (s Store) Search(query string, n int) ([]Hit, error) {
	l, err := s.library()
	if err != nil {
		return nil, err
	}
	return l.search(query, n), nil
}

// Get returns the page ref names, or ErrNotFound.
func (s Store) Get(ref string) (Page, error) {
	l, err := s.library()
	if err != nil {
		return Page{}, err
	}
	return l.get(ref)
}

func (s Store) library() (*library, error) {
	pages, snap, err := s.synced()
	if err != nil {
		return nil, err
	}
	if pages == nil {
		return embedded()
	}
	return load(pages, snap)
}

// synced returns the copy this build synced, or a nil fs.FS when it synced none.
func (s Store) synced() (fs.FS, []byte, error) {
	if s.Dir == "" {
		return nil, nil, nil
	}
	stamp, err := os.ReadFile(filepath.Join(s.Dir, buildFile))
	if errors.Is(err, fs.ErrNotExist) || (err == nil && string(stamp) != s.Build) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, fmt.Errorf("read the synced docs: %w", err)
	}
	snap, err := os.ReadFile(filepath.Join(s.Dir, "snapshot.json"))
	if err != nil {
		return nil, nil, fmt.Errorf("read the synced docs: %w", err)
	}
	// The root stays open for the command's reads, and keeps them inside the copy.
	root, err := os.OpenRoot(filepath.Join(s.Dir, "pages"))
	if err != nil {
		return nil, nil, fmt.Errorf("read the synced docs: %w", err)
	}
	return root.FS(), snap, nil
}

// Sync brings the copy docs commands read up to date from the help centre, asking only for the pages that changed.
func (s Store) Sync(ctx context.Context, c *http.Client) (SyncReport, error) {
	if s.Dir == "" {
		return SyncReport{}, errors.New("no directory to keep the synced docs in")
	}
	parent := filepath.Dir(s.Dir)
	if err := os.MkdirAll(parent, 0o700); err != nil {
		return SyncReport{}, fmt.Errorf("create %s: %w", parent, err)
	}
	// The sync works on a copy beside the old one and swaps it in whole, so no reader sees half a sync.
	tmp, err := os.MkdirTemp(parent, ".docs-sync-*")
	if err != nil {
		return SyncReport{}, fmt.Errorf("create a directory for the sync: %w", err)
	}
	defer os.RemoveAll(tmp)
	if err := s.seed(tmp); err != nil {
		return SyncReport{}, err
	}
	n := &counter{next: c.Transport}
	if n.next == nil {
		n.next = http.DefaultTransport
	}
	counted := *c
	counted.Transport = n
	r, err := docsfetch.Sync(ctx, &counted, docsfetch.Options{Origin: Origin, Dir: tmp})
	if err != nil {
		return SyncReport{}, err
	}
	r.Requests = n.n.Load()
	snap, err := os.ReadFile(filepath.Join(tmp, "snapshot.json")) //nolint:gosec // the directory the sync created
	if err != nil {
		return SyncReport{}, fmt.Errorf("read the synced snapshot: %w", err)
	}
	date, err := dated(snap)
	if err != nil {
		return SyncReport{}, err
	}
	if err := os.WriteFile(filepath.Join(tmp, buildFile), []byte(s.Build), 0o600); err != nil {
		return SyncReport{}, fmt.Errorf("stamp the synced docs: %w", err)
	}
	if err := swap(tmp, s.Dir); err != nil {
		return SyncReport{}, err
	}
	return SyncReport{Report: r, Snapshot: date}, nil
}

// seed writes the copy docs commands read now into dir, for the sync to bring up to date.
func (s Store) seed(dir string) error {
	pages, snap, err := s.synced()
	if err != nil {
		return err
	}
	if pages == nil {
		if pages, err = fs.Sub(pagesFS, "pages"); err != nil {
			return fmt.Errorf("open the built-in docs: %w", err)
		}
		snap = snapshotJSON
	}
	if err := os.CopyFS(filepath.Join(dir, "pages"), pages); err != nil {
		return fmt.Errorf("copy the docs: %w", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "snapshot.json"), snap, 0o600); err != nil {
		return fmt.Errorf("copy the docs snapshot: %w", err)
	}
	return nil
}

// swap moves the copy in tmp to dir; a reader between the renames finds no copy and reads the built-in one.
func swap(tmp, dir string) error {
	old := tmp + ".old"
	if err := os.Rename(dir, old); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("move the old docs aside: %w", err)
	}
	defer os.RemoveAll(old)
	if err := os.Rename(tmp, dir); err != nil {
		_ = os.Rename(old, dir)
		return fmt.Errorf("move the synced docs in place: %w", err)
	}
	return nil
}

// counter counts a sync's requests, so docs sync reports its load on the help centre.
type counter struct {
	next http.RoundTripper
	n    atomic.Int64
}

func (c *counter) RoundTrip(r *http.Request) (*http.Response, error) {
	c.n.Add(1)
	return c.next.RoundTrip(r) //nolint:wrapcheck // http.Client wraps transport errors itself
}
