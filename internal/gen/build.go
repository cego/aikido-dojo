package gen

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/cego/aikido-dojo/internal/api"
	"github.com/cego/aikido-dojo/internal/ops"
	"github.com/cego/aikido-dojo/internal/overlay"
)

// Model is the catalog before rendering.
type Model struct {
	Ops     []ops.Op                 // sorted by command
	Schemas map[string]ops.SchemaSet // by operation ID
}

// destructiveVerbs are the verbs whose every command is destructive.
var destructiveVerbs = map[string]bool{"delete": true, "deactivate": true, "rotate": true}

// Build joins the spec and the overlay into one descriptor per operation. It
// fails when they disagree, so a spec change can't add, drop or silently
// change a command.
func Build(s *Spec, overlays map[string]overlay.Op) (*Model, error) {
	var errs []error
	inSpec := map[string]bool{}
	for _, sop := range s.Ops {
		inSpec[sop.ID] = true
	}
	for _, id := range slices.Sorted(maps.Keys(overlays)) {
		if !inSpec[id] {
			errs = append(errs, fmt.Errorf("overlay entry %s has no operation in the spec", id))
		}
	}
	m := &Model{Schemas: map[string]ops.SchemaSet{}}
	for _, sop := range s.Ops {
		o, ok := overlays[sop.ID]
		if !ok {
			errs = append(errs, fmt.Errorf("operation %s (%s %s) has no overlay entry", sop.ID, sop.Method, sop.Path))
			continue
		}
		op, schemas, err := build(sop, o)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", sop.ID, err))
			continue
		}
		m.Ops = append(m.Ops, op)
		m.Schemas[op.ID] = schemas
	}
	if err := errors.Join(errs...); err != nil {
		return nil, err
	}
	slices.SortFunc(m.Ops, func(a, b ops.Op) int { return strings.Compare(a.Command, b.Command) })
	return m, nil
}

func build(sop SpecOp, o overlay.Op) (ops.Op, ops.SchemaSet, error) {
	_, verb, ok := strings.Cut(o.Cmd, " ")
	if !ok || strings.HasPrefix(o.Cmd, " ") || verb == "" || strings.Contains(verb, " ") {
		return ops.Op{}, ops.SchemaSet{}, fmt.Errorf("command %q is not <resource> <verb>", o.Cmd)
	}
	op := ops.Op{ID: sop.ID, Command: o.Cmd, Method: sop.Method, Path: sop.Path, Summary: sop.Summary,
		Description: strings.TrimSpace(sop.Description), Help: o.Help, Scope: sop.Scope, BadRequestHint: o.BadRequestHint}
	if o.Scope != "" {
		if o.Scope == sop.Scope {
			return ops.Op{}, ops.SchemaSet{}, fmt.Errorf("the overlay's Scope %s repeats the spec", o.Scope)
		}
		op.Scope = o.Scope
	}
	if destructiveVerbs[verb] && o.Destructive {
		return ops.Op{}, ops.SchemaSet{}, fmt.Errorf("the overlay's Destructive is implied by the verb %s", verb)
	}
	op.Destructive = destructiveVerbs[verb] || o.Destructive
	paging, owned, err := pagingOf(sop, o)
	if err != nil {
		return ops.Op{}, ops.SchemaSet{}, err
	}
	op.Paging = paging
	args, argsSchema, err := argsOf(sop)
	if err != nil {
		return ops.Op{}, ops.SchemaSet{}, err
	}
	flags, flagsSchema, err := flagsOf(sop, owned)
	if err != nil {
		return ops.Op{}, ops.SchemaSet{}, err
	}
	body, err := bodyOf(sop, o)
	if err != nil {
		return ops.Op{}, ops.SchemaSet{}, err
	}
	op.Args, op.Flags, op.Body = args, flags, body
	if err := uniqueFlags(op); err != nil {
		return ops.Op{}, ops.SchemaSet{}, err
	}
	resp, err := responseOf(sop, o, paging)
	if err != nil {
		return ops.Op{}, ops.SchemaSet{}, err
	}
	schemas := ops.SchemaSet{Command: o.Cmd, Method: sop.Method, Path: sop.Path, Scope: op.Scope,
		Args: argsSchema, Flags: flagsSchema, Response: resp, ResponseTypes: sop.Response.Types}
	if sop.Body != nil {
		schemas.Body = sop.Body.Schema
	}
	return op, schemas, nil
}

