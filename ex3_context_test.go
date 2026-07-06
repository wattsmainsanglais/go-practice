package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestFetchWithTimeout_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	got, err := FetchWithTimeout(ctx, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 12 {
		t.Fatalf("got %d, want 12", got)
	}
}

func TestFetchWithTimeout_TimesOut(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := FetchWithTimeout(ctx, 4)
	elapsed := time.Since(start)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("got err %v, want context.DeadlineExceeded", err)
	}
	// Should return right around the 10ms deadline, not wait for slowWork's
	// full 50ms — if this takes ~50ms, the select isn't actually racing.
	if elapsed > 30*time.Millisecond {
		t.Fatalf("took %v to return — should return near the 10ms deadline, not wait on slowWork", elapsed)
	}
}
