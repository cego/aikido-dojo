package cli

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/cego/aikido-dojo/internal/catalog"
	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/config"
)

var writes = [][]string{
	{"team", "create", "--name", "x"},
	{"team", "delete", "1"},
	{"api", "POST", "/teams"},
	{"api", "DELETE", "/teams/1"},
	{"auth", "login", "--client-id", "x"},
	{"--profile", "p", "auth", "logout"},
}

func TestReadOnlyRefusesEveryWrite(t *testing.T) {
	for _, viaEnv := range []bool{false, true} {
		for _, args := range writes {
			f, env, vars := newFake(t, respond(`{}`))
			if viaEnv {
				vars[config.EnvReadOnly] = "1"
			} else {
				args = append([]string{"--read-only"}, args...)
			}
			_, stderr, code := run(t, env, args...)
			e := errorOf(t, stderr)
			if code != clierr.ExitRefused || e.Code != "read_only" || len(f.seen())+len(f.seenLogins()) != 0 {
				t.Errorf("%q (env %v): exit %d, error %+v; want 7 read_only before any call", args, viaEnv, code, e)
			}
			// The mode is a guard handed to agents, so the hint must not say how to lift it.
			if strings.Contains(e.Hint, "unset") || strings.Contains(e.Hint, "drop") {
				t.Errorf("hint = %q, want no way out", e.Hint)
			}
		}
	}
}

func TestReadOnlyAllowsReads(t *testing.T) {
	_, env, _ := newFake(t, pages(`[1]`, `[]`))
	mustRun(t, env, "--read-only", "repo", "list")
	_, env2, _ := newFake(t, respond(`{}`))
	mustRun(t, env2, "--read-only", "api", "GET", "/workspace")
}

func TestReadOnlyEnvMustBeABoolean(t *testing.T) {
	f, env, vars := newFake(t, pages(`[1]`, `[]`))
	vars[config.EnvReadOnly] = "maybe"
	_, stderr, code := run(t, env, "repo", "list")
	if e := errorOf(t, stderr); code != clierr.ExitUsage || !strings.Contains(e.Message, config.EnvReadOnly) || len(f.seen()) != 0 {
		t.Errorf("exit %d, error %+v; want a usage error naming %s", code, e, config.EnvReadOnly)
	}
}

// The global --read-only takes the flag name of the body field read_only.
func TestReadOnlyFlagMakesTheBodyFieldBodyOnly(t *testing.T) {
	env, _ := testEnv(t)
	if out := mustRun(t, env, "user-role", "update", "--help"); !strings.Contains(out, "read_only share a name with a built-in flag") {
		t.Errorf("help = %s, want read_only listed as body-only", out)
	}
}

func TestDryRunPrintsTheRequest(t *testing.T) {
	f, env, vars := newFake(t, respond(`{}`))
	vars[config.EnvRegion] = "us"
	want := `{"dry_run":true,"method":"POST","url":"https://app.us.aikido.dev/api/public/v1/teams",` +
		`"headers":{"Authorization":"Bearer [redacted]","Content-Type":"application/json"},"body":{"name":"x"}}` + "\n"
	if out := mustRun(t, env, "team", "create", "--name", "x", "--dry-run"); out != want {
		t.Errorf("stdout = %s, want %s", out, want)
	}
	if n := len(f.seen()) + len(f.seenLogins()); n != 0 {
		t.Errorf("%d requests, want none", n)
	}
}

func TestDryRunRedactsCredentials(t *testing.T) {
	for _, args := range [][]string{
		{"code-scanning-token", "update", "--body-file", "-", "--dry-run"},
		{"api", "POST", "/access-tokens/code-scanning", "--input", "-", "--dry-run"},
		// The wrong method, or a path no command calls: a credential is still a credential.
		{"api", "PUT", "/access-tokens/code-scanning", "--input", "-", "--dry-run"},
		{"api", "POST", "/nosuch", "--input", "-", "--dry-run"},
	} {
		_, env, _ := newFake(t, respond(`{}`))
		env.Stdin = strings.NewReader(`{"access_token":"s3cr3t"}`)
		out := mustRun(t, env, args...)
		if strings.Contains(out, "s3cr3t") || !strings.Contains(out, `"access_token":"[redacted]"`) {
			t.Errorf("%q: stdout = %s, want the credential redacted", args, out)
		}
	}
}

