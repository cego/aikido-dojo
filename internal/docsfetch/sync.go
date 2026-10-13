package docsfetch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/cego/aikido-dojo/internal/specfetch"
)

// Gentle on GitBook: a full run is ~750 pages.
const concurrency = 2

// Snapshot is snapshot.json: what the last run saw, so the next can tell what changed.
type Snapshot struct {
	FullSHA256 string            `json:"full_sha256"`
	LastMod    map[string]string `json:"lastmod"`
}

// Options says where to fetch from and where the copy lives.
type Options struct {
	Origin string // the help centre
	Dir    string // holds pages/ and snapshot.json
	Full   bool   // refetch every page
}

// Report says what a run did; the caller fills in Requests.
type Report struct {
	Unchanged bool     `json:"unchanged"`
	Full      bool     `json:"full"`
	Fetched   int      `json:"fetched"`
	Changed   []string `json:"changed"`
	Removed   []string `json:"removed"`
	Requests  int64    `json:"requests"`
}

// Sync brings the copy in o.Dir up to date, fetching pages only when llms-full.txt changed.
func Sync(ctx context.Context, c *http.Client, o Options) (Report, error) {
	origin, err := url.Parse(o.Origin)
	if err != nil {
		return Report{}, fmt.Errorf("parse origin: %w", err)
	}
	// A redirect would fetch a URL no index listed.
	noRedirects := *c
	noRedirects.CheckRedirect = func(req *http.Request, _ []*http.Request) error {
		return fmt.Errorf("refusing redirect to %s", req.URL)
	}
	c = &noRedirects
	snapPath := filepath.Join(o.Dir, "snapshot.json")
	old, err := readSnapshot(snapPath)
	if err != nil {
		return Report{}, err
	}
	digest, err := FullDigest(ctx, c, o.Origin+"/llms-full.txt", origin)
	if err != nil {
		return Report{}, err
	}
	if !o.Full && digest == old.FullSHA256 {
		return Report{Unchanged: true, Changed: []string{}, Removed: []string{}}, nil
	}
	index, err := specfetch.Get(ctx, c, o.Origin+"/llms.txt")
	if err != nil {
		return Report{}, err
	}
	paths, err := PagePaths(string(index), origin)
	if err != nil {
		return Report{}, fmt.Errorf("llms.txt: %w", err)
	}
	sitemap, err := specfetch.Get(ctx, c, o.Origin+"/sitemap-pages.xml")
	if err != nil {
		return Report{}, err
	}
	mods, err := LastMods(sitemap, origin)
	if err != nil {
		return Report{}, err
	}
	pagesDir := filepath.Join(o.Dir, "pages")
	if err := os.MkdirAll(pagesDir, 0o750); err != nil {
		return Report{}, fmt.Errorf("create %s: %w", pagesDir, err)
	}
	root, err := os.OpenRoot(pagesDir)
	if err != nil {
		return Report{}, fmt.Errorf("open %s: %w", pagesDir, err)
	}
	defer root.Close()
	held, err := heldPages(root)
	if err != nil {
		return Report{}, err
	}
	if len(paths) < len(held)/2 {
		return Report{}, fmt.Errorf("llms.txt lists %d pages but %d are held; refusing to delete most of the copy", len(paths), len(held))
	}
	gone := unlisted(held, paths)
	r := Report{Full: o.Full || old.FullSHA256 == ""}
	want := paths
	if !r.Full {
		want = stale(paths, mods, old.LastMod, held)
	}
	changed, err := update(ctx, c, o.Origin, root, want)
	if err != nil {
		return Report{}, err
	}
	r.Fetched = len(want)
	// The digest moved but no lastmod did, so one lagged: check every page.
	if !r.Full && len(changed) == 0 && len(gone) == 0 && !moved(paths, mods, old.LastMod) {
		rest := slices.DeleteFunc(slices.Clone(paths), func(p string) bool { return slices.Contains(want, p) })
		more, err := update(ctx, c, o.Origin, root, rest)
		if err != nil {
			return Report{}, err
		}
		changed = append(changed, more...)
		r.Full, r.Fetched = true, r.Fetched+len(rest)
	}
	if err := remove(root, gone); err != nil {
		return Report{}, err
	}
	keep := map[string]string{}
	for _, p := range paths {
		if m, ok := mods[p]; ok {
			keep[p] = m
		}
	}
	if err := writeSnapshot(snapPath, Snapshot{FullSHA256: digest, LastMod: keep}); err != nil {
		return Report{}, err
	}
	slices.Sort(changed)
	r.Changed, r.Removed = append([]string{}, changed...), gone
	return r, nil
}

