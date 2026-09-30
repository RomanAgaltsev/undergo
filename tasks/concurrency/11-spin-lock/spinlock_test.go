package spinlock

import (
	"runtime"
	"sync"
	"testing"
	"time"
)

// within runs f and fails the test if it has not returned by limit. A broken
// lock can deadlock, and a hang must fail with a message rather than wait for
// the ten-minute timeout of go test.
func within(t *testing.T, limit time.Duration, what string, f func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		f()
	}()
	timer := time.NewTimer(limit)
	defer timer.Stop()
	select {
	case <-done:
	case <-timer.C:
		t.Fatalf("%s did not finish within %v: does Unlock release the lock?", what, limit)
	}
}

// TestMutualExclusion has eight goroutines add one to a plain int under the
// lock, fifty thousand times each, thirty times over. A lock that ever lets two
// holders in loses increments, and the count comes out short.
func TestMutualExclusion(t *testing.T) {
	const rounds, workers, per = 30, 8, 50_000
	for round := range rounds {
		var l SpinLock
		count := 0
		within(t, 30*time.Second, "a round of increments", func() {
			var wg sync.WaitGroup
			for range workers {
				wg.Go(func() {
					for range per {
						l.Lock()
						count++
						l.Unlock()
					}
				})
			}
			wg.Wait()
		})
		if count != workers*per {
			t.Fatalf("round %d: counted %d increments, want %d: two goroutines held the lock at once",
				round, count, workers*per)
		}
	}
}

// TestYieldsWhileWaiting runs on one P. A holder takes the lock and parks; a
// waiter calls Lock; and only the test goroutine can release the holder, but
// the test goroutine needs the single P to do it. A Lock that retries without
// yielding keeps that P until the runtime preempts it, about ten milliseconds
// later, every round.
func TestYieldsWhileWaiting(t *testing.T) {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(1))
	const rounds = 50
	const budget = 250 * time.Millisecond

	start := time.Now()
	within(t, 30*time.Second, "the handoff rounds", func() {
		for range rounds {
			var l SpinLock
			held := make(chan struct{})
			release := make(chan struct{})
			acquired := make(chan struct{})
			go func() {
				l.Lock()
				close(held)
				<-release
				l.Unlock()
			}()
			<-held
			go func() {
				l.Lock()
				close(acquired)
				l.Unlock()
			}()
			runtime.Gosched() // let the waiter reach Lock
			close(release)
			<-acquired
		}
	})
	if elapsed := time.Since(start); elapsed > budget {
		t.Fatalf("%d lock handoffs on one P took %v, budget %v: does Lock let other goroutines run while it waits?",
			rounds, elapsed.Round(time.Millisecond), budget)
	}
}
