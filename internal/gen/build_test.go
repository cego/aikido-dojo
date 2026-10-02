package gen

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/cego/aikido-dojo/internal/api"
	"github.com/cego/aikido-dojo/internal/ops"
	"github.com/cego/aikido-dojo/internal/overlay"
)

const fixtureSpec = `{"paths": {
	"/things": {
		"get": {"operationId": "listThings", "summary": "List things", "security": [{"oauth": ["things:read"]}],
			"parameters": [
				{"in": "query", "name": "page", "schema": {"type": "integer"}},
				{"in": "query", "name": "per_page", "schema": {"type": "integer", "maximum": 50}},
				{"in": "query", "name": "filter_status", "description": "Status to keep.\nMore text.", "schema": {"type": "string", "enum": ["open", "closed"]}},
				{"in": "query", "name": "ids", "required": true, "schema": {"type": "array", "items": {"type": "integer"}}}],
			"responses": {"200": {"content": {"application/json": {"schema": {"type": "array", "items": {"type": "object"}}}}}}},
		"post": {"operationId": "addThing", "summary": "Add a thing", "security": [{"oauth": ["things:write"]}],
			"requestBody": {"required": true, "content": {"application/json": {"schema": {"type": "object", "required": ["name"],
				"properties": {"name": {"type": "string", "description": "Its name."}, "api_token": {"type": "string"},
					"size": {"type": ["integer", "null"]}, "tags": {"type": "array", "items": {"type": "string"}}}}}}},
			"responses": {"201": {"content": {"application/json": {"schema": {"type": "object"}}}}}}},
	"/things/{thing_id}": {
		"delete": {"operationId": "deleteThing", "security": [{"oauth": ["things:write"]}],
			"parameters": [{"in": "path", "name": "thing_id", "required": true, "schema": {"type": "integer"}}],
			"responses": {"204": {}}}},
	"/things/export": {
		"get": {"operationId": "exportThings",
			"parameters": [{"in": "query", "name": "page", "schema": {"type": "integer"}}, {"in": "query", "name": "per_page", "schema": {"type": "integer"}}],
			"responses": {"200": {"content": {"text/csv": {}}}}}}
}}`

func fixtureOverlay() map[string]overlay.Op {
	return map[string]overlay.Op{
		"listThings":   {Cmd: "thing list", Page: overlay.UntilEmpty},
		"addThing":     {Cmd: "thing create", Secret: []string{"api_token"}},
		"deleteThing":  {Cmd: "thing delete"},
		"exportThings": {Cmd: "thing export", NotPaged: true},
	}
}

func buildFixture(t *testing.T, mutate func(map[string]overlay.Op)) (*Model, error) {
	t.Helper()
	s, err := LoadSpec([]byte(fixtureSpec))
	if err != nil {
		t.Fatal(err)
	}
	o := fixtureOverlay()
	if mutate != nil {
		mutate(o)
	}
	return Build(s, o)
}

