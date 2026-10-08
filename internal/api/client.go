// Package api sends requests to Aikido's public REST API. It adds the access
// token, paces calls under the rate limit, retries what is safe to retry and
// turns failures into the CLI's error model.
package api

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/cego/aikido-dojo/internal/auth"
	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/ratelimit"
)

type Client struct {
	http    *http.Client
	host    string
	tokens  *auth.Source
	limiter *ratelimit.Limiter
	now     func() time.Time
	sleep   func(context.Context, time.Duration) error
	backoff func(attempt int) time.Duration
}

// New returns a client for one region's host. Leaving Accept-Encoding to the
// transport makes it ask for gzip and decompress while streaming, which cut
// a 305 MB issue export to 9 MB on the wire.
func New(c *http.Client, host string, tokens *auth.Source, limiter *ratelimit.Limiter) *Client {
	return &Client{
		http:    c,
		host:    host,
		tokens:  tokens,
		limiter: limiter,
		now:     time.Now,
		sleep:   sleep,
		backoff: backoff,
	}
}

// Request is one API call. Scope is the scope the operation needs, named in
// the hint when Aikido answers 403.
type Request struct {
	Method string
	Path   string
	Query  url.Values
	Body   []byte
	Scope  string
}

// URL is where req goes on host; a dry run prints the same URL a call sends to.
func URL(host string, req Request) string {
	u := "https://" + host + "/api/public/v1" + req.Path
	if len(req.Query) > 0 {
		u += "?" + req.Query.Encode()
	}
	return u
}

// Do sends req and returns the response when it is 2xx; the caller closes its
// body. Every attempt waits for the rate limiter first.
func (c *Client) Do(ctx context.Context, req Request) (*http.Response, error) {
	var refreshed bool
	var rateRetries, serverRetries int
	for {
		if err := c.limiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("%s %s: %w", req.Method, req.Path, err)
		}
		resp, err := c.send(ctx, req)
		if err != nil {
			return nil, err
		}
		status := resp.StatusCode
		switch {
		case status >= 200 && status < 300:
			return resp, nil
		case status == http.StatusUnauthorized && !refreshed:
			// The cached token may have been revoked or its secret rotated; one fresh token decides.
			discard(resp)
			refreshed = true
			if err := c.tokens.Invalidate(); err != nil {
				return nil, err
			}
			continue
		case status == http.StatusTooManyRequests && rateRetries < maxRateRetries:
			wait := retryAfter(resp.Header.Get("Retry-After"), c.now())
			discard(resp)
			rateRetries++
			if err := c.sleep(ctx, wait); err != nil {
				return nil, err
			}
			continue
		case status >= 500 && req.Method == http.MethodGet && serverRetries < maxServerRetries:
			// Only GETs are retried: a write may have taken effect before the server failed.
			discard(resp)
			wait := c.backoff(serverRetries)
			serverRetries++
			if err := c.sleep(ctx, wait); err != nil {
				return nil, err
			}
			continue
		}
		return nil, failure(req, resp)
	}
}

func (c *Client) send(ctx context.Context, req Request) (*http.Response, error) {
	tok, err := c.tokens.Token(ctx)
	if err != nil {
		return nil, err
	}
	u := URL(c.host, req)
	var body io.Reader
	if req.Body != nil {
		body = bytes.NewReader(req.Body)
	}
	hr, err := http.NewRequestWithContext(ctx, req.Method, u, body)
	if err != nil {
		return nil, fmt.Errorf("build %s %s: %w", req.Method, req.Path, err)
	}
	hr.Header.Set("Authorization", "Bearer "+tok)
	if req.Body != nil {
		hr.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.http.Do(hr)
	if err != nil {
		return nil, fmt.Errorf("%s %s: %w", req.Method, req.Path, err)
	}
	return resp, nil
}

// discard reads a little of a body we don't use, so the connection can be reused.
func discard(resp *http.Response) {
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 64<<10))
	_ = resp.Body.Close()
}

func failure(req Request, resp *http.Response) error {
	defer resp.Body.Close()
	status := resp.StatusCode
	e := &clierr.Error{Message: fmt.Sprintf("%s %s: %s", req.Method, req.Path, apiMessage(resp)), HTTPStatus: status}
	switch {
	case status == http.StatusUnauthorized:
		e.Code, e.Exit = "auth_failed", clierr.ExitAuth
		e.Hint = "Aikido rejected a freshly issued token; check the API client in Aikido's workspace settings"
	case status == http.StatusForbidden && req.Scope != "":
		e.Code, e.Exit = "missing_scope", clierr.ExitForbidden
		e.Hint = "grant the " + req.Scope + " scope on the API client in Aikido's workspace settings"
	case status == http.StatusForbidden:
		e.Code, e.Exit = "forbidden", clierr.ExitForbidden
		e.Hint = "check the API client's scopes in Aikido's workspace settings"
	case status == http.StatusNotFound:
		e.Code, e.Exit = "not_found", clierr.ExitNotFound
		e.Hint = "check the ID; list commands show the valid ones"
	case status == http.StatusTooManyRequests:
		e.Code, e.Exit = "rate_limited", clierr.ExitRateLimited
		e.Hint = "Aikido allows 20 calls per minute per workspace; wait a minute and retry"
	case status >= 400 && status < 500:
		e.Code, e.Exit = "bad_request", clierr.ExitUsage
		e.Hint = "check the flags and body against: aikido-dojo schema <command>"
	default:
		e.Code, e.Exit = "server_error", clierr.ExitUnexpected
		e.Hint = "Aikido failed to handle the request; retry later"
	}
	return e
}

// apiMessage extracts Aikido's error text ({"error": …}, {"message": …}, or
// {"reason_phrase": …} as seen live on 2026-09-30), else a short plain-text
// body, else the status text.
func apiMessage(resp *http.Response) string {
	fallback := http.StatusText(resp.StatusCode)
	if strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html") {
		return fallback
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
	if err != nil {
		return fallback
	}
	var body struct {
		Error        string `json:"error"`
		Message      string `json:"message"`
		ReasonPhrase string `json:"reason_phrase"`
	}
	if json.Unmarshal(data, &body) == nil {
		if m := cmp.Or(body.Error, body.Message, body.ReasonPhrase); m != "" {
			return m
		}
	}
	if text := strings.TrimSpace(string(data)); text != "" && utf8.ValidString(text) && len(text) <= 200 {
		return text
	}
	return fallback
}
