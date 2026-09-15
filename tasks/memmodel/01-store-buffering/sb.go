// Package sb runs the store-buffering shape, in two variants.
//
// Two goroutines each write one variable and then read the other:
//
//	goroutine A: x = 1; r1 = y
//	goroutine B: y = 1; r2 = x
//
// Three outcomes are obviously possible. The fourth, r1 == 0 && r2 == 0, says
// both goroutines read the value from before the other's write — which no
// interleaving of the four statements can produce.
//
// The plain variant is a data race on purpose. That is the subject, and it is
// why this task is pinned to the default build.
package sb

import (
	"sync"
	"sync/atomic"
)

// PlainSB runs the shape with ordinary variables and counts how many iterations
// produced the (0, 0) outcome.
func PlainSB(iters int) int {
	var x, y int
	seen := 0
	for range iters {
		x, y = 0, 0
		var r1, r2 int
		var wg sync.WaitGroup
		wg.Add(2)
		go func() { x = 1; r1 = y; wg.Done() }()
		go func() { y = 1; r2 = x; wg.Done() }()
		wg.Wait()
		if r1 == 0 && r2 == 0 {
			seen++
		}
	}
	return seen
}

// AtomicSB runs the same shape with sequentially consistent atomics.
func AtomicSB(iters int) int {
	var x, y atomic.Int32
	seen := 0
	for range iters {
		x.Store(0)
		y.Store(0)
		var r1, r2 int32
		var wg sync.WaitGroup
		wg.Add(2)
		go func() { x.Store(1); r1 = y.Load(); wg.Done() }()
		go func() { y.Store(1); r2 = x.Load(); wg.Done() }()
		wg.Wait()
		if r1 == 0 && r2 == 0 {
			seen++
		}
	}
	return seen
}
