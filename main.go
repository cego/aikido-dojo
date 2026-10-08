// Command aikido-dojo is an unofficial command-line client for the Aikido
// Security public API.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/term"

	"github.com/cego/aikido-dojo/internal/cli"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	code := cli.Run(ctx, os.Args[1:], cli.Env{
		Stdin:      os.Stdin,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
		Getenv:     os.Getenv,
		CacheDir:   os.UserCacheDir,
		ReadSecret: cli.TerminalSecret(os.Stdin, os.Stderr),
		StdoutTTY:  term.IsTerminal(int(os.Stdout.Fd())),
	})
	stop()
	os.Exit(code)
}
