package specfetch

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const IndexURL = "https://apidocs.aikido.dev/llms.txt"

const (
	fetchConcurrency = 4
	// Far above any real page; the cap only stops a misbehaving server from
	// exhausting memory.
	maxBodyBytes = 8 << 20
)

// Assemble fetches the index and every page it lists, and returns the merged
// spec and the newest page updatedAt.
func Assemble(ctx context.Context, c *http.Client, indexURL string) (map[string]any, time.Time, error) {
	origin, err := url.Parse(indexURL)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("parse index url: %w", err)
	}
	// A redirect would fetch a URL the index never listed and bypass PageURLs' origin check.
	noRedirects := *c
	noRedirects.CheckRedirect = func(req *http.Request, _ []*http.Request) error {
		return fmt.Errorf("refusing redirect to %s", req.URL)
	}
	c = &noRedirects
	index, err := get(ctx, c, indexURL)
	if err != nil {
		return nil, time.Time{}, err
	}
	urls, err := PageURLs(string(index), origin)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("%s: %w", indexURL, err)
	}
	pages, err := getAll(ctx, c, urls)
	if err != nil {
		return nil, time.Time{}, err
	}
	docs := make(map[string][]byte, len(pages))
	var newest time.Time
	for u, page := range pages {
		doc, err := ExtractSpec(page)
		if err != nil {
			return nil, time.Time{}, fmt.Errorf("%s: %w", u, err)
		}
		updated, err := UpdatedAt(page)
		if err != nil {
			return nil, time.Time{}, fmt.Errorf("%s: %w", u, err)
		}
		docs[u] = doc
		if updated.After(newest) {
			newest = updated
		}
	}
	spec, err := Merge(docs)
	if err != nil {
		return nil, time.Time{}, err
	}
	return spec, newest, nil
}

// getAll fetches with bounded concurrency and cancels the rest at the first failure.
func getAll(ctx context.Context, c *http.Client, urls []string) (map[string][]byte, error) {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)
	var (
		mu    sync.Mutex
		wg    sync.WaitGroup
		pages = make(map[string][]byte, len(urls))
		sem   = make(chan struct{}, fetchConcurrency)
	)
	for _, u := range urls {
		wg.Go(func() {
			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				return
			}
			defer func() { <-sem }()
			body, err := get(ctx, c, u)
			if err != nil {
				cancel(err)
				return
			}
			mu.Lock()
			pages[u] = body
			mu.Unlock()
		})
	}
	wg.Wait()
	if err := context.Cause(ctx); err != nil {
		return nil, fmt.Errorf("fetch pages: %w", err)
	}
	return pages, nil
}

func get(ctx context.Context, c *http.Client, u string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("get %s: %w", u, err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get %s: %w", u, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get %s: HTTP %d", u, resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes+1))
	if err != nil {
		return nil, fmt.Errorf("get %s: read body: %w", u, err)
	}
	if len(body) > maxBodyBytes {
		return nil, fmt.Errorf("get %s: body exceeds %d bytes", u, maxBodyBytes)
	}
	return body, nil
}
