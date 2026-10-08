package cli

import (
	"bufio"
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
		read := func() ([]byte, error) {
			b, err := term.ReadPassword(fd)
			fmt.Fprintln(w) // the Enter wasn't echoed either
			return b, err
		}
		return ask(ctx, w, prompt, "the client secret", read, func() error { return term.Restore(fd, state) })
	}
}

// TerminalConfirm asks prompt on w and reads the answer typed at the
// terminal on in. Off a terminal it fails with no_terminal, so the caller can
// refuse instead of asking.
func TerminalConfirm(in *os.File, w io.Writer) func(ctx context.Context, prompt string) (bool, error) {
	return func(ctx context.Context, prompt string) (bool, error) {
		if !term.IsTerminal(int(in.Fd())) {
			return false, &clierr.Error{Code: "no_terminal", Message: "no terminal to confirm on", Exit: clierr.ExitUsage}
		}
		read := func() ([]byte, error) {
			b, err := bufio.NewReader(in).ReadBytes('\n')
			if errors.Is(err, io.EOF) {
				return b, nil // Ctrl-D ends the answer, which is then a no
			}
			return b, err
		}
		answer, err := ask(ctx, w, prompt, "the answer", read, nil)
		if err != nil {
			return false, err
		}
		return affirmative(answer), nil
	}
}

func affirmative(answer string) bool {
	a := strings.ToLower(answer)
	return a == "y" || a == "yes"
}

// ask writes prompt and returns what read gets, trimmed. read runs aside so
// that a cancelled ctx ends the wait at once. A hidden read turns echo off
// and waits for Enter until the process exits, so restore, when given,
// gives the terminal its echo back. After an earlier cancel no read starts;
// only a cancel in the instant before read turns echo off can still leave it
// off.
func ask(ctx context.Context, w io.Writer, prompt, what string, read func() ([]byte, error), restore func() error) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("read %s: %w", what, err)
	}
	fmt.Fprint(w, prompt)
	type result struct {
		answer []byte
		err    error
	}
	done := make(chan result, 1)
	go func() {
		b, err := read()
		done <- result{b, err}
	}()
	select {
	case r := <-done:
		if r.err != nil {
			return "", fmt.Errorf("read %s: %w", what, r.err)
		}
		return strings.TrimSpace(string(r.answer)), nil
	case <-ctx.Done():
		var restoreErr error
		if restore != nil {
			restoreErr = restore()
		}
		fmt.Fprintln(w)
		return "", errors.Join(fmt.Errorf("read %s: %w", what, ctx.Err()), restoreErr)
	}
}
