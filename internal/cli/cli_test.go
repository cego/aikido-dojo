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
