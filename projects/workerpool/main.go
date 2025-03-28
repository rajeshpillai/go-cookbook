// Question 1: Implementing a Worker Pool
// Implement a worker pool in Go. The pool should receive jobs through a channel,
// process each job concurrently, and send the results back through another
// channel. Make sure to gracefully shut down the workers when all jobs have
// been processed."

package main

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID    int
	Value int
}

type Result struct {
	JobID  int
	Output int
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		// Simulate work with a sleep
		time.Sleep(100 * time.Millisecond)
		result := Result{JobID: job.ID, Output: job.Value * 2}
		fmt.Printf("Worker %d processed job %d\n", id, job.ID)
		results <- result
	}
}

func main() {
	jobs := make(chan Job, 10)
	results := make(chan Result, 10)

	var wg sync.WaitGroup
	numWorkers := 3

	// Start the worker goroutines
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Dispatch jobs
	numJobs := 5
	for j := 1; j <= numJobs; j++ {
		jobs <- Job{ID: j, Value: j * 10}
	}
	close(jobs)

	// Close results channel once all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Retrieve and print results
	for res := range results {
		fmt.Printf("Result: JobID %d, Output %d\n", res.JobID, res.Output)
	}
}
