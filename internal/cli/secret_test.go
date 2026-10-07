package cli

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/config"
)

func TestTerminalSecretRefusesAPipe(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()
	var prompt bytes.Buffer
	_, err = TerminalSecret(r, &prompt)(t.Context(), "Client secret: ")
	var e *clierr.Error
	if !errors.As(err, &e) || e.Code != "no_terminal" || e.Exit != clierr.ExitUsage || !strings.Contains(e.Hint, config.EnvClientSecret) {
		t.Errorf("err = %v, want no_terminal pointing at %s", err, config.EnvClientSecret)
	}
	if prompt.Len() != 0 {
		t.Errorf("prompted %q off a terminal", prompt.String())
	}
}
