package main

import (
	"sync"
	"testing"
)

// Throwaway demo — proves concurrent map WRITES crash even with zero
// readers involved. Delete this file once you've seen it happen.
//
// Run with: go test -run DemoConcurrentMapWrite -v
// (no -race needed — this isn't a "maybe" race, Go's runtime detects
// concurrent map writes unconditionally and panics.)
func TestDemoConcurrentMapWrite(t *testing.T) {
	m := make(map[string]int)
	start := make(chan struct{})

	var wg sync.WaitGroup
	for range 500 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start // every goroutine blocks here until we release them all at once
			for range 1000 {
				m["hits"]++ // no lock at all — only Inc-style writers, no readers
			}
		}()
	}
	close(start) // release all 500 goroutines simultaneously — maximises collisions
	wg.Wait()

	t.Logf("final value: %d", m["hits"])
}