func TestBuild(t *testing.T) {
	m, err := buildFixture(t, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := []ops.Op{
		{ID: "addThing", Command: "thing create", Method: "POST", Path: "/things", Summary: "Add a thing", Scope: "things:write",
			Body: &ops.Body{Required: true, Object: true, Secret: []string{"api_token"}, Fields: []ops.Param{
				{Name: "name", Kind: ops.String, Required: true, Usage: "Its name. (required)"},
				{Name: "size", Kind: ops.Integer},
			}}},
		{ID: "deleteThing", Command: "thing delete", Method: "DELETE", Path: "/things/{thing_id}", Scope: "things:write", Destructive: true,
			Args: []ops.Param{{Name: "thing_id", Kind: ops.Integer, Required: true}}},
		{ID: "exportThings", Command: "thing export", Method: "GET", Path: "/things/export",
			Flags: []ops.Param{{Name: "page", Kind: ops.Integer}, {Name: "per_page", Kind: ops.Integer}}},
		{ID: "listThings", Command: "thing list", Method: "GET", Path: "/things", Summary: "List things", Scope: "things:read",
			Flags: []ops.Param{
				{Name: "filter_status", Kind: ops.String, Usage: "Status to keep. (one of: open, closed)"},
				{Name: "ids", Kind: ops.IntegerList, Required: true, Usage: "(required)"},
			},
			Paging: &api.Paging{End: api.EndEmpty, SizeParam: "per_page", Size: 50}},
	}
	if !reflect.DeepEqual(m.Ops, want) {
		t.Errorf("ops =\n%+v\nwant\n%+v", m.Ops, want)
	}
	list := m.Schemas["listThings"]
	if list.Command != "thing list" || list.Response["type"] != "array" {
		t.Errorf("listThings schemas = %+v, want the command and an array response", list)
	}
	flags, _ := list.Flags["properties"].(map[string]any)
	if _, owned := flags["per_page"]; owned || flags["filter_status"] == nil || flags["ids"] == nil {
		t.Errorf("flag schema properties = %v, want filter_status and ids but not the paginator's per_page", flags)
	}
	if list.Flags["required"] == nil {
		t.Errorf("flag schema = %v, want ids required", list.Flags)
	}
	if body := m.Schemas["addThing"].Body; body["type"] != "object" {
		t.Errorf("addThing body schema = %v, want the spec's", body)
	}
}

func TestBuildRejects(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(map[string]overlay.Op)
		want   string
	}{
		{"an operation without an overlay entry", func(o map[string]overlay.Op) { delete(o, "addThing") }, "addThing (POST /things) has no overlay entry"},
		{"an overlay entry without an operation", func(o map[string]overlay.Op) { o["ghost"] = overlay.Op{Cmd: "ghost get"} }, "overlay entry ghost has no operation"},
		{"a command that isn't resource and verb", func(o map[string]overlay.Op) { o["deleteThing"] = overlay.Op{Cmd: "thing"} }, "is not <resource> <verb>"},
		{"a scope correction that repeats the spec", func(o map[string]overlay.Op) {
			o["deleteThing"] = overlay.Op{Cmd: "thing delete", Scope: "things:write"}
		}, "repeats the spec"},
		{"Destructive on a delete", func(o map[string]overlay.Op) { o["deleteThing"] = overlay.Op{Cmd: "thing delete", Destructive: true} }, "implied by the verb"},
		{"a page parameter without Page", func(o map[string]overlay.Op) { o["listThings"] = overlay.Op{Cmd: "thing list"} }, "set Page, or NotPaged"},
		{"Page without a page parameter", func(o map[string]overlay.Op) {
			o["deleteThing"] = overlay.Op{Cmd: "thing delete", Page: overlay.UntilEmpty}
		}, "takes no page parameter"},
		{"NotPaged with Page", func(o map[string]overlay.Op) {
			o["exportThings"] = overlay.Op{Cmd: "thing export", NotPaged: true, Page: overlay.UntilEmpty}
		}, "NotPaged needs"},
		{"PageSize over a spec maximum", func(o map[string]overlay.Op) {
			o["listThings"] = overlay.Op{Cmd: "thing list", Page: overlay.UntilEmpty, PageSize: 10}
		}, "overrides the spec's maximum 50"},
		{"no page size anywhere", func(o map[string]overlay.Op) {
			o["exportThings"] = overlay.Op{Cmd: "thing export", Page: overlay.UntilEmpty}
		}, "no maximum"},
		{"PageSize without Page", func(o map[string]overlay.Op) { o["deleteThing"] = overlay.Op{Cmd: "thing delete", PageSize: 5} }, "PageSize without Page"},
		{"More without EndField", func(o map[string]overlay.Op) {
			o["listThings"] = overlay.Op{Cmd: "thing list", Page: overlay.Page{End: api.EndEmpty, More: "hasMore"}}
		}, "More is needed exactly"},
		{"Page without End", func(o map[string]overlay.Op) {
			o["listThings"] = overlay.Op{Cmd: "thing list", Page: overlay.Page{Items: "x"}}
		}, "Page.End is not set"},
		{"Items that aren't an array field", func(o map[string]overlay.Op) {
			o["listThings"] = overlay.Op{Cmd: "thing list", Page: overlay.Page{End: api.EndEmpty, Items: "things"}}
		}, "Page.Items things"},
		{"ResponseArray on an array", func(o map[string]overlay.Op) {
			o["listThings"] = overlay.Op{Cmd: "thing list", Page: overlay.UntilEmpty, ResponseArray: true}
		}, "isn't one object"},
		{"a Secret the body doesn't have", func(o map[string]overlay.Op) {
			o["addThing"] = overlay.Op{Cmd: "thing create", Secret: []string{"api_token", "nope"}}
		}, "nope, which the body doesn't have"},
		{"a credential-looking field not in Secret", func(o map[string]overlay.Op) { o["addThing"] = overlay.Op{Cmd: "thing create"} }, "api_token looks like a credential"},
		{"Secret without a body", func(o map[string]overlay.Op) {
			o["deleteThing"] = overlay.Op{Cmd: "thing delete", Secret: []string{"x"}}
		}, "has no body"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := buildFixture(t, tt.mutate)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("err = %v, want it to contain %q", err, tt.want)
			}
		})
	}
}

