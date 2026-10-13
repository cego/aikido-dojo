package docsfetch

import (
	"net/url"
	"slices"
	"strings"
	"testing"
)

func mustURL(t *testing.T, s string) *url.URL {
	t.Helper()
	u, err := url.Parse(s)
	if err != nil {
		t.Fatal(err)
	}
	return u
}

const testIndex = `# Aikido

## Docs

- [Aikido Docs Overview](https://help.test/readme.md)
- [Setting Up Your Account](https://help.test/getting-started/setting-up-your-account.md)
- [Change Language](https://help.test/getting-started/setting-up-your-account/change-language.md): Change the language.
- [[Beta] Node.js ESM](https://help.test/zen/node.js-esm.md)
- [GCP + API](https://help.test/cloud/get-values-+-connect-gcp.md)

## Changelog

- [Changelog](https://help.test/changelog/readme.md)

## Querying This Documentation
`

func TestPagePaths(t *testing.T) {
	origin := mustURL(t, "https://help.test/llms.txt")
	tests := []struct {
		name, index string
		want        []string
		wantErr     string
	}{
		{name: "keeps docs pages in index order", index: testIndex, want: []string{
			"readme",
			"getting-started/setting-up-your-account",
			"getting-started/setting-up-your-account/change-language",
			"zen/node.js-esm",
			"cloud/get-values-+-connect-gcp",
		}},
		{name: "refuses a link off origin", index: "## Docs\n- [X](https://evil.test/x.md)\n", wantErr: "is not on https://help.test"},
		{name: "refuses a path that climbs", index: "## Docs\n- [X](https://help.test/a/../../x.md)\n", wantErr: "unexpected path"},
		{name: "refuses an escaped climb", index: "## Docs\n- [X](https://help.test/%2e%2e/x.md)\n", wantErr: "unexpected path"},
		{name: "refuses an entry it can't read", index: "## Docs\n- [X](https://help.test/x)\n", wantErr: "unrecognised index entry"},
		{name: "needs a docs page", index: "## Changelog\n- [C](https://help.test/changelog/readme.md)\n", wantErr: "no docs pages"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PagePaths(tt.index, origin)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("err = %v, want it to contain %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("PagePaths = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLastMods(t *testing.T) {
	origin := mustURL(t, "https://help.test/llms.txt")
	sitemap := []byte(`<?xml version="1.0" encoding="utf-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
<url><loc>https://help.test</loc><lastmod>2026-08-28T13:44:11.701Z</lastmod></url>
<url><loc>https://help.test/getting-started/setting-up-your-account</loc><lastmod>2025-06-10T18:59:14.563Z</lastmod></url>
<url><loc>https://help.test/zen/use-cases</loc></url>
</urlset>`)
	got, err := LastMods(sitemap, origin)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"readme": "2026-08-28T13:44:11.701Z", "getting-started/setting-up-your-account": "2025-06-10T18:59:14.563Z"}
	if len(got) != len(want) || got["readme"] != want["readme"] || got["getting-started/setting-up-your-account"] != want["getting-started/setting-up-your-account"] {
		t.Errorf("LastMods = %v, want %v", got, want)
	}
	if _, err := LastMods([]byte(`<urlset><url><loc>https://evil.test/x</loc></url></urlset>`), origin); err == nil || !strings.Contains(err.Error(), "is not on") {
		t.Errorf("err = %v, want a loc off origin refused", err)
	}
}

func TestNormalize(t *testing.T) {
	page := "> For the complete documentation index, see [llms.txt](https://help.test/llms.txt).\n\n" +
		"# Change Language\n\nPick one.\n\n" +
		"<figure><img src=\"https://files.test/a.png?alt=media\" alt=\"Language &amp; region dropdown\"><figcaption></figcaption></figure>\n\n" +
		"<figure><img src=\"https://files.test/b.png\" alt=\"\"><figcaption><p>The <strong>profile</strong> page</p></figcaption></figure>\n\n" +
		"<figure><img src=\"https://files.test/c.png\" alt=\"\"><figcaption></figcaption></figure>\n\n" +
		"Inline <img src=\"https://files.test/d.png\" alt=\"logo\"> here.\n\n" +
		"---\n\n# Agent Instructions\nQuery https://help.test/index.md?ask=x\n"
	want := "# Change Language\n\nPick one.\n\n[image: Language & region dropdown]\n\n[image: The profile page]\n\nInline [image: logo] here.\n"
	if got := string(Normalize([]byte(page))); got != want {
		t.Errorf("Normalize =\n%q\nwant\n%q", got, want)
	}
	md := "# Shots\n\n![Ignore dialog](https://files.test/e.png?alt=media)\n\n![](https://files.test/f.png)\n\nSee ![logo](https://files.test/g.png \"Logo\") and [a link](https://help.test/x).\n"
	if got, want := string(Normalize([]byte(md))), "# Shots\n\n[image: Ignore dialog]\n\nSee [image: logo] and [a link](https://help.test/x).\n"; got != want {
		t.Errorf("markdown images: Normalize =\n%q\nwant\n%q", got, want)
	}
	quoted := "# Footers\n\nGitBook ends pages with:\n\n---\n\n# Agent Instructions\nquoted\n\n---\n\n# Agent Instructions\nthe real footer\n"
	if got, want := string(Normalize([]byte(quoted))), "# Footers\n\nGitBook ends pages with:\n\n---\n\n# Agent Instructions\nquoted\n"; got != want {
		t.Errorf("a quoted footer: Normalize =\n%q\nwant\n%q", got, want)
	}
	if got := string(Normalize([]byte("# Plain\n\nText.\n"))); got != "# Plain\n\nText.\n" {
		t.Errorf("a page without GitBook's header and footer changed: %q", got)
	}
}
