// Package gen turns the vendored OpenAPI spec and the overlay into the
// command catalog: a descriptor per operation, the input and response
// schemas, and the search index. It runs under go generate only.
package gen

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
)

// JSONSchema is a JSON Schema object as decoded from the spec.
type JSONSchema = map[string]any

type Spec struct {
	Ops []SpecOp // sorted by ID
}

type SpecOp struct {
	ID          string
	Method      string
	Path        string
	Summary     string
	Description string
	Scope       string // "" when the spec names none
	Params      []SpecParam
	Body        *SpecBody
	Response    SpecResponse
}

type SpecParam struct {
	Name        string
	In          string // "path" or "query"
	Required    bool
	Description string
	Schema      JSONSchema
}

type SpecBody struct {
	Required bool
	Schema   JSONSchema
}

type SpecResponse struct {
	Types  []string   // content types of the first 2xx response, sorted
	Schema JSONSchema // its application/json schema, or nil
}

type rawParam struct {
	Ref         string     `json:"$ref"`
	Name        string     `json:"name"`
	In          string     `json:"in"`
	Required    bool       `json:"required"`
	Description string     `json:"description"`
	Schema      JSONSchema `json:"schema"`
}

type rawMedia map[string]struct {
	Schema JSONSchema `json:"schema"`
}

type rawOp struct {
	OperationID string     `json:"operationId"`
	Summary     string     `json:"summary"`
	Description string     `json:"description"`
	Parameters  []rawParam `json:"parameters"`
	RequestBody *struct {
		Required bool     `json:"required"`
		Content  rawMedia `json:"content"`
	} `json:"requestBody"`
	Responses map[string]struct {
		Content rawMedia `json:"content"`
	} `json:"responses"`
	Security []map[string][]string `json:"security"`
}

var methods = map[string]bool{"get": true, "post": true, "put": true, "patch": true, "delete": true}

