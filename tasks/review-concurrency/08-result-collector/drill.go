// Package drill — C1/08 result-collector.
package drill

import "sync"

// Collect runs fn over every input concurrently and returns the results. It
// signals on `out` the first time fn yields a zero, so callers can react.
func Collect(inputs []int, fn func(int) int) []int {
	var wg sync.WaitGroup
	results := make([]int, 0, len(inputs))
	out := make(chan struct{})

	for _, in := range inputs {
		go func(x int) {
			wg.Add(1)
			defer wg.Done()
			r := fn(x)
			results = append(results, r)
			if r == 0 {
				close(out)
			}
		}(in)
	}

	wg.Wait()
	select {
	case <-out:
	default:
	}
	return results
}
