package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/cego/aikido-dojo/internal/clierr"
)

// sequence answers the i-th request with statuses[i], and 200 after the list runs out.
func sequence(attempts *atomic.Int32, statuses ...int) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		i := int(attempts.Add(1)) - 1
		if i < len(statuses) {
			if statuses[i] == http.StatusTooManyRequests {
				w.Header().Set("Retry-After", "2")
			}
			w.WriteHeader(statuses[i])
			return
		}
		fmt.Fprint(w, `{}`)
	}
}

func TestDoHonoursRetryAfter(t *testing.T) {
	var attempts atomic.Int32
	ta := newTestAPI(t, sequence(&attempts, 429, 429), nil)
	resp, err := ta.client.Do(t.Context(), Request{Method: http.MethodGet, Path: "/workspace"})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if want := []time.Duration{2 * time.Second, 2 * time.Second}; !slices.Equal(ta.slept, want) {
		t.Errorf("slept %v, want %v", ta.slept, want)
	}
}

func TestDoGivesUpAfterFourRateLimits(t *testing.T) {
	var attempts atomic.Int32
	ta := newTestAPI(t, sequence(&attempts, 429, 429, 429, 429, 429), nil)
	err := doErr(t, ta.client, Request{Method: http.MethodGet, Path: "/workspace"})
	wantCode(t, err, "rate_limited", clierr.ExitRateLimited)
	if got := attempts.Load(); got != 5 {
		t.Errorf("attempts = %d, want 5", got)
	}
}

func TestDoRetriesAGetOn5xxWithGrowingWaits(t *testing.T) {
	var attempts atomic.Int32
	ta := newTestAPI(t, sequence(&attempts, 502, 503), nil)
	ta.client.backoff = func(attempt int) time.Duration { return time.Duration(attempt+1) * time.Second }
	resp, err := ta.client.Do(t.Context(), Request{Method: http.MethodGet, Path: "/workspace"})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if want := []time.Duration{time.Second, 2 * time.Second}; !slices.Equal(ta.slept, want) {
		t.Errorf("slept %v, want %v", ta.slept, want)
	}
}

func TestDoStopsRetryingAGetAfterThreeAttempts(t *testing.T) {
	var attempts atomic.Int32
	ta := newTestAPI(t, sequence(&attempts, 500, 500, 500, 500), nil)
	err := doErr(t, ta.client, Request{Method: http.MethodGet, Path: "/workspace"})
	wantCode(t, err, "server_error", clierr.ExitUnexpected)
	if got := attempts.Load(); got != 3 {
		t.Errorf("attempts = %d, want 3", got)
	}
}

func TestDoResendsTheBodyAfterA429(t *testing.T) {
	var mu sync.Mutex
	var bodies []string
	var attempts atomic.Int32
	next := sequence(&attempts, 429)
	ta := newTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		mu.Lock()
		bodies = append(bodies, string(b))
		mu.Unlock()
		next(w, r)
	}, nil)
	resp, err := ta.client.Do(t.Context(), Request{Method: http.MethodPost, Path: "/teams", Body: []byte(`{"name":"x"}`)})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	mu.Lock()
	defer mu.Unlock()
	if want := []string{`{"name":"x"}`, `{"name":"x"}`}; !slices.Equal(bodies, want) {
		t.Errorf("bodies = %q, want the same body twice", bodies)
	}
}

func TestRetryAfter(t *testing.T) {
	now := time.Date(2026, 9, 30, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name, header string
		want         time.Duration
	}{
		{"seconds", "7", 7 * time.Second},
		{"an HTTP date", now.Add(30 * time.Second).Format(http.TimeFormat), 30 * time.Second},
		{"a date in the past", now.Add(-time.Minute).Format(http.TimeFormat), 0},
		{"missing", "", 10 * time.Second},
		{"garbage", "soon", 10 * time.Second},
		{"negative", "-5", 10 * time.Second},
		{"longer than the window is capped", "3600", time.Minute},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := retryAfter(tt.header, now); got != tt.want {
				t.Errorf("retryAfter(%q) = %s, want %s", tt.header, got, tt.want)
			}
		})
	}
}

func TestBackoffGrowsWithJitter(t *testing.T) {
	for attempt, bounds := range [][2]time.Duration{{250 * time.Millisecond, 500 * time.Millisecond}, {500 * time.Millisecond, time.Second}} {
		for range 50 {
			if d := backoff(attempt); d < bounds[0] || d >= bounds[1] {
				t.Fatalf("backoff(%d) = %s, want in [%s, %s)", attempt, d, bounds[0], bounds[1])
			}
		}
	}
}

func TestSleepStopsWhenTheContextEnds(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if err := sleep(ctx, time.Hour); !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want context.Canceled", err)
	}
}

func TestDoReadsAnHTTPDateRetryAfter(t *testing.T) {
	fixed := time.Date(2026, 9, 30, 12, 0, 0, 0, time.UTC)
	var attempts atomic.Int32
	ta := newTestAPI(t, func(w http.ResponseWriter, _ *http.Request) {
		if attempts.Add(1) == 1 {
			w.Header().Set("Retry-After", fixed.Add(5*time.Second).Format(http.TimeFormat))
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		fmt.Fprint(w, `{}`)
	}, nil)
	ta.client.now = func() time.Time { return fixed }
	resp, err := ta.client.Do(t.Context(), Request{Method: http.MethodGet, Path: "/workspace"})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if want := []time.Duration{5 * time.Second}; !slices.Equal(ta.slept, want) {
		t.Errorf("slept %v, want %v", ta.slept, want)
	}
}

func TestDoStopsWaitingForARetryWhenTheContextEnds(t *testing.T) {
	for _, status := range []int{http.StatusTooManyRequests, http.StatusServiceUnavailable} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			var attempts atomic.Int32
			ta := newTestAPI(t, sequence(&attempts, status, status), nil)
			ta.client.sleep = sleep
			ta.client.backoff = func(int) time.Duration { return time.Hour }
			ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
			defer cancel()
			resp, err := ta.client.Do(ctx, Request{Method: http.MethodGet, Path: "/workspace"})
			if err == nil {
				defer resp.Body.Close()
				t.Fatal("Do returned a response instead of stopping")
			}
			if !errors.Is(err, context.DeadlineExceeded) {
				t.Errorf("err = %v, want context.DeadlineExceeded", err)
			}
		})
	}
}

func TestSleepReturnsAfterTheDelay(t *testing.T) {
	if err := sleep(t.Context(), time.Millisecond); err != nil {
		t.Error(err)
	}
}
