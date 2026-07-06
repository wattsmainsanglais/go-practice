package main

import (
	"sort"
	"testing"
	"time"
)

func TestProcessJobs_Correctness(t *testing.T) {
	jobs := []int{1, 2, 3, 4, 5}
	got := ProcessJobs(jobs, 2)
	sort.Ints(got)
	want := []int{2, 4, 6, 8, 10}

	if len(got) != len(want) {
		t.Fatalf("got %d results, want %d: %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestProcessJobs_ActuallyConcurrent(t *testing.T) {
	jobs := make([]int, 10) // 10 jobs, each takes 50ms via doWork
	for i := range jobs {
		jobs[i] = i
	}

	start := time.Now()
	ProcessJobs(jobs, 5) // 5 workers
	elapsed := time.Since(start)

	// Sequential would take ~500ms (10 jobs * 50ms each). With 5 workers
	// running concurrently, it should take roughly 2 batches: ~100ms.
	// 300ms is a generous upper bound to avoid flakiness on a slow machine —
	// if your implementation is secretly sequential (e.g. one worker doing
	// everything), this will fail loudly.
	if elapsed > 300*time.Millisecond {
		t.Fatalf("took %v — doesn't look concurrent (expected well under 300ms)", elapsed)
	}
}
