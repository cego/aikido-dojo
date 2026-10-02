package gen

import (
	"maps"
	"slices"

	"github.com/cego/aikido-dojo/internal/search"
)

// buildIndex makes one search document per command, from its command words,
// summary, description, and argument, flag and body field names.
func buildIndex(m *Model) search.Index {
	idx := search.Index{DF: map[string]int{}}
	total := 0
	for _, op := range m.Ops {
		texts := []string{op.Command, op.Summary, op.Description}
		sc := m.Schemas[op.ID]
		for _, s := range []map[string]any{sc.Args, sc.Flags, sc.Body} {
			texts = append(texts, propertyNames(s)...)
		}
		doc := search.Doc{Command: op.Command, Summary: op.Summary, TF: map[string]int{}}
		for _, text := range texts {
			for _, tok := range search.Tokens(text) {
				doc.TF[tok]++
				doc.Len++
			}
		}
		for tok := range doc.TF {
			idx.DF[tok]++
		}
		idx.Docs = append(idx.Docs, doc)
		total += doc.Len
	}
	if len(idx.Docs) > 0 {
		idx.AvgLen = float64(total) / float64(len(idx.Docs))
	}
	return idx
}

// propertyNames lists the property names of s and of the shapes it allows.
func propertyNames(s map[string]any) []string {
	if s == nil {
		return nil
	}
	var names []string
	if props, ok := s["properties"].(map[string]any); ok {
		names = slices.Sorted(maps.Keys(props))
	}
	if branches, ok := s["oneOf"].([]any); ok {
		for _, b := range branches {
			if bs, ok := b.(map[string]any); ok {
				names = append(names, propertyNames(bs)...)
			}
		}
	}
	return names
}