func TestBuildRejectsSpecShapes(t *testing.T) {
	const page = `{"in": "query", "name": "page", "schema": {"type": "integer"}}, {"in": "query", "name": "per_page", "schema": {"type": "integer", "maximum": 5}}`
	tests := []struct {
		name, paths string
		op          overlay.Op
		want        string
	}{
		{"a path placeholder the spec doesn't define", `{"/a/{x}": {"get": {"operationId": "a"}}}`, overlay.Op{Cmd: "a get"}, "names {x}"},
		{"a path parameter missing from the path", `{"/a": {"get": {"operationId": "a", "parameters": [{"in": "path", "name": "x", "required": true, "schema": {"type": "integer"}}]}}}`,
			overlay.Op{Cmd: "a get"}, "don't appear in the path"},
		{"a boolean path parameter", `{"/a/{x}": {"get": {"operationId": "a", "parameters": [{"in": "path", "name": "x", "required": true, "schema": {"type": "boolean"}}]}}}`,
			overlay.Op{Cmd: "a get"}, "want an integer or string"},
		{"an object query parameter", `{"/a": {"get": {"operationId": "a", "parameters": [{"in": "query", "name": "q", "schema": {"type": "object"}}]}}}`,
			overlay.Op{Cmd: "a get"}, "query parameter q"},
		{"two names with one flag", `{"/a": {"post": {"operationId": "a", "parameters": [{"in": "query", "name": "a_b", "schema": {"type": "string"}}],
			"requestBody": {"content": {"application/json": {"schema": {"type": "object", "properties": {"a_b": {"type": "string"}}}}}}}}}`,
			overlay.Op{Cmd: "a create"}, "--a-b is used twice"},
		{"a paged response that isn't an array", `{"/a": {"get": {"operationId": "a", "parameters": [` + page + `],
			"responses": {"200": {"content": {"application/json": {"schema": {"type": "object"}}}}}}}}`,
			overlay.Op{Cmd: "a list", Page: overlay.UntilEmpty}, "isn't an array"},
		{"a camelCase credential field", `{"/a": {"post": {"operationId": "a",
			"requestBody": {"content": {"application/json": {"schema": {"type": "object", "properties": {"apiKey": {"type": "string"}}}}}}}}}`,
			overlay.Op{Cmd: "a create"}, "apiKey looks like a credential"},
		{"a credentials field", `{"/a": {"post": {"operationId": "a",
			"requestBody": {"content": {"application/json": {"schema": {"type": "object", "properties": {"credentials": {"type": "string"}}}}}}}}}`,
			overlay.Op{Cmd: "a create"}, "credentials looks like a credential"},
		{"a More field that isn't a boolean", `{"/a": {"get": {"operationId": "a", "parameters": [` + page + `],
			"responses": {"200": {"content": {"application/json": {"schema": {"type": "object", "properties": {"items": {"type": "array"}}}}}}}}}}`,
			overlay.Op{Cmd: "a list", Page: overlay.Page{End: api.EndField, Items: "items", More: "hasMore"}}, "More hasMore is not a boolean field"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := LoadSpec(specWith(tt.paths))
			if err != nil {
				t.Fatal(err)
			}
			_, err = Build(s, map[string]overlay.Op{"a": tt.op})
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("err = %v, want it to contain %q", err, tt.want)
			}
		})
	}
}

