package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/zalando/go-keyring"

	"github.com/cego/aikido-dojo/internal/config"
)

// Login checks r's client ID and secret with one token request. Only once
// Aikido accepts them does it run record, which saves the profile, and then
// store the secret for r.Profile and cache the token, replacing any issued
// under the profile's previous secret. If record fails, nothing is stored, so
// the profile's old secret still matches the client ID the config names.
func Login(ctx context.Context, c *http.Client, r config.Resolved, record func() error) error {
	s := &Source{http: c, tokenURL: tokenURL(r.Host), clientID: r.ClientID, secret: r.Secret, profile: r.Profile, now: time.Now, login: true}
	tok, err := s.fetch(ctx)
	if err != nil {
		return err
	}
	if err := record(); err != nil {
		return err
	}
	if err := keyring.Set(keychainService, account(r.Profile, "client_secret"), r.Secret); err != nil {
		return fmt.Errorf("store the client secret in the keychain: %w", err)
	}
	return s.writeCache(tok)
}

// Forget deletes the profile's stored secret and cached token. removed
// reports whether there was a secret to delete.
func Forget(profile string) (removed bool, err error) {
	err = keyring.Delete(keychainService, account(profile, "access_token"))
	if err != nil && !errors.Is(err, keyring.ErrNotFound) {
		return false, fmt.Errorf("delete the cached access token: %w", err)
	}
	err = keyring.Delete(keychainService, account(profile, "client_secret"))
	if errors.Is(err, keyring.ErrNotFound) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("delete the client secret: %w", err)
	}
	return true, nil
}
