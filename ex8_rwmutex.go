package main

import "sync"

// Exercise 8: sync.RWMutex
//
// Same shape of problem as Exercise 6 (SafeCounter) — many goroutines
// touching a shared map — but this time most callers only ever read, and
// writes are rare. A plain sync.Mutex forces every reader to wait behind
// every other reader too, even though two reads can never corrupt each
// other. sync.RWMutex splits the lock in two:
//
//   - RLock() / RUnlock() — a "read lock". Any number of goroutines can
//     hold this at the same time, as long as nobody holds the write lock.
//   - Lock() / Unlock() — a "write lock". Exclusive, same as a plain Mutex:
//     while one goroutine holds it, nobody else (reader or writer) can
//     proceed.
//
// Rule of thumb: RWMutex only pays for itself when reads vastly outnumber
// writes. If writes are frequent, the extra bookkeeping RWMutex does makes
// it slower than a plain Mutex — don't reach for it by default.
//
// The bug to watch for: using RLock() in a method that WRITES to the map.
// Two goroutines could both hold the read lock and both mutate the map at
// the same instant — a data race, exactly what the lock exists to prevent.
// Think through which lock each method below needs BEFORE you write it,
// don't just copy Inc's pattern into all three.
//
// RWCounter tracks hit counts per key — identical job to ex6's SafeCounter,
// different lock.
type RWCounter struct {
	mu     sync.RWMutex
	counts map[string]int
}

func NewRWCounter() *RWCounter {
	return &RWCounter{counts: make(map[string]int)}
}

// YOUR TASK 1: Inc
//
// Increment counts[key] by 1. This WRITES to the map.
func (c *RWCounter) Inc(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counts[key]++
}

// YOUR TASK 2: Value
//
// Return counts[key] (0 if the key has never been seen). This only READS.
func (c *RWCounter) Value(key string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.counts[key]
}

// YOUR TASK 3: Len
//
// Return the number of distinct keys currently tracked. Also read-only.
func (c *RWCounter) Len() int {
	// TODO
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.counts)
}
