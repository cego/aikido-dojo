package auth

import (
	"errors"

	"github.com/zalando/go-keyring"

	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/config"
)

const (
	keychainService = "aikido-dojo"
	clientSecret    = "client_secret"
	accessToken     = "access_token"
)

// store keeps one stored profile's client secret and cached access token.
type store interface {
	get(item string) (value string, ok bool, err error)
	set(item, value string) error
	remove(item string) (removed bool, err error)
	// failed is err from doing action, with what to do about it.
	failed(action string, err error) error
}

func storeFor(r config.Resolved) store {
	if r.Storage == config.StorageFile {
		return credentialsFile{path: r.Credentials, profile: r.Profile}
	}
	return keychain{profile: r.Profile}
}

func otherStore(r config.Resolved) store {
	if r.Storage == config.StorageFile {
		return keychain{profile: r.Profile}
	}
	return credentialsFile{path: r.Credentials, profile: r.Profile}
}

type keychain struct{ profile string }

func account(profile, item string) string { return profile + "/" + item }

func (k keychain) get(item string) (string, bool, error) {
	v, err := keyring.Get(keychainService, account(k.profile, item))
	if errors.Is(err, keyring.ErrNotFound) {
		return "", false, nil
	}
	return v, err == nil, err //nolint:wrapcheck // failed adds the context
}

func (k keychain) set(item, value string) error {
	return keyring.Set(keychainService, account(k.profile, item), value) //nolint:wrapcheck // failed adds the context
}

func (k keychain) remove(item string) (bool, error) {
	err := keyring.Delete(keychainService, account(k.profile, item))
	if errors.Is(err, keyring.ErrNotFound) {
		return false, nil
	}
	return err == nil, err //nolint:wrapcheck // failed adds the context
}

func (k keychain) failed(action string, err error) error {
	return &clierr.Error{Code: "keychain_unavailable", Message: action, Err: err, Exit: clierr.ExitUnexpected,
		Hint: "on a machine with no keychain, such as a headless Linux server, run aikido-dojo auth login --profile " + k.profile +
			" --insecure-storage to keep the secret in a file only you can read, or set " + config.EnvClientID + " and " +
			config.EnvClientSecret + " instead"}
}