// pagingOf checks the overlay's paging against the spec's parameters. It
// returns the parameters the paginator owns, which get no flag.
func pagingOf(sop SpecOp, o overlay.Op) (*api.Paging, map[string]bool, error) {
	hasPage := slices.ContainsFunc(sop.Params, func(p SpecParam) bool { return p.In == "query" && p.Name == "page" })
	paged := o.Page != (overlay.Page{})
	switch {
	case o.NotPaged && (!hasPage || paged):
		return nil, nil, errors.New("the overlay's NotPaged needs a page parameter and no Page")
	case o.NotPaged:
		return nil, nil, nil
	case hasPage && !paged:
		return nil, nil, errors.New("the spec takes a page parameter: set Page, or NotPaged with evidence")
	case !hasPage && paged:
		return nil, nil, errors.New("the overlay sets Page but the spec takes no page parameter")
	case !paged && o.PageSize != 0:
		return nil, nil, errors.New("the overlay sets PageSize without Page")
	case !paged:
		return nil, nil, nil
	case o.Page.End != api.EndEmpty && o.Page.End != api.EndHeader && o.Page.End != api.EndField:
		return nil, nil, errors.New("the overlay's Page.End is not set")
	case (o.Page.End == api.EndField) != (o.Page.More != ""):
		return nil, nil, errors.New("the overlay's Page.More is needed exactly when the list ends by field")
	}
	var size *SpecParam
	for i, p := range sop.Params {
		if p.In != "query" || (p.Name != "per_page" && p.Name != "limit") {
			continue
		}
		if size != nil {
			return nil, nil, errors.New("the spec takes both per_page and limit")
		}
		size = &sop.Params[i]
	}
	if size == nil {
		return nil, nil, errors.New("the spec takes no per_page or limit parameter")
	}
	maxSize, hasMax := size.Schema["maximum"].(float64)
	switch {
	case hasMax && o.PageSize != 0:
		return nil, nil, fmt.Errorf("the overlay's PageSize %d overrides the spec's maximum %v", o.PageSize, maxSize)
	case !hasMax && o.PageSize == 0:
		return nil, nil, fmt.Errorf("the spec gives %s no maximum: set PageSize with evidence", size.Name)
	}
	n := o.PageSize
	if hasMax {
		n = int(maxSize)
	}
	p := &api.Paging{End: o.Page.End, SizeParam: size.Name, Size: n, Items: o.Page.Items, More: o.Page.More}
	return p, map[string]bool{"page": true, size.Name: true}, nil
}

// argsOf returns the path parameters in the order the path names them.
func argsOf(sop SpecOp) ([]ops.Param, map[string]any, error) {
	byName := map[string]SpecParam{}
	for _, p := range sop.Params {
		if p.In == "path" {
			byName[p.Name] = p
		}
	}
	var args []ops.Param
	var obj object
	for _, name := range pathNames(sop.Path) {
		p, ok := byName[name]
		if !ok {
			return nil, nil, fmt.Errorf("the path names {%s} but the spec doesn't define it", name)
		}
		delete(byName, name)
		kind, ok := kindOf(p.Schema)
		if !ok || (kind != ops.Integer && kind != ops.String) {
			return nil, nil, fmt.Errorf("path parameter %s: want an integer or string", name)
		}
		args = append(args, ops.Param{Name: name, Kind: kind, Required: true, Usage: usage(p.Description, p.Schema, false)})
		obj.add(name, described(p), true)
	}
	if len(byName) > 0 {
		return nil, nil, fmt.Errorf("path parameters %v don't appear in the path", slices.Sorted(maps.Keys(byName)))
	}
	return args, obj.schema(), nil
}

