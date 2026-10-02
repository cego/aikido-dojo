// Package search finds commands from a request in plain words. go generate
// builds the index from the spec's summaries, descriptions and parameter
// names; ranking runs locally and sends nothing anywhere.
package search

import (
	"strings"
	"unicode"
)

// Index holds the term statistics BM25 ranks with.
type Index struct {
	AvgLen float64        `json:"avg_len"`
	DF     map[string]int `json:"df"`
	Docs   []Doc          `json:"docs"`
}

type Doc struct {
	Command string         `json:"command"`
	Summary string         `json:"summary"`
	Len     int            `json:"len"`
	TF      map[string]int `json:"tf"`
}

var stopwords = map[string]bool{
	"a": true, "an": true, "the": true, "of": true, "for": true, "to": true, "in": true, "on": true,
	"and": true, "or": true, "by": true, "with": true, "is": true, "are": true, "be": true, "this": true,
	"that": true, "it": true, "as": true, "from": true, "your": true, "you": true, "can": true,
}

// Tokens splits s into lowercase words, drops stopwords and reduces plurals,
// so "List the repositories" and "list repository" meet.
func Tokens(s string) []string {
	var out []string
	for _, w := range strings.FieldsFunc(strings.ToLower(s), func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) }) {
		if !stopwords[w] {
			out = append(out, singular(w))
		}
	}
	return out
}

func singular(w string) string {
	switch {
	case len(w) > 4 && strings.HasSuffix(w, "ies"):
		return w[:len(w)-3] + "y"
	case len(w) > 3 && strings.HasSuffix(w, "s") && !strings.HasSuffix(w, "ss") && !strings.HasSuffix(w, "us") && !strings.HasSuffix(w, "is"):
		return w[:len(w)-1]
	}
	return w
}
