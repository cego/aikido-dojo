// Package ratelimit keeps calls under Aikido's per-workspace limit across
// every process on the machine. Agents run the CLI as many short processes,
// so a limiter inside one process would never engage.
package ratelimit

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"syscall"
	"time"
)

// Aikido allows 20 calls per rolling minute per workspace (rate-limiting.md, checked 2026-09-29).
const (
	limit  = 20
	window = time.Minute
)

type Limiter struct {
	path  string
	now   func() time.Time
	sleep func(context.Context, time.Duration) error
}

// New returns the limiter for one API client. The state file is keyed by a
// hash of the client ID, so every profile and process using that client
// shares one window.
func New(dir, clientID string) *Limiter {
	sum := sha256.Sum256([]byte(clientID))
	return &Limiter{
		path:  filepath.Join(dir, "ratelimit-"+hex.EncodeToString(sum[:8])+".json"),
		now:   time.Now,
		sleep: sleep,
	}
}

// Wait blocks until a call fits in the window, then records it.
func (l *Limiter) Wait(ctx context.Context) error {
	for {
		wait, err := l.reserve()
		if err != nil {
			return err
		}
		if wait <= 0 {
			return nil
		}
		if err := l.sleep(ctx, wait); err != nil {
			return err
		}
	}
}

// reserve records a call and returns 0 when one fits in the window, or how
// long until the oldest call leaves it.
func (l *Limiter) reserve() (time.Duration, error) {
	if err := os.MkdirAll(filepath.Dir(l.path), 0o700); err != nil {
		return 0, fmt.Errorf("create the rate-limit directory: %w", err)
	}
	f, err := os.OpenFile(l.path, os.O_RDWR|os.O_CREATE, 0o600)
	if err != nil {
		return 0, fmt.Errorf("open the rate-limit state: %w", err)
	}
	defer f.Close() // closing also releases the lock
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		return 0, fmt.Errorf("lock the rate-limit state: %w", err)
	}
	calls, err := readCalls(f)
	if err != nil {
		return 0, err
	}
	now := l.now()
	calls = slices.DeleteFunc(calls, func(t int64) bool { return now.Sub(time.Unix(0, t)) >= window })
	for i, t := range calls {
		// A call dated after now means the clock stepped back; counting it as
		// made now keeps every wait within one window.
		calls[i] = min(t, now.UnixNano())
	}
	if len(calls) >= limit {
		// Saved even when full, so the clamped times age out instead of being clamped again.
		return time.Unix(0, calls[0]).Add(window).Sub(now), writeCalls(f, calls)
	}
	calls = append(calls, now.UnixNano())
	slices.Sort(calls)
	return 0, writeCalls(f, calls)
}

func readCalls(f *os.File) ([]int64, error) {
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("read the rate-limit state: %w", err)
	}
	var calls []int64
	if err := json.Unmarshal(data, &calls); err != nil {
		// The file only paces calls, so a damaged one is rebuilt rather than blocking every call.
		calls = nil
	}
	return calls, nil
}

func writeCalls(f *os.File, calls []int64) error {
	data, _ := json.Marshal(calls) // a []int64 always encodes
	if err := f.Truncate(0); err != nil {
		return fmt.Errorf("reset the rate-limit state: %w", err)
	}
	if _, err := f.WriteAt(data, 0); err != nil {
		return fmt.Errorf("write the rate-limit state: %w", err)
	}
	return nil
}

func sleep(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("wait for the rate limit: %w", context.Cause(ctx))
	}
}