func TestBuildDropsBoundsTheOverlayMarksWrong(t *testing.T) {
	s, err := LoadSpec(specWith(`{"/a/{x}": {"get": {"operationId": "a", "parameters": [
		{"in": "path", "name": "x", "required": true, "schema": {"type": "integer", "minimum": 0, "maximum": 1}},
		{"in": "query", "name": "q", "schema": {"type": "integer", "minimum": 0, "maximum": 1}}]}}}`))
	if err != nil {
		t.Fatal(err)
	}
	m, err := Build(s, map[string]overlay.Op{"a": {Cmd: "a get", Unbounded: []string{"x"}}})
	if err != nil {
		t.Fatal(err)
	}
	x := m.Schemas["a"].Args["properties"].(map[string]any)["x"].(map[string]any)
	q := m.Schemas["a"].Flags["properties"].(map[string]any)["q"].(map[string]any)
	if _, ok := x["maximum"]; ok || x["minimum"] != nil || q["maximum"] == nil {
		t.Errorf("x = %v, q = %v; want x's bounds dropped and q's kept", x, q)
	}
	for _, tt := range []struct {
		unbounded []string
		want      string
	}{
		{[]string{"nope"}, "Unbounded names nope, which is no parameter"},
	} {
		_, err := Build(s, map[string]overlay.Op{"a": {Cmd: "a get", Unbounded: tt.unbounded}})
		if err == nil || !strings.Contains(err.Error(), tt.want) {
			t.Errorf("Unbounded %v: err = %v, want %q", tt.unbounded, err, tt.want)
		}
	}
	s2, err := LoadSpec(specWith(`{"/a/{x}": {"get": {"operationId": "a", "parameters": [{"in": "path", "name": "x", "required": true, "schema": {"type": "integer"}}]}}}`))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Build(s2, map[string]overlay.Op{"a": {Cmd: "a get", Unbounded: []string{"x"}}}); err == nil || !strings.Contains(err.Error(), "x, which has no bounds") {
		t.Errorf("err = %v, want Unbounded on an unbounded parameter refused", err)
	}
}

func TestBuildVendoredSpec(t *testing.T) {
	data, err := os.ReadFile("../../spec/openapi.json")
	if err != nil {
		t.Fatal(err)
	}
	s, err := LoadSpec(data)
	if err != nil {
		t.Fatal(err)
	}
	m, err := Build(s, overlay.Ops)
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Ops) != 201 || len(m.Schemas) != 201 {
		t.Fatalf("ops = %d, schemas = %d, want 201 each", len(m.Ops), len(m.Schemas))
	}
	byCmd := map[string]ops.Op{}
	for _, op := range m.Ops {
		byCmd[op.Command] = op
	}
	checks := []struct {
		cmd  string
		ok   func(ops.Op) bool
		want string
	}{
		{"repo list", func(op ops.Op) bool {
			return reflect.DeepEqual(op.Paging, &api.Paging{End: api.EndEmpty, SizeParam: "per_page", Size: 200})
		}, "per_page pages of 200 until empty"},
		{"issue-group list", func(op ops.Op) bool { return op.Paging != nil && op.Paging.End == api.EndHeader }, "header paging"},
		{"cloud-asset list", func(op ops.Op) bool {
			return reflect.DeepEqual(op.Paging, &api.Paging{End: api.EndField, SizeParam: "limit", Size: 100, Items: "assets", More: "hasMore"})
		}, "limit pages of 100 in assets, ending by hasMore"},
		{"code-quality-finding list", func(op ops.Op) bool {
			return op.Paging != nil && op.Paging.End == api.EndHeader && op.Paging.Items == "findings"
		}, "findings ending by header"},
		{"issue export", func(op ops.Op) bool { return op.Paging == nil && len(op.Flags) > 0 }, "no paging, page flags kept"},
		{"repo delete", func(op ops.Op) bool { return op.Destructive }, "destructive"},
		{"pr-check-config-all update", func(op ops.Op) bool { return op.Destructive }, "destructive from the overlay"},
		{"container-connectivity update", func(op ops.Op) bool { return op.Scope == "containers:write" }, "the corrected scope"},
		{"code-scanning-token update", func(op ops.Op) bool {
			return op.Body != nil && len(op.Body.Fields) == 0 && reflect.DeepEqual(op.Body.Secret, []string{"access_token"})
		}, "the token only through --body"},
	}
	for _, c := range checks {
		if op, found := byCmd[c.cmd]; !found || !c.ok(op) {
			t.Errorf("%s = %+v, want %s", c.cmd, op, c.want)
		}
	}
	if r := m.Schemas["listCodeRepos"].Response; r["type"] != "array" {
		t.Errorf("repo list response = %v, want the corrected array", r)
	}
}
