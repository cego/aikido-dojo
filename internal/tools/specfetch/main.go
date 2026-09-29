// Command specfetch rebuilds spec/ from Aikido's public API reference pages.
// Run it from the repository root: go run ./internal/tools/specfetch
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cego/aikido-dojo/internal/specfetch"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "specfetch:", err)
		os.Exit(1)
	}
}

func run() error {
	// spec/ is resolved against the working directory, so refuse to scatter it elsewhere.
	if _, err := os.Stat("go.mod"); err != nil {
		return fmt.Errorf("run from the repository root: %w", err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	client := &http.Client{Timeout: 30 * time.Second}
	spec, updated, err := specfetch.Assemble(ctx, client, specfetch.IndexURL)
	if err != nil {
		return fmt.Errorf("assemble: %w", err)
	}
	if err := specfetch.Write("spec", spec, updated); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
