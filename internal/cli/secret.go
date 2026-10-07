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
		fd := int(in.Fd()) //nolint:gosec // a file descriptor fits in an int
		if !term.IsTerminal(fd) {
			return "", &clierr.Error{Code: "no_terminal", Message: "no terminal to read the client secret from",
				Hint: "set " + config.EnvClientSecret + " for this one command instead", Exit: clierr.ExitUsage}
		}
		state, err := term.GetState(fd)
		if err != nil {
			return "", fmt.Errorf("read the terminal state: %w", err)
		}
		fmt.Fprint(w, prompt)
		type result struct {
			secret []byte
			err    error
		}
		done := make(chan result, 1)
		go func() {
			b, err := term.ReadPassword(fd)
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
			// ReadPassword turned echo off and waits for Enter until the process
			// exits, which follows; give the terminal its echo back first.
			restoreErr := term.Restore(fd, state)
			fmt.Fprintln(w)
			return "", errors.Join(fmt.Errorf("read the client secret: %w", ctx.Err()), restoreErr)
		}
	}
}
