package overlay

import (
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
	"testing"
)

// verbs is the closed verb set: the handover's list, minus sync (no operation
// uses it), plus the four the spec needs (design note, approved 2026-09-29).
var verbs = map[string]bool{
	"list": true, "get": true, "create": true, "update": true, "delete": true,
	"ignore": true, "unignore": true, "snooze": true, "unsnooze": true, "solve": true,
	"scan": true, "activate": true, "deactivate": true, "export": true,
	"clone": true, "import": true, "rotate": true, "review": true,
}

var resourceName = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

func TestCommandsAreResourceThenVerb(t *testing.T) {
	for id, op := range Ops {
		resource, verb, ok := strings.Cut(op.Cmd, " ")
		if !ok || !resourceName.MatchString(resource) {
			t.Errorf("%s: %q is not <kebab-case resource> <verb>", id, op.Cmd)
			continue
		}
		if !verbs[verb] {
			t.Errorf("%s: verb %q is outside the closed verb set", id, verb)
		}
	}
}

func TestCommandsAreUnique(t *testing.T) {
	seen := map[string]string{}
	for id, op := range Ops {
		if other, dup := seen[op.Cmd]; dup {
			t.Errorf("%s and %s are both %q", other, id, op.Cmd)
		}
		seen[op.Cmd] = id
	}
}

// needsEvidence are the fields that correct or add to the spec, or mark an
// operation destructive beyond its verb. Each needs a reason a reviewer can check.
var needsEvidence = map[string]bool{
	"Destructive": true, "Scope": true, "ResponseArray": true, "PageSize": true,
	"NotPaged": true, "Secret": true, "Help": true, "BadRequestHint": true,
}

func TestCorrectionsCiteEvidence(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "overlay.go", nil, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	commentEnds := map[int]bool{}
	for _, cg := range f.Comments {
		commentEnds[fset.Position(cg.End()).Line] = true
	}
	lit := opsLiteral(t, f)
	for _, elt := range lit.Elts {
		entry, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			t.Fatalf("Ops element at %s is not key: value", fset.Position(elt.Pos()))
		}
		value, ok := entry.Value.(*ast.CompositeLit)
		if !ok {
			t.Fatalf("Ops value at %s is not a literal", fset.Position(entry.Pos()))
		}
		for _, field := range value.Elts {
			kv, ok := field.(*ast.KeyValueExpr)
			if !ok {
				t.Fatalf("Op field at %s is not key: value", fset.Position(field.Pos()))
			}
			name, _ := kv.Key.(*ast.Ident)
			if name != nil && needsEvidence[name.Name] && !commentEnds[fset.Position(entry.Pos()).Line-1] {
				t.Errorf("line %d sets %s without an evidence comment on the line above", fset.Position(entry.Pos()).Line, name.Name)
			}
		}
	}
}

func opsLiteral(t *testing.T, f *ast.File) *ast.CompositeLit {
	t.Helper()
	for _, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range gen.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok || len(vs.Names) != 1 || vs.Names[0].Name != "Ops" || len(vs.Values) != 1 {
				continue
			}
			if lit, ok := vs.Values[0].(*ast.CompositeLit); ok {
				return lit
			}
		}
	}
	t.Fatal("no var Ops = map[string]Op{...} in overlay.go")
	return nil
}
