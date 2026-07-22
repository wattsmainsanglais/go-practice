package main

import "sync"

// Exercise 9: Fan-in (N producers -> one channel)
//
// You've done fan-out (ex1: one goroutine per input) and a worker pool
// (ex2: N workers pulling from ONE shared channel). This is the reverse
// direction: N *separate* channels, each fed by its own producer, and you
// need to merge everything into a single output channel a consumer can
// range over.
//
// Concepts:
//   - A WaitGroup only tells you when a set of goroutines is DONE — it
//     doesn't merge values by itself. You still need somewhere for the
//     values to go while producers are running.
//   - range over a channel keeps receiving until the channel is closed.
//     If nobody ever closes the output channel, a consumer ranging over
//     it blocks forever, even after every producer has finished.
//   - Closing a channel that something might still send on panics. So
//     whoever closes the output channel must be certain no more sends
//     are coming — that means waiting for ALL producer goroutines to
//     finish first.
//
// YOUR TASK: FanIn
//
// Given a variable number of input channels, return ONE output channel.
// Every value sent on any input channel should eventually appear on the
// output channel (order across producers doesn't matter). Once every
// input channel has been closed AND drained, close the output channel.
//
//	STEP 1: create the output channel.
//	STEP 2: for each input channel, launch a goroutine that ranges over
//	        it, sending everything it receives to the output channel.
//	        Use a WaitGroup to track these goroutines.
//	STEP 3: launch ONE more goroutine that waits on the WaitGroup, then
//	        closes the output channel.
//	        Question to think about before you write this: why must this
//	        run in its own goroutine, rather than you calling wg.Wait()
//	        and close() directly in FanIn before it returns?
//	STEP 4: return the output channel immediately — don't block in FanIn
//	        itself waiting for anything.
func FanIn(channels ...<-chan int) <-chan int {
	output := make(chan int, len(channels))

	var wg sync.WaitGroup

	for _, n := range channels {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for w := range n {
				output <- w
			}
		}()

	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}