// Nothing is sent, so a dry run is how an agent in read-only mode shows a
// person the write it would make.
func TestDryRunWorksInReadOnlyMode(t *testing.T) {
	_, env, _ := newFake(t, respond(`{}`))
	if out := mustRun(t, env, "--read-only", "team", "delete", "1", "--dry-run"); !strings.Contains(out, `"method":"DELETE"`) {
		t.Errorf("stdout = %s, want the DELETE request", out)
	}
}

func TestDryRunIsOnlyOnWrites(t *testing.T) {
	env, _ := testEnv(t)
	if _, stderr, code := run(t, env, "repo", "list", "--dry-run"); code != clierr.ExitUsage || !strings.Contains(stderr, "unknown flag: --dry-run") {
		t.Errorf("repo list --dry-run: exit %d, %s; want an unknown flag", code, stderr)
	}
	_, env2, _ := newFake(t, respond(`{}`))
	out := mustRun(t, env2, "api", "GET", "/workspace", "--dry-run")
	if !strings.Contains(out, `"method":"GET"`) || strings.Contains(out, "Content-Type") || strings.Contains(out, `"body"`) {
		t.Errorf("api GET --dry-run = %s, want a GET with no body", out)
	}
}

func TestDestructiveNeedsYesOffATerminal(t *testing.T) {
	all, ok := catalog.Command("pr-check-config-all update")
	if !ok || !all.Destructive {
		t.Fatal("pr-check-config-all update is not a destructive command")
	}
	for _, args := range [][]string{
		{"team", "delete", "1"},
		{"api", "DELETE", "/teams/1"},
		// Not a DELETE, but a path a destructive command calls.
		{"api", all.Method, all.Path},
	} {
		f, env, _ := newFake(t, respond(`{}`))
		_, stderr, code := run(t, env, args...)
		e := errorOf(t, stderr)
		if code != clierr.ExitRefused || e.Code != "confirmation_required" || !strings.Contains(e.Hint, "--yes") || len(f.seen())+len(f.seenLogins()) != 0 {
			t.Errorf("%q: exit %d, error %+v; want 7 confirmation_required before any call", args, code, e)
		}
	}
}

func TestDestructiveWithYes(t *testing.T) {
	for _, args := range [][]string{{"team", "delete", "1", "--yes"}, {"api", "DELETE", "/teams/1", "--yes"}} {
		f, env, _ := newFake(t, respond(`{}`))
		mustRun(t, env, args...)
		if n := len(f.seen()); n != 1 {
			t.Errorf("%q: %d calls, want 1", args, n)
		}
	}
}

func TestDestructiveAsksOnATerminal(t *testing.T) {
	for _, answer := range []bool{true, false} {
		f, env, _ := newFake(t, respond(`{}`))
		var prompt string
		env.Confirm = func(_ context.Context, p string) (bool, error) {
			prompt = p
			return answer, nil
		}
		_, stderr, code := run(t, env, "team", "delete", "1")
		if !strings.Contains(prompt, "team delete 1") {
			t.Errorf("prompt = %q, want it to name the command", prompt)
		}
		if answer && (code != clierr.ExitOK || len(f.seen()) != 1) {
			t.Errorf("answered yes: exit %d, %d calls, stderr %s", code, len(f.seen()), stderr)
		}
		if !answer && (code != clierr.ExitRefused || errorOf(t, stderr).Code != "not_confirmed" || len(f.seen()) != 0) {
			t.Errorf("answered no: exit %d, %d calls, stderr %s", code, len(f.seen()), stderr)
		}
	}
}

func TestDryRunNeedsNoConfirmation(t *testing.T) {
	_, env, _ := newFake(t, respond(`{}`))
	mustRun(t, env, "team", "delete", "1", "--dry-run")
}

// Ctrl-C at the question is a no, and a no says how to get a yes.
func TestDestructiveCancelledIsNotConfirmed(t *testing.T) {
	for _, answer := range []error{nil, fmt.Errorf("read the answer: %w", context.Canceled)} {
		f, env, _ := newFake(t, respond(`{}`))
		env.Confirm = func(context.Context, string) (bool, error) { return false, answer }
		_, stderr, code := run(t, env, "team", "delete", "1")
		if e := errorOf(t, stderr); code != clierr.ExitRefused || e.Code != "not_confirmed" || !strings.Contains(e.Hint, "--yes") || len(f.seen()) != 0 {
			t.Errorf("answer error %v: exit %d, error %+v; want 7 not_confirmed with a hint", answer, code, e)
		}
	}
}