// LoadSpec decodes the spec and normalises its schema quirks. It fails on
// anything the generator would otherwise have to guess about.
func LoadSpec(data []byte) (*Spec, error) {
	var raw struct {
		Paths      map[string]map[string]json.RawMessage `json:"paths"`
		Components struct {
			Parameters map[string]rawParam `json:"parameters"`
		} `json:"components"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("decode the spec: %w", err)
	}
	s := &Spec{}
	seen := map[string]bool{}
	for _, path := range slices.Sorted(maps.Keys(raw.Paths)) {
		for _, key := range slices.Sorted(maps.Keys(raw.Paths[path])) {
			if !methods[key] {
				return nil, fmt.Errorf("%s: unsupported path item key %q", path, key)
			}
			method := strings.ToUpper(key)
			var ro rawOp
			if err := json.Unmarshal(raw.Paths[path][key], &ro); err != nil {
				return nil, fmt.Errorf("%s %s: %w", method, path, err)
			}
			op, err := specOp(method, path, ro, raw.Components.Parameters)
			if err != nil {
				return nil, fmt.Errorf("%s %s: %w", method, path, err)
			}
			if seen[op.ID] {
				return nil, fmt.Errorf("operation ID %s is used twice", op.ID)
			}
			seen[op.ID] = true
			s.Ops = append(s.Ops, op)
		}
	}
	slices.SortFunc(s.Ops, func(a, b SpecOp) int { return strings.Compare(a.ID, b.ID) })
	return s, nil
}

func specOp(method, path string, ro rawOp, shared map[string]rawParam) (SpecOp, error) {
	if ro.OperationID == "" {
		return SpecOp{}, errors.New("no operationId")
	}
	scope, err := scopeOf(ro.Security)
	if err != nil {
		return SpecOp{}, err
	}
	op := SpecOp{ID: ro.OperationID, Method: method, Path: path, Summary: ro.Summary, Description: ro.Description, Scope: scope}
	for _, rp := range ro.Parameters {
		p, err := param(rp, shared)
		if err != nil {
			return SpecOp{}, err
		}
		op.Params = append(op.Params, p)
	}
	if rb := ro.RequestBody; rb != nil {
		media, ok := rb.Content["application/json"]
		if len(rb.Content) != 1 || !ok || media.Schema == nil {
			return SpecOp{}, fmt.Errorf("request body types %v: only application/json is supported", slices.Sorted(maps.Keys(rb.Content)))
		}
		if err := input(media.Schema, "request body"); err != nil {
			return SpecOp{}, err
		}
		op.Body = &SpecBody{Required: rb.Required, Schema: media.Schema}
	}
	op.Response = response(ro)
	return op, nil
}

// scopeOf returns the one oauth scope an operation names. The missing-scope
// hint can name only one.
func scopeOf(security []map[string][]string) (string, error) {
	if len(security) == 0 {
		return "", nil
	}
	scopes, ok := security[0]["oauth"]
	if len(security) > 1 || len(security[0]) != 1 || !ok || len(scopes) > 1 {
		return "", fmt.Errorf("security %v: want at most one oauth scope", security)
	}
	if len(scopes) == 0 {
		return "", nil
	}
	return scopes[0], nil
}

func param(rp rawParam, shared map[string]rawParam) (SpecParam, error) {
	if rp.Ref != "" {
		name, isParam := strings.CutPrefix(rp.Ref, "#/components/parameters/")
		p, found := shared[name]
		if !isParam || !found {
			return SpecParam{}, fmt.Errorf("unresolvable parameter $ref %q", rp.Ref)
		}
		rp = p
	}
	if rp.In != "path" && rp.In != "query" {
		return SpecParam{}, fmt.Errorf("parameter %s is in %q: only path and query parameters are supported", rp.Name, rp.In)
	}
	if rp.Schema == nil {
		return SpecParam{}, fmt.Errorf("parameter %s has no schema", rp.Name)
	}
	if err := input(rp.Schema, "parameter "+rp.Name); err != nil {
		return SpecParam{}, err
	}
	return SpecParam{Name: rp.Name, In: rp.In, Required: rp.Required, Description: rp.Description, Schema: rp.Schema}, nil
}

func response(ro rawOp) SpecResponse {
	for _, status := range slices.Sorted(maps.Keys(ro.Responses)) {
		if !strings.HasPrefix(status, "2") {
			continue
		}
		content := ro.Responses[status].Content
		r := SpecResponse{Types: slices.Sorted(maps.Keys(content))}
		if media, ok := content["application/json"]; ok && media.Schema != nil {
			normalize(media.Schema)
			r.Schema = media.Schema
		}
		return r
	}
	return SpecResponse{}
}

// inputKeywords are what an input schema may use: the ten keywords the
// validator enforces, plus annotations it can safely ignore.
var inputKeywords = map[string]bool{
	"type": true, "properties": true, "required": true, "enum": true, "const": true,
	"minimum": true, "maximum": true, "items": true, "minItems": true, "oneOf": true,
	"description": true, "title": true, "default": true, "example": true, "format": true,
}

// input normalises an input schema and refuses a keyword the validator
// doesn't enforce, which would otherwise let bad input through unchecked.
func input(s JSONSchema, where string) error {
	normalize(s)
	return checkKeywords(s, where)
}

func checkKeywords(s JSONSchema, where string) error {
	for _, k := range slices.Sorted(maps.Keys(s)) {
		if !inputKeywords[k] {
			return fmt.Errorf("%s: schema keyword %q is not supported in inputs", where, k)
		}
	}
	for _, sub := range subschemas(s) {
		if err := checkKeywords(sub, where); err != nil {
			return err
		}
	}
	return nil
}

// normalize fixes two spec quirks in place, at every depth: many schemas give
// their type only as a format ("object", "Object" or "array"), and some type
// lists hold a JSON null where the string "null" is meant.
func normalize(s JSONSchema) {
	if _, typed := s["type"]; !typed {
		if f, ok := s["format"].(string); ok && (strings.EqualFold(f, "object") || f == "array") {
			s["type"] = strings.ToLower(f)
			delete(s, "format")
		}
	}
	if types, ok := s["type"].([]any); ok {
		for i, t := range types {
			if t == nil {
				types[i] = "null"
			}
		}
	}
	for _, sub := range subschemas(s) {
		normalize(sub)
	}
}

// subschemas returns the schemas nested directly in s.
func subschemas(s JSONSchema) []JSONSchema {
	var out []JSONSchema
	if props, ok := s["properties"].(map[string]any); ok {
		for _, name := range slices.Sorted(maps.Keys(props)) {
			if p, ok := props[name].(map[string]any); ok {
				out = append(out, p)
			}
		}
	}
	if items, ok := s["items"].(map[string]any); ok {
		out = append(out, items)
	}
	for _, key := range []string{"oneOf", "anyOf", "allOf"} {
		if list, ok := s[key].([]any); ok {
			for _, b := range list {
				if bs, ok := b.(map[string]any); ok {
					out = append(out, bs)
				}
			}
		}
	}
	return out
}
