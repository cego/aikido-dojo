// Package specfetch rebuilds the Aikido OpenAPI spec from the public API
// reference pages, which need no credentials, so CI can track spec drift
// without holding an Aikido secret.
package specfetch

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// These llms.txt sections are prose or the OAuth server, which lives outside
// /api/public/v1. Every other "API Reference" section documents operations.
var skippedSections = map[string]bool{
	"Documentation": true,
	"Authorization": true,
}

var (
	entryRE = regexp.MustCompile(`^\s*- \[`)
	// The lazy title match lets titles contain brackets, as in "[Beta] List things".
	linkRE = regexp.MustCompile(`^\s*- \[.*?\]\(([^)\s]+\.md)\)`)
)

// PageURLs returns the reference pages llms.txt lists, in index order. Links
// must stay on the index's origin so a tampered index can't send the fetcher
// elsewhere.
func PageURLs(index string, origin *url.URL) ([]string, error) {
	var urls []string
	section := ""
	sc := bufio.NewScanner(strings.NewReader(index))
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "## ") {
			name, ok := strings.CutPrefix(line, "## API Reference: ")
			name = strings.TrimSpace(name)
			section = ""
			if ok && !skippedSections[name] {
				section = name
			}
			continue
		}
		if section == "" {
			continue
		}
		if !entryRE.MatchString(line) {
			continue
		}
		// Skipping an entry would silently drop an operation from the spec.
		m := linkRE.FindStringSubmatch(line)
		if m == nil {
			return nil, fmt.Errorf("unrecognised index entry %q", line)
		}
		u, err := url.Parse(m[1])
		if err != nil {
			return nil, fmt.Errorf("parse link %q: %w", m[1], err)
		}
		if u.Scheme != origin.Scheme || u.Host != origin.Host {
			return nil, fmt.Errorf("link %s is not on %s://%s", u, origin.Scheme, origin.Host)
		}
		urls = append(urls, u.String())
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("scan index: %w", err)
	}
	if len(urls) == 0 {
		return nil, errors.New("index lists no API reference pages")
	}
	return urls, nil
}
