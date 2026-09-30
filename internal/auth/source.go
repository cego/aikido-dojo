// Package auth gets access tokens with the client_credentials grant and
// caches them in the OS keychain, so each short-lived CLI process doesn't
// spend a call on the token endpoint.
package auth

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/zalando/go-keyring"

	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/config"
)

const keychainService = "aikido-dojo"

// A token this close to expiry is treated as expired, so it can't lapse mid-request.
const expiryMargin = time.Minute

type Source struct {
	http     *http.Client
	tokenURL string
	clientID string
	secret   string
	profile  string
	now      func() time.Time
	token    cachedToken
}

// cachedToken records which client and token endpoint issued it: a profile
// can be repointed at another client or region while its cache entry lives on.
type cachedToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	ClientID  string    `json:"client_id"`
	TokenURL  string    `json:"token_url"`
}

// NewSource builds the token source for r. A stored profile's secret comes
// from the keychain; the environment profile carries its own.
func NewSource(c *http.Client, r config.Resolved) (*Source, error) {
	s := &Source{
		http:     c,
		tokenURL: "https://" + r.Host + "/api/oauth/token",
		clientID: r.ClientID,
		secret:   r.Secret,
		profile:  r.Profile,
		now:      time.Now,
	}
	if r.Profile == "" {
		return s, nil
	}
	secret, err := keyring.Get(keychainService, account(r.Profile, "client_secret"))
	if errors.Is(err, keyring.ErrNotFound) {
		return nil, &clierr.Error{Code: "no_secret", Message: fmt.Sprintf("no client secret stored for profile %q", r.Profile),
			Hint: "run aikido-dojo auth login --profile " + r.Profile, Exit: clierr.ExitAuth}
	}
	if err != nil {
		return nil, fmt.Errorf("read the client secret from the keychain: %w", err)
	}
	s.secret = secret
	return s, nil
}

func account(profile, item string) string { return profile + "/" + item }

// Token returns a valid access token from memory, then the keychain cache,
// then the token endpoint.
func (s *Source) Token(ctx context.Context) (string, error) {
	if s.valid(s.token) {
		return s.token.Token, nil
	}
	if s.profile != "" {
		cached, err := s.readCache()
		if err != nil {
			return "", err
		}
		if s.valid(cached) {
			s.token = cached
			return cached.Token, nil
		}
	}
	tok, err := s.fetch(ctx)
	if err != nil {
		return "", err
	}
	s.token = tok
	if s.profile != "" {
		if err := s.writeCache(tok); err != nil {
			return "", err
		}
	}
	return tok.Token, nil
}

// Invalidate drops the cached token after the API rejected it.
func (s *Source) Invalidate() error {
	s.token = cachedToken{}
	if s.profile == "" {
		return nil
	}
	err := keyring.Delete(keychainService, account(s.profile, "access_token"))
	if err != nil && !errors.Is(err, keyring.ErrNotFound) {
		return fmt.Errorf("drop the cached access token: %w", err)
	}
	return nil
}

func (s *Source) valid(t cachedToken) bool {
	return t.Token != "" && t.ClientID == s.clientID && t.TokenURL == s.tokenURL &&
		s.now().Add(expiryMargin).Before(t.ExpiresAt)
}

func (s *Source) readCache() (cachedToken, error) {
	raw, err := keyring.Get(keychainService, account(s.profile, "access_token"))
	if errors.Is(err, keyring.ErrNotFound) {
		return cachedToken{}, nil
	}
	if err != nil {
		return cachedToken{}, fmt.Errorf("read the cached access token: %w", err)
	}
	var t cachedToken
	if err := json.Unmarshal([]byte(raw), &t); err != nil {
		// A damaged entry only costs one token call, so it is replaced rather than reported.
		t = cachedToken{}
	}
	return t, nil
}

func (s *Source) writeCache(t cachedToken) error {
	raw, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("encode the access token cache: %w", err)
	}
	if err := keyring.Set(keychainService, account(s.profile, "access_token"), string(raw)); err != nil {
		return fmt.Errorf("cache the access token in the keychain: %w", err)
	}
	return nil
}

type tokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (s *Source) fetch(ctx context.Context) (cachedToken, error) {
	form := url.Values{"grant_type": {"client_credentials"}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return cachedToken{}, fmt.Errorf("build the token request: %w", err)
	}
	req.SetBasicAuth(s.clientID, s.secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := s.http.Do(req)
	if err != nil {
		return cachedToken{}, fmt.Errorf("request an access token: %w", err)
	}
	defer resp.Body.Close()
	var body tokenResponse
	decodeErr := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&body)
	if resp.StatusCode != http.StatusOK {
		return cachedToken{}, s.failure(resp.StatusCode, body)
	}
	if decodeErr != nil || body.AccessToken == "" {
		return cachedToken{}, &clierr.Error{Code: "auth_failed", Message: "the token endpoint returned no access token",
			HTTPStatus: resp.StatusCode, Err: decodeErr, Hint: "check the API client in Aikido's workspace settings", Exit: clierr.ExitAuth}
	}
	return cachedToken{
		Token:     body.AccessToken,
		ExpiresAt: s.now().Add(time.Duration(body.ExpiresIn) * time.Second),
		ClientID:  s.clientID,
		TokenURL:  s.tokenURL,
	}, nil
}

// Aikido reports every client and grant problem as 401 invalid_client and
// tells them apart only in error_description (auth spike, 2026-09-29).
func (s *Source) failure(status int, body tokenResponse) error {
	e := &clierr.Error{
		Code:       "auth_failed",
		Message:    "token request failed: " + cmp.Or(body.ErrorDescription, body.Error, http.StatusText(status)),
		HTTPStatus: status,
		Exit:       clierr.ExitAuth,
	}
	switch {
	case status == http.StatusTooManyRequests:
		e.Code, e.Exit, e.Hint = "rate_limited", clierr.ExitRateLimited, "Aikido allows 20 calls per minute per workspace; wait a minute and retry"
	case status >= 500:
		e.Code, e.Exit, e.Hint = "server_error", clierr.ExitUnexpected, "Aikido's token endpoint failed; retry later"
	case strings.Contains(body.ErrorDescription, "public"):
		e.Hint = "this API client is a Public app, which only allows browser login; create a non-public API client in Aikido's workspace settings"
	case s.profile == "":
		e.Hint = "check " + config.EnvClientID + " and " + config.EnvClientSecret
	default:
		e.Hint = "check the client ID, or store the secret again with aikido-dojo auth login --profile " + s.profile
	}
	return e
}
