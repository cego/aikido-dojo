package api

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Debug wraps next so every request is logged to w as method, URL, status and
// timing. Headers and bodies are never logged, so the Authorization header,
// the client secret and tokens can't reach the log.
func Debug(next http.RoundTripper, w io.Writer) http.RoundTripper {
	return debugTransport{next: next, w: w}
}

type debugTransport struct {
	next http.RoundTripper
	w    io.Writer
}

func (t debugTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	start := time.Now()
	resp, err := t.next.RoundTrip(r)
	took := time.Since(start).Round(time.Millisecond)
	if err != nil {
		fmt.Fprintf(t.w, "debug: %s %s -> %v (%s)\n", r.Method, r.URL.Redacted(), err, took)
		return nil, err //nolint:wrapcheck // http.Client wraps transport errors itself
	}
	fmt.Fprintf(t.w, "debug: %s %s -> %d (%s)\n", r.Method, r.URL.Redacted(), resp.StatusCode, took)
	return resp, nil
}
