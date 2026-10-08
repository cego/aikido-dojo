package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"
	"testing"
	"unicode/utf8"
)

var tenantPatterns = []struct {
	name string
	re   *regexp.Regexp
}{
	{"a JWT", regexp.MustCompile(`eyJ[A-Za-z0-9_-]{8,}\.eyJ[A-Za-z0-9_-]{8,}`)},
	{"an API client ID", regexp.MustCompile(`AIK_CLIENT_[A-Za-z0-9]{16,}`)},
	{"an internal host", regexp.MustCompile(`\bcego\.dk\b`)},
	{"a private key", regexp.MustCompile(`-----BEGIN [A-Z ]*PRIVATE KEY-----`)},
}

var (
	email = regexp.MustCompile(`[A-Za-z0-9._%+-]+@((?:[A-Za-z0-9-]+\.)+[A-Za-z]{2,})`)
	// A cego repository other than this one and its tap names an internal project.
	cegoRepo = regexp.MustCompile(`github\.com[/:]cego/([A-Za-z0-9._-]+)`)
)

type finding struct {
	line    int
	pattern string
}

// tenantData lists the lines of s that hold a credential or a tenant's
// identifier, and what each looks like. An address is allowed only on a
// domain no tenant has, the RFC 2606 examples and github.com; a git@host:
// remote isn't an address.
func tenantData(s string) []finding {
	var found []finding
	for i, line := range strings.Split(s, "\n") {
		for _, p := range tenantPatterns {
			if p.re.MatchString(line) {
				found = append(found, finding{i + 1, p.name})
			}
		}
		for _, m := range email.FindAllStringSubmatchIndex(line, -1) {
			domain := strings.ToLower(line[m[2]:m[3]])
			remote := m[1] < len(line) && line[m[1]] == ':'
			if remote || domain == "github.com" || strings.HasSuffix("."+domain, ".example.com") || strings.HasSuffix("."+domain, ".example.org") {
				continue
			}
			found = append(found, finding{i + 1, "an e-mail address"})
		}
		for _, m := range cegoRepo.FindAllStringSubmatch(line, -1) {
			if repo := strings.TrimSuffix(m[1], ".git"); repo != "aikido-dojo" && repo != "homebrew-tap" {
				found = append(found, finding{i + 1, "a cego repository"})
			}
		}
	}
	return found
}

// The samples are built from pieces, so this file passes its own scan.
func TestTenantDataPatterns(t *testing.T) {
	hits := map[string]string{
		"token " + "eyJ" + "hbGciOiJSUzI1NiJ9.eyJ" + "zY29wZSI6Imlzc3VlczpyZWFkIn0.sig": "a JWT",
		"AIK_" + "CLIENT_0123456789abcdef0123456789abcdef":                              "an API client ID",
		"https://cego" + ".dk/group/repo.git":                                           "an internal host",
		"-----BEGIN RSA PRIVATE " + "KEY-----":                                          "a private key",
		"mail someone@" + "company.io about it":                                         "an e-mail address",
		"see https://github.com/" + "cego/internal-tool":                                "a cego repository",
	}
	for s, want := range hits {
		if got := tenantData(s); len(got) != 1 || got[0].pattern != want {
			t.Errorf("tenantData(%q) = %v, want %s", s, got, want)
		}
	}
	misses := []string{
		"AIK_CLIENT_cego and AIK_CLIENT_…",
		"git@github.com:cego/aikido-dojo.git and git@gitlab.example.com:g/r.git",
		"git@gitlab.com:group/repo.git and git@bitbucket.org:team/repo.git",
		"user@example.com, ops@example.org",
		"41898282+github-actions[bot]@users.noreply.github.com",
		"github.com/cego/aikido-dojo and github.com/cego/homebrew-tap",
	}
	for _, s := range misses {
		if got := tenantData(s); len(got) != 0 {
			t.Errorf("tenantData(%q) = %v, want nothing", s, got)
		}
	}
}

// Tests use hand-written fakes, not recordings, so nothing is scrubbed on
// the way in; this fails if tenant data reaches a tracked file anyway. It
// names the file, line and kind, never the text, which CI logs would keep.
// Skipped: the vendored spec and the catalog generated from it, which are
// Aikido's own document.
func TestRepositoryHoldsNoTenantData(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	out, err := exec.CommandContext(t.Context(), "git", "ls-files", "-z").Output()
	if err != nil {
		t.Fatalf("list the tracked files: %v", err)
	}
	generated := []string{"internal/catalog/schemas.json", "internal/catalog/search.json", "internal/catalog/zz_ops.go"}
	for _, path := range strings.Split(strings.TrimSuffix(string(out), "\x00"), "\x00") {
		if strings.HasPrefix(path, "spec/") || slices.Contains(generated, path) {
			continue
		}
		data, err := os.ReadFile(path) //nolint:gosec // a file git tracks here
		if err != nil {
			t.Fatal(fmt.Errorf("read %s: %w", path, err))
		}
		if !utf8.Valid(data) || bytes.IndexByte(data, 0) >= 0 {
			continue
		}
		for _, f := range tenantData(string(data)) {
			t.Errorf("%s:%d holds %s", path, f.line, f.pattern)
		}
	}
}
