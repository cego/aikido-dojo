package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"

	"github.com/cego/aikido-dojo/internal/clierr"
)

// credentialsFile keeps a profile's secret and cached token in a file shared
// by every profile with file storage, for machines with no keychain.
type credentialsFile struct {
	path    string
	profile string
}

type credentials struct {
	Profiles map[string]map[string]string `json:"profiles"`
}

func (c credentialsFile) get(item string) (string, bool, error) {
	creds, err := c.read()
	if err != nil {
		return "", false, err
	}
	v, ok := creds.Profiles[c.profile][item]
	return v, ok, nil
}

func (c credentialsFile) set(item, value string) error {
	return c.update(func(items map[string]string) { items[item] = value })
}

func (c credentialsFile) remove(item string) (bool, error) {
	// Nothing to remove needn't create the directory and its lock.
	if _, err := os.Lstat(c.path); errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	var removed bool
	err := c.update(func(items map[string]string) {
		_, removed = items[item]
		delete(items, item)
	})
	return removed, err
}

func (credentialsFile) failed(action string, err error) error {
	return fmt.Errorf("%s: %w", action, err)
}

// read decodes the file once check passes; a missing file holds nothing.
func (c credentialsFile) read() (credentials, error) {
	if err := c.check(os.Getuid()); err != nil {
		return credentials{}, err
	}
	data, err := os.ReadFile(c.path)
	if errors.Is(err, fs.ErrNotExist) {
		return credentials{}, nil
	}
	if err != nil {
		return credentials{}, fmt.Errorf("read the credentials file: %w", err)
	}
	var creds credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return credentials{}, &clierr.Error{Code: "bad_credentials_file", Message: "parse " + c.path, Err: err, Exit: clierr.ExitUsage,
			Hint: "fix the file, which holds profiles with client_secret and access_token, or remove it and run " + c.login()}
	}
	return creds, nil
}

// update applies change to the profile's items and replaces the file through
// a temp file and a rename, so a failed write keeps the old one. It holds a
// lock on a file beside it, which the rename doesn't replace, so a process
// changing another profile at the same time can't lose this change.
func (c credentialsFile) update(change func(items map[string]string)) error {
	dir := filepath.Dir(c.path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("create %s: %w", dir, err)
	}
	if err := c.check(os.Getuid()); err != nil {
		return err
	}
	lock, err := os.OpenFile(filepath.Join(dir, ".credentials.lock"), os.O_RDWR|os.O_CREATE|syscall.O_NOFOLLOW, 0o600) //nolint:gosec // the user's own config directory
	if err != nil {
		return fmt.Errorf("open the credentials lock: %w", err)
	}
	defer lock.Close() // closing also releases the lock
	if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX); err != nil {
		return fmt.Errorf("lock the credentials file: %w", err)
	}
	creds, err := c.read()
	if err != nil {
		return err
	}
	items := creds.Profiles[c.profile]
	if items == nil {
		items = map[string]string{}
	}
	change(items)
	if creds.Profiles == nil {
		creds.Profiles = map[string]map[string]string{}
	}
	creds.Profiles[c.profile] = items
	if len(items) == 0 {
		delete(creds.Profiles, c.profile)
	}
	if len(creds.Profiles) == 0 {
		if err := os.Remove(c.path); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("remove %s: %w", c.path, err)
		}
		return nil
	}
	return c.write(creds)
}

func (c credentialsFile) write(creds credentials) error {
	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("encode the credentials file: %w", err)
	}
	tmp, err := os.CreateTemp(filepath.Dir(c.path), ".credentials-*.json")
	if err != nil {
		return fmt.Errorf("write %s: %w", c.path, err)
	}
	defer os.Remove(tmp.Name()) // fails harmlessly once the rename has moved it
	if _, err := tmp.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("write %s: %w", c.path, errors.Join(err, tmp.Close()))
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("write %s: %w", c.path, err)
	}
	if err := os.Rename(tmp.Name(), c.path); err != nil {
		return fmt.Errorf("write %s: %w", c.path, err)
	}
	return nil
}

// check refuses the file, as ssh's StrictModes does a key, unless only uid
// can read it or replace it. A missing file passes: it holds no secret.
func (c credentialsFile) check(uid int) error {
	dir := filepath.Dir(c.path)
	info, err := os.Stat(dir)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("check %s: %w", dir, err)
	}
	if info.Mode().Perm()&0o022 != 0 {
		return insecure(fmt.Sprintf("%s holds the credentials file but others can write to it (mode %04o)", dir, info.Mode().Perm()),
			"run chmod go-w "+dir)
	}
	info, err = os.Lstat(c.path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("check %s: %w", c.path, err)
	}
	st, ok := info.Sys().(*syscall.Stat_t)
	switch {
	case !info.Mode().IsRegular():
		return insecure(c.path+" is not a regular file", "remove "+c.path+", which may be a link to another file, then run "+c.login())
	case !ok || int64(st.Uid) != int64(uid):
		return insecure(c.path+" belongs to another user", "remove "+c.path+", then run "+c.login())
	case info.Mode().Perm()&0o077 != 0:
		return insecure(fmt.Sprintf("others can read or write %s (mode %04o)", c.path, info.Mode().Perm()),
			"run chmod 600 "+c.path+", and rotate the secret in Aikido's workspace settings if someone else may have read it")
	}
	return nil
}

func (c credentialsFile) login() string {
	return "aikido-dojo auth login --profile " + c.profile + " --insecure-storage"
}

func insecure(msg, hint string) error {
	return &clierr.Error{Code: "insecure_credentials_file", Message: msg, Hint: hint, Exit: clierr.ExitUsage}
}