// stale lists the pages to refetch: new, without a lastmod, or with a lastmod that moved.
func stale(paths []string, mods, old map[string]string, held map[string]bool) []string {
	var out []string
	for _, p := range paths {
		if !held[p] || mods[p] == "" || mods[p] != old[p] {
			out = append(out, p)
		}
	}
	return out
}

// moved reports whether the sitemap dates any page differently from the last run.
func moved(paths []string, mods, old map[string]string) bool {
	for _, p := range paths {
		if mods[p] != "" && mods[p] != old[p] {
			return true
		}
	}
	return false
}

// update fetches paths and writes those whose text changed; it writes nothing unless every page is markdown.
func update(ctx context.Context, c *http.Client, origin string, root *os.Root, paths []string) ([]string, error) {
	urls := make([]string, len(paths))
	for i, p := range paths {
		urls[i] = origin + "/" + p + ".md"
	}
	bodies, err := specfetch.GetAll(ctx, c, urls, concurrency)
	if err != nil {
		return nil, err
	}
	texts := make([][]byte, len(paths))
	for i, p := range paths {
		texts[i] = Normalize(bodies[urls[i]])
		if !bytes.HasPrefix(texts[i], []byte("# ")) {
			return nil, fmt.Errorf("%s: not a markdown page", p)
		}
	}
	var changed []string
	for i, p := range paths {
		name := filepath.FromSlash(p + ".md")
		if held, err := root.ReadFile(name); err == nil && bytes.Equal(held, texts[i]) {
			continue
		}
		if err := root.MkdirAll(filepath.Dir(name), 0o750); err != nil {
			return nil, fmt.Errorf("create %s: %w", filepath.Dir(name), err)
		}
		if err := root.WriteFile(name, texts[i], 0o600); err != nil {
			return nil, fmt.Errorf("write %s: %w", name, err)
		}
		changed = append(changed, p)
	}
	return changed, nil
}

func heldPages(root *os.Root) (map[string]bool, error) {
	held := map[string]bool{}
	err := fs.WalkDir(root.FS(), ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(p, ".md") {
			held[strings.TrimSuffix(p, ".md")] = true
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("list the held pages: %w", err)
	}
	return held, nil
}

func unlisted(held map[string]bool, paths []string) []string {
	gone := []string{}
	for p := range held {
		if !slices.Contains(paths, p) {
			gone = append(gone, p)
		}
	}
	slices.Sort(gone)
	return gone
}

// remove deletes pages and the directories their removal empties.
func remove(root *os.Root, gone []string) error {
	for _, p := range gone {
		name := filepath.FromSlash(p + ".md")
		if err := root.Remove(name); err != nil {
			return fmt.Errorf("remove %s: %w", name, err)
		}
		for dir := filepath.Dir(name); dir != "."; dir = filepath.Dir(dir) {
			entries, err := fs.ReadDir(root.FS(), filepath.ToSlash(dir))
			if err != nil {
				return fmt.Errorf("read %s: %w", dir, err)
			}
			if len(entries) > 0 {
				break
			}
			if err := root.Remove(dir); err != nil {
				return fmt.Errorf("remove %s: %w", dir, err)
			}
		}
	}
	return nil
}

func readSnapshot(path string) (Snapshot, error) {
	data, err := os.ReadFile(path) //nolint:gosec // the copy's own snapshot, at a path the caller chose
	if errors.Is(err, fs.ErrNotExist) {
		return Snapshot{}, nil
	}
	if err != nil {
		return Snapshot{}, fmt.Errorf("read %s: %w", path, err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return Snapshot{}, fmt.Errorf("decode %s: %w", path, err)
	}
	return s, nil
}

func writeSnapshot(path string, s Snapshot) error {
	data, err := specfetch.Canonical(s)
	if err != nil {
		return fmt.Errorf("render snapshot: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
