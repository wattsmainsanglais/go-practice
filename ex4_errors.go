package main

import "errors"

// Exercise 4: custom error types + errors.As
//
// Concepts:
//   - A custom error type is just a struct with an `Error() string` method
//     — that's the entire interface. Once it satisfies `error`, it works
//     anywhere an error is expected.
//   - errors.As(err, &target) walks an error's Unwrap() chain looking for
//     an error whose concrete type matches target (target is a pointer to
//     a pointer, e.g. `var nf *NotFoundError; errors.As(err, &nf)`). It's
//     how you recover a specific error type after it's been wrapped with
//     fmt.Errorf("...: %w", err) one or more times.
//   - This is exactly how a real CLI decides "should I retry this HTTP
//     call, or give up?" — see internal/apierr/errors.go and
//     internal/client/retry.go in the andamio-cli repo if you want to see
//     the production version this exercise (and the next one) is modeled
//     on. isRetryable() there does the same errors.As dance you're about
//     to write, just against more error types.
//
// Two error types are provided below — don't modify them.

// NotFoundError means the resource doesn't exist. Never worth retrying.
type NotFoundError struct {
	Resource string
}

func (e *NotFoundError) Error() string {
	return "not found: " + e.Resource
}

// RateLimitError means the caller is being throttled. Worth retrying.
type RateLimitError struct {
	RetryAfterSeconds int
}

func (e *RateLimitError) Error() string {
	return "rate limited"
}

// YOUR TASK: delete the panic() line and implement IsRetryable so that:
//   - it returns false for a nil error
//   - it returns false if err is (or wraps) a *NotFoundError
//   - it returns true if err is (or wraps) a *RateLimitError
//   - it returns false for any other error
//
// Suggested structure:
//
//	STEP 1: handle the nil case first — return false immediately.
//
//	STEP 2: declare `var nf *NotFoundError` and check `errors.As(err, &nf)`.
//	        If true, return false.
//
//	STEP 3: declare `var rl *RateLimitError` and check `errors.As(err, &rl)`.
//	        If true, return true.
//
//	STEP 4: anything else — return false.
//
// Run the tests once you've got something written. One of them wraps a
// RateLimitError with fmt.Errorf("...: %w", err) before passing it in —
// errors.As has to see through that wrap for the test to pass.
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	var nf *NotFoundError

	if errors.As(err, &nf) {
		return false
	}

	var rl *RateLimitError

	if errors.As(err, &rl) {
		return true
	}

	return false

}
