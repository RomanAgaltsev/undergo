// Package drill — C1/03 worker-pipeline.
package drill

import "sync"

// Process fans jobs out across n workers, doubles each, and collects the results.
func Process(jobs []int, n int) []int {
	var wg sync.WaitGroup
	jobCh := make(chan int)
	results := make([]int, 0, len(jobs))

	for i := 0; i < n; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			for j := range jobCh {
				results = append(results, j*2)
			}
		}()
	}

	for _, j := range jobs {
		jobCh <- j
	}
	close(jobCh)
	wg.Wait()
	return results
}

// drain consumes events until the channel is closed, then signals done.
func drain(events <-chan int, done chan<- struct{}) {
	for range events {
	}
	done <- struct{}{}
}

// Warm starts a background drainer for a warm-up stream, then returns once the
// stream is closed.
func Warm(stream []int) {
	events := make(chan int)
	done := make(chan struct{})
	go drain(events, done)

	if len(stream) == 0 {
		return
	}
	for _, e := range stream {
		events <- e
	}
	close(events)
	<-done
}