func pathNames(path string) []string {
	var names []string
	rest := path
	for {
		_, after, ok := strings.Cut(rest, "{")
		if !ok {
			return names
		}
		name, tail, ok := strings.Cut(after, "}")
		if !ok {
			return names
		}
		names = append(names, name)
		rest = tail
	}
}

func flagsOf(sop SpecOp, owned map[string]bool) ([]ops.Param, map[string]any, error) {
	var flags []ops.Param
	var obj object
	for _, p := range sop.Params {
		if p.In != "query" || owned[p.Name] {
			continue
		}
		kind, ok := kindOf(p.Schema)
		if !ok {
			return nil, nil, fmt.Errorf("query parameter %s: unsupported schema type %v", p.Name, p.Schema["type"])
		}
		flags = append(flags, ops.Param{Name: p.Name, Kind: kind, Required: p.Required, Usage: usage(p.Description, p.Schema, p.Required)})
		obj.add(p.Name, described(p), p.Required)
	}
	slices.SortFunc(flags, func(a, b ops.Param) int { return strings.Compare(a.Name, b.Name) })
	return flags, obj.schema(), nil
}

// bodyOf offers each flat body field as a flag, except credentials, which the
// overlay must list so a new one in the spec can't become a flag unnoticed.
func bodyOf(sop SpecOp, o overlay.Op) (*ops.Body, error) {
	if sop.Body == nil {
		if len(o.Secret) > 0 {
			return nil, errors.New("the overlay sets Secret but the operation has no body")
		}
		return nil, nil
	}
	s := sop.Body.Schema
	props, _ := s["properties"].(map[string]any)
	for _, name := range o.Secret {
		if _, ok := props[name]; !ok {
			return nil, fmt.Errorf("the overlay's Secret names %s, which the body doesn't have", name)
		}
	}
	required := map[string]bool{}
	if list, ok := s["required"].([]any); ok {
		for _, r := range list {
			if name, ok := r.(string); ok {
				required[name] = true
			}
		}
	}
	b := &ops.Body{Required: sop.Body.Required, Object: s["type"] == "object"}
	if len(o.Secret) > 0 {
		b.Secret = slices.Sorted(slices.Values(o.Secret))
	}
	for _, name := range slices.Sorted(maps.Keys(props)) {
		field, _ := props[name].(map[string]any)
		kind, ok := kindOf(field)
		if !ok || kind == ops.StringList || kind == ops.IntegerList || slices.Contains(o.Secret, name) {
			continue
		}
		if kind == ops.String && looksSecret(name) {
			return nil, fmt.Errorf("body field %s looks like a credential: list it in Secret, with evidence", name)
		}
		desc, _ := field["description"].(string)
		b.Fields = append(b.Fields, ops.Param{Name: name, Kind: kind, Required: required[name], Usage: usage(desc, field, required[name])})
	}
	return b, nil
}

func looksSecret(name string) bool {
	for _, w := range []string{"token", "secret", "password", "key"} {
		if strings.Contains(name, w) {
			return true
		}
	}
	return false
}

func uniqueFlags(op ops.Op) error {
	seen := map[string]bool{}
	names := make([]string, 0, len(op.Flags))
	for _, p := range op.Flags {
		names = append(names, p.Name)
	}
	if op.Body != nil {
		for _, p := range op.Body.Fields {
			names = append(names, p.Name)
		}
	}
	for _, name := range names {
		flag := ops.FlagName(name)
		if seen[flag] {
			return fmt.Errorf("flag --%s is used twice", flag)
		}
		seen[flag] = true
	}
	return nil
}

