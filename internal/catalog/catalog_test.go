package catalog

import (
	"regexp"
	"strings"
	"testing"
)

func TestEveryOperationHasSchemas(t *testing.T) {
	if len(All) != 206 {
		t.Errorf("All = %d ops, want 206", len(All))
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

func TestMatch(t *testing.T) {
	tests := []struct{ method, path, want string }{
		{"GET", "/repositories/code", "repo list"},
		{"GET", "/repositories/code/", "repo list"},
		{"PUT", "/issues/groups/12/ignore", "issue-group ignore"},
		{"GET", "/issues/12", "issue get"},
		// /issues/{issue_id} fits too; the literal segment wins.
		{"GET", "/issues/export", "issue export"},
		// Spellings a server may route alike match alike, so api's guard can't be dodged.
		{"POST", "/repositories/code/Deactivate", "repo deactivate"},
		{"POST", "/repositories/code/%64eactivate", "repo deactivate"},
		{"POST", "/repositories//code/deactivate", "repo deactivate"},
		{"POST", "/repositories/code%2Fdeactivate", "repo deactivate"},
		{"POST", "/repositories%2Fcode%2Fdeactivate", "repo deactivate"},
		{"POST", `/repositories\code\deactivate`, "repo deactivate"},
		{"DELETE", "/repositories/code", ""},
		{"GET", "/nosuch", ""},
	}
	for _, tt := range tests {
		op, ok := Match(tt.method, tt.path)
		if op.Command != tt.want || ok != (tt.want != "") {
			t.Errorf("Match(%s %s) = %q, %v; want %q", tt.method, tt.path, op.Command, ok, tt.want)
		}
	}
}

// Normalising spellings must not move a real path to another operation.
func TestEveryOperationMatchesItsOwnPath(t *testing.T) {
	placeholder := regexp.MustCompile(`\{[^}]+\}`)
	for _, op := range All {
		path := placeholder.ReplaceAllString(op.Path, "1")
		if got, ok := Match(op.Method, path); !ok || got.ID != op.ID {
			t.Errorf("Match(%s %s) = %s, want %s", op.Method, path, got.Command, op.Command)
		}
	}
}
