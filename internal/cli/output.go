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
	jq     *gojq.Code
}

// value prints v, which encoding/json encodes.
func (o *output) value(ctx context.Context, v any) error {
	raw, err := compact(v)
	if err != nil {
		return err
	}
	if o.jq != nil {
		return o.filter(ctx, raw)
	}
	return o.write(raw)
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
// gh --jq prints it, anything else as JSON.
func (o *output) filter(ctx context.Context, raw []byte) error {
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
			if _, err := io.WriteString(o.w, r+"\n"); err != nil {
				return fmt.Errorf("write the output: %w", err)
			}
		default:
			b, err := compact(r)
			if err != nil {
				return err
			}
			if err := o.write(b); err != nil {
				return err
			}
		}
	}
}

// check refuses, before any call, an output flag the command's response
// can't satisfy: --jq needs JSON, which a CSV or PDF download isn't.
func (o *output) check(op ops.Op, sc ops.SchemaSet) error {
	if o.jq != nil && op.Paging == nil && !slices.Contains(sc.ResponseTypes, "application/json") {
		return notJSON(op.Command + " prints " + strings.Join(sc.ResponseTypes, " or "))
	}
	return nil
}

func notJSON(what string) error {
	return &clierr.Error{Code: "invalid_input", Message: "--jq needs JSON, but " + what,
		Hint: "drop --jq for this command", Exit: clierr.ExitUsage}
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
	if o.jq != nil {
		if !isJSON(contentType) {
			return notJSON("the response is " + contentType)
		}
		data, err := io.ReadAll(r)
		if err != nil {
			return fmt.Errorf("read the response: %w", err)
		}
		return o.filter(ctx, data)
	}
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
