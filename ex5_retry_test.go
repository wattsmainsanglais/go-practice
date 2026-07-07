package main

import (
	"errors"
	"testing"
	"time"
)

func TestRetryWithBackoff_SucceedsFirstTry(t *testing.T) {
	calls := 0
	err := RetryWithBackoff(3, time.Millisecond, func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestRetryWithBackoff_SucceedsAfterRetries(t *testing.T) {
	calls := 0
	err := RetryWithBackoff(5, time.Millisecond, func() error {
		calls++
		if calls < 3 {
			return &RateLimitError{RetryAfterSeconds: 1}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestRetryWithBackoff_GivesUpOnNonRetryable(t *testing.T) {
	calls := 0
	err := RetryWithBackoff(5, time.Millisecond, func() error {
		calls++
		return &NotFoundError{Resource: "x"}
	})
	if calls != 1 {
		t.Fatalf("expected exactly 1 call (no retries on NotFoundError), got %d", calls)
	}
	var nf *NotFoundError
	if !errors.As(err, &nf) {
		t.Fatalf("expected a NotFoundError back, got %v", err)
	}
}

func TestRetryWithBackoff_ExhaustsAttempts(t *testing.T) {
	calls := 0
	err := RetryWithBackoff(3, time.Millisecond, func() error {
		calls++
		return &RateLimitError{RetryAfterSeconds: 1}
	})
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
	if err == nil {
		t.Fatal("expected an error after exhausting all attempts")
	}
}

func TestRetryWithBackoff_ActuallyBacksOff(t *testing.T) {
	// base=10ms, 3 attempts, all fail retryable: sleeps after attempt 1
	// (10ms) and attempt 2 (20ms), none after attempt 3. 30ms minimum.
	start := time.Now()
	RetryWithBackoff(3, 10*time.Millisecond, func() error {
		return &RateLimitError{RetryAfterSeconds: 1}
	})
	elapsed := time.Since(start)
	if elapsed < 25*time.Millisecond {
		t.Fatalf("elapsed %v — doesn't look like it backed off between attempts", elapsed)
	}
}
