// Package search finds commands from a request in plain words. go generate
// builds the index from the spec's summaries, descriptions and parameter
// names; ranking runs locally and sends nothing anywhere.
package search

import (
	"cmp"
	"math"
	"slices"
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

// Hit is a command search found. Score is the BM25 score, rounded; only the
// order it gives is meaningful.
type Hit struct {
	Command string  `json:"command"`
	Summary string  `json:"summary"`
	Score   float64 `json:"score"`
}

// BM25's usual constants: k1 limits how much a repeated word counts, b how
// much a long description is discounted.
const (
	k1 = 1.2
	b  = 0.75
)

// Corpus holds what BM25 needs to know of all the documents it ranks.
type Corpus struct {
	Docs   int
	AvgLen float64
	DF     map[string]int
}

// QueryTerms are a query's tokens, each once.
func QueryTerms(query string) []string {
	return slices.Compact(slices.Sorted(slices.Values(Tokens(query))))
}

// Score is the BM25 score of a document of length n with term frequencies tf.
func (c Corpus) Score(terms []string, tf map[string]int, n int) float64 {
	if c.AvgLen == 0 {
		return 0
	}
	docs := float64(c.Docs)
	var s float64
	for _, t := range terms {
		f := float64(tf[t])
		if f == 0 {
			continue
		}
		df := float64(c.DF[t])
		idf := math.Log(1 + (docs-df+0.5)/(df+0.5))
		s += idf * f * (k1 + 1) / (f + k1*(1-b+b*float64(n)/c.AvgLen))
	}
	return s
}

// Rank scores every command against query with BM25 and returns the best n,
// best first. A command sharing no word with the query is left out.
func (idx *Index) Rank(query string, n int) []Hit {
	terms := QueryTerms(query)
	c := Corpus{Docs: len(idx.Docs), AvgLen: idx.AvgLen, DF: idx.DF}
	type scored struct {
		doc   *Doc
		score float64
	}
	var all []scored
	for i := range idx.Docs {
		d := &idx.Docs[i]
		if s := c.Score(terms, d.TF, d.Len); s > 0 {
			all = append(all, scored{d, s})
		}
	}
	slices.SortFunc(all, func(x, y scored) int {
		return cmp.Or(cmp.Compare(y.score, x.score), strings.Compare(x.doc.Command, y.doc.Command))
	})
	hits := []Hit{}
	for _, s := range all[:min(n, len(all))] {
		hits = append(hits, Hit{Command: s.doc.Command, Summary: s.doc.Summary, Score: math.Round(s.score*100) / 100})
	}
	return hits
}
