package cli

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestOutputIndentsOurJSONOnATerminal(t *testing.T) {
	env, _ := testEnv(t)
	env.StdoutTTY = true
	if out := mustRun(t, env, "version"); !strings.HasPrefix(out, "{\n  \"version\": ") {
		t.Errorf("stdout = %q, want indented JSON", out)
	}
}

func TestOutputPassesAikidoJSONThroughOffATerminal(t *testing.T) {
	_, env, _ := newFake(t, respond(`{"a":"x\/y"}`))
	if out := mustRun(t, env, "workspace", "get"); out != `{"a":"x\/y"}` {
		t.Errorf("stdout = %q, want Aikido's bytes unchanged", out)
	}
}

func TestOutputIndentsAikidoJSONOnATerminal(t *testing.T) {
	_, env, _ := newFake(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"a":"x\/y"}`)
	})
	env.StdoutTTY = true
	if out := mustRun(t, env, "workspace", "get"); out != "{\n  \"a\": \"x/y\"\n}\n" {
		t.Errorf("stdout = %q, want it indented", out)
	}
}

func TestOutputCopiesCSVOnATerminal(t *testing.T) {
	_, env, _ := newFake(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		fmt.Fprint(w, "a,b\n1,2\n")
	})
	env.StdoutTTY = true
	if out := mustRun(t, env, "issue", "export", "--format", "csv"); out != "a,b\n1,2\n" {
		t.Errorf("stdout = %q, want the CSV unchanged", out)
	}
}

func TestListIndentsOnATerminal(t *testing.T) {
	_, env, _ := newFake(t, pages(`[{"id":1}]`, `[]`))
	env.StdoutTTY = true
	if out := mustRun(t, env, "repo", "list"); out != "[\n  {\n    \"id\": 1\n  }\n]\n" {
		t.Errorf("stdout = %q, want an indented array", out)
	}
}
