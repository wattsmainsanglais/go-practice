package main

import "sync"

// Exercise 1: goroutines + channels basics
//
// Concepts:
//   - `go f()` starts f running concurrently, doesn't block.
//   - a channel (`chan int`) is a typed pipe: one goroutine sends with `ch <- v`,
//     another receives with `v := <-ch`. Receiving blocks until something is sent.
//   - `sync.WaitGroup` lets you wait for a set of goroutines to finish:
//     wg.Add(1) before starting one, wg.Done() when it finishes (usually via defer),
//     wg.Wait() blocks until all Done() calls have happened.
//
// YOUR TASK: delete the panic() line below and write real code, following
// these steps in order. Nothing here is filled in for you.
//
//	STEP 1: create a channel of int: `ch := make(chan int)`
//
//	STEP 2: create a WaitGroup: `var wg sync.WaitGroup`
//
//	STEP 3: loop over `nums`. For each `n`:
//	          - call `wg.Add(1)`
//	          - start a goroutine (`go func() { ... }()`) that computes n*n,
//	            sends it on ch, then calls `wg.Done()` (use `defer wg.Done()`
//	            as the first line inside the goroutine)
//
//	STEP 4: start ONE more goroutine that calls `wg.Wait()` and then
//	        `close(ch)`. This has to run concurrently with STEP 5 below —
//	        if you call wg.Wait() directly in the main function body before
//	        reading from ch, you'll deadlock (nothing is reading the channel
//	        yet, so the sends in STEP 3 can never complete, so Wait() never
//	        returns).
//
//	STEP 5: collect results: `for v := range ch { results = append(results, v) }`
//	        This loop keeps pulling values until the channel is closed.
//
//	STEP 6: return the collected slice.
//
// Gotcha to watch for: closures over loop variables. In Go 1.22+ (you're on
// 1.24) each loop iteration gets its own copy of `n`, so `n` is safe to use
// directly inside the goroutine closure in STEP 3. In older Go you'd need to
// pass it as a parameter instead.
func Squares(nums []int) []int {
	ch := make(chan int)

	var wg sync.WaitGroup

	for _, n := range nums {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := n * n
			ch <- result

		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var results []int

	for v := range ch {
		results = append(results, v)

	}
	return results

}
