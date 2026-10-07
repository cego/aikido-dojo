// Package catalog holds the generated command catalog: zz_ops.go,
// schemas.json and search.json. Regenerate with go generate ./internal/catalog
// after changing spec/openapi.json or internal/overlay.
package catalog

//go:generate go run ../tools/gen -spec ../../spec/openapi.json -snapshot ../../spec/snapshot.json -dir .

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/cego/aikido-dojo/internal/ops"
	"github.com/cego/aikido-dojo/internal/search"
)

//go:embed schemas.json
var schemasJSON []byte

// loadSchemas decodes schemas.json on first use only, so a command that
// needs no schema doesn't pay for it. Numbers stay json.Number for the validator.
var loadSchemas = sync.OnceValues(func() (map[string]ops.SchemaSet, error) {
	var m map[string]ops.SchemaSet
	dec := json.NewDecoder(bytes.NewReader(schemasJSON))
	dec.UseNumber()
	if err := dec.Decode(&m); err != nil {
		return nil, fmt.Errorf("decode the embedded schemas: %w", err)
	}
	return m, nil
})

func Schemas(id string) (ops.SchemaSet, error) {
	m, err := loadSchemas()
	if err != nil {
		return ops.SchemaSet{}, err
	}
	s, ok := m[id]
	if !ok {
		return ops.SchemaSet{}, fmt.Errorf("no schemas for operation %s", id)
	}
	return s, nil
}

//go:embed search.json
var searchJSON []byte

// SearchIndex decodes the embedded index. Only search needs it, so it isn't
// decoded at startup.
func SearchIndex() (*search.Index, error) {
	var idx search.Index
	if err := json.Unmarshal(searchJSON, &idx); err != nil {
		return nil, fmt.Errorf("decode the embedded search index: %w", err)
	}
	return &idx, nil
}

// Command finds a generated command by name, such as "repo list".
func Command(name string) (ops.Op, bool) {
	for _, op := range All {
		if op.Command == name {
			return op, true
		}
	}
	return ops.Op{}, false
}
