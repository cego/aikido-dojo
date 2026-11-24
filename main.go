// Command aikido-dojo is an unofficial command-line client for the Aikido
// Security public API.
package main

import (
	"context"
	_ "embed"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/term"

	"github.com/cego/aikido-dojo/internal/cli"
)

//go:embed AGENTS.md
var primer []byte

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	code := cli.Run(ctx, os.Args[1:], cli.Env{
		Stdin:      os.Stdin,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
		Getenv:     os.Getenv,
		CacheDir:   os.UserCacheDir,
		Home:       os.UserHomeDir,
		Primer:     primer,
		ReadSecret: cli.TerminalSecret(os.Stdin, os.Stderr),
		StdoutTTY:  term.IsTerminal(int(os.Stdout.Fd())),
		Confirm:    cli.TerminalConfirm(os.Stdin, os.Stderr),
	})
	stop()
	os.Exit(code)
}
