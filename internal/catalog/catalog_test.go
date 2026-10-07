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

func TestSearchIndexHasEveryCommand(t *testing.T) {
	idx, err := SearchIndex()
	if err != nil {
		t.Fatal(err)
	}
	indexed := map[string]bool{}
	for _, d := range idx.Docs {
		indexed[d.Command] = true
	}
	for _, op := range All {
		if !indexed[op.Command] {
			t.Errorf("%s is not in the search index", op.Command)
		}
	}
	if len(idx.Docs) != len(All) || idx.AvgLen <= 0 {
		t.Errorf("index has %d docs and avg_len %v, want %d and above 0", len(idx.Docs), idx.AvgLen, len(All))
	}
}

func TestCommand(t *testing.T) {
	if op, ok := Command("repo list"); !ok || op.ID != "listCodeRepos" {
		t.Errorf("Command(repo list) = %s, %v", op.ID, ok)
	}
	if _, ok := Command("repo nosuch"); ok {
		t.Error("found a command that doesn't exist")
	}
}
