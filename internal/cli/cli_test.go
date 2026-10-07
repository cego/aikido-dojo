package cli

import (
	"strings"
	"testing"

	"github.com/cego/aikido-dojo/internal/clierr"
)

func TestNoArgumentsPrintsHelp(t *testing.T) {
	env, _ := testEnv(t)
	stdout, _, code := run(t, env)
	if code != clierr.ExitOK || !strings.Contains(stdout, "not affiliated with Aikido Security") {
		t.Errorf("exit %d, stdout %q; want 0 and the root help", code, stdout)
	}
}

func TestCommandLineMistakesAreUsageErrors(t *testing.T) {
	tests := []struct {
		args []string
		code string
	}{
		{[]string{"nosuch"}, "unknown_command"},
		{[]string{"--nosuch"}, "usage"},
	}
	for _, tt := range tests {
		env, _ := testEnv(t)
		stdout, stderr, code := run(t, env, tt.args...)
		if code != clierr.ExitUsage || stdout != "" {
			t.Errorf("%q: exit %d, stdout %q; want exit 2 and no stdout", tt.args, code, stdout)
		}
		if e := errorOf(t, stderr); e.Code != tt.code || e.Hint == "" {
			t.Errorf("%q: error = %+v, want code %s with a hint", tt.args, e, tt.code)
		}
	}
}

func TestRootHelpEndsWithSearchAndSchema(t *testing.T) {
	env, _ := testEnv(t)
	stdout, _, code := run(t, env, "--help")
	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	if n := len(lines); code != clierr.ExitOK || n < 2 || !strings.Contains(lines[n-2], "aikido-dojo search") || !strings.Contains(lines[n-1], "aikido-dojo schema") {
		t.Errorf("exit %d, root help ends %q", code, lines[max(0, len(lines)-3):])
	}
	if sub, _, _ := run(t, env, "repo", "--help"); strings.Contains(sub, "aikido-dojo search") {
		t.Errorf("repo help repeats the root's pointer:\n%s", sub)
	}
}
