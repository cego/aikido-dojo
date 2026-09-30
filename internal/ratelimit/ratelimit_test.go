package ratelimit

import (
	"context"
	"errors"
	"os"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type fakeClock struct {
	t     time.Time
	slept []time.Duration
}

func (c *fakeClock) now() time.Time { return c.t }

func (c *fakeClock) sleep(_ context.Context, d time.Duration) error {
	c.slept = append(c.slept, d)
	c.t = c.t.Add(d)
	return nil
}

func limiter(dir, clientID string, c *fakeClock) *Limiter {
	l := New(dir, clientID)
	l.now, l.sleep = c.now, c.sleep
	return l
}

func TestWaitAllowsTwentyCallsPerMinute(t *testing.T) {
	c := &fakeClock{t: time.Unix(1_800_000_000, 0)}
	l := limiter(t.TempDir(), "AIK_CLIENT_a", c)
	for range 20 {
		if err := l.Wait(t.Context()); err != nil {
			t.Fatal(err)
		}
		c.t = c.t.Add(time.Second)
	}
	if len(c.slept) != 0 {
		t.Fatalf("slept %v during the first 20 calls", c.slept)
	}
	// The 21st call comes 20 s after the first, which leaves the window 40 s later.
	if err := l.Wait(t.Context()); err != nil {
		t.Fatal(err)
	}
	if want := []time.Duration{40 * time.Second}; !slices.Equal(c.slept, want) {
		t.Errorf("slept %v, want %v", c.slept, want)
	}
}

func TestWaitSharesTheWindowAcrossLimiters(t *testing.T) {
	dir := t.TempDir()
	c := &fakeClock{t: time.Unix(1_800_000_000, 0)}
	first, second := limiter(dir, "AIK_CLIENT_a", c), limiter(dir, "AIK_CLIENT_a", c)
	for range 20 {
		if err := first.Wait(t.Context()); err != nil {
			t.Fatal(err)
		}
	}
	if err := second.Wait(t.Context()); err != nil {
		t.Fatal(err)
	}
	if want := []time.Duration{time.Minute}; !slices.Equal(c.slept, want) {
		t.Errorf("slept %v, want %v", c.slept, want)
	}
}

func TestWaitKeepsClientsApart(t *testing.T) {
	dir := t.TempDir()
	c := &fakeClock{t: time.Unix(1_800_000_000, 0)}
	busy, other := limiter(dir, "AIK_CLIENT_a", c), limiter(dir, "AIK_CLIENT_b", c)
	for range 20 {
		if err := busy.Wait(t.Context()); err != nil {
			t.Fatal(err)
		}
	}
	if err := other.Wait(t.Context()); err != nil {
		t.Fatal(err)
	}
	if len(c.slept) != 0 {
		t.Errorf("slept %v, want no wait for another client", c.slept)
	}
}

func TestWaitRebuildsADamagedStateFile(t *testing.T) {
	c := &fakeClock{t: time.Unix(1_800_000_000, 0)}
	l := limiter(t.TempDir(), "AIK_CLIENT_a", c)
	if err := os.WriteFile(l.path, []byte("not json"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := l.Wait(t.Context()); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(l.path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "[1800000000000000000]" {
		t.Errorf("state = %s, want one recorded call", data)
	}
}

func TestWaitStopsWhenTheContextEnds(t *testing.T) {
	l := New(t.TempDir(), "AIK_CLIENT_a")
	fixed := time.Unix(1_800_000_000, 0)
	l.now = func() time.Time { return fixed }
	for range 20 {
		if err := l.Wait(t.Context()); err != nil {
			t.Fatal(err)
		}
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if err := l.Wait(ctx); !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want context.Canceled", err)
	}
}

// Two limiters on one state file stand in for two processes: each has its own
// file descriptor, so only flock keeps their reservations from interleaving.
func TestReserveIsExclusiveAcrossProcesses(t *testing.T) {
	dir := t.TempDir()
	fixed := time.Unix(1_800_000_000, 0)
	var granted atomic.Int32
	var wg sync.WaitGroup
	for range 2 {
		l := New(dir, "AIK_CLIENT_a")
		l.now = func() time.Time { return fixed }
		wg.Go(func() {
			for range 15 {
				wait, err := l.reserve()
				if err != nil {
					t.Error(err)
					return
				}
				if wait == 0 {
					granted.Add(1)
				}
			}
		})
	}
	wg.Wait()
	if got := granted.Load(); got != 20 {
		t.Errorf("granted %d calls in one window, want exactly 20", got)
	}
}
