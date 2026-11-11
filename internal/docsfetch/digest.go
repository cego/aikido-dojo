package docsfetch

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/cego/aikido-dojo/internal/specfetch"
)

// The changelog is the last page and changes with every post; the copy leaves it out.
var changelog = []byte("\n# Changelog\n")

// FullDigest returns the SHA-256 of llms-full.txt without its changelog, so one digest tells whether any page changed.
func FullDigest(ctx context.Context, c *http.Client, u string) (string, error) {
	body, err := specfetch.Get(ctx, c, u)
	if err != nil {
		return "", err
	}
	// A challenge or error page isn't markdown; hashing it would trigger a sync of every page.
	if !bytes.HasPrefix(body, []byte("# ")) {
		return "", fmt.Errorf("%s: not the help centre's llms-full.txt", u)
	}
	if j := bytes.LastIndex(body, changelog); j >= 0 {
		body = body[:j]
	}
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:]), nil
}
