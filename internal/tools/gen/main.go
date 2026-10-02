// Command gen writes the command catalog (descriptors, schemas and search
// index) from the vendored spec and the overlay. go generate runs it.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cego/aikido-dojo/internal/gen"
	"github.com/cego/aikido-dojo/internal/overlay"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "gen:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("gen", flag.ContinueOnError)
	specPath := fs.String("spec", "", "the vendored OpenAPI spec")
	dir := fs.String("dir", "", "the catalog package directory to write into")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}
	if *specPath == "" || *dir == "" {
		return errors.New("usage: gen -spec <openapi.json> -dir <catalog dir>")
	}
	data, err := os.ReadFile(*specPath)
	if err != nil {
		return fmt.Errorf("read the spec: %w", err)
	}
	out, err := gen.Generate(data, overlay.Ops)
	if err != nil {
		return err
	}
	for _, f := range []struct {
		name string
		data []byte
	}{{"zz_ops.go", out.Ops}, {"schemas.json", out.Schemas}, {"search.json", out.Search}} {
		if err := os.WriteFile(filepath.Join(*dir, f.name), f.data, 0o600); err != nil {
			return fmt.Errorf("write %s: %w", f.name, err)
		}
	}
	return nil
}
