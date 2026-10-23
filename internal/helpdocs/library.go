package helpdocs

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"math"
	"net/url"
	"regexp"
	"slices"
	"strings"

	"github.com/cego/aikido-dojo/internal/search"
)

const snippetRunes = 200

var linkRE = regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`)

type pageInfo struct {
	path  string
	title string
}

// section is a page's text under one "## " heading, or before the first.
type section struct {
	page    *pageInfo
	heading string
	snippet string
	tf      map[string]int
	n       int
}

type library struct {
	fsys     fs.FS
	snapshot string
	secs     []section
	corpus   search.Corpus
}

func load(fsys fs.FS, snapshotJSON []byte) (*library, error) {
	date, err := dated(snapshotJSON)
	if err != nil {
		return nil, err
	}
	l := &library{fsys: fsys, snapshot: date, corpus: search.Corpus{DF: map[string]int{}}}
	total := 0
	err = fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(p, ".md") {
			return nil
		}
		md, err := fs.ReadFile(fsys, p)
		if err != nil {
			return fmt.Errorf("read %s: %w", p, err)
		}
		for _, s := range split(&pageInfo{path: strings.TrimSuffix(p, ".md")}, string(md)) {
			for t := range s.tf {
				l.corpus.DF[t]++
			}
			total += s.n
			l.secs = append(l.secs, s)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("index the docs: %w", err)
	}
	l.corpus.Docs = len(l.secs)
	if len(l.secs) > 0 {
		l.corpus.AvgLen = float64(total) / float64(len(l.secs))
	}
	return l, nil
}

// dated dates a copy by its snapshot's newest lastmod.
func dated(snapshotJSON []byte) (string, error) {
	var snap struct {
		LastMod map[string]string `json:"lastmod"`
	}
	if err := json.Unmarshal(snapshotJSON, &snap); err != nil {
		return "", fmt.Errorf("decode the docs snapshot: %w", err)
	}
	return newest(snap.LastMod), nil
}

// newest dates the copy by its newest lastmod; RFC 3339 UTC strings sort by time.
func newest(mods map[string]string) string {
	latest := ""
	for _, m := range mods {
		latest = max(latest, m)
	}
	return latest[:min(len(latest), len("2006-01-02"))]
}

// split cuts a page at its "## " headings outside code fences and sets info.title from its "# " line.
func split(info *pageInfo, md string) []section {
	var (
		secs    []section
		heading string
		body    []string
		fenced  bool
	)
	flush := func() {
		if heading != "" || strings.TrimSpace(strings.Join(body, "")) != "" {
			secs = append(secs, newSection(info, heading, body))
		}
	}
	for _, line := range strings.Split(md, "\n") {
		if t := strings.TrimSpace(line); strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~") {
			fenced = !fenced
		}
		switch {
		case !fenced && info.title == "" && strings.HasPrefix(line, "# "):
			info.title = strings.TrimSpace(line[2:])
		case !fenced && strings.HasPrefix(line, "## "):
			flush()
			heading, body = strings.TrimSpace(line[3:]), nil
		default:
			body = append(body, line)
		}
	}
	flush()
	if info.title == "" {
		info.title = info.path
	}
	return secs
}

// newSection counts the page title and section heading twice, so they outweigh a passing mention.
func newSection(info *pageInfo, heading string, body []string) section {
	s := section{page: info, heading: heading, snippet: snippet(body), tf: map[string]int{}}
	for _, text := range []string{info.title, info.title, heading, heading, strings.Join(body, "\n")} {
		for _, t := range search.Tokens(text) {
			s.tf[t]++
			s.n++
		}
	}
	return s
}

// snippet is a section's first line of prose, without link targets or emphasis.
func snippet(body []string) string {
	fenced := false
	for _, line := range body {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~") {
			fenced = !fenced
			continue
		}
		if fenced || t == "" || t == "---" || strings.ContainsAny(t[:1], "#<{|") {
			continue
		}
		t = strings.ReplaceAll(linkRE.ReplaceAllString(t, "$1"), "**", "")
		if r := []rune(t); len(r) > snippetRunes {
			t = string(r[:snippetRunes]) + "…"
		}
		return t
	}
	return ""
}

func (l *library) search(query string, n int) []Hit {
	terms := search.QueryTerms(query)
	type scored struct {
		s     *section
		score float64
	}
	var all []scored
	for i := range l.secs {
		if v := l.corpus.Score(terms, l.secs[i].tf, l.secs[i].n); v > 0 {
			all = append(all, scored{&l.secs[i], v})
		}
	}
	slices.SortStableFunc(all, func(x, y scored) int {
		return cmp.Or(cmp.Compare(y.score, x.score), strings.Compare(x.s.page.path, y.s.page.path))
	})
	hits := []Hit{}
	seen := map[string]bool{}
	for _, a := range all {
		if len(hits) == n {
			break
		}
		if seen[a.s.page.path] {
			continue
		}
		seen[a.s.page.path] = true
		hits = append(hits, Hit{Path: a.s.page.path, Title: a.s.page.title, Section: a.s.heading,
			URL: pageURL(a.s.page.path), Snippet: a.s.snippet, Score: math.Round(a.score*100) / 100})
	}
	return hits
}

func (l *library) get(ref string) (Page, error) {
	p, ok := pagePath(ref)
	if !ok {
		return Page{}, ErrNotFound
	}
	md, err := fs.ReadFile(l.fsys, p+".md")
	if errors.Is(err, fs.ErrNotExist) {
		return Page{}, ErrNotFound
	}
	if err != nil {
		return Page{}, fmt.Errorf("read %s: %w", p, err)
	}
	info := &pageInfo{path: p}
	split(info, string(md))
	return Page{Path: p, Title: info.title, URL: pageURL(p), Snapshot: l.snapshot, Markdown: string(md)}, nil
}

// pagePath turns a path, a path with ".md" or a help-centre URL into a page path of the copy.
func pagePath(ref string) (string, bool) {
	ref = strings.TrimSpace(ref)
	if strings.Contains(ref, "://") {
		u, err := url.Parse(ref)
		if err != nil || u.Scheme+"://"+u.Host != Origin {
			return "", false
		}
		ref = u.Path
	}
	ref, _, _ = strings.Cut(ref, "#")
	ref, _, _ = strings.Cut(ref, "?")
	p := strings.TrimSuffix(strings.Trim(ref, "/"), ".md")
	if p == "" {
		p = "readme"
	}
	return p, fs.ValidPath(p + ".md")
}

func pageURL(p string) string {
	if p == "readme" {
		return Origin
	}
	return Origin + "/" + p
}
