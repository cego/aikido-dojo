package docsfetch

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
)

// fakeHelp is a help centre a test can change between runs; it counts requests per path.
type fakeHelp struct {
	srv   *httptest.Server
	mu    sync.Mutex
	order []string
	pages map[string]string
	mods  map[string]string
	hits  map[string]int
}

func newFakeHelp(t *testing.T) *fakeHelp {
	t.Helper()
	f := &fakeHelp{pages: map[string]string{}, mods: map[string]string{}, hits: map[string]int{}}
	f.srv = httptest.NewTLSServer(http.HandlerFunc(f.serve))
	t.Cleanup(f.srv.Close)
	f.set("readme", "# Overview\n\nWelcome.", "2026-10-01T00:00:00.000Z")
	f.set("guide/setup", "# Setup\n\nSteps.", "2026-10-01T00:00:00.000Z")
	f.set("guide/setup/language", "# Language\n\nPick one.", "2026-10-01T00:00:00.000Z")
	f.set("guide/group", "# Group\n\nChildren.", "") // group pages have no lastmod
	return f
}

// set adds or changes a page as GitBook serves it; an empty lastmod leaves it out of the sitemap's dates.
func (f *fakeHelp) set(path, body, lastmod string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.pages[path]; !ok {
		f.order = append(f.order, path)
	}
	f.pages[path] = "> For the complete documentation index, see llms.txt.\n\n" + body + "\n\n---\n\n# Agent Instructions\nAsk here.\n"
	if lastmod == "" {
		delete(f.mods, path)
	} else {
		f.mods[path] = lastmod
	}
}

func (f *fakeHelp) drop(path string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.pages, path)
	delete(f.mods, path)
	f.order = slices.DeleteFunc(f.order, func(p string) bool { return p == path })
}

