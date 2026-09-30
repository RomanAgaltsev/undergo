// Package atomiccost measures three ways to add one to an int64.
//
// A plain increment, an atomic.Int64.Add and a mutex-guarded increment leave
// the same number behind. They do not cost the same, and the gap between them
// changes shape once more than one core wants the same word.
package atomiccost

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// Plain adds one to *p, n times. The compiler keeps the loop: each iteration
// is one memory increment through p.
//
//go:noinline
func Plain(p *int64, n int) {
	for range n {
		*p++
	}
}

// Atomic adds one to v, n times, with Add.
//
//go:noinline
func Atomic(v *atomic.Int64, n int) {
	for range n {
		v.Add(1)
	}
}

// Guarded is an int64 behind a mutex.
type Guarded struct {
	mu sync.Mutex
	v  int64
}

// Locked adds one to g, n times, taking and releasing the lock each time.
//
//go:noinline
func Locked(g *Guarded, n int) {
	for range n {
		g.mu.Lock()
		g.v++
		g.mu.Unlock()
	}
}

// Costs holds nanoseconds per increment. The test compares them and never
// prints them.
type Costs struct {
	Plain, Atomic, Locked      float64 // one goroutine
	AtomicShared, LockedShared float64 // GOMAXPROCS goroutines, one shared value
}

// trials is how many times each arm runs. The cheapest run counts, because
// noise on a shared machine only ever adds time.
const trials = 5

// perOp runs body(n) in workers goroutines released together and returns the
// best wall time over trials divided by n: the cost of one increment as each
// goroutine experienced it.
func perOp(workers, n int, body func(int)) float64 {
	best := 0.0
	for i := range trials {
		start := make(chan struct{})
		var wg sync.WaitGroup
		for range workers {
			wg.Go(func() {
				<-start
				body(n)
			})
		}
		t0 := time.Now()
		close(start)
		wg.Wait()
		ns := float64(time.Since(t0).Nanoseconds()) / float64(n)
		if i == 0 || ns < best {
			best = ns
		}
	}
	return best
}

// Measure times all five arms. The iteration counts keep every run above
// about ten milliseconds, well clear of the coarsest clock this runs on.
func Measure() Costs {
	procs := runtime.GOMAXPROCS(0)
	var x int64
	var a atomic.Int64
	var g Guarded
	return Costs{
		Plain:        perOp(1, 50_000_000, func(n int) { Plain(&x, n) }),
		Atomic:       perOp(1, 5_000_000, func(n int) { Atomic(&a, n) }),
		Locked:       perOp(1, 2_500_000, func(n int) { Locked(&g, n) }),
		AtomicShared: perOp(procs, 200_000, func(n int) { Atomic(&a, n) }),
		LockedShared: perOp(procs, 50_000, func(n int) { Locked(&g, n) }),
	}
}
