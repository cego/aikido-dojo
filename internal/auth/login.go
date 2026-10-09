package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/cego/aikido-dojo/internal/config"
)

// Login checks r's client ID and secret with one token request. Only once
// Aikido accepts them does it run record, which saves the profile, and then
// store the secret for r.Profile and cache the token, replacing any issued
// under the profile's previous secret. If record fails, nothing is stored, so
// the profile's old secret still matches the client ID the config names.
// moved says the profile's storage changed: once the new copy is stored, the
// old one is deleted, and a failure there is ignored, since a machine moved
// to the file often has no keychain at all.
func Login(ctx context.Context, c *http.Client, r config.Resolved, moved bool, record func() error) error {
	st := storeFor(r)
	s := &Source{http: c, tokenURL: tokenURL(r.Host), clientID: r.ClientID, secret: r.Secret, profile: r.Profile, store: st, now: time.Now, login: true}
	tok, err := s.fetch(ctx)
	if err != nil {
		return err
	}
	if err := record(); err != nil {
		return err
	}
	if err := st.set(clientSecret, r.Secret); err != nil {
		return st.failed("store the client secret", err)
	}
	if err := s.writeCache(tok); err != nil {
		return err
	}
	if moved {
		_, _ = forget(otherStore(r))
	}
	return nil
}

// Forget deletes the profile's stored secret and cached token from the
// keychain and the credentials file. removed reports whether either held a
// secret. Errors from the store the profile doesn't use are ignored.
func Forget(r config.Resolved) (removed bool, err error) {
	removed, err = forget(storeFor(r))
	if err != nil {
		return false, err
	}
	other, _ := forget(otherStore(r))
	return removed || other, nil
}

func forget(st store) (bool, error) {
	if _, err := st.remove(accessToken); err != nil {
		return false, st.failed("delete the cached access token", err)
	}
	removed, err := st.remove(clientSecret)
	if err != nil {
		return false, st.failed("delete the client secret", err)
	}
	return removed, nil
}
