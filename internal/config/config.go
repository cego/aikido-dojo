// Package config loads the config file and resolves which profile, API
// client and region a call runs as.
package config

import (
	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/cego/aikido-dojo/internal/clierr"
)

const (
	EnvConfig       = "AIKIDO_DOJO_CONFIG"
	EnvProfile      = "AIKIDO_DOJO_PROFILE"
	EnvClientID     = "AIKIDO_DOJO_CLIENT_ID"
	EnvClientSecret = "AIKIDO_DOJO_CLIENT_SECRET" //nolint:gosec // the variable's name, not a credential
	EnvRegion       = "AIKIDO_DOJO_REGION"
)

type File struct {
	DefaultProfile string             `json:"default_profile,omitempty"`
	Profiles       map[string]Profile `json:"profiles,omitempty"`
}

type Profile struct {
	ClientID string `json:"client_id"`
	Region   string `json:"region,omitempty"`
}

// Flags holds the global flags that take part in resolution; "" means unset.
type Flags struct {
	Config  string
	Profile string
}

// Resolved is who a call runs as. Profile is "" for the environment profile,
// the only one that carries its Secret: stored profiles keep theirs in the
// keychain.
type Resolved struct {
	Profile  string
	ClientID string
	Secret   string `json:"-"` // never encoded, so no output built from a Resolved can carry it
	Region   string
	Host     string
}

func Path(flags Flags, getenv func(string) string) (string, error) {
	if p := cmp.Or(flags.Config, getenv(EnvConfig)); p != "" {
		return p, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("find the home directory for the config file: %w", err)
	}
	return filepath.Join(home, ".config", "aikido-dojo", "config.json"), nil
}

// Load reads the config file. A missing file is an empty config, so a CI run
// that only sets environment variables needs none.
func Load(path string) (File, error) {
	data, err := os.ReadFile(path) //nolint:gosec // the path is the user's own config file
	if errors.Is(err, fs.ErrNotExist) {
		return File{}, nil
	}
	if err != nil {
		return File{}, fmt.Errorf("read config %s: %w", path, err)
	}
	// Checked before the strict decode, so a pasted secret gets this message
	// rather than "unknown field".
	if key := secretKey(data); key != "" {
		return File{}, usage("secret_in_config", fmt.Sprintf("config %s contains %q", path, key),
			"remove it: secrets belong in the OS keychain (aikido-dojo auth login) or in "+EnvClientSecret)
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var f File
	err = dec.Decode(&f)
	if err == nil {
		// A second value would otherwise be ignored unchecked, secrets included.
		if _, next := dec.Token(); !errors.Is(next, io.EOF) {
			err = errors.New("unexpected data after the JSON object")
		}
	}
	if err != nil {
		return File{}, &clierr.Error{Code: "bad_config", Message: "parse config " + path, Err: err,
			Hint: "fix the file; it holds default_profile and profiles with client_id and region", Exit: clierr.ExitUsage}
	}
	return f, nil
}

// secretKey returns the first secret-shaped key at the top level or inside a
// profile, in sorted order, or "" if there is none. Profile names are the
// user's to choose, so they aren't checked, and anything nested deeper is
// rejected by the strict decode. Invalid JSON returns "": Load reports it.
func secretKey(data []byte) string {
	var top map[string]json.RawMessage
	if json.Unmarshal(data, &top) != nil {
		return ""
	}
	if k := firstSecret(top); k != "" {
		return k
	}
	var profiles map[string]map[string]json.RawMessage
	if json.Unmarshal(top["profiles"], &profiles) != nil {
		return ""
	}
	for _, name := range slices.Sorted(maps.Keys(profiles)) {
		if k := firstSecret(profiles[name]); k != "" {
			return k
		}
	}
	return ""
}

func firstSecret(m map[string]json.RawMessage) string {
	for _, k := range slices.Sorted(maps.Keys(m)) {
		if looksSecret(k) {
			return k
		}
	}
	return ""
}

func looksSecret(key string) bool {
	k := strings.ToLower(key)
	return strings.Contains(k, "secret") || strings.Contains(k, "password") || strings.HasSuffix(k, "token")
}

// Resolve picks the profile and region. A stored profile's client ID comes only
// from its entry: taking one from the environment would send the profile's
// secret to another client.
func Resolve(f File, flags Flags, getenv func(string) string) (Resolved, error) {
	r, err := pick(f, flags, getenv)
	if err != nil {
		return Resolved{}, err
	}
	r.Region = cmp.Or(getenv(EnvRegion), r.Region, DefaultRegion)
	host, err := RegionHost(r.Region)
	if err != nil {
		return Resolved{}, err
	}
	r.Host = host
	return r, nil
}

// Save writes f to path through a temp file and a rename, so a failed write
// keeps the old file. A symlink to an existing config, as a dotfiles checkout
// keeps one, is written through to its target when only its owner can write
// its directory; elsewhere, as in /tmp, a planted link could aim the write at
// any of the user's files, so the link itself is replaced. A new file is 0600
// and an existing one keeps its mode.
func Save(path string, f File) error {
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return fmt.Errorf("encode the config: %w", err)
	}
	if dir, err := os.Stat(filepath.Dir(path)); err == nil && dir.Mode().Perm()&0o022 == 0 {
		if target, err := filepath.EvalSymlinks(path); err == nil {
			path = target
		}
	}
	mode := os.FileMode(0o600)
	if info, err := os.Stat(path); err == nil {
		mode = info.Mode().Perm()
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("create %s: %w", dir, err)
	}
	tmp, err := os.CreateTemp(dir, ".config-*.json")
	if err != nil {
		return fmt.Errorf("write config %s: %w", path, err)
	}
	defer os.Remove(tmp.Name()) // fails harmlessly once the rename has moved it
	if _, err := tmp.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("write config %s: %w", path, errors.Join(err, tmp.Close()))
	}
	if err := tmp.Chmod(mode); err != nil {
		return fmt.Errorf("write config %s: %w", path, errors.Join(err, tmp.Close()))
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("write config %s: %w", path, err)
	}
	if err := os.Rename(tmp.Name(), path); err != nil {
		return fmt.Errorf("write config %s: %w", path, err)
	}
	return nil
}

