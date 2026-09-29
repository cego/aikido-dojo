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

// Write replaces dir/openapi.json and dir/snapshot.json. Each goes through a
// temp file and a rename, so a failed run leaves the previous snapshot.
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
	if err := writeAtomic(filepath.Join(dir, "openapi.json"), specJSON); err != nil {
		return err
	}
	return writeAtomic(filepath.Join(dir, "snapshot.json"), snapJSON)
}

func writeAtomic(path string, data []byte) error {
	f, err := os.CreateTemp(filepath.Dir(path), ".specfetch-*")
	if err != nil {
		return fmt.Errorf("create temp for %s: %w", path, err)
	}
	defer os.Remove(f.Name()) // fails harmlessly once the rename has moved it
	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("write %s: %w", path, errors.Join(err, f.Close()))
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("close %s: %w", path, err)
	}
	if err := os.Rename(f.Name(), path); err != nil {
		return fmt.Errorf("rename to %s: %w", path, err)
	}
	return nil
}
