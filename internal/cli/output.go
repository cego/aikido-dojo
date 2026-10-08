package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/itchyny/gojq"

	"github.com/cego/aikido-dojo/internal/clierr"
	"github.com/cego/aikido-dojo/internal/ops"
)

// output prints a command's data to stdout: indented for a person when stdout
// is a terminal, compact for a program otherwise, and through --jq when set.
type output struct {
	w      io.Writer
	pretty bool
	ndjson bool // one array item per line
	jq     *gojq.Code
}

// value prints v, which encoding/json encodes.
func (o *output) value(ctx context.Context, v any) error {
	raw, err := compact(v)
	if err != nil {
		return err
	}
	if o.ndjson {
		var items []json.RawMessage
		if json.Unmarshal(raw, &items) != nil {
			return notAnArray("this command prints one object")
		}
		for _, item := range items {
			if err := o.item(ctx, item); err != nil {
				return err
			}
		}
		return nil
	}
	if o.jq != nil {
		return o.filter(ctx, raw, o.write)
	}
	return o.write(raw)
}

// item prints one array item on a line of its own, or each --jq result of it.
func (o *output) item(ctx context.Context, raw json.RawMessage) error {
	if o.jq != nil {
		return o.filter(ctx, raw, o.line)
	}
	return o.line(raw)
}

// line prints raw compact, on one line.
func (o *output) line(raw []byte) error {
	var b bytes.Buffer
	if err := json.Compact(&b, raw); err != nil {
		return fmt.Errorf("compact the output: %w", err)
	}
	b.WriteByte('\n')
	if _, err := o.w.Write(b.Bytes()); err != nil {
		return fmt.Errorf("write the output: %w", err)
	}
	return nil
}

func compact(v any) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, fmt.Errorf("encode the output: %w", err)
	}
	return bytes.TrimSuffix(b.Bytes(), []byte("\n")), nil
}

// filter runs --jq over raw and prints each result: a string as it is, as
// gh --jq prints it, anything else as JSON through emit.
func (o *output) filter(ctx context.Context, raw []byte, emit func([]byte) error) error {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return fmt.Errorf("decode the output for --jq: %w", err)
	}
	iter := o.jq.RunWithContext(ctx, v)
	for {
		result, ok := iter.Next()
		if !ok {
			return nil
		}
		switch r := result.(type) {
		case error:
			return &clierr.Error{Code: "jq_failed", Message: "--jq: " + r.Error(), Hint: "check the filter against: aikido-dojo schema <command>", Exit: clierr.ExitUsage}
		case string:
			if o.pretty {
				r = escapeControls(r)
			}
			if _, err := io.WriteString(o.w, r+"\n"); err != nil {
				return fmt.Errorf("write the output: %w", err)
			}
		default:
			b, err := compact(r)
			if err != nil {
				return err
			}
			if err := emit(b); err != nil {
				return err
			}
		}
	}
}

// escapeControls writes control characters other than newline and tab as
// \u escapes: a raw string bound for a terminal comes from Aikido's data,
// such as names in scanned repos, and mustn't carry escape sequences to it.
func escapeControls(s string) string {
	var b strings.Builder
	for _, r := range s {
		if (r < 0x20 && r != '\n' && r != '\t') || (r >= 0x7f && r <= 0x9f) {
			fmt.Fprintf(&b, "\\u%04x", r)
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// check refuses, before any call, an output flag the command's response
// can't satisfy: --jq needs JSON, which a CSV or PDF download isn't.
func (o *output) check(op ops.Op, sc ops.SchemaSet) error {
	if op.Paging != nil {
		return nil
	}
	if (o.jq != nil || o.ndjson) && !slices.Contains(sc.ResponseTypes, "application/json") {
		if len(sc.ResponseTypes) == 0 {
			return notJSON(op.Command + " prints no JSON")
		}
		return notJSON(op.Command + " prints " + strings.Join(sc.ResponseTypes, " or "))
	}
	if o.ndjson && sc.Response["type"] != "array" {
		return notAnArray(op.Command + " prints one object")
	}
	return nil
}

func notAnArray(what string) error {
	return &clierr.Error{Code: "invalid_input", Message: "--ndjson needs a list, but " + what,
		Hint: "drop --ndjson; it applies to commands that print an array", Exit: clierr.ExitUsage}
}

func notJSON(what string) error {
	return &clierr.Error{Code: "invalid_input", Message: "--jq and --ndjson need JSON, but " + what,
		Hint: "drop --jq and --ndjson for this command", Exit: clierr.ExitUsage}
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
func (o *output) response(ctx context.Context, r io.Reader, contentType string) error {
	if (o.jq != nil || o.ndjson) && !isJSON(contentType) {
		return notJSON("the response is " + contentType)
	}
	if o.ndjson {
		return o.items(ctx, r)
	}
	if o.jq != nil {
		data, err := io.ReadAll(r)
		if err != nil {
			return fmt.Errorf("read the response: %w", err)
		}
		return o.filter(ctx, data, o.write)
	}
	if o.pretty && isJSON(contentType) {
		return prettyCopy(o.w, r)
	}
	if _, err := io.Copy(o.w, r); err != nil {
		return fmt.Errorf("read the response: %w", err)
	}
	return nil
}

// items prints each element of the JSON array in r as it arrives, so an
// export streams instead of being held in memory.
func (o *output) items(ctx context.Context, r io.Reader) error {
	dec := json.NewDecoder(r)
	tok, err := dec.Token()
	if err != nil {
		return fmt.Errorf("read the response: %w", err)
	}
	if tok != json.Delim('[') {
		return notAnArray("the response is not an array")
	}
	for dec.More() {
		var item json.RawMessage
		if err := dec.Decode(&item); err != nil {
			return fmt.Errorf("read the response: %w", err)
		}
		if err := o.item(ctx, item); err != nil {
			return err
		}
	}
	if _, err := dec.Token(); err != nil {
		return fmt.Errorf("read the response: %w", err)
	}
	return nil
}

// transforms reports whether printing a response of contentType changes it,
// rather than copying it through.
func (o *output) transforms(contentType string) bool {
	return o.jq != nil || o.ndjson || (o.pretty && isJSON(contentType))
}

func isJSON(contentType string) bool { return strings.HasPrefix(contentType, "application/json") }

func (a *app) printJSON(ctx context.Context, v any) error { return a.out.value(ctx, v) }
