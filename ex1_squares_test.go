package main

import (
	"sort"
	"testing"
)

func TestSquares(t *testing.T) {
	got := Squares([]int{1, 2, 3, 4, 5})
	sort.Ints(got)
	want := []int{1, 4, 9, 16, 25}

	if len(got) != len(want) {
		t.Fatalf("got %d results, want %d: %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestSquaresEmpty(t *testing.T) {
	got := Squares(nil)
	if len(got) != 0 {
		t.Fatalf("got %v, want empty", got)
	}
}
