package gen

import (
	"testing"

	"github.com/cego/aikido-dojo/internal/ops"
)

func TestBuildIndex(t *testing.T) {
	m := &Model{
		Ops: []ops.Op{
			{ID: "a", Command: "repo list", Summary: "List repos", Description: "Every repo."},
			{ID: "b", Command: "issue ignore", Summary: "Ignore an issue"},
		},
		Schemas: map[string]ops.SchemaSet{
			"a": {Flags: map[string]any{"properties": map[string]any{"filter_name": map[string]any{}}}},
			"b": {Args: map[string]any{"properties": map[string]any{"issue_id": map[string]any{}}},
				Body: map[string]any{"oneOf": []any{map[string]any{"properties": map[string]any{"reason": map[string]any{}}}}}},
		},
	}
	idx := buildIndex(m)
	if len(idx.Docs) != 2 {
		t.Fatalf("docs = %d, want 2", len(idx.Docs))
	}
	a := idx.Docs[0]
	// repo list + list repo + every repo + filter name
	if a.Command != "repo list" || a.Summary != "List repos" || a.Len != 8 || a.TF["repo"] != 3 || a.TF["list"] != 2 || a.TF["filter"] != 1 {
		t.Errorf("doc a = %+v", a)
	}
	b := idx.Docs[1]
	// issue ignore + ignore issue + issue id + reason
	if b.Len != 7 || b.TF["issue"] != 3 || b.TF["reason"] != 1 {
		t.Errorf("doc b = %+v", b)
	}
	if idx.DF["repo"] != 1 || idx.DF["issue"] != 1 || idx.AvgLen != 7.5 {
		t.Errorf("df = %v, avg_len = %v", idx.DF, idx.AvgLen)
	}
}
