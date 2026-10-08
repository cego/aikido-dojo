package gen

import (
	"encoding/json"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
)

// specWith returns a spec holding the given path items, plus one shared
// parameter that tests can $ref.
func specWith(paths string) []byte {
	return []byte(`{"paths": ` + paths + `, "components": {"parameters": {"PageParam": {"in": "query", "name": "page", "schema": {"type": "integer", "minimum": 0}}}}}`)
}

func loadOne(t *testing.T, paths string) SpecOp {
	t.Helper()
	s, err := LoadSpec(specWith(paths))
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Ops) != 1 {
		t.Fatalf("ops = %d, want 1", len(s.Ops))
	}
	return s.Ops[0]
}

func TestLoadSpecNormalisesSchemas(t *testing.T) {
	op := loadOne(t, `{"/things": {"post": {"operationId": "addThing",
		"requestBody": {"required": true, "content": {"application/json": {"schema":
			{"format": "Object", "properties": {"tags": {"format": "array", "items": {"type": ["string", null]}}}}}}},
		"responses": {"200": {"content": {"application/json": {"schema": {"format": "object"}}}}}}}}`)
	got, err := json.Marshal(op.Body.Schema)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"properties":{"tags":{"items":{"type":["string","null"]},"type":"array"}},"type":"object"}`
	if string(got) != want {
		t.Errorf("body schema = %s, want %s", got, want)
	}
	if op.Response.Schema["type"] != "object" {
		t.Errorf("response schema = %v, want type object", op.Response.Schema)
	}
}

func TestLoadSpecReadsAnOperation(t *testing.T) {
	op := loadOne(t, `{"/things/{thing_id}": {"get": {"operationId": "getThing", "summary": "Get a thing", "description": "All of it.",
		"security": [{"oauth": ["things:read"]}],
		"parameters": [{"$ref": "#/components/parameters/PageParam"},
			{"in": "path", "name": "thing_id", "required": true, "description": "The thing.", "schema": {"type": "integer"}}],
		"responses": {"400": {}, "201": {"content": {"text/csv": {}, "application/json": {"schema": {"type": "array"}}}}}}}}`)
	if op.ID != "getThing" || op.Method != http.MethodGet || op.Path != "/things/{thing_id}" || op.Scope != "things:read" ||
		op.Summary != "Get a thing" || op.Description != "All of it." {
		t.Errorf("op = %+v", op)
	}
	if len(op.Params) != 2 || op.Params[0].Name != "page" || op.Params[0].In != "query" ||
		op.Params[1].Name != "thing_id" || !op.Params[1].Required {
		t.Errorf("params = %+v, want the shared page parameter resolved, then thing_id", op.Params)
	}
	if !reflect.DeepEqual(op.Response.Types, []string{"application/json", "text/csv"}) || op.Response.Schema["type"] != "array" {
		t.Errorf("response = %+v, want the 201's types and schema", op.Response)
	}
}

func TestLoadSpecRejects(t *testing.T) {
	tests := []struct {
		name, paths, want string
	}{
		{"an unresolvable $ref", `{"/x": {"get": {"operationId": "a", "parameters": [{"$ref": "#/components/parameters/Nope"}]}}}`, "unresolvable parameter $ref"},
		{"a header parameter", `{"/x": {"get": {"operationId": "a", "parameters": [{"in": "header", "name": "h", "schema": {"type": "string"}}]}}}`, `is in "header"`},
		{"a parameter without a schema", `{"/x": {"get": {"operationId": "a", "parameters": [{"in": "query", "name": "q"}]}}}`, "has no schema"},
		{"a non-JSON body", `{"/x": {"post": {"operationId": "a", "requestBody": {"content": {"text/plain": {"schema": {"type": "string"}}}}}}}`, "only application/json"},
		{"an input keyword the validator doesn't enforce", `{"/x": {"get": {"operationId": "a", "parameters": [{"in": "query", "name": "q", "schema": {"type": "string", "pattern": "^a"}}]}}}`, `keyword "pattern"`},
		{"a nested unsupported keyword", `{"/x": {"post": {"operationId": "a", "requestBody": {"content": {"application/json": {"schema": {"type": "object", "properties": {"n": {"type": "integer", "multipleOf": 2}}}}}}}}}`, `keyword "multipleOf"`},
		{"two scopes", `{"/x": {"get": {"operationId": "a", "security": [{"oauth": ["a:read", "b:read"]}]}}}`, "at most one oauth scope"},
		{"a non-oauth scheme", `{"/x": {"get": {"operationId": "a", "security": [{"apiKey": []}]}}}`, "at most one oauth scope"},
		{"a path-level key", `{"/x": {"parameters": []}}`, `unsupported path item key "parameters"`},
		{"a missing operationId", `{"/x": {"get": {}}}`, "no operationId"},
		{"a duplicate operationId", `{"/x": {"get": {"operationId": "a"}}, "/y": {"get": {"operationId": "a"}}}`, "used twice"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoadSpec(specWith(tt.paths))
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("err = %v, want it to contain %q", err, tt.want)
			}
		})
	}
}

func TestLoadVendoredSpec(t *testing.T) {
	data, err := os.ReadFile("../../spec/openapi.json")
	if err != nil {
		t.Fatal(err)
	}
	s, err := LoadSpec(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Ops) != 206 {
		t.Errorf("ops = %d, want the 206 of the 2026-10-07 snapshot", len(s.Ops))
	}
}
