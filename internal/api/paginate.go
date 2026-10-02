package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"maps"
	"net/http"
	"net/url"
	"strconv"
)

// End is how a paged list marks its last page. Pages are 0-indexed in every style.
type End int

const (
	// EndEmpty: the list ends on an empty page.
	EndEmpty End = iota + 1
	// EndHeader: only X-Has-Next-Page ends the list. Aikido filters after
	// paging, so a short or empty page is not the end.
	EndHeader
	// EndField: the list ends when the More field is false.
	EndField
)

type Paging struct {
	End       End
	SizeParam string // "per_page" or "limit"
	Size      int
	Items     string // the field holding the items when the body is an object; "" when the body is the array
	More      string // EndField: the boolean field saying more pages exist
}

// Items fetches page after page and yields each item. Stopping the loop stops
// the fetching, so a caller that wants N items spends only the calls it needs.
func (c *Client) Items(ctx context.Context, req Request, p Paging) iter.Seq2[json.RawMessage, error] {
	return func(yield func(json.RawMessage, error) bool) {
		for n := 0; ; n++ {
			q := url.Values{}
			maps.Copy(q, req.Query)
			q.Set("page", strconv.Itoa(n))
			q.Set(p.SizeParam, strconv.Itoa(p.Size))
			pageReq := req
			pageReq.Query = q
			resp, err := c.Do(ctx, pageReq)
			if err != nil {
				yield(nil, err)
				return
			}
			items, more, err := decodePage(resp, p)
			if err != nil {
				yield(nil, fmt.Errorf("%s %s: decode page %d: %w", req.Method, req.Path, n, err))
				return
			}
			for _, item := range items {
				if !yield(item, nil) {
					return
				}
			}
			if !more {
				return
			}
		}
	}
}

func decodePage(resp *http.Response, p Paging) ([]json.RawMessage, bool, error) {
	defer resp.Body.Close()
	items, env, err := decodeItems(resp.Body, p.Items)
	if err != nil {
		return nil, false, err
	}
	switch p.End {
	case EndHeader:
		return items, resp.Header.Get("X-Has-Next-Page") == "true", nil
	case EndField:
		var more bool
		if raw, ok := env[p.More]; ok {
			if err := json.Unmarshal(raw, &more); err != nil {
				return nil, false, fmt.Errorf("field %q is not a boolean: %w", p.More, err)
			}
		}
		return items, more, nil
	}
	return items, len(items) > 0, nil
}

// decodeItems returns a page's items and, when they sit in an object, that
// object's fields.
func decodeItems(r io.Reader, field string) ([]json.RawMessage, map[string]json.RawMessage, error) {
	dec := json.NewDecoder(r)
	var items []json.RawMessage
	if field == "" {
		if err := dec.Decode(&items); err != nil {
			return nil, nil, fmt.Errorf("want a JSON array: %w", err)
		}
		return items, nil, nil
	}
	var env map[string]json.RawMessage
	if err := dec.Decode(&env); err != nil {
		return nil, nil, fmt.Errorf("want a JSON object: %w", err)
	}
	raw, ok := env[field]
	if !ok {
		return nil, nil, fmt.Errorf("no %q field in the response", field)
	}
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, nil, fmt.Errorf("field %q: %w", field, err)
	}
	return items, env, nil
}