func (f *fakeHelp) serve(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.hits[r.URL.Path]++
	half := len(f.order) / 2
	switch r.URL.Path {
	case "/llms.txt":
		fmt.Fprint(w, "# Aikido\n\n## Docs\n\n")
		for _, p := range f.order {
			fmt.Fprintf(w, "- [%s](%s/%s.md)\n", p, f.srv.URL, p)
		}
		fmt.Fprintf(w, "\n## Changelog\n\n- [Changelog](%s/changelog/readme.md)\n", f.srv.URL)
	case "/sitemap-pages.xml":
		fmt.Fprint(w, `<?xml version="1.0" encoding="utf-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
		for _, p := range f.order {
			loc := f.srv.URL + "/" + p
			if p == "readme" {
				loc = f.srv.URL
			}
			fmt.Fprintf(w, "<url><loc>%s</loc>", loc)
			if m, ok := f.mods[p]; ok {
				fmt.Fprintf(w, "<lastmod>%s</lastmod>", m)
			}
			fmt.Fprint(w, "</url>")
		}
		fmt.Fprint(w, "</urlset>")
	case "/llms-full.txt":
		fmt.Fprintf(w, "> This is part 1 of 2.\n> Next part: %s/llms-full.txt/1\n\n", f.srv.URL)
		for _, p := range f.order[:half] {
			fmt.Fprintln(w, f.pages[p])
		}
	case "/llms-full.txt/1":
		fmt.Fprint(w, "> This is part 2 of 2.\n\n")
		for _, p := range f.order[half:] {
			fmt.Fprintln(w, f.pages[p])
		}
	default:
		body, ok := f.pages[strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/"), ".md")]
		if !ok || !strings.HasSuffix(r.URL.Path, ".md") {
			http.NotFound(w, r)
			return
		}
		fmt.Fprint(w, body)
	}
}

func (f *fakeHelp) reset() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.hits = map[string]int{}
}

// requests copies the request counts under the lock, so -race sees no unsynchronised read.
func (f *fakeHelp) requests() map[string]int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return maps.Clone(f.hits)
}

// fetched lists the pages requested since the last reset.
func (f *fakeHelp) fetched() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	var got []string
	for p := range f.hits {
		if strings.HasSuffix(p, ".md") {
			got = append(got, strings.TrimSuffix(strings.TrimPrefix(p, "/"), ".md"))
		}
	}
	slices.Sort(got)
	return got
}

func (f *fakeHelp) sync(t *testing.T, dir string, full bool) (Report, error) {
	t.Helper()
	return Sync(t.Context(), f.srv.Client(), Options{Origin: f.srv.URL, Dir: dir, Full: full})
}

func (f *fakeHelp) mustSync(t *testing.T, dir string, full bool) Report {
	t.Helper()
	r, err := f.sync(t, dir, full)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func mustRead(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(name) //nolint:gosec // a file in the test's own temp dir
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func readPage(t *testing.T, dir, path string) string {
	t.Helper()
	return string(mustRead(t, filepath.Join(dir, "pages", filepath.FromSlash(path)+".md")))
}

func TestSyncFirstRunFetchesEveryPage(t *testing.T) {
	f, dir := newFakeHelp(t), t.TempDir()
	r := f.mustSync(t, dir, false)
	if !r.Full || r.Fetched != 4 || len(r.Changed) != 4 {
		t.Errorf("report = %+v, want a full run that wrote all 4 pages", r)
	}
	if got := readPage(t, dir, "guide/setup/language"); got != "# Language\n\nPick one.\n" {
		t.Errorf("language page = %q, want it normalised", got)
	}
	var snap Snapshot
	if err := json.Unmarshal(mustRead(t, filepath.Join(dir, "snapshot.json")), &snap); err != nil {
		t.Fatal(err)
	}
	if snap.FullSHA256 == "" || len(snap.LastMod) != 3 || snap.LastMod["readme"] != "2026-10-01T00:00:00.000Z" {
		t.Errorf("snapshot = %+v, want the digest and the 3 lastmods the sitemap gave", snap)
	}
}

func TestSyncUnchangedAsksOnlyForLlmsFull(t *testing.T) {
	f, dir := newFakeHelp(t), t.TempDir()
	f.mustSync(t, dir, false)
	f.reset()
	r := f.mustSync(t, dir, false)
	if !r.Unchanged || r.Fetched != 0 {
		t.Errorf("report = %+v, want unchanged", r)
	}
	if got, want := f.requests(), map[string]int{"/llms-full.txt": 1, "/llms-full.txt/1": 1}; !maps.Equal(got, want) {
		t.Errorf("requests = %v, want only %v", got, want)
	}
}

func TestSyncRefetchesOnlyChangedPages(t *testing.T) {
	f, dir := newFakeHelp(t), t.TempDir()
	f.mustSync(t, dir, false)
	f.set("guide/setup/language", "# Language\n\nPick one of six.", "2026-10-12T16:15:53.385Z")
	f.reset()
	r := f.mustSync(t, dir, false)
	if want := []string{"guide/group", "guide/setup/language"}; !slices.Equal(f.fetched(), want) {
		t.Errorf("fetched %q, want only the changed page and the page without a lastmod %q", f.fetched(), want)
	}
	if r.Full || !slices.Equal(r.Changed, []string{"guide/setup/language"}) {
		t.Errorf("report = %+v, want one changed page", r)
	}
	if got := readPage(t, dir, "guide/setup/language"); got != "# Language\n\nPick one of six.\n" {
		t.Errorf("language page = %q", got)
	}
}

func TestSyncChecksEveryPageWhenNoLastmodMoved(t *testing.T) {
	f, dir := newFakeHelp(t), t.TempDir()
	f.mustSync(t, dir, false)
	f.set("guide/setup", "# Setup\n\nNew steps.", "2026-10-01T00:00:00.000Z")
	f.reset()
	r := f.mustSync(t, dir, false)
	if !r.Full || !slices.Equal(r.Changed, []string{"guide/setup"}) || len(f.fetched()) != 4 {
		t.Errorf("report = %+v, fetched %q; want a full check that found guide/setup", r, f.fetched())
	}
}

func TestSyncRemovesUnlistedPages(t *testing.T) {
	f, dir := newFakeHelp(t), t.TempDir()
	f.mustSync(t, dir, false)
	f.drop("guide/setup/language")
	r := f.mustSync(t, dir, false)
	if !slices.Equal(r.Removed, []string{"guide/setup/language"}) {
		t.Errorf("removed = %q, want the dropped page", r.Removed)
	}
	if _, err := os.Stat(filepath.Join(dir, "pages", "guide", "setup")); !os.IsNotExist(err) {
		t.Errorf("the emptied directory is still there: %v", err)
	}
	if got := readPage(t, dir, "guide/setup"); got != "# Setup\n\nSteps.\n" {
		t.Errorf("the sibling page changed: %q", got)
	}
}

func TestSyncRefusesAnIndexThatLostMostPages(t *testing.T) {
	f, dir := newFakeHelp(t), t.TempDir()
	f.mustSync(t, dir, false)
	f.drop("guide/setup")
	f.drop("guide/setup/language")
	f.drop("guide/group")
	if _, err := f.sync(t, dir, false); err == nil || !strings.Contains(err.Error(), "refusing") {
		t.Fatalf("err = %v, want the shrunken index refused", err)
	}
	if got := readPage(t, dir, "guide/setup"); got != "# Setup\n\nSteps.\n" {
		t.Errorf("a held page was touched: %q", got)
	}
}

func TestSyncRefusesAPageThatIsNotMarkdown(t *testing.T) {
	f, dir := newFakeHelp(t), t.TempDir()
	f.mustSync(t, dir, false)
	before := mustRead(t, filepath.Join(dir, "snapshot.json"))
	f.set("guide/setup", "<html>Just a moment...</html>", "2026-10-12T00:00:00.000Z")
	if _, err := f.sync(t, dir, false); err == nil || !strings.Contains(err.Error(), "guide/setup: not a markdown page") {
		t.Fatalf("err = %v, want the challenge page refused", err)
	}
	after := mustRead(t, filepath.Join(dir, "snapshot.json"))
	if string(before) != string(after) || readPage(t, dir, "guide/setup") != "# Setup\n\nSteps.\n" {
		t.Error("a refused run changed the snapshot or a page; the next run wouldn't retry")
	}
}

func TestSyncFullRefetchesEveryPage(t *testing.T) {
	f, dir := newFakeHelp(t), t.TempDir()
	f.mustSync(t, dir, false)
	f.reset()
	r := f.mustSync(t, dir, true)
	if !r.Full || len(f.fetched()) != 4 || len(r.Changed) != 0 {
		t.Errorf("report = %+v, fetched %q; want all 4 fetched and none changed", r, f.fetched())
	}
}