// responseOf is the schema of what the command prints; for a paged list, that
// is the array of every page's items.
func responseOf(sop SpecOp, o overlay.Op, paging *api.Paging) (map[string]any, error) {
	s := sop.Response.Schema
	if o.ResponseArray {
		if s == nil || s["type"] != "object" {
			return nil, errors.New("the overlay sets ResponseArray but the spec's response isn't one object")
		}
		s = map[string]any{"type": "array", "items": s}
	}
	if paging == nil {
		return s, nil
	}
	if s == nil {
		return nil, errors.New("a paged operation has no JSON response")
	}
	props, _ := s["properties"].(map[string]any)
	if paging.More != "" {
		if more, _ := props[paging.More].(map[string]any); more["type"] != "boolean" {
			return nil, fmt.Errorf("the overlay's Page.More %s is not a boolean field of the response", paging.More)
		}
	}
	if paging.Items != "" {
		field, _ := props[paging.Items].(map[string]any)
		if field["type"] != "array" {
			return nil, fmt.Errorf("the overlay's Page.Items %s is not an array field of the response", paging.Items)
		}
		s = field
	}
	if s["type"] != "array" {
		return nil, errors.New("the paged response isn't an array: set ResponseArray or Page.Items, with evidence")
	}
	return s, nil
}

// kindOf maps a schema to a flag kind. A type list holding "null" counts as
// its other type, since a flag can't send null anyway.
func kindOf(s map[string]any) (ops.Kind, bool) {
	switch nonNullType(s) {
	case "string":
		return ops.String, true
	case "integer":
		return ops.Integer, true
	case "boolean":
		return ops.Boolean, true
	case "array":
		items, _ := s["items"].(map[string]any)
		switch nonNullType(items) {
		case "string":
			return ops.StringList, true
		case "integer":
			return ops.IntegerList, true
		}
	}
	return "", false
}

func nonNullType(s map[string]any) string {
	switch t := s["type"].(type) {
	case string:
		return t
	case []any:
		var others []string
		for _, x := range t {
			if name, ok := x.(string); ok && name != "null" {
				others = append(others, name)
			}
		}
		if len(others) == 1 {
			return others[0]
		}
	}
	return ""
}

// usage is a flag's help line: the first line of its description, then its
// allowed values.
func usage(desc string, s map[string]any, required bool) string {
	line, _, _ := strings.Cut(strings.TrimSpace(desc), "\n")
	var notes []string
	enum, ok := s["enum"].([]any)
	if items, isList := s["items"].(map[string]any); !ok && isList {
		enum, ok = items["enum"].([]any)
	}
	if ok {
		var values []string
		for _, v := range enum {
			if v != nil {
				values = append(values, fmt.Sprint(v))
			}
		}
		notes = append(notes, "one of: "+strings.Join(values, ", "))
	}
	if required {
		notes = append(notes, "required")
	}
	if len(notes) > 0 {
		line = strings.TrimSpace(line + " (" + strings.Join(notes, "; ") + ")")
	}
	return line
}

// described copies a parameter's schema with its description, because
// shared parameters share one schema map.
func described(p SpecParam) map[string]any {
	s := maps.Clone(p.Schema)
	if _, ok := s["description"]; !ok && p.Description != "" {
		s["description"] = p.Description
	}
	return s
}

// object builds the object schema of an operation's arguments or flags.
type object struct {
	props    map[string]any
	required []string
}

func (o *object) add(name string, s map[string]any, required bool) {
	if o.props == nil {
		o.props = map[string]any{}
	}
	o.props[name] = s
	if required {
		o.required = append(o.required, name)
	}
}

func (o *object) schema() map[string]any {
	s := map[string]any{"type": "object", "properties": map[string]any{}}
	if o.props != nil {
		s["properties"] = o.props
	}
	if len(o.required) > 0 {
		slices.Sort(o.required)
		s["required"] = o.required
	}
	return s
}
