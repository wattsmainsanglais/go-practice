package main

import "time"

// Exercise 5: retry with backoff
//
// Builds directly on Exercise 4's IsRetryable. Concepts:
//   - Exponential backoff: wait longer after each failed attempt so you
//     don't hammer a struggling server. Doubling is simplest: attempt 1
//     waits `base`, attempt 2 waits `base*2`, attempt 3 waits `base*4`.
//   - `base << n` (bit shift) computes base * 2^n — a fast way to write
//     exponential growth. The andamio-cli retry.go you looked at does
//     exactly this: `cfg.InitialBackoff << (attempt - 1)`.
//   - A retry loop needs two exit conditions, not one: success (return
//     immediately), and "this specific error will never succeed no matter
//     how many times we try" (also return immediately — don't burn
//     attempts sleeping between retries of something that can't work).
//
// YOUR TASK: delete the panic() line and implement RetryWithBackoff so
// that:
//   - it calls fn() up to maxAttempts times
//   - if fn() returns nil, return nil immediately (success)
//   - if fn() returns an error where IsRetryable(err) is false, return
//     that error immediately — don't retry it, don't sleep
//   - if fn() returns a retryable error and attempts remain, sleep then
//     try again. Sleep duration doubles each time: base, base*2, base*4...
//   - if fn() returns a retryable error on the LAST attempt, return that
//     error without sleeping afterward (nothing left to wait for)
//
// Suggested structure:
//
//	STEP 1: loop `for attempt := 0; attempt < maxAttempts; attempt++`
//
//	STEP 2: inside the loop, call fn() and store the result in `err`.
//
//	STEP 3: if err == nil, return nil.
//
//	STEP 4: if !IsRetryable(err), return err.
//
//	STEP 5: if this is the last attempt (attempt == maxAttempts-1),
//	        return err.
//
//	STEP 6: otherwise, sleep for `base << attempt` (a time.Duration), then
//	        let the loop continue to the next attempt.
//
//	STEP 7: after the loop — Go needs a return statement here even though
//	        it's never reached in practice. Return nil.
func RetryWithBackoff(maxAttempts int, base time.Duration, fn func() error) error {

	for attempt := 0; attempt < maxAttempts; attempt++ {

		result := fn()

		if result == nil {
			return nil
		}

		if !IsRetryable(result) {
			return result
		}

		if IsRetryable(result) {
			if maxAttempts-attempt == 1 {
				return result
			}

			time.Sleep(base * time.Duration(2^attempt))
		}

	}

	return nil
}
