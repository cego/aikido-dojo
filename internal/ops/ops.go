// Package ops describes generated commands: what each one calls, its
// arguments and flags, and how it pages. The generator fills these in from the
// spec and the overlay; the CLI builds its commands from them.
package ops

import (
	"strings"

	"github.com/cego/aikido-dojo/internal/api"
)

type Op struct {
	ID             string // the spec's operationId
	Command        string // "<resource> <verb>"
	Method         string
	Path           string // with {name} placeholders for Args
	Summary        string
	Description    string
	Help           string // the overlay's addition to the help
	Scope          string
	Args           []Param // path parameters, in path order
	Flags          []Param // query parameters; the paginator owns page and page size
	Body           *Body
	Paging         *api.Paging
	Destructive    bool
	BadRequestHint string
}

type Param struct {
	Name     string // the spec's name; FlagName gives the flag
	Kind     Kind
	Required bool
	Usage    string
}

type Body struct {
	Required bool
	Object   bool     // a JSON object: flags can set its fields, and an empty body is {}
	Fields   []Param  // flat fields offered as flags
	Secret   []string // credential fields, set only through --body or --body-file
}

type Kind string

const (
	String      Kind = "string"
	Integer     Kind = "integer"
	Boolean     Kind = "boolean"
	StringList  Kind = "[]string"
	IntegerList Kind = "[]integer"
)

// SchemaSet is one operation's entry in schemas.json: the JSON Schemas of its
// arguments, flags and body, and of what the command prints.
type SchemaSet struct {
	Command       string         `json:"command"`
	Method        string         `json:"method"`
	Path          string         `json:"path"`
	Scope         string         `json:"scope,omitempty"`
	Args          map[string]any `json:"args"`
	Flags         map[string]any `json:"flags"`
	Body          map[string]any `json:"body,omitempty"`
	Response      map[string]any `json:"response,omitempty"`
	ResponseTypes []string       `json:"response_types,omitempty"`
}

// FlagName is the flag for a spec name: filter_code_repo_id is --filter-code-repo-id.
func FlagName(name string) string { return strings.ReplaceAll(name, "_", "-") }
