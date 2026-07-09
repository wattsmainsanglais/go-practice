package main

import (
	"fmt"
	"time"
)

// Exercise 7: select
//
// Every exercise so far has worked with one channel at a time (or a
// WaitGroup fanning into one). select lets a goroutine wait on MULTIPLE
// channel operations at once, and proceed with whichever one is ready
// first. It reads like a switch, but each case is a channel send or
// receive, not a value comparison.
//
// Concepts:
//   - select { case v := <-a: ... case v := <-b: ... } blocks until EITHER
//     a or b has something ready, then runs that one case. If both are
//     ready at the same instant, Go picks one at random — never assume
//     ordering between cases.
//   - time.After(d) returns a <-chan time.Time that fires once, after d has
//     elapsed. Putting it in a select case is the standard way to add a
//     timeout to a wait without needing context at all.
//   - select with a `default:` case never blocks — if no other case is
//     ready RIGHT NOW, it runs default instead of waiting. That's how you
//     do a non-blocking channel check.

// YOUR TASK 1: FirstString
//
// Two goroutines (already running, you don't write them) are each trying
// to produce a string on their own channel — think of it as racing two
// API mirrors and taking whichever answers first. Wait on BOTH a and b at
// once. Return whichever value arrives first. If neither arrives within
// timeout, return an error instead.
//
//	STEP 1: build a select with three cases: receive from a, receive from
//	        b, and receive from time.After(timeout).
//	STEP 2: whichever case fires, that's your return. For the two channel
//	        cases return (value, nil). For the timeout case return
//	        ("", an error).
//
// Question to think about before you write this: once select picks a
// case and this function returns, what happens to the OTHER channel if it
// produces a value later? Is anything listening?
func FirstString(a, b <-chan string, timeout time.Duration) (string, error) {

	var result string

	select {
	case result = <-a:

		return result, nil
	case result = <-b:

		return result, nil
	case <-time.After(timeout):
		return "", fmt.Errorf("timer elapsed")
	}
}

// YOUR TASK 2: TryReceive
//
// Given a channel that some other goroutine may or may not have put a
// value on yet, check ONCE, right now, without blocking. If a value was
// ready, return (value, true). If nothing was ready, return (0, false)
// immediately — do not wait around for one to show up.
//
//	STEP 1: select with one case receiving from ch, and a default case.
//	STEP 2: channel case returns (v, true). default case returns (0, false).
func TryReceive(ch <-chan int) (int, bool) {

	var result int
	select {
	case result = <-ch:
		return result, true
	default:
		return 0, false

	}
}
