package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/config"
)

// TerminalSecret reads a secret typed at the terminal on in, without echo,
// after writing prompt to w. Off a terminal it refuses: the hidden prompt
// and AIKIDO_DOJO_CLIENT_SECRET are the only ways a secret comes in.
func TerminalSecret(in *os.File, w io.Writer) func(ctx context.Context, prompt string) (string, error) {
	return func(ctx context.Context, prompt string) (string, error) {
		fd := int(in.Fd())
		if !term.IsTerminal(fd) {
			return "", &clierr.Error{Code: "no_terminal", Message: "no terminal to read the client secret from",
				Hint: "set " + config.EnvClientSecret + " for this one command instead", Exit: clierr.ExitUsage}
		}
		state, err := term.GetState(fd)
		if err != nil {
			return "", fmt.Errorf("read the terminal state: %w", err)
		}
		return readHidden(ctx, w, prompt, func() ([]byte, error) { return term.ReadPassword(fd) },
			func() error { return term.Restore(fd, state) })
	}
}

// readHidden writes prompt and returns what read gets, trimmed. read runs
// aside so that a cancelled ctx ends the wait at once: read has turned echo
// off and waits for Enter until the process exits, so restore gives the
// terminal its echo back first. After an earlier cancel no read starts.
func readHidden(ctx context.Context, w io.Writer, prompt string, read func() ([]byte, error), restore func() error) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("read the client secret: %w", err)
	}
	fmt.Fprint(w, prompt)
	type result struct {
		secret []byte
		err    error
	}
	done := make(chan result, 1)
	go func() {
		b, err := read()
		done <- result{b, err}
	}()
	select {
	case r := <-done:
		fmt.Fprintln(w)
		if r.err != nil {
			return "", fmt.Errorf("read the client secret: %w", r.err)
		}
		return strings.TrimSpace(string(r.secret)), nil
	case <-ctx.Done():
		restoreErr := restore()
		fmt.Fprintln(w)
		return "", errors.Join(fmt.Errorf("read the client secret: %w", ctx.Err()), restoreErr)
	}
}
