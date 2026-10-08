package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"unicode/utf8"
)

var tenantPatterns = []*regexp.Regexp{
	regexp.MustCompile(`eyJ[A-Za-z0-9_-]{8,}\.eyJ[A-Za-z0-9_-]{8,}`), // a JWT, such as an access token
	regexp.MustCompile(`AIK_CLIENT_[A-Za-z0-9]{16,}`),                // a real API client ID
	regexp.MustCompile(`\bcego\.dk\b`),                               // the company's own hosts
	regexp.MustCompile(`-----BEGIN [A-Z ]*PRIVATE KEY-----`),
}

var email = regexp.MustCompile(`[A-Za-z0-9._%+-]+@((?:[A-Za-z0-9-]+\.)+[A-Za-z]{2,})`)

// tenantData lists what in s looks like a credential or a tenant's
// identifier. An address is allowed only on a domain no tenant has: the
// RFC 2606 examples, and github.com in git@github.com remotes.
func tenantData(s string) []string {
	var hits []string
	for _, p := range tenantPatterns {
		hits = append(hits, p.FindAllString(s, -1)...)
	}
	for _, m := range email.FindAllStringSubmatch(s, -1) {
		domain := strings.ToLower(m[1])
		if domain == "github.com" || strings.HasSuffix("."+domain, ".example.com") || strings.HasSuffix("."+domain, ".example.org") {
			continue
		}
		hits = append(hits, m[0])
	}
	return hits
}

func TestTenantDataPatterns(t *testing.T) {
	hits := []string{
		"token eyJhbGciOiJSUzI1NiJ9.eyJzY29wZSI6Imlzc3VlczpyZWFkIn0.sig",
		"AIK_CLIENT_0123456789abcdef0123456789abcdef",
		"mail someone@company.io about it",
		"https://cego.dk/group/repo.git",
		"-----BEGIN RSA PRIVATE KEY-----",
	}
	for _, s := range hits {
		if len(tenantData(s)) == 0 {
			t.Errorf("tenantData(%q) found nothing", s)
		}
	}
	misses := []string{
		"AIK_CLIENT_cego and AIK_CLIENT_…",
		"git@github.com:cego/aikido-dojo.git and git@gitlab.example.com:g/r.git",
		"user@example.com, ops@example.org",
		"41898282+github-actions[bot]@users.noreply.github.com",
		"github.com/cego/aikido-dojo",
	}
	for _, s := range misses {
		if got := tenantData(s); len(got) != 0 {
			t.Errorf("tenantData(%q) = %q, want nothing", s, got)
		}
	}
}

// Tests use hand-written fakes, not recordings, so nothing is scrubbed on
// the way in; this fails if tenant data reaches a file anyway. Skipped: the
// vendored spec and the catalog generated from it, which are Aikido's own
// document, and this file, which holds the samples above.
func TestRepositoryHoldsNoTenantData(t *testing.T) {
	skip := map[string]bool{
		".git": true, "dist": true, "spec": true,
		"internal/catalog/schemas.json": true, "internal/catalog/search.json": true, "internal/catalog/zz_ops.go": true,
		"main_test.go": true,
	}
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if skip[filepath.ToSlash(path)] {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path) //nolint:gosec // the repository's own files
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		if !utf8.Valid(data) {
			return nil
		}
		for _, hit := range tenantData(string(data)) {
			t.Errorf("%s holds tenant data: %s", path, hit)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
