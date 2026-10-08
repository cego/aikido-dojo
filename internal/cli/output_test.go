package cli

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/cego/aikido-dojo/internal/catalog"
	"github.com/cego/aikido-dojo/internal/clierr"
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

func TestJQFiltersTheOutput(t *testing.T) {
	two := pages(`[{"name":"a","id":1},{"name":"b","id":2}]`, `[]`)
	for _, tt := range []struct {
		filter string
		tty    bool
		want   string
	}{
		{".[0].name", false, "a\n"},
		{".[] | .name", false, "a\nb\n"},
		{"length", false, "2\n"},
		{".[0]", false, `{"id":1,"name":"a"}` + "\n"},
		{".[0]", true, "{\n  \"id\": 1,\n  \"name\": \"a\"\n}\n"},
	} {
		_, env, _ := newFake(t, two)
		env.StdoutTTY = tt.tty
		if out := mustRun(t, env, "repo", "list", "--jq", tt.filter); out != tt.want {
			t.Errorf("--jq %q (tty %v) = %q, want %q", tt.filter, tt.tty, out, tt.want)
		}
	}
}

func TestJQOnOurJSONAndOnAPIResponses(t *testing.T) {
	env, _ := testEnv(t)
	if out := mustRun(t, env, "version", "--jq", ".spec.updated_at"); out != catalog.SpecUpdatedAt+"\n" {
		t.Errorf("version --jq = %q", out)
	}
	_, env2, _ := newFake(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"name":"ws","id":12345678901234567890}`)
	})
	if out := mustRun(t, env2, "workspace", "get", "--jq", ".id"); out != "12345678901234567890\n" {
		t.Errorf("workspace get --jq .id = %q, want the integer exactly", out)
	}
}

func TestJQMistakes(t *testing.T) {
	f, env, _ := newFake(t, respond(`[1]`))
	_, stderr, code := run(t, env, "repo", "list", "--jq", ".[")
	if e := errorOf(t, stderr); code != 2 || e.Code != "invalid_jq" || len(f.seen())+len(f.seenLogins()) != 0 {
		t.Errorf("exit %d, error %+v; want invalid_jq before any call", code, e)
	}
	_, stderr, code = run(t, env, "workspace", "get", "--jq", ".a")
	if e := errorOf(t, stderr); code != 2 || e.Code != "jq_failed" {
		t.Errorf("exit %d, error %+v; want jq_failed for indexing an array", code, e)
	}
}

// --jq and --ndjson need JSON; a command whose schema says it prints
// something else is refused before the call.
func TestOutputMistakesCostNoCall(t *testing.T) {
	for _, args := range [][]string{
		{"report", "export", "--jq", "."},
	} {
		f, env, _ := newFake(t, respond(`{}`))
		_, stderr, code := run(t, env, args...)
		if e := errorOf(t, stderr); code != 2 || e.Hint == "" || len(f.seen())+len(f.seenLogins()) != 0 {
			t.Errorf("%q: exit %d, error %+v; want a usage error before any call", args, code, e)
		}
	}
}

func TestNDJSONStreamsAList(t *testing.T) {
	_, env, _ := newFake(t, pages(`[{"id":1},{"id":2}]`, `[{"id":3},{"id":4}]`, `[]`))
	if out := mustRun(t, env, "repo", "list", "--ndjson", "--limit", "3"); out != "{\"id\":1}\n{\"id\":2}\n{\"id\":3}\n" {
		t.Errorf("stdout = %q, want three lines", out)
	}
	// A failure part-way leaves the lines already printed; the error says the rest is missing.
	_, env2, _ := newFake(t, pages(`[1,2]`))
	stdout, stderr, code := run(t, env2, "repo", "list", "--ndjson")
	if stdout != "1\n2\n" || code != 5 || errorOf(t, stderr).Code != "not_found" {
		t.Errorf("stdout %q, exit %d, stderr %s; want the first page's lines and the error", stdout, code, stderr)
	}
}

func TestNDJSONStreamsAnArrayResponse(t *testing.T) {
	_, env, _ := newFake(t, respond(`[{"a":1},{"a":"x\/y"}]`))
	env.StdoutTTY = true
	// Each line is Aikido's own bytes, as everything is off a terminal.
	if out := mustRun(t, env, "issue", "export", "--ndjson"); out != "{\"a\":1}\n{\"a\":\"x\\/y\"}\n" {
		t.Errorf("stdout = %q, want one compact line per item, even on a terminal", out)
	}
}

func TestNDJSONOnOurArrays(t *testing.T) {
	env, _ := testEnv(t)
	out := mustRun(t, env, "search", "list repositories", "--ndjson")
	if lines := strings.Split(strings.TrimSuffix(out, "\n"), "\n"); len(lines) != 10 || !strings.HasPrefix(lines[0], `{"command":`) {
		t.Errorf("stdout = %q, want 10 lines of hits", out)
	}
}

func TestNDJSONWithJQ(t *testing.T) {
	_, env, _ := newFake(t, pages(`[{"name":"a"},{"name":"b"}]`, `[]`))
	if out := mustRun(t, env, "repo", "list", "--ndjson", "--jq", ".name"); out != "a\nb\n" {
		t.Errorf("stdout = %q, want one name per line", out)
	}
}

func TestNDJSONNeedsAnArray(t *testing.T) {
	for _, args := range [][]string{
		{"workspace", "get", "--ndjson"},
		{"version", "--ndjson"},
	} {
		f, env, _ := newFake(t, respond(`{}`))
		_, stderr, code := run(t, env, args...)
		if e := errorOf(t, stderr); code != 2 || e.Hint == "" || len(f.seen())+len(f.seenLogins()) != 0 {
			t.Errorf("%q: exit %d, error %+v; want a usage error before any call", args, code, e)
		}
	}
}

// Once a write has gone through, an output mistake must not read as "nothing
// happened": an agent would repeat the write to see its result.
func TestOutputFailureAfterAWriteSaysItTookEffect(t *testing.T) {
	f, env, _ := newFake(t, respond(`{"id":7}`))
	stdout, stderr, code := run(t, env, "team", "create", "--name", "x", "--jq", ".[0]")
	e := errorOf(t, stderr)
	if code != clierr.ExitUnexpected || e.Code != "output_failed" || !strings.Contains(e.Message, "succeeded") || !strings.Contains(e.Hint, "don't repeat") {
		t.Errorf("exit %d, error %+v; want output_failed saying the write took effect", code, e)
	}
	if len(f.seen()) != 1 || !strings.Contains(stdout, `{"id":7}`) {
		t.Errorf("%d calls, stdout %q; want one call and its response kept on stdout", len(f.seen()), stdout)
	}
}

// api checks the output flags against the command its path matches, before the call, as generated commands do.
func TestAPIChecksOutputFlagsBeforeTheCall(t *testing.T) {
	for _, args := range [][]string{
		{"api", "GET", "/workspace", "--ndjson"},
		{"api", "DELETE", "/teams/1", "--yes", "--jq", "."},
	} {
		f, env, _ := newFake(t, respond(`{}`))
		_, stderr, code := run(t, env, args...)
		if code != clierr.ExitUsage || len(f.seen())+len(f.seenLogins()) != 0 {
			t.Errorf("%q: exit %d, %d calls, stderr %s; want a usage error before any call", args, code, len(f.seen()), stderr)
		}
	}
}

func TestOutputMistakeMessageNamesNoJSON(t *testing.T) {
	env, _ := testEnv(t)
	_, stderr, _ := run(t, env, "team", "delete", "1", "--yes", "--jq", ".")
	if e := errorOf(t, stderr); !strings.HasSuffix(e.Message, "team delete prints no JSON") {
		t.Errorf("message = %q", e.Message)
	}
}

// A dry run prints JSON of its own, whatever the call would return.
func TestDryRunTakesJQ(t *testing.T) {
	_, env, _ := newFake(t, respond(`{}`))
	if out := mustRun(t, env, "team", "delete", "1", "--dry-run", "--jq", ".method"); out != "DELETE\n" {
		t.Errorf("stdout = %q", out)
	}
}

// Aikido's data includes names from scanned repos; on a terminal a raw string
// mustn't carry escape sequences to it. A program off a terminal gets it exact.
func TestJQStringsCarryNoControlCharactersToATerminal(t *testing.T) {
	for _, tty := range []bool{true, false} {
		_, env, _ := newFake(t, respond(`{"name":"a\u001b[31mb\tc"}`))
		env.StdoutTTY = tty
		out := mustRun(t, env, "workspace", "get", "--jq", ".name")
		want := "a\x1b[31mb\tc\n"
		if tty {
			want = "a\\u001b[31mb\tc\n"
		}
		if out != want {
			t.Errorf("tty %v: stdout = %q, want %q", tty, out, want)
		}
	}
}

// A response that breaks off is a failed read, not a usage mistake.
func TestNDJSONReportsABrokenResponseAsARead(t *testing.T) {
	_, env, _ := newFake(t, respond(``))
	_, stderr, code := run(t, env, "issue", "export", "--ndjson")
	if e := errorOf(t, stderr); code != clierr.ExitUnexpected || e.Code == "invalid_input" {
		t.Errorf("exit %d, error %+v; want exit 1 for a response that ended", code, e)
	}
}

// A command that can print JSON or CSV is checked after the call too.
func TestJQOnACSVResponse(t *testing.T) {
	_, env, _ := newFake(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/csv;charset=UTF-8")
		fmt.Fprint(w, "a,b\n")
	})
	_, stderr, code := run(t, env, "issue", "export", "--format", "csv", "--jq", ".")
	if e := errorOf(t, stderr); code != clierr.ExitUsage || !strings.Contains(e.Message, "text/csv") {
		t.Errorf("exit %d, error %+v; want a usage error naming text/csv", code, e)
	}
}

// On a path no command matches, api applies --jq and --ndjson to whatever comes back.
func TestAPIOutputFlagsOnAnUnknownPath(t *testing.T) {
	_, env, _ := newFake(t, respond(`[{"a":1},{"a":2}]`))
	if out := mustRun(t, env, "api", "GET", "/nosuch", "--jq", ".[1].a"); out != "2\n" {
		t.Errorf("--jq: stdout = %q", out)
	}
	if out := mustRun(t, env, "api", "GET", "/nosuch", "--ndjson"); out != "{\"a\":1}\n{\"a\":2}\n" {
		t.Errorf("--ndjson: stdout = %q", out)
	}
}

// After a filter fails part-way, stdout holds Aikido's response alone, not
// some results followed by it.
func TestWriteResultReprintsOnlyTheRawResponse(t *testing.T) {
	_, env, _ := newFake(t, respond(`{"id":7}`))
	stdout, stderr, code := run(t, env, "team", "create", "--name", "x", "--jq", ".id, .[0]")
	if stdout != "{\"id\":7}\n" || code != clierr.ExitUnexpected || errorOf(t, stderr).Code != "output_failed" {
		t.Errorf("stdout %q, exit %d, stderr %s; want only the raw response and output_failed", stdout, code, stderr)
	}
}

// A write's response that breaks off is still a write that went through.
func TestWriteResultThatBreaksOff(t *testing.T) {
	_, env, _ := newFake(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", "100")
		fmt.Fprint(w, `{"id":`)
	})
	_, stderr, code := run(t, env, "team", "create", "--name", "x")
	if e := errorOf(t, stderr); code != clierr.ExitUnexpected || e.Code != "output_failed" || !strings.Contains(e.Message, "team create succeeded") {
		t.Errorf("exit %d, error %+v; want output_failed", code, e)
	}
}
