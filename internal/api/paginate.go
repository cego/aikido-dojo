package api

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"maps"
	"net/http"
	"net/url"
	"strconv"
)

// Style is how an operation pages. Pages are 0-indexed in every style.
type Style int

const (
	// PageArray: the body is an array and the list ends on an empty page.
	PageArray Style = iota + 1
	// PageHeader: the body is an array and only X-Has-Next-Page ends the
	// list. Aikido filters after paging, so a short or empty page is not the end.
	PageHeader
	// PageEnvelope: the body is an object holding the items. The list ends
	// when More is false, or on an empty page when there is no More field.
	PageEnvelope
)

type Paging struct {
	Style     Style
	SizeParam string // "per_page" or "limit"
	Size      int
	Items     string // PageEnvelope: the field holding the items
	More      string // PageEnvelope: the field saying more pages exist, or ""
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
	dec := json.NewDecoder(resp.Body)
	if p.Style != PageEnvelope {
		var items []json.RawMessage
		if err := dec.Decode(&items); err != nil {
			return nil, false, fmt.Errorf("want a JSON array: %w", err)
		}
		if p.Style == PageHeader {
			return items, resp.Header.Get("X-Has-Next-Page") == "true", nil
		}
		return items, len(items) > 0, nil
	}
	var env map[string]json.RawMessage
	if err := dec.Decode(&env); err != nil {
		return nil, false, fmt.Errorf("want a JSON object: %w", err)
	}
	raw, ok := env[p.Items]
	if !ok {
		return nil, false, fmt.Errorf("no %q field in the response", p.Items)
	}
	var items []json.RawMessage
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, false, fmt.Errorf("field %q: %w", p.Items, err)
	}
	if p.More == "" {
		return items, len(items) > 0, nil
	}
	var more bool
	if raw, ok := env[p.More]; ok {
		if err := json.Unmarshal(raw, &more); err != nil {
			return nil, false, fmt.Errorf("field %q is not a boolean: %w", p.More, err)
		}
	}
	return items, more, nil
}
