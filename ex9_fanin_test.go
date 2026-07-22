package main

import (
	"testing"
	"time"
)

func TestFanInMergesAllValues(t *testing.T) {
	a := make(chan int)
	b := make(chan int)
	c := make(chan int)

	go func() {
		for _, v := range []int{1, 2, 3} {
			a <- v
		}
		close(a)
	}()
	go func() {
		for _, v := range []int{4, 5} {
			b <- v
		}
		close(b)
	}()
	go func() {
		c <- 6
		close(c)
	}()

	out := FanIn(a, b, c)

	seen := make(map[int]bool)
	for v := range out {
		seen[v] = true
	}

	want := []int{1, 2, 3, 4, 5, 6}
	for _, v := range want {
		if !seen[v] {
			t.Errorf("missing value %d in output", v)
		}
	}
	if len(seen) != len(want) {
		t.Errorf("got %d unique values, want %d", len(seen), len(want))
	}
}

func TestFanInClosesOutput(t *testing.T) {
	a := make(chan int)
	close(a)

	out := FanIn(a)

	select {
	case _, ok := <-out:
		if ok {
			t.Fatal("expected output channel to be empty and closed")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("output channel was never closed")
	}
}

func TestFanInNoChannels(t *testing.T) {
	out := FanIn()

	select {
	case _, ok := <-out:
		if ok {
			t.Fatal("expected empty output for zero input channels")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("output channel was never closed with zero inputs")
	}
}
