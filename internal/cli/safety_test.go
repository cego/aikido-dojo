package cli

import (
	"strings"
	"testing"

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
