package validate

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"
)

func decode(t *testing.T, s string) any {
	t.Helper()
	dec := json.NewDecoder(strings.NewReader(s))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		t.Fatalf("decode %s: %v", s, err)
	}
	return v
}

func render(issues []Issue) []string {
	out := make([]string, 0, len(issues))
	for _, i := range issues {
		out = append(out, Join("body", i.Path)+": "+i.Msg)
	}
	return out
}

const shapes = `{"oneOf": [
	{"type": "object", "required": ["enabled"], "properties": {"enabled": {"const": true}, "ids": {"type": "array"}}},
	{"type": "object", "required": ["enabled"], "properties": {"enabled": {"const": false}}}]}`

func TestCheck(t *testing.T) {
	tests := []struct {
		name, schema, value string
		errs, warns         []string
	}{
		{name: "a type mismatch", schema: `{"type": "string"}`, value: `1`, errs: []string{"body: want string, got number"}},
		{name: "a fraction is not an integer", schema: `{"type": "integer"}`, value: `1.5`, errs: []string{"body: want integer, got number"}},
		{name: "an integer", schema: `{"type": "integer"}`, value: `3`},
		{name: "null allowed by a type list", schema: `{"type": ["string", "null"]}`, value: `null`},
		{name: "a missing required field", schema: `{"type": "object", "required": ["a", "b"], "properties": {"a": {"type": "string"}}}`, value: `{"a": "x"}`,
			errs: []string{"body.b: required"}},
		{name: "a nested field", schema: `{"type": "object", "properties": {"a": {"type": "object", "properties": {"n": {"type": "integer"}}}}}`,
			value: `{"a": {"n": "x"}}`, errs: []string{"body.a.n: want integer, got string"}},
		{name: "an array item", schema: `{"type": "array", "minItems": 2, "items": {"type": "integer"}}`, value: `[1, "x"]`,
			errs: []string{"body[1]: want integer, got string"}},
		{name: "too few items", schema: `{"type": "array", "minItems": 2}`, value: `[]`, errs: []string{"body: want at least 2 items, got 0"}},
		{name: "an unknown enum value warns", schema: `{"type": "string", "enum": ["open", "closed"]}`, value: `"auto-ignored"`,
			warns: []string{`body: "auto-ignored" is not one of "open", "closed" in the spec; sending it anyway`}},
		{name: "null in an enum", schema: `{"enum": ["low", null]}`, value: `null`},
		{name: "numbers compare by value", schema: `{"enum": [1, 2]}`, value: `1.0`},
		{name: "a const", schema: `{"const": true}`, value: `false`, errs: []string{"body: want true, got false"}},
		{name: "below the minimum", schema: `{"type": "integer", "minimum": 10, "maximum": 20}`, value: `5`, errs: []string{"body: want at least 10, got 5"}},
		{name: "above the maximum", schema: `{"type": "integer", "minimum": 10, "maximum": 20}`, value: `25`, errs: []string{"body: want at most 20, got 25"}},
		{name: "unknown properties pass", schema: `{"type": "object", "properties": {}}`, value: `{"x": 1}`},
		{name: "one shape matches", schema: shapes, value: `{"enabled": true, "ids": []}`},
		{name: "no shape matches", schema: shapes, value: `{"enabled": "yes"}`,
			errs: []string{`body: matches none of the allowed shapes (shape 1: enabled: want true, got "yes"; shape 2: enabled: want false, got "yes")`}},
		{name: "two shapes match", schema: `{"oneOf": [{"type": "object"}, {"type": "object"}]}`, value: `{}`,
			errs: []string{"body: matches 2 of the allowed shapes; exactly one must match"}},
		{name: "the matching shape's warnings are kept", schema: `{"oneOf": [{"type": "object", "properties": {"s": {"enum": ["a"]}}}, {"type": "array"}]}`,
			value: `{"s": "b"}`, warns: []string{`body.s: "b" is not one of "a" in the spec; sending it anyway`}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema, ok := decode(t, tt.schema).(map[string]any)
			if !ok {
				t.Fatal("schema is not an object")
			}
			r := Check(schema, decode(t, tt.value))
			if got := render(r.Errors); !slices.Equal(got, tt.errs) && (len(got) > 0 || len(tt.errs) > 0) {
				t.Errorf("errors = %q, want %q", got, tt.errs)
			}
			if got := render(r.Warnings); !slices.Equal(got, tt.warns) && (len(got) > 0 || len(tt.warns) > 0) {
				t.Errorf("warnings = %q, want %q", got, tt.warns)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	if got := Join("--cloud-id", []string{"[1]"}); got != "--cloud-id[1]" {
		t.Errorf("Join = %q", got)
	}
	if got := Join("", []string{"a", "[0]", "b"}); got != "a[0].b" {
		t.Errorf("Join = %q", got)
	}
}
