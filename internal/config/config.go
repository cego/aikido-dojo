// Package config loads the config file and resolves which profile, API
// client and region a call runs as.
package config

import (
	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
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
	Region  string
}

// Resolved is who a call runs as. Profile is "" for the environment profile,
// the only one that carries its Secret: stored profiles keep theirs in the
// keychain.
type Resolved struct {
	Profile  string
	ClientID string
	Secret   string
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
	if err := dec.Decode(&f); err != nil {
		return File{}, &clierr.Error{Code: "bad_config", Message: "parse config " + path, Err: err,
			Hint: "fix the file; it holds default_profile and profiles with client_id and region", Exit: clierr.ExitUsage}
	}
	return f, nil
}

// secretKey returns the first secret-shaped key, visiting keys in sorted
// order, or "" if there is none. Invalid JSON returns "": the strict decode
// in Load reports it.
func secretKey(data []byte) string {
	var v any
	if json.Unmarshal(data, &v) != nil {
		return ""
	}
	return findSecretKey(v)
}

func findSecretKey(v any) string {
	switch v := v.(type) {
	case map[string]any:
		for _, k := range slices.Sorted(maps.Keys(v)) {
			if looksSecret(k) {
				return k
			}
			if found := findSecretKey(v[k]); found != "" {
				return found
			}
		}
	case []any:
		for _, e := range v {
			if found := findSecretKey(e); found != "" {
				return found
			}
		}
	}
	return ""
}

func looksSecret(key string) bool {
	k := strings.ToLower(key)
	return strings.Contains(k, "secret") || strings.Contains(k, "password") || strings.HasSuffix(k, "token")
}

// Resolve picks the profile and region. A stored profile always uses its own
// client ID, so the secret the keychain holds for it still matches.
func Resolve(f File, flags Flags, getenv func(string) string) (Resolved, error) {
	r, err := pick(f, flags, getenv)
	if err != nil {
		return Resolved{}, err
	}
	r.Region = cmp.Or(flags.Region, getenv(EnvRegion), r.Region, DefaultRegion)
	host, err := regionHost(r.Region)
	if err != nil {
		return Resolved{}, err
	}
	r.Host = host
	return r, nil
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
