package main

import (
	"errors"
	"fmt"
	"testing"
)

func TestIsRetryable_Nil(t *testing.T) {
	if IsRetryable(nil) {
		t.Fatal("nil error should not be retryable")
	}
}

func TestIsRetryable_NotFound(t *testing.T) {
	err := &NotFoundError{Resource: "course/123"}
	if IsRetryable(err) {
		t.Fatal("NotFoundError should not be retryable")
	}
}

func TestIsRetryable_RateLimit(t *testing.T) {
	err := &RateLimitError{RetryAfterSeconds: 5}
	if !IsRetryable(err) {
		t.Fatal("RateLimitError should be retryable")
	}
}

func TestIsRetryable_WrappedRateLimit(t *testing.T) {
	inner := &RateLimitError{RetryAfterSeconds: 2}
	wrapped := fmt.Errorf("posting task: %w", inner)
	if !IsRetryable(wrapped) {
		t.Fatal("wrapped RateLimitError should still be retryable via errors.As")
	}
}

func TestIsRetryable_OtherError(t *testing.T) {
	if IsRetryable(errors.New("something else")) {
		t.Fatal("a plain error should not be retryable")
	}
}
