package main

import (
	"testing"
	"time"
)

// Throwaway demo — not part of the exercise, delete when done.
//
// Question: does select+default on an UNBUFFERED channel ever succeed the
// receive case, or does it always fall to default since there's no buffer
// to "store" a value?

func TestUnbufferedDefault_SenderAlreadyWaiting(t *testing.T) {
	ch := make(chan int) // unbuffered — zero capacity, no storage at all

	go func() {
		ch <- 99 // blocks here until someone receives
	}()

	time.Sleep(50 * time.Millisecond) // let the sender actually reach the blocked send

	v, ok := TryReceive(ch)
	t.Logf("with a sender already blocked and waiting: v=%d ok=%v", v, ok)
	if !ok {
		t.Error("expected ok=true — a sender was ready to hand off right now")
	}
}

func TestUnbufferedDefault_NoSenderWaiting(t *testing.T) {
	ch := make(chan int) // unbuffered, and nobody is trying to send at all

	v, ok := TryReceive(ch)
	t.Logf("with nobody sending: v=%d ok=%v", v, ok)
	if ok {
		t.Error("expected ok=false — nothing was ready")
	}
}
