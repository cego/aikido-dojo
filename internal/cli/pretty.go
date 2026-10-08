package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

// prettyCopy indents the JSON in r as it reads it, two spaces a level, as
// json.Indent would. It holds one token at a time, so an export of any size
// streams, and it keeps keys in their order. Strings are written anew, which
// drops the \/ escapes Aikido sends; control characters stay escaped.
func prettyCopy(w io.Writer, r io.Reader) error {
	dec := json.NewDecoder(r)
	dec.UseNumber()
	// Written to a buffer, whose writes can't fail, and flushed now and then,
	// so memory stays small and each write to w is checked.
	var bw bytes.Buffer
	flush := func() error {
		if _, err := w.Write(bw.Bytes()); err != nil {
			return fmt.Errorf("write the output: %w", err)
		}
		bw.Reset()
		return nil
	}
	type frame struct {
		object bool
		count  int
		key    bool // an object expects a key next
	}
	var stack []frame
	newline := func() {
		bw.WriteByte('\n')
		bw.WriteString(strings.Repeat("  ", len(stack)))
	}
	for {
		if bw.Len() > 32<<10 {
			if err := flush(); err != nil {
				return err
			}
		}
		tok, err := dec.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("read the JSON response: %w", err)
		}
		if d, ok := tok.(json.Delim); ok && (d == '}' || d == ']') {
			closed := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if closed.count > 0 {
				newline()
			}
			bw.WriteRune(rune(d))
			if len(stack) == 0 {
				bw.WriteByte('\n')
			} else if top := &stack[len(stack)-1]; top.object {
				top.key = true
			}
			continue
		}
		if len(stack) > 0 {
			top := &stack[len(stack)-1]
			switch {
			case top.object && top.key:
				if top.count > 0 {
					bw.WriteByte(',')
				}
				newline()
				bw.Write(encodeString(tok.(string)))
				bw.WriteString(": ")
				top.count++
				top.key = false
				continue
			case !top.object:
				if top.count > 0 {
					bw.WriteByte(',')
				}
				newline()
				top.count++
			}
		}
		if d, ok := tok.(json.Delim); ok {
			bw.WriteRune(rune(d))
			stack = append(stack, frame{object: d == '{', key: d == '{'})
			continue
		}
		bw.Write(encodeScalar(tok))
		if len(stack) == 0 {
			bw.WriteByte('\n')
		} else if top := &stack[len(stack)-1]; top.object {
			top.key = true
		}
	}
	if err := flush(); err != nil {
		return err
	}
	if len(stack) > 0 {
		return errors.New("read the JSON response: it ended early")
	}
	return nil
}

func encodeScalar(tok json.Token) []byte {
	switch v := tok.(type) {
	case string:
		return encodeString(v)
	case json.Number:
		return []byte(v)
	case bool:
		if v {
			return []byte("true")
		}
		return []byte("false")
	}
	return []byte("null")
}

// encodeString writes s as JSON without HTML escaping: nothing renders the output as HTML.
func encodeString(s string) []byte {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(s) // a string always encodes
	return bytes.TrimSuffix(b.Bytes(), []byte("\n"))
}
