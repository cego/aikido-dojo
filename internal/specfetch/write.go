package specfetch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot is spec/snapshot.json. UpdatedAt is the newest page updatedAt
// rather than the fetch time, so a run over unchanged docs changes no file.
type Snapshot struct {
	UpdatedAt time.Time `json:"updated_at"`
}

// Canonical renders v with sorted keys, two-space indent and no HTML
// escaping, so the vendored spec diffs cleanly against Aikido's own.
func Canonical(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return nil, fmt.Errorf("encode: %w", err)
	}
	return buf.Bytes(), nil
}

// Write replaces dir/openapi.json and dir/snapshot.json. Both are written in
// full to temp files before either is renamed, so a failure while writing
// keeps the previous pair. Only a failed second rename can leave a new spec
// beside the old snapshot, and Write then returns that error.
func Write(dir string, spec map[string]any, updated time.Time) error {
	specJSON, err := Canonical(spec)
	if err != nil {
		return fmt.Errorf("render spec: %w", err)
	}
	snapJSON, err := Canonical(Snapshot{UpdatedAt: updated.UTC()})
	if err != nil {
		return fmt.Errorf("render snapshot: %w", err)
	}
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("create %s: %w", dir, err)
	}
	specTmp, err := stage(dir, specJSON)
	if err != nil {
		return err
	}
	defer os.Remove(specTmp) // fails harmlessly once the rename has moved it
	snapTmp, err := stage(dir, snapJSON)
	if err != nil {
		return err
	}
	defer os.Remove(snapTmp)
	if err := os.Rename(specTmp, filepath.Join(dir, "openapi.json")); err != nil {
		return fmt.Errorf("rename to %s: %w", filepath.Join(dir, "openapi.json"), err)
	}
	if err := os.Rename(snapTmp, filepath.Join(dir, "snapshot.json")); err != nil {
		return fmt.Errorf("rename to %s: %w", filepath.Join(dir, "snapshot.json"), err)
	}
	return nil
}

// stage writes data to a new temp file in dir and returns its path. It
// removes the file itself when it fails.
func stage(dir string, data []byte) (name string, err error) {
	f, err := os.CreateTemp(dir, ".specfetch-*")
	if err != nil {
		return "", fmt.Errorf("create temp in %s: %w", dir, err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, os.Remove(f.Name()))
		}
	}()
	if _, err := f.Write(data); err != nil {
		return "", fmt.Errorf("write %s: %w", f.Name(), errors.Join(err, f.Close()))
	}
	if err := f.Close(); err != nil {
		return "", fmt.Errorf("close %s: %w", f.Name(), err)
	}
	return f.Name(), nil
}
