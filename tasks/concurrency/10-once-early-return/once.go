// Package earlyreturn asks what "exactly once" promises the second caller.
//
// Naive runs f exactly once, and it is still not sync.Once. Race puts each of
// them in the one situation where the difference shows: a second caller
// arriving while the first is still inside f.
package earlyreturn

import (
	"sync/atomic"
	"time"
)

// Doer is anything with a Do method: a *sync.Once, or a *Naive.
type Doer interface{ Do(f func()) }

// Naive is the Once that everybody writes first. The first caller to flip done
// runs f; every later caller finds done already set and moves on.
type Naive struct{ done atomic.Bool }

// Do calls f if and only if no earlier call to Do has.
func (o *Naive) Do(f func()) {
	if o.done.CompareAndSwap(false, true) {
		f()
	}
}

// wait is how long Race gives the second caller to come back while the first
// is still inside f. A caller that returns early does so in microseconds; a
// caller that waits cannot return at all until f finishes. So this only has
// to be long, not precise.
const wait = time.Second

// Result is what Race observed.
type Result struct {
	Calls         int32 // how many of the functions passed to Do ran
	ReturnedEarly bool  // the second caller returned while the first was inside f
}

// Race has a first caller enter o.Do with an f that blocks, then has a second
// caller call o.Do while the first is still inside it.
func Race(o Doer) Result {
	var calls atomic.Int32
	entered := make(chan struct{})
	release := make(chan struct{})
	firstDone := make(chan struct{})
	go func() {
		defer close(firstDone)
		o.Do(func() {
			calls.Add(1)
			close(entered)
			<-release
		})
	}()
	<-entered

	secondDone := make(chan struct{})
	go func() {
		defer close(secondDone)
		o.Do(func() { calls.Add(1) })
	}()

	timer := time.NewTimer(wait)
	defer timer.Stop()
	var r Result
	select {
	case <-secondDone:
		r.ReturnedEarly = true
	case <-timer.C:
	}

	close(release)
	<-firstDone
	<-secondDone
	r.Calls = calls.Load()
	return r
}
