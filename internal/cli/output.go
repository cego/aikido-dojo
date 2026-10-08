package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// output prints a command's data to stdout: indented for a person when stdout
// is a terminal, compact for a program otherwise.
type output struct {
	w      io.Writer
	pretty bool
}

// value prints v, which encoding/json encodes.
func (o *output) value(_ context.Context, v any) error {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("encode the output: %w", err)
	}
	return o.write(bytes.TrimSuffix(b.Bytes(), []byte("\n")))
}

// write prints one compact JSON value, indented on a terminal.
func (o *output) write(raw []byte) error {
	var b bytes.Buffer
	if o.pretty {
		if err := json.Indent(&b, raw, "", "  "); err != nil {
			return fmt.Errorf("indent the output: %w", err)
		}
	} else {
		b.Write(raw)
	}
	b.WriteByte('\n')
	if _, err := o.w.Write(b.Bytes()); err != nil {
		return fmt.Errorf("write the output: %w", err)
	}
	return nil
}

// response prints an API response as it arrives. Off a terminal, Aikido's
// JSON is already compact (live, 2026-10-08), so it passes through unchanged;
// CSV and PDF always do.
func (o *output) response(_ context.Context, r io.Reader, contentType string) error {
	if o.pretty && isJSON(contentType) {
		return prettyCopy(o.w, r)
	}
	if _, err := io.Copy(o.w, r); err != nil {
		return fmt.Errorf("read the response: %w", err)
	}
	return nil
}

func isJSON(contentType string) bool { return strings.HasPrefix(contentType, "application/json") }

func (a *app) printJSON(v any) error { return a.out.value(context.Background(), v) }
