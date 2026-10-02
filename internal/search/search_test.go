package search

import (
	"slices"
	"testing"
)

func TestTokens(t *testing.T) {
	tests := map[string][]string{
		"List code repositories":        {"list", "code", "repository"},
		"Ignore an issue group":         {"ignore", "issue", "group"},
		"filter_code_repo_id":           {"filter", "code", "repo", "id"},
		"Get the CVE-2024-1234 details": {"get", "cve", "2024", "1234", "detail"},
		"repos findings rules":          {"repo", "finding", "rule"},
		"status address analysis gas":   {"status", "address", "analysis", "gas"},
		"issue-group-note create":       {"issue", "group", "note", "create"},
	}
	for in, want := range tests {
		if got := Tokens(in); !slices.Equal(got, want) {
			t.Errorf("Tokens(%q) = %q, want %q", in, got, want)
		}
	}
}
