package gen

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/cego/aikido-dojo/internal/overlay"
)

// The committed catalog must be what the spec and overlay generate, so a
// hand edit or a forgotten go generate fails CI.
func TestGeneratedFilesAreCurrent(t *testing.T) {
	spec, err := os.ReadFile("../../spec/openapi.json")
	if err != nil {
		t.Fatal(err)
	}
	out, err := Generate(spec, overlay.Ops)
	if err != nil {
		t.Fatal(err)
	}
	for name, want := range map[string][]byte{"zz_ops.go": out.Ops, "schemas.json": out.Schemas, "search.json": out.Search} {
		got, err := os.ReadFile(filepath.Join("..", "catalog", name)) //nolint:gosec // name is one of the three generated files
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(got, want) {
			t.Errorf("internal/catalog/%s is stale: run go generate ./internal/catalog", name)
		}
	}
}
