package helpdocs

import (
	"errors"
	"slices"
	"strings"
	"testing"
	"testing/fstest"
)

func testLibrary(t *testing.T) *library {
	t.Helper()
	fsys := fstest.MapFS{
		"readme.md": {Data: []byte("# Aikido Docs Overview\n\nStart here.\n")},
		"getting-started/change-language.md": {Data: []byte("# Change Language\n\nSet the app language from your profile.\n\n" +
			"## Change your language\n\n1. Go to [**My Profile**](https://app.aikido.dev/my-profile).\n2. Pick a language.\n")},
		"zen/install-node.md": {Data: []byte("# Install Zen for Node.js\n\nZen protects your app.\n\n## Install\n\n" +
			"```bash\n## not a heading\nnpm install @aikidosec/firewall\n```\n\n## Configure the token\n\nSet AIKIDO_TOKEN to your firewall token.\n")},
	}
	snap := []byte(`{"full_sha256":"x","lastmod":{"readme":"2026-10-01T00:00:00.000Z","zen/install-node":"2026-10-12T16:15:53.385Z"}}`)
	l, err := load(fsys, snap)
	if err != nil {
		t.Fatal(err)
	}
	return l
}

func TestSplitIgnoresHeadingsInCode(t *testing.T) {
	info := &pageInfo{path: "zen/install-node"}
	secs := split(info, "# Install Zen\n\n## Install\n\n```bash\n## not a heading\n```\n\n## Configure\n\nText.\n")
	var headings []string
	for _, s := range secs {
		headings = append(headings, s.heading)
	}
	if want := []string{"Install", "Configure"}; !slices.Equal(headings, want) || info.title != "Install Zen" {
		t.Errorf("headings %q, title %q; want %q and Install Zen", headings, info.title, want)
	}
}

func TestSearchFindsTheBestSectionOfEachPage(t *testing.T) {
	hits := testLibrary(t).search("firewall token", 10)
	if len(hits) != 1 {
		t.Fatalf("hits = %+v, want only the Zen page, once", hits)
	}
	want := Hit{Path: "zen/install-node", Title: "Install Zen for Node.js", Section: "Configure the token",
		URL: "https://help.aikido.dev/zen/install-node", Snippet: "Set AIKIDO_TOKEN to your firewall token.", Score: hits[0].Score}
	if hits[0] != want || hits[0].Score <= 0 {
		t.Errorf("hit = %+v, want %+v", hits[0], want)
	}
}

func TestSearchSnippetDropsLinkTargets(t *testing.T) {
	hits := testLibrary(t).search("change language", 10)
	if len(hits) == 0 || hits[0].Path != "getting-started/change-language" {
		t.Fatalf("hits = %+v, want the language page first", hits)
	}
	if strings.Contains(hits[0].Snippet, "](") || strings.Contains(hits[0].Snippet, "**") {
		t.Errorf("snippet = %q, want plain text", hits[0].Snippet)
	}
}

func TestSearchWithoutAMatchIsEmpty(t *testing.T) {
	if hits := testLibrary(t).search("zzqqxx", 10); hits == nil || len(hits) != 0 {
		t.Errorf("hits = %#v, want an empty, non-nil slice", hits)
	}
}

func TestGetAcceptsEveryWayOfNamingAPage(t *testing.T) {
	l := testLibrary(t)
	for _, ref := range []string{
		"getting-started/change-language",
		"getting-started/change-language.md",
		"/getting-started/change-language",
		"https://help.aikido.dev/getting-started/change-language",
		"https://help.aikido.dev/getting-started/change-language.md#change-your-language",
		"getting-started/change-language#change-your-language",
	} {
		p, err := l.get(ref)
		if err != nil {
			t.Errorf("get(%q): %v", ref, err)
			continue
		}
		if p.Path != "getting-started/change-language" || p.Title != "Change Language" || p.Snapshot != "2026-10-12" ||
			p.URL != "https://help.aikido.dev/getting-started/change-language" || !strings.HasPrefix(p.Markdown, "# Change Language\n") {
			t.Errorf("get(%q) = %+v", ref, p)
		}
	}
	for _, ref := range []string{"", "https://help.aikido.dev", "https://help.aikido.dev/"} {
		if p, err := l.get(ref); err != nil || p.Path != "readme" || p.URL != "https://help.aikido.dev" {
			t.Errorf("get(%q) = %+v, %v; want the overview", ref, p, err)
		}
	}
}

func TestGetRefusesWhatIsNotAPage(t *testing.T) {
	l := testLibrary(t)
	for _, ref := range []string{
		"nope",
		"getting-started",
		"../../etc/passwd",
		"getting-started/../readme",
		"https://evil.example/getting-started/change-language",
		"http://help.aikido.dev/readme",
	} {
		if _, err := l.get(ref); !errors.Is(err, ErrNotFound) {
			t.Errorf("get(%q) err = %v, want ErrNotFound", ref, err)
		}
	}
}
