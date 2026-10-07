package search

import (
	"maps"
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

// index builds an Index the way the generator does, from one text per command.
func index(texts map[string]string) *Index {
	idx := &Index{DF: map[string]int{}}
	total := 0
	for _, command := range slices.Sorted(maps.Keys(texts)) {
		d := Doc{Command: command, Summary: "summary of " + command, TF: map[string]int{}}
		for _, tok := range Tokens(texts[command]) {
			d.TF[tok]++
			d.Len++
		}
		for tok := range d.TF {
			idx.DF[tok]++
		}
		idx.Docs = append(idx.Docs, d)
		total += d.Len
	}
	idx.AvgLen = float64(total) / float64(len(idx.Docs))
	return idx
}

func TestRank(t *testing.T) {
	idx := index(map[string]string{
		"repo list":          "list code repositories",
		"repo scan":          "scan a code repository",
		"issue list":         "list issues",
		"issue ignore":       "ignore an issue",
		"issue-group ignore": "ignore an issue group and every issue in the group at once",
	})
	tests := []struct {
		name  string
		query string
		n     int
		want  []string
	}{
		{"a rare word outweighs a common one", "scan repositories", 10, []string{"repo scan", "repo list"}},
		{"a short description beats a long one", "ignore", 10, []string{"issue ignore", "issue-group ignore"}},
		{"a repeated word counts once", "ignore ignore", 10, []string{"issue ignore", "issue-group ignore"}},
		{"ties go by command name", "issue", 10, []string{"issue ignore", "issue list", "issue-group ignore"}},
		{"n caps the hits", "list", 1, []string{"issue list"}},
		{"no shared word, no hit", "nothing here", 10, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hits := idx.Rank(tt.query, tt.n)
			got := []string{}
			for _, h := range hits {
				got = append(got, h.Command)
				if h.Summary != "summary of "+h.Command || h.Score <= 0 {
					t.Errorf("hit %+v lacks its summary or a score", h)
				}
			}
			if hits == nil || !slices.Equal(got, tt.want) {
				t.Errorf("Rank(%q) = %q, want %q", tt.query, got, tt.want)
			}
		})
	}
	if a, b := idx.Rank("ignore", 10), idx.Rank("ignore ignore", 10); a[0].Score != b[0].Score {
		t.Errorf("repeating a word changed the score: %v vs %v", a[0].Score, b[0].Score)
	}
}
