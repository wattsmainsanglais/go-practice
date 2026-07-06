package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Squares:", Squares([]int{1, 2, 3, 4, 5}))

	jobs := make([]int, 10)
	for i := range jobs {
		jobs[i] = i
	}

	start := time.Now()
	results := ProcessJobs(jobs, 5)
	elapsed := time.Since(start)

	fmt.Println("ProcessJobs results:", results)
	fmt.Println("ProcessJobs took:", elapsed)
}
