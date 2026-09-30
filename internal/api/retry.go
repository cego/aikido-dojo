package api

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"
)

const (
	// Retry-After is honoured this many times before giving up with exit 6.
	maxRateRetries = 4
	// A GET is sent at most this many extra times after a 5xx.
	maxServerRetries = 2
	// Aikido's window is a rolling minute, so waiting longer never frees a call sooner.
	maxRetryAfter = time.Minute
	// Used when a 429 carries no usable Retry-After.
	defaultRetryAfter = 10 * time.Second
	backoffBase       = 500 * time.Millisecond
)

// retryAfter reads Retry-After as seconds or an HTTP date.
func retryAfter(v string, now time.Time) time.Duration {
	d := defaultRetryAfter
	if secs, err := strconv.Atoi(v); err == nil && secs >= 0 {
		d = time.Duration(secs) * time.Second
	} else if at, err := http.ParseTime(v); err == nil {
		d = at.Sub(now)
	}
	return min(max(d, 0), maxRetryAfter)
}

// backoff doubles per attempt with equal jitter, so agents that hit the same
// 5xx don't retry in lockstep.
func backoff(attempt int) time.Duration {
	d := backoffBase << attempt
	return d/2 + rand.N(d/2) //nolint:gosec // jitter needs no cryptographic randomness
}

func sleep(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("wait before retrying: %w", context.Cause(ctx))
	}
}
