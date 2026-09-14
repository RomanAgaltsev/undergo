// Package runnext observes the order in which goroutines spawned from a single
// goroutine actually run, on a single P, with no preemption points between the
// spawns.
package runnext

import (
	"runtime"
	"sync"
)

// SpawnOrder starts n goroutines numbered 1..n from one goroutine under
// GOMAXPROCS=1, and returns the order in which they recorded themselves.
//
// The spawner does not yield between the `go` statements, and none of the
// spawned functions blocks before recording. So the order returned is the
// order the scheduler chose, not a race between them.
func SpawnOrder(n int) []int {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(1))

	var mu sync.Mutex
	order := make([]int, 0, n)

	var wg sync.WaitGroup
	wg.Add(n)
	for i := 1; i <= n; i++ {
		go func(id int) {
			defer wg.Done()
			mu.Lock()
			order = append(order, id)
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	return order
}
