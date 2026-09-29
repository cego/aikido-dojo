package specfetch

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"
)

func TestWriteIsCanonical(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "spec")
	spec := decoded(t, `{"paths":{"/b":{},"/a":{"get":{"description":"<b> & c","maximum":9007199254740993}}},"openapi":"3.1.0"}`)
	updated := time.Date(2026, 6, 1, 8, 0, 0, 0, time.FixedZone("CEST", 2*3600))

	obj, err := object(spec, "spec")
	if err != nil {
		t.Fatal(err)
	}
	if err := Write(dir, obj, updated); err != nil {
		t.Fatal(err)
	}

	wantSpec := `{
  "openapi": "3.1.0",
  "paths": {
    "/a": {
      "get": {
        "description": "<b> & c",
        "maximum": 9007199254740993
      }
    },
    "/b": {}
  }
}
`
	if got := readFile(t, filepath.Join(dir, "openapi.json")); got != wantSpec {
		t.Errorf("openapi.json =\n%s\nwant\n%s", got, wantSpec)
	}
	if got, want := readFile(t, filepath.Join(dir, "snapshot.json")), "{\n  \"updated_at\": \"2026-06-01T06:00:00Z\"\n}\n"; got != want {
		t.Errorf("snapshot.json = %q, want %q", got, want)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	if want := []string{"openapi.json", "snapshot.json"}; !slices.Equal(names, want) {
		t.Errorf("dir holds %q, want %q", names, want)
	}
}

func TestWriteFailureLeavesNoTempFiles(t *testing.T) {
	dir := t.TempDir()
	// A non-empty directory where openapi.json should be makes the rename fail.
	if err := os.MkdirAll(filepath.Join(dir, "openapi.json", "x"), 0o750); err != nil {
		t.Fatal(err)
	}
	err := Write(dir, map[string]any{"openapi": "3.1.0"}, time.Now())
	wantErr(t, err, "rename")
	matches, err := filepath.Glob(filepath.Join(dir, ".specfetch-*"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 0 {
		t.Errorf("temp files left behind: %q", matches)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path) //nolint:gosec // path comes from t.TempDir
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
