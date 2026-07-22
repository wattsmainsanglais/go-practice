package main

import (
	"sync"
	"testing"
)

func TestRWCounterSingleKey(t *testing.T) {
	c := NewRWCounter()

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

func TestRWCounterMissingKey(t *testing.T) {
	c := NewRWCounter()

	if got := c.Value("nope"); got != 0 {
		t.Errorf("Value on missing key = %d, want 0", got)
	}
}

func TestRWCounterLen(t *testing.T) {
	c := NewRWCounter()
	c.Inc("a")
	c.Inc("b")
	c.Inc("a")

	if got := c.Len(); got != 2 {
		t.Errorf("Len() = %d, want 2", got)
	}
}

// Concurrent readers and writers at once — the point of this exercise.
// -race must pass here, not just the value assertions.
func TestRWCounterConcurrentReadsAndWrites(t *testing.T) {
	c := NewRWCounter()
	keys := []string{"a", "b", "c"}

	var wg sync.WaitGroup

	// writers
	for _, k := range keys {
		for range 50 {
			wg.Add(1)
			go func(key string) {
				defer wg.Done()
				c.Inc(key)
			}(k)
		}
	}

	// readers, running concurrently with the writers above
	for range 200 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = c.Value("a")
			_ = c.Len()
		}()
	}

	wg.Wait()

	for _, k := range keys {
		if got := c.Value(k); got != 50 {
			t.Errorf("Value(%q) = %d, want 50", k, got)
		}
	}
}
