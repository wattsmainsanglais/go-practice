package main

import (
	"context"
	"time"
)

// Exercise 3: context + select
//
// Concepts:
//   - context.Context carries a deadline/cancellation signal through a call
//     chain. context.WithTimeout(parent, d) returns a ctx that automatically
//     fires after duration d, plus a cancel() func you should always defer.
//   - ctx.Done() returns a channel that's closed when the context expires or
//     is cancelled. ctx.Err() tells you why (context.DeadlineExceeded, etc).
//   - `select` waits on multiple channel operations at once and proceeds with
//     whichever is ready first. It's how you race "did the work finish" against
//     "did the deadline expire" without polling.
//
// slowWork is provided — don't modify it. Same idea as doWork in Exercise 2,
// just a different name so the two exercises don't collide.
func slowWork(job int) int {
	time.Sleep(50 * time.Millisecond)
	return job * 3
}

// YOUR TASK: delete the panic() line and implement FetchWithTimeout so that:
//   - it runs slowWork(job) concurrently (in a goroutine)
//   - if slowWork finishes before ctx's deadline, return (result, nil)
//   - if ctx's deadline fires first, return (0, ctx.Err()) — and don't block
//     forever waiting on slowWork in that case
//
// Suggested structure:
//
//	STEP 1: create a buffered result channel: `resultChan := make(chan int, 1)`
//	        Buffered with capacity 1 matters here: if the timeout wins the
//	        select below, nothing will ever receive from resultChan again.
//	        Without buffering, the goroutine in STEP 2 would block forever on
//	        its send — a goroutine leak. With buffer 1, the send always
//	        succeeds immediately and the goroutine exits cleanly either way.

//	STEP 2: start a goroutine that computes `slowWork(job)` and sends the
//	        result into resultChan.
//
//	STEP 3: `select` on two cases:
//	          - receiving from resultChan -> return (result, nil)
//	          - receiving from ctx.Done() -> return (0, ctx.Err())
//
// Try it, and we'll walk through it line by line once you've got something
// written — same as last time.
func FetchWithTimeout(ctx context.Context, job int) (int, error) {
	resultChan := make(chan int, 1)

	go func() {
		resultChan <- slowWork(job)

	}()

	select {
	case result := <-resultChan:
		return result, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}
