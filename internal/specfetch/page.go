package specfetch

import (
	"bytes"
	"errors"
	"fmt"
	"time"
)

var (
	definitionHeading = []byte("\n# OpenAPI definition\n")
	jsonFence         = []byte("```json\n")
	// JSON strings can't hold a raw newline, so the first "\n```" after the
	// opening fence always ends the document.
	closingFence = []byte("\n```")
)

// ExtractSpec returns the OpenAPI document a reference page embeds. Fences
// before the "# OpenAPI definition" heading are prose examples.
func ExtractSpec(page []byte) ([]byte, error) {
	_, rest, ok := bytes.Cut(page, definitionHeading)
	if !ok {
		return nil, errors.New(`no "# OpenAPI definition" heading`)
	}
	_, rest, ok = bytes.Cut(rest, jsonFence)
	if !ok {
		return nil, errors.New("no json block after the OpenAPI definition heading")
	}
	doc, _, ok := bytes.Cut(rest, closingFence)
	if !ok {
		return nil, errors.New("unterminated json block")
	}
	return doc, nil
}

// UpdatedAt reads the updatedAt field of a page's front matter.
func UpdatedAt(page []byte) (time.Time, error) {
	front, ok := bytes.CutPrefix(page, []byte("---\n"))
	if !ok {
		return time.Time{}, errors.New("no front matter")
	}
	front, _, ok = bytes.Cut(front, []byte("\n---\n"))
	if !ok {
		return time.Time{}, errors.New("unterminated front matter")
	}
	for line := range bytes.Lines(front) {
		v, ok := bytes.CutPrefix(bytes.TrimSpace(line), []byte("updatedAt:"))
		if !ok {
			continue
		}
		t, err := time.Parse(time.RFC3339, string(bytes.TrimSpace(v)))
		if err != nil {
			return time.Time{}, fmt.Errorf("parse updatedAt: %w", err)
		}
		return t, nil
	}
	return time.Time{}, errors.New("front matter has no updatedAt")
}
