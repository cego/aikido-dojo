// Command docsfetch refreshes internal/helpdocs from help.aikido.dev. Run it from the repository root: go run ./internal/tools/docsfetch [-full]
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/cego/aikido-dojo/internal/docsfetch"
)

// Names the fetcher to Aikido, and Cloudflare refuses some default agents.
const userAgent = "aikido-dojo-docsfetch (+https://github.com/cego/aikido-dojo)"

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "docsfetch:", err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("docsfetch", flag.ContinueOnError)
	full := flags.Bool("full", false, "refetch every page, not only those the sitemap says changed")
	if err := flags.Parse(args); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}
	// The copy is resolved against the working directory, so refuse to scatter it elsewhere.
	if _, err := os.Stat("go.mod"); err != nil {
		return fmt.Errorf("run from the repository root: %w", err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	counter := &counting{rt: http.DefaultTransport}
	client := &http.Client{Timeout: 30 * time.Second, Transport: counter}
	r, err := docsfetch.Sync(ctx, client, docsfetch.Options{Origin: docsfetch.Origin, Dir: filepath.Join("internal", "helpdocs"), Full: *full})
	if err != nil {
		return fmt.Errorf("sync: %w", err)
	}
	r.Requests = counter.n.Load()
	if err := json.NewEncoder(stdout).Encode(r); err != nil {
		return fmt.Errorf("print the report: %w", err)
	}
	return nil
}

// counting sets the User-Agent and counts requests, so each run reports its load on the site.
type counting struct {
	rt http.RoundTripper
	n  atomic.Int64
}

func (c *counting) RoundTrip(r *http.Request) (*http.Response, error) {
	c.n.Add(1)
	r = r.Clone(r.Context())
	r.Header.Set("User-Agent", userAgent)
	return c.rt.RoundTrip(r) //nolint:wrapcheck // a RoundTripper passes its transport's error through
}
