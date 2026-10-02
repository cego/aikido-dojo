package cli

import (
	"errors"
	"net/http"
	"time"

	"github.com/cego/aikido-dojo/internal/api"
	"github.com/cego/aikido-dojo/internal/clierr"
)

// headerTimeout bounds only the wait for a response's headers. Exports stream
// bodies that take minutes (a 305 MB issue export took 117 s), so there is no
// overall timeout; Ctrl-C cancels a stuck call.
const headerTimeout = 5 * time.Minute

func newHTTPClient(env Env, host string, debug bool) *http.Client {
	rt := env.Transport
	if rt == nil {
		t := http.DefaultTransport.(*http.Transport).Clone()
		t.ResponseHeaderTimeout = headerTimeout
		rt = t
	}
	if debug {
		rt = api.Debug(rt, env.Stderr)
	}
	return &http.Client{Transport: rt, CheckRedirect: sameHost(host)}
}

// sameHost refuses a redirect off the API host: aikido-dojo talks to no other
// host, and the bearer token must not follow a redirect elsewhere.
func sameHost(host string) func(*http.Request, []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		if req.URL.Scheme != "https" || req.URL.Host != host {
			return &clierr.Error{Code: "redirect_refused", Message: "refused a redirect to " + req.URL.Redacted(),
				Hint: "aikido-dojo only talks to " + host + "; retry, and report it if Aikido keeps redirecting", Exit: clierr.ExitUnexpected}
		}
		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}
		return nil
	}
}
