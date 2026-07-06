package main

import (
	"sync"
	"time"
)

// Exercise 2: worker pool
//
// Exercise 1 spawned one goroutine PER INPUT — fine for 5 items, a bad idea
// for 5 million (or if each item hits a database with a limited connection
// pool). A worker pool instead uses a FIXED number of goroutines that pull
// work from a shared job channel, so concurrency is capped regardless of how
// many jobs there are — same idea as the semaphore we discussed, structured
// slightly differently.
//
// doWork is provided — don't modify it. It simulates something slow (an API
// call, a DB query) by sleeping.
func doWork(job int) int {
	time.Sleep(50 * time.Millisecond)
	return job * 2
}

// YOUR TASK: implement ProcessJobs so that:
//   - every job in `jobs` gets passed to doWork(), and all results are
//     collected and returned (order doesn't matter)
//   - exactly `workers` goroutines run concurrently to do this — not one
//     goroutine per job like Exercise 1
//
// Suggested structure:
//
//	STEP 1: create a channel for jobs: `jobsChan := make(chan int, len(jobs))`
//	        then loop over `jobs`, sending each one into jobsChan, then
//	        close(jobsChan). (Buffering it to len(jobs) means all the sends
//	        succeed immediately without needing a separate goroutine for
//	        this part — feel free to ask if that's unclear.)
//
//	STEP 2: create a channel for results: `resultsChan := make(chan int)`
//
//	STEP 3: create a WaitGroup. Start exactly `workers` goroutines (a loop
//	        from 0 to workers, NOT from 0 to len(jobs)). Each one:
//	          - wg.Add(1) before starting it
//	          - defer wg.Done() inside it
//	          - loops with `for job := range jobsChan { ... }`, calling
//	            doWork(job) and sending the result into resultsChan
//	        Each worker's range loop naturally ends once jobsChan is closed
//	        AND drained — that's why STEP 1 closes it after loading all jobs.
//
//	STEP 4: same trick as Exercise 1 — one more goroutine that does
//	        wg.Wait() then close(resultsChan), so the collection loop below
//	        knows when to stop.
//
//	STEP 5: collect with `for r := range resultsChan { ... }`, return the
//	        slice.
func ProcessJobs(jobs []int, workers int) []int {
	jobsChan := make(chan int, len(jobs))

	for _, j := range jobs {
		jobsChan <- j

	}

	close(jobsChan)

	resultsChan := make(chan int)
	var wg sync.WaitGroup

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobsChan {
				result := doWork(job)
				resultsChan <- result
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var results []int

	for r := range resultsChan {
		results = append(results, r)
	}
	return results
}
