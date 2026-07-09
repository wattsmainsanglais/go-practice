package main

import (
	"testing"
	"time"
)

func TestFirstStringFastWins(t *testing.T) {
	a := make(chan string, 1)
	b := make(chan string, 1)

	go func() {
		time.Sleep(10 * time.Millisecond)
		a <- "from a"
	}()
	go func() {
		time.Sleep(50 * time.Millisecond)
		b <- "from b"
	}()

	got, err := FirstString(a, b, 200*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "from a" {
		t.Errorf("FirstString = %q, want %q", got, "from a")
	}
}

func TestFirstStringTimeout(t *testing.T) {
	a := make(chan string)
	b := make(chan string)

	_, err := FirstString(a, b, 20*time.Millisecond)
	if err == nil {
		t.Fatal("expected an error on timeout, got nil")
	}
}

func TestTryReceiveHasValue(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 42

	v, ok := TryReceive(ch)
	if !ok {
		t.Fatal("TryReceive ok = false, want true")
	}
	if v != 42 {
		t.Errorf("TryReceive value = %d, want 42", v)
	}
}

func TestTryReceiveEmpty(t *testing.T) {
	ch := make(chan int)

	start := time.Now()
	_, ok := TryReceive(ch)
	elapsed := time.Since(start)

	if ok {
		t.Fatal("TryReceive ok = true on empty channel, want false")
	}
	if elapsed > 5*time.Millisecond {
		t.Errorf("TryReceive took %v — should return immediately, not block", elapsed)
	}
}
