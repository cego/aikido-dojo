// Package docsfetch keeps a copy of Aikido's help centre in the repository, asking the site for as little as it can.
package docsfetch

import (
	"bufio"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strings"
)

// Origin is the help centre; every URL a run fetches must be on it.
const Origin = "https://help.aikido.dev"

const (
	indexNote = "> For the complete documentation index"
	// GitBook appends instructions for agents to every page; they aren't Aikido's docs.
	agentFooter = "\n---\n\n# Agent Instructions"
)

var (
	entryRE = regexp.MustCompile(`^\s*- \[`)
	// The lazy title match lets titles contain brackets.
	linkRE = regexp.MustCompile(`^\s*- \[.*?\]\(([^)\s]+\.md)\)`)
	// A segment can't start with a dot, so no path climbs out of the copy.
	pathRE    = regexp.MustCompile(`^[a-z0-9][a-z0-9._+-]*(/[a-z0-9][a-z0-9._+-]*)*$`)
	figureRE  = regexp.MustCompile(`(?s)<figure>.*?</figure>`)
	imgRE     = regexp.MustCompile(`<img\b[^>]*>`)
	mdImageRE = regexp.MustCompile(`!\[([^\]]*)\]\([^)\s]+(?:\s+"[^"]*")?\)`)
	altRE     = regexp.MustCompile(`\balt="([^"]*)"`)
	captionRE = regexp.MustCompile(`(?s)<figcaption>(.*?)</figcaption>`)
	tagRE     = regexp.MustCompile(`<[^>]+>`)
	blanksRE  = regexp.MustCompile(`\n{3,}`)
)

// PagePaths returns the pages llms.txt lists under "## Docs", in index order, as paths without ".md".
func PagePaths(index string, origin *url.URL) ([]string, error) {
	var paths []string
	seen := map[string]bool{}
	inDocs := false
	sc := bufio.NewScanner(strings.NewReader(index))
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "## ") {
			inDocs = strings.TrimSpace(line) == "## Docs"
			continue
		}
		if !inDocs || !entryRE.MatchString(line) {
			continue
		}
		m := linkRE.FindStringSubmatch(line)
		if m == nil {
			return nil, fmt.Errorf("unrecognised index entry %q", line)
		}
		p, err := pathOf(m[1], origin)
		if err != nil {
			return nil, err
		}
		if !seen[p] {
			seen[p] = true
			paths = append(paths, p)
		}
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("scan index: %w", err)
	}
	if len(paths) == 0 {
		return nil, errors.New("index lists no docs pages")
	}
	return paths, nil
}

func pathOf(link string, origin *url.URL) (string, error) {
	u, err := url.Parse(link)
	if err != nil {
		return "", fmt.Errorf("parse link %q: %w", link, err)
	}
	if u.Scheme != origin.Scheme || u.Host != origin.Host {
		return "", fmt.Errorf("link %s is not on %s://%s", u, origin.Scheme, origin.Host)
	}
	p := strings.TrimSuffix(strings.TrimPrefix(u.Path, "/"), ".md")
	if !pathRE.MatchString(p) {
		return "", fmt.Errorf("link %s has an unexpected path", u)
	}
	return p, nil
}

type urlset struct {
	URLs []struct {
		Loc     string `xml:"loc"`
		LastMod string `xml:"lastmod"`
	} `xml:"url"`
}

// LastMods maps each sitemap page to its lastmod, leaving out pages without one; the site root is "readme".
func LastMods(sitemap []byte, origin *url.URL) (map[string]string, error) {
	var set urlset
	if err := xml.Unmarshal(sitemap, &set); err != nil {
		return nil, fmt.Errorf("parse sitemap: %w", err)
	}
	mods := map[string]string{}
	for _, e := range set.URLs {
		u, err := url.Parse(strings.TrimSpace(e.Loc))
		if err != nil {
			return nil, fmt.Errorf("parse sitemap loc %q: %w", e.Loc, err)
		}
		if u.Scheme != origin.Scheme || u.Host != origin.Host {
			return nil, fmt.Errorf("sitemap loc %s is not on %s://%s", u, origin.Scheme, origin.Host)
		}
		p := strings.Trim(u.Path, "/")
		if p == "" {
			p = "readme"
		}
		if m := strings.TrimSpace(e.LastMod); m != "" {
			mods[p] = m
		}
	}
	return mods, nil
}

// Normalize keeps only Aikido's text: no GitBook header or footer, images as their description.
func Normalize(page []byte) []byte {
	s := string(page)
	if strings.HasPrefix(s, indexNote) {
		_, s, _ = strings.Cut(s, "\n")
	}
	s, _, _ = strings.Cut(s, agentFooter)
	s = figureRE.ReplaceAllStringFunc(s, func(fig string) string {
		desc := text(altRE, fig)
		if desc == "" {
			desc = text(captionRE, fig)
		}
		return image(desc)
	})
	s = imgRE.ReplaceAllStringFunc(s, func(img string) string { return image(text(altRE, img)) })
	s = mdImageRE.ReplaceAllStringFunc(s, func(img string) string { return image(text(mdImageRE, img)) })
	s = blanksRE.ReplaceAllString(s, "\n\n")
	return []byte(strings.TrimSpace(s) + "\n")
}

func text(re *regexp.Regexp, s string) string {
	m := re.FindStringSubmatch(s)
	if m == nil {
		return ""
	}
	return strings.TrimSpace(html.UnescapeString(tagRE.ReplaceAllString(m[1], "")))
}

func image(desc string) string {
	if desc == "" {
		return ""
	}
	return "[image: " + desc + "]"
}
