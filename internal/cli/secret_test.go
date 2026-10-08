package cli

import (
	"bytes"
	"context"
	"errors"
	"io"
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

func TestAskReturnsTheTrimmedSecret(t *testing.T) {
	var w bytes.Buffer
	got, err := ask(t.Context(), &w, "Secret: ", "the client secret", func() ([]byte, error) { return []byte(" s3cr3t \n"), nil }, func() error {
		t.Error("restore ran after a read that ended on its own")
		return nil
	})
	if err != nil || got != "s3cr3t" || w.String() != "Secret: " {
		t.Errorf("ask = %q, %v, wrote %q", got, err, w.String())
	}
}

func TestAskReportsAReadError(t *testing.T) {
	_, err := ask(t.Context(), io.Discard, "Secret: ", "the client secret", func() ([]byte, error) { return nil, errors.New("tty gone") }, func() error { return nil })
	if err == nil || !strings.Contains(err.Error(), "read the client secret: tty gone") {
		t.Errorf("err = %v", err)
	}
}

// Ctrl-C while the read waits must give the terminal its echo back at once.
func TestAskRestoresTheTerminalOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	unblock := make(chan struct{})
	defer close(unblock)
	restored := 0
	read := func() ([]byte, error) {
		cancel()
		<-unblock
		return nil, nil
	}
	_, err := ask(ctx, io.Discard, "Secret: ", "the client secret", read, func() error { restored++; return nil })
	if !errors.Is(err, context.Canceled) || restored != 1 {
		t.Errorf("err = %v, restored %d times; want context.Canceled and one restore", err, restored)
	}
}

// A Ctrl-C before the prompt starts no read, which could turn echo off
// after the process had given up on it.
func TestAskAfterACancelReadsNothing(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	var w bytes.Buffer
	_, err := ask(ctx, &w, "Secret: ", "the client secret", func() ([]byte, error) {
		t.Error("read started after the cancel")
		return nil, nil
	}, func() error { return nil })
	if !errors.Is(err, context.Canceled) || w.Len() != 0 {
		t.Errorf("err = %v, wrote %q; want context.Canceled and no prompt", err, w.String())
	}
}

// A visible question has no terminal state to restore.
func TestAskWithoutRestore(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	unblock := make(chan struct{})
	defer close(unblock)
	_, err := ask(ctx, io.Discard, "Continue? ", "the answer", func() ([]byte, error) {
		cancel()
		<-unblock
		return nil, nil
	}, nil)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want context.Canceled", err)
	}
}

func TestTerminalConfirmRefusesAPipe(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()
	_, err = TerminalConfirm(r, w)(t.Context(), "Continue? ")
	var e *clierr.Error
	if !errors.As(err, &e) || e.Code != "no_terminal" {
		t.Errorf("err = %v, want no_terminal", err)
	}
}

func TestAffirmative(t *testing.T) {
	for in, want := range map[string]bool{"y": true, "Y": true, "yes": true, "YES": true, "": false, "n": false, "no": false, "yeah": false} {
		if got := affirmative(in); got != want {
			t.Errorf("affirmative(%q) = %v, want %v", in, got, want)
		}
	}
}
