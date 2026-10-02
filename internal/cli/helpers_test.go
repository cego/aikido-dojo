package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/cego/aikido-dojo/internal/config"
)

// testEnv can't reach the user's config, keychain or network: the config
// path points into a temp dir, and any request fails until a test routes it.
// Tests set more variables in the returned map.
func testEnv(t *testing.T) (Env, map[string]string) {
	t.Helper()
	vars := map[string]string{config.EnvConfig: filepath.Join(t.TempDir(), "config.json")}
	cache := t.TempDir()
	return Env{
		Stdin:     strings.NewReader(""),
		Getenv:    func(k string) string { return vars[k] },
		CacheDir:  func() (string, error) { return cache, nil },
		Transport: noNetwork{},
	}, vars
}

type noNetwork struct{}

func (noNetwork) RoundTrip(r *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("the test tried to reach %s", r.URL.Redacted())
}

func run(t *testing.T, env Env, args ...string) (stdout, stderr string, code int) {
	t.Helper()
	var out, errOut bytes.Buffer
	env.Stdout, env.Stderr = &out, &errOut
	code = Run(t.Context(), args, env)
	return out.String(), errOut.String(), code
}

type reported struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Hint       string `json:"hint"`
	HTTPStatus int    `json:"http_status"`
}

// errorOf decodes the error line on stderr, skipping any warnings.
func errorOf(t *testing.T, stderr string) reported {
	t.Helper()
	for _, line := range strings.Split(strings.TrimSpace(stderr), "\n") {
		var v struct {
			Error *reported `json:"error"`
		}
		if json.Unmarshal([]byte(line), &v) == nil && v.Error != nil {
			return *v.Error
		}
	}
	t.Fatalf("no error on stderr: %q", stderr)
	return reported{}
}

// rewrite sends every request to a test server and records the host the CLI
// meant, so tests can check which region it picked.
type rewrite struct {
	to    string
	next  http.RoundTripper
	mu    sync.Mutex
	hosts []string
}

func (r *rewrite) RoundTrip(req *http.Request) (*http.Response, error) {
	r.mu.Lock()
	r.hosts = append(r.hosts, req.URL.Host)
	r.mu.Unlock()
	out := req.Clone(req.Context())
	out.URL.Host = r.to
	return r.next.RoundTrip(out) //nolint:wrapcheck // http.Client wraps transport errors itself
}

func (r *rewrite) seenHosts() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return slices.Clone(r.hosts)
}
