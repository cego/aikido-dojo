package docsfetch

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/cego/aikido-dojo/internal/specfetch"
)

// The site has 8 parts of 100 pages; the cap only stops a link loop.
const maxParts = 50

var nextPartRE = regexp.MustCompile(`(?m)^> Next part: (\S+)$`)

// FullDigest returns the SHA-256 of llms-full.txt and every part it links to, so one digest tells whether any page changed.
func FullDigest(ctx context.Context, c *http.Client, first string, origin *url.URL) (string, error) {
	h := sha256.New()
	next := first
	for range maxParts {
		body, err := specfetch.Get(ctx, c, next)
		if err != nil {
			return "", err
		}
		h.Write(body)
		m := nextPartRE.FindSubmatch(body[:min(len(body), 1024)])
		if m == nil {
			return hex.EncodeToString(h.Sum(nil)), nil
		}
		u, err := url.Parse(string(m[1]))
		if err != nil {
			return "", fmt.Errorf("parse next part %q: %w", m[1], err)
		}
		if u.Scheme != origin.Scheme || u.Host != origin.Host {
			return "", fmt.Errorf("next part %s is not on %s://%s", u, origin.Scheme, origin.Host)
		}
		next = u.String()
	}
	return "", fmt.Errorf("llms-full.txt has more than %d parts", maxParts)
}
