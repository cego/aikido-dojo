package api

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"
	"time"
)

// fill takes n of the window's 20 slots on the real clock.
func fill(t *testing.T, ta *testAPI, n int) {
	t.Helper()
	for range n {
		if err := ta.client.limiter.Wait(t.Context()); err != nil {
			t.Fatal(err)
		}
	}
}

func TestDoWaitsForTheRateLimiter(t *testing.T) {
	var attempts atomic.Int32
	ta := newTestAPI(t, sequence(&attempts), nil)
	fill(t, ta, 20)
	ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
	defer cancel()
	resp, err := ta.client.Do(ctx, Request{Method: http.MethodGet, Path: "/workspace"})
	if err == nil {
		defer resp.Body.Close()
		t.Fatal("Do sent a request while the window was full")
	}
	if !errors.Is(err, context.DeadlineExceeded) || attempts.Load() != 0 {
		t.Errorf("err = %v, attempts = %d; want a deadline while waiting and no request", err, attempts.Load())
	}
}

func TestDoWaitsForTheRateLimiterBeforeARetry(t *testing.T) {
	var attempts atomic.Int32
	ta := newTestAPI(t, sequence(&attempts, 429), nil)
	fill(t, ta, 19)
	ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
	defer cancel()
	resp, err := ta.client.Do(ctx, Request{Method: http.MethodGet, Path: "/workspace"})
	if err == nil {
		defer resp.Body.Close()
		t.Fatal("the retry went out while the window was full")
	}
	if !errors.Is(err, context.DeadlineExceeded) || attempts.Load() != 1 {
		t.Errorf("err = %v, attempts = %d; want the retry to wait for the limiter", err, attempts.Load())
	}
}
