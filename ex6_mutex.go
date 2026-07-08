package main

import "sync"

// Exercise 6: sync.Mutex
//
// Every exercise so far has synchronized goroutines by passing data through
// a channel (WaitGroup + channel, or the buffered-channel semaphore in
// Exercise 2). A Mutex is a different tool for a different situation: it
// doesn't pass anything between goroutines, it just makes sure only ONE
// goroutine at a time can read/modify some shared piece of memory. You
// already used one indirectly — TaskStore in http-challenge/store.go guards
// its map with a sync.Mutex, but that store was given to you fully written.
// This time you write the locking yourself.
//
// Concepts:
//   - mu.Lock() / mu.Unlock() bracket a "critical section" — code that
//     touches shared state. Anything between Lock and Unlock runs as if
//     no other goroutine exists; Go halts other goroutines' Lock() calls
//     until Unlock() happens.
//   - defer mu.Unlock() right after mu.Lock() is the idiomatic pattern —
//     guarantees the unlock happens even if the function returns early or
//     panics, same reasoning as defer wg.Done() in Exercise 2.
//   - Rule of thumb for when to reach for Mutex instead of a channel: a
//     channel is for handing off a value or a signal between goroutines
//     (a pipeline, a job queue, "I'm done"). A Mutex is for protecting a
//     piece of state that many goroutines read AND write directly, where
//     there's no natural "hand-off" — a shared map or counter is the
//     classic case.
//
// SafeCounter tracks hit counts per key (imagine counting API requests per
// endpoint) from many goroutines at once.
type SafeCounter struct {
	mu     sync.Mutex
	counts map[string]int
}

func NewSafeCounter() *SafeCounter {
	return &SafeCounter{counts: make(map[string]int)}
}

// YOUR TASK 1: Inc
//
// Increment counts[key] by 1. Must be safe to call from many goroutines
// at once for the same key AND for different keys.
//
//	STEP 1: lock the mutex, defer the unlock.
//	STEP 2: c.counts[key]++
func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counts[key]++
}

// YOUR TASK 2: Value
//
// Return the current count for key (0 if never incremented — that's the
// zero value of int, which is what you get reading a missing map key).
//
// Same lock/defer-unlock pattern as Inc. Question to think about before
// you write this: Value only reads, it doesn't write — does it still need
// the lock? What could go wrong if a read ran with no lock while another
// goroutine's Inc was mid-increment on the same map?
func (c *SafeCounter) Value(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.counts[key]
}

// YOUR TASK 3: Snapshot
//
// Return a copy of the entire counts map — a plain map[string]int, safe
// for the caller to read/range over after this function returns even if
// more Inc calls happen concurrently on the original.
//
//	STEP 1: lock, defer unlock.
//	STEP 2: make a new map, copy every key/value from c.counts into it.
//	STEP 3: return the copy.
//
// Why a copy, and not just returning c.counts directly? Think about what
// the caller could do with the map you hand back, and when your Unlock
// actually runs relative to that.
func (c *SafeCounter) Snapshot() map[string]int {

	c.mu.Lock()
	defer c.mu.Unlock()

	newCounter := make(map[string]int)
	for key, val := range c.counts {
		newCounter[key] = val
	}

	return newCounter

}
