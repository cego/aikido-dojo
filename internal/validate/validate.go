// Package validate checks a decoded JSON value against the part of JSON
// Schema the spec's inputs use: type, properties, required, enum, const,
// minimum, maximum, items, minItems and oneOf. A value outside an enum only
// warns, because the spec's enums are known to be incomplete.
package validate

import (
	"encoding/json"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"
)

type Issue struct {
	Path []string // property names, and [i] for array items
	Msg  string
}

type Result struct {
	Errors   []Issue
	Warnings []Issue
}

// Check validates v, decoded with json.Decoder.UseNumber, against schema.
func Check(schema map[string]any, v any) Result {
	var r Result
	check(schema, v, nil, &r)
	return r
}

// Join renders path after root: Join("body", ["a", "[0]", "b"]) is body.a[0].b.
func Join(root string, path []string) string {
	var b strings.Builder
	b.WriteString(root)
	for _, seg := range path {
		if b.Len() > 0 && !strings.HasPrefix(seg, "[") {
			b.WriteByte('.')
		}
		b.WriteString(seg)
	}
	return b.String()
}

func (r *Result) fail(path []string, msg string) {
	r.Errors = append(r.Errors, Issue{Path: slices.Clone(path), Msg: msg})
}

func (r *Result) warn(path []string, msg string) {
	r.Warnings = append(r.Warnings, Issue{Path: slices.Clone(path), Msg: msg})
}

func check(s map[string]any, v any, path []string, r *Result) {
	if s == nil {
		return
	}
	if types := typeList(s["type"]); len(types) > 0 && !slices.ContainsFunc(types, func(t string) bool { return is(v, t) }) {
		r.fail(path, fmt.Sprintf("want %s, got %s", strings.Join(types, " or "), kindName(v)))
		return
	}
	if c, ok := s["const"]; ok && !equal(c, v) {
		r.fail(path, "want "+show(c)+", got "+show(v))
		return
	}
	if enum, ok := s["enum"].([]any); ok && !slices.ContainsFunc(enum, func(e any) bool { return equal(e, v) }) {
		shown := make([]string, len(enum))
		for i, e := range enum {
			shown[i] = show(e)
		}
		r.warn(path, fmt.Sprintf("%s is not one of %s in the spec; sending it anyway", show(v), strings.Join(shown, ", ")))
	}
	switch x := v.(type) {
	case json.Number:
		checkRange(s, x, path, r)
	case map[string]any:
		checkObject(s, x, path, r)
	case []any:
		checkArray(s, x, path, r)
	}
	if branches, ok := s["oneOf"].([]any); ok {
		checkOneOf(branches, v, path, r)
	}
}

func checkRange(s map[string]any, n json.Number, path []string, r *Result) {
	f, err := n.Float64()
	if err != nil {
		return
	}
	if lo, ok := number(s["minimum"]); ok && f < lo {
		r.fail(path, "want at least "+show(s["minimum"])+", got "+n.String())
	}
	if hi, ok := number(s["maximum"]); ok && f > hi {
		r.fail(path, "want at most "+show(s["maximum"])+", got "+n.String())
	}
}

func checkObject(s map[string]any, obj map[string]any, path []string, r *Result) {
	if required, ok := s["required"].([]any); ok {
		for _, name := range required {
			if n, ok := name.(string); ok {
				if _, present := obj[n]; !present {
					r.fail(append(path, n), "required")
				}
			}
		}
	}
	props, _ := s["properties"].(map[string]any)
	for _, name := range slices.Sorted(maps.Keys(obj)) {
		if sub, ok := props[name].(map[string]any); ok {
			check(sub, obj[name], append(path, name), r)
		}
	}
}

func checkArray(s map[string]any, arr []any, path []string, r *Result) {
	if lo, ok := number(s["minItems"]); ok && float64(len(arr)) < lo {
		r.fail(path, fmt.Sprintf("want at least %s items, got %d", show(s["minItems"]), len(arr)))
	}
	items, _ := s["items"].(map[string]any)
	for i, v := range arr {
		check(items, v, append(path, fmt.Sprintf("[%d]", i)), r)
	}
}

// checkOneOf wants exactly one shape to match. When none does, it reports
// each shape's first error, so the fix for the shape the caller meant is visible.
func checkOneOf(branches []any, v any, path []string, r *Result) {
	var matched []Result
	var misses []string
	for i, b := range branches {
		bs, _ := b.(map[string]any)
		var br Result
		check(bs, v, path, &br)
		if len(br.Errors) == 0 {
			matched = append(matched, br)
			continue
		}
		e := br.Errors[0]
		misses = append(misses, fmt.Sprintf("shape %d: %s: %s", i+1, Join("", e.Path), e.Msg))
	}
	switch len(matched) {
	case 1:
		r.Warnings = append(r.Warnings, matched[0].Warnings...)
	case 0:
		r.fail(path, "matches none of the allowed shapes ("+strings.Join(misses, "; ")+")")
	default:
		r.fail(path, fmt.Sprintf("matches %d of the allowed shapes; exactly one must match", len(matched)))
	}
}

func typeList(t any) []string {
	switch x := t.(type) {
	case string:
		return []string{x}
	case []any:
		var out []string
		for _, e := range x {
			if s, ok := e.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}

func is(v any, t string) bool {
	switch t {
	case "null":
		return v == nil
	case "boolean":
		_, ok := v.(bool)
		return ok
	case "string":
		_, ok := v.(string)
		return ok
	case "object":
		_, ok := v.(map[string]any)
		return ok
	case "array":
		_, ok := v.([]any)
		return ok
	case "number":
		_, ok := v.(json.Number)
		return ok
	case "integer":
		n, ok := v.(json.Number)
		if !ok {
			return false
		}
		_, err := n.Int64()
		return err == nil
	}
	return false
}

func kindName(v any) string {
	switch v.(type) {
	case nil:
		return "null"
	case bool:
		return "boolean"
	case string:
		return "string"
	case map[string]any:
		return "object"
	case []any:
		return "array"
	}
	return "number"
}

// equal compares JSON values; numbers compare by value, so 1 and 1.0 match.
func equal(a, b any) bool {
	an, aok := number(a)
	bn, bok := number(b)
	if aok || bok {
		return aok && bok && an == bn
	}
	return reflect.DeepEqual(a, b)
}

func number(v any) (float64, bool) {
	switch n := v.(type) {
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	case float64:
		return n, true
	}
	return 0, false
}

func show(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprint(v)
	}
	return string(b)
}
