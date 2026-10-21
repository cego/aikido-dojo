package docsfetch

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/cego/aikido-dojo/internal/specfetch"
)

// The site has 8 parts of 100 pages; the cap only stops a link loop.
const maxParts = 50

var (
	partRE     = regexp.MustCompile(`^> This is part (\d+) of (\d+)\b`)
	nextPartRE = regexp.MustCompile(`(?m)^> Next part: (\S+)$`)
	// The changelog is the last page and changes with every post; the copy leaves it out.
	changelog = []byte("\n# Changelog\n")
)

// FullDigest returns the SHA-256 of llms-full.txt and every part it links to, so one digest tells whether any page changed.
func FullDigest(ctx context.Context, c *http.Client, first string, origin *url.URL) (string, error) {
	h := sha256.New()
	next := first
	for i := 1; i <= maxParts; i++ {
		body, err := specfetch.Get(ctx, c, next)
		if err != nil {
			return "", err
		}
		head := body[:min(len(body), 1024)]
		m := partRE.FindSubmatch(head)
		if m == nil || string(m[1]) != strconv.Itoa(i) {
			return "", fmt.Errorf("%s: not part %d of GitBook's llms-full.txt", next, i)
		}
		link := nextPartRE.FindSubmatch(head)
		if link == nil {
			if total := string(m[2]); total != strconv.Itoa(i) {
				return "", fmt.Errorf("llms-full.txt ends at part %d of %s", i, total)
			}
			if j := bytes.LastIndex(body, changelog); j >= 0 {
				body = body[:j]
			}
			h.Write(body)
			return hex.EncodeToString(h.Sum(nil)), nil
		}
		h.Write(body)
		u, err := url.Parse(string(link[1]))
		if err != nil {
			return "", fmt.Errorf("parse next part %q: %w", link[1], err)
		}
		if u.Scheme != origin.Scheme || u.Host != origin.Host {
			return "", fmt.Errorf("next part %s is not on %s://%s", u, origin.Scheme, origin.Host)
		}
		next = u.String()
	}
	return "", fmt.Errorf("llms-full.txt has more than %d parts", maxParts)
}
