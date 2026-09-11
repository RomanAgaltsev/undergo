// Package drill — C1/05 waitgroup-negative.
package drill

import "sync"

// RunAll runs each task in its own goroutine and waits for all of them.
func RunAll(tasks []func()) {
	var wg sync.WaitGroup
	for _, t := range tasks {
		go func(fn func()) {
			defer wg.Done()
			fn()
		}(t)
	}
	wg.Wait()
}