func pick(f File, flags Flags, getenv func(string) string) (Resolved, error) {
	if name := cmp.Or(flags.Profile, getenv(EnvProfile)); name != "" {
		return stored(f, name)
	}
	id, secret := getenv(EnvClientID), getenv(EnvClientSecret)
	switch {
	case id != "" && secret != "":
		return Resolved{ClientID: id, Secret: secret}, nil
	case id != "" || secret != "":
		// Half a pair is almost always a CI mistake; falling back to a stored profile would hide it.
		missing := EnvClientSecret
		if id == "" {
			missing = EnvClientID
		}
		return Resolved{}, usage("incomplete_credentials", missing+" is not set",
			"set both "+EnvClientID+" and "+EnvClientSecret+", or neither")
	case f.DefaultProfile != "":
		return stored(f, f.DefaultProfile)
	}
	return Resolved{}, &clierr.Error{Code: "no_credentials", Message: "no profile or credentials configured",
		Hint: "run aikido-dojo auth login, or set " + EnvClientID + " and " + EnvClientSecret, Exit: clierr.ExitAuth}
}

func stored(f File, name string) (Resolved, error) {
	p, ok := f.Profiles[name]
	if !ok {
		return Resolved{}, usage("unknown_profile", fmt.Sprintf("no profile %q in the config file", name),
			"create it with aikido-dojo auth login --profile "+name)
	}
	if p.ClientID == "" {
		return Resolved{}, usage("bad_config", fmt.Sprintf("profile %q has no client_id", name),
			"set its client_id, or recreate it with aikido-dojo auth login --profile "+name)
	}
	return Resolved{Profile: name, ClientID: p.ClientID, Region: p.Region}, nil
}

func usage(code, msg, hint string) *clierr.Error {
	return &clierr.Error{Code: code, Message: msg, Hint: hint, Exit: clierr.ExitUsage}
}
