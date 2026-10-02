package catalog

import (
	"strings"
	"testing"
)

func TestEveryOperationHasSchemas(t *testing.T) {
	if len(All) != 201 {
		t.Errorf("All = %d ops, want 201", len(All))
	}
	for _, op := range All {
		s, err := Schemas(op.ID)
		if err != nil {
			t.Errorf("%s: %v", op.Command, err)
			continue
		}
		if s.Command != op.Command || s.Args == nil || s.Flags == nil {
			t.Errorf("%s: schemas = %+v, want its command, args and flags", op.Command, s)
		}
		if op.Body != nil && s.Body == nil {
			t.Errorf("%s has a body but no body schema", op.Command)
		}
	}
}

func TestSchemasOfAnUnknownOperation(t *testing.T) {
	if _, err := Schemas("nope"); err == nil || !strings.Contains(err.Error(), "nope") {
		t.Errorf("err = %v, want one naming the operation", err)
	}
}
