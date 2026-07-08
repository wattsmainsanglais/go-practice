package main

import (
	"sync"
	"testing"
)

func TestSafeCounterSingleKey(t *testing.T) {
	c := NewSafeCounter()

	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Inc("hits")
		}()
	}
	wg.Wait()

	if got := c.Value("hits"); got != 100 {
		t.Errorf("Value(%q) = %d, want 100", "hits", got)
	}
}

func TestSafeCounterMissingKey(t *testing.T) {
	c := NewSafeCounter()

	if got := c.Value("nope"); got != 0 {
		t.Errorf("Value on missing key = %d, want 0", got)
	}
}

func TestSafeCounterMultipleKeys(t *testing.T) {
	c := NewSafeCounter()
	keys := []string{"a", "b", "c"}

	var wg sync.WaitGroup
	for _, k := range keys {
		for range 50 {
			wg.Add(1)
			go func(key string) {
				defer wg.Done()
				c.Inc(key)
			}(k)
		}
	}
	wg.Wait()

	snap := c.Snapshot()
	for _, k := range keys {
		if snap[k] != 50 {
			t.Errorf("Snapshot()[%q] = %d, want 50", k, snap[k])
		}
	}
	if len(snap) != len(keys) {
		t.Errorf("Snapshot() has %d keys, want %d", len(snap), len(keys))
	}
}

func TestSafeCounterSnapshotIsIndependentCopy(t *testing.T) {
	c := NewSafeCounter()
	c.Inc("x")

	snap := c.Snapshot()
	snap["x"] = 999 // mutating the returned map must not affect the counter

	if got := c.Value("x"); got != 1 {
		t.Errorf("Value(%q) after mutating snapshot = %d, want 1 (Snapshot should return a copy)", "x", got)
	}
}
