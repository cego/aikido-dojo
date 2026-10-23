package helpdocs

import (
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/cego/aikido-dojo/internal/helpdocs/helpdocstest"
)

const (
	language = "getting-started/setting-up-your-account/change-language"
	ips      = "getting-started/setting-up-your-account/limit-aikido-access-to-specific-ips"
)

// counted is a client for a fake help centre that counts the requests it sends.
type counted struct {
	next http.RoundTripper
	n    int
}

func (c *counted) RoundTrip(r *http.Request) (*http.Response, error) {
	c.n++
	return c.next.RoundTrip(r) //nolint:wrapcheck // http.Client wraps transport errors itself
}

func fakeHelp(t *testing.T, edits map[string]string) (*http.Client, *counted) {
	t.Helper()
	c := &counted{next: helpdocstest.Transport(t, helpdocstest.Handler(".", edits))}
	return &http.Client{Transport: c}, c
}

// languageIn is the change-language page of the copy s reads.
func languageIn(t *testing.T, s Store) Page {
	t.Helper()
	p, err := s.Get(language)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func mustSync(t *testing.T, s Store, c *http.Client) SyncReport {
	t.Helper()
	r, err := s.Sync(t.Context(), c)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestStoreReadsTheBuiltInCopyUntilASync(t *testing.T) {
	s := Store{Dir: filepath.Join(t.TempDir(), "docs"), Build: "v1"}
	if got := languageIn(t, s).Markdown; got != languageIn(t, Store{}).Markdown {
		t.Errorf("markdown = %q, want the built-in page", got)
	}
}

func TestSyncFetchesWhatAikidoChanged(t *testing.T) {
	edited := "# Change Language\n\nPick a language under zebracorn settings.\n"
	c, _ := fakeHelp(t, map[string]string{language: edited})
	s := Store{Dir: filepath.Join(t.TempDir(), "docs"), Build: "v1"}
	r := mustSync(t, s, c)
	if r.Unchanged || !slices.Equal(r.Changed, []string{language}) || r.Snapshot != "2099-01-01" || r.Requests < 4 {
		t.Errorf("report = %+v, want %s changed, snapshot 2099-01-01 and its requests counted", r, language)
	}
	if p := languageIn(t, s); p.Markdown != edited || p.Snapshot != "2099-01-01" {
		t.Errorf("docs get = %q at %s, want the edit at 2099-01-01", p.Markdown, p.Snapshot)
	}
	hits, err := s.Search("zebracorn", 1)
	if err != nil || len(hits) != 1 || hits[0].Path != language {
		t.Errorf("search zebracorn = %+v, %v; want %s", hits, err, language)
	}
	if got := languageIn(t, Store{}).Markdown; got == edited {
		t.Error("the built-in copy changed")
	}
}

func TestSyncWhenNothingChangedAsksOnlyForLlmsFull(t *testing.T) {
	edited := "# Change Language\n\nPick a language under zebracorn settings.\n"
	c, n := fakeHelp(t, map[string]string{language: edited})
	s := Store{Dir: filepath.Join(t.TempDir(), "docs"), Build: "v1"}
	mustSync(t, s, c)
	n.n = 0
	if r := mustSync(t, s, c); !r.Unchanged || r.Requests != 1 || n.n != 1 || r.Snapshot != "2099-01-01" {
		t.Errorf("report = %+v after %d requests; want unchanged after the one part of llms-full.txt", r, n.n)
	}
	if got := languageIn(t, s).Markdown; got != edited {
		t.Errorf("markdown = %q, want the edit kept", got)
	}
}

// An update of aikido-dojo brings back its built-in copy, and its next sync starts from that copy.
func TestAnUpdateBringsBackTheBuiltInCopy(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "docs")
	c1, _ := fakeHelp(t, map[string]string{language: "# Change Language\n\nThe v1 edit.\n"})
	mustSync(t, Store{Dir: dir, Build: "v1"}, c1)

	v2 := Store{Dir: dir, Build: "v2"}
	if got := languageIn(t, v2).Markdown; got != languageIn(t, Store{}).Markdown {
		t.Errorf("after the update, markdown = %q; want the built-in page", got)
	}
	c2, _ := fakeHelp(t, map[string]string{ips: "# Limit Aikido Access to Specific IPs\n\nThe v2 edit.\n"})
	if r := mustSync(t, v2, c2); !slices.Equal(r.Changed, []string{ips}) {
		t.Errorf("changed = %q, want only %s: the sync must start from the built-in copy", r.Changed, ips)
	}
	if got := languageIn(t, v2).Markdown; got != languageIn(t, Store{}).Markdown {
		t.Errorf("markdown = %q, want the built-in page, not v1's edit", got)
	}
}

func TestAFailedSyncKeepsTheCopy(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "docs")
	s := Store{Dir: dir, Build: "v1"}
	c, _ := fakeHelp(t, map[string]string{language: "# Change Language\n\nKept.\n"})
	mustSync(t, s, c)
	broken, _ := fakeHelp(t, map[string]string{ips: "Not markdown."})
	if _, err := s.Sync(t.Context(), broken); err == nil {
		t.Fatal("synced a page that is not markdown")
	}
	if got := languageIn(t, s).Markdown; got != "# Change Language\n\nKept.\n" {
		t.Errorf("markdown = %q, want the copy from before the failed sync", got)
	}
	entries, err := os.ReadDir(filepath.Dir(dir))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Errorf("the cache holds %v, want only docs", entries)
	}
}
