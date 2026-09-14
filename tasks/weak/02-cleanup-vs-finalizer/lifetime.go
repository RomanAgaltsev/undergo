// Package lifetime runs the same six experiments against runtime.AddCleanup and
// runtime.SetFinalizer.
//
// You do not implement anything here. Work out the answers, write them into
// prediction.yaml, then verify.
//
// Every function below waits generously for a result rather than reading a flag
// straight after runtime.GC(). Since Go 1.25 cleanups run concurrently, so the
// obvious version of this code races the runtime and answers differently on a
// busy machine.
package lifetime

import (
	"runtime"
	"sync/atomic"
	"time"
)

// Obj is an ordinary heap object with nothing special about it.
type Obj struct {
	data [512]byte //nolint:unused // ballast, so the object is an ordinary heap allocation
}

// settle is how long each experiment waits after a collection before deciding
// nothing is going to happen. It is deliberately generous: a slow machine must
// not change an answer.
const settle = 200 * time.Millisecond

// maxGCs bounds every experiment, so a negative answer takes a known time
// rather than hanging.
const maxGCs = 5

// gcsUntil forces collections until done fires, and reports how many it took.
// It returns -1 if nothing happened within maxGCs.
func gcsUntil(done <-chan struct{}) int {
	for i := 1; i <= maxGCs; i++ {
		runtime.GC()
		select {
		case <-done:
			return i
		case <-time.After(settle):
		}
	}
	return -1
}

// CleanupGCsNeeded reports how many collections pass before a cleanup runs.
func CleanupGCsNeeded() int {
	done := make(chan struct{})
	func() {
		o := &Obj{}
		runtime.AddCleanup(o, func(struct{}) { close(done) }, struct{}{})
	}()
	return gcsUntil(done)
}

// FinalizerGCsNeeded reports how many collections pass before a finalizer runs.
func FinalizerGCsNeeded() int {
	done := make(chan struct{})
	func() {
		o := &Obj{}
		runtime.SetFinalizer(o, func(*Obj) { close(done) })
	}()
	return gcsUntil(done)
}

// CleanupSelfReferencePanics registers a cleanup whose argument is the object
// being cleaned up, and reports whether that panics.
func CleanupSelfReferencePanics() (panicked bool) {
	defer func() {
		if recover() != nil {
			panicked = true
		}
	}()

	o := &Obj{}
	runtime.AddCleanup(o, func(*Obj) {}, o)
	return false
}

// FinalizerRunsOnSelfReference sets a finalizer whose closure captures the
// object it is finalizing, and reports whether it ever runs.
func FinalizerRunsOnSelfReference() bool {
	done := make(chan struct{})
	func() {
		o := &Obj{}
		runtime.SetFinalizer(o, func(*Obj) {
			// The closure captures o, not just the argument.
			_ = o
			close(done)
		})
	}()
	return gcsUntil(done) > 0
}

// resurrected is where FinalizerCanResurrect parks the object it was handed.
var resurrected *Obj

// FinalizerCanResurrect reports whether a finalizer can make its object
// reachable again by storing the pointer it receives.
func FinalizerCanResurrect() bool {
	done := make(chan struct{})
	func() {
		o := &Obj{}
		runtime.SetFinalizer(o, func(self *Obj) {
			resurrected = self
			close(done)
		})
	}()
	if gcsUntil(done) < 0 {
		return false
	}
	got := resurrected != nil
	resurrected = nil
	return got
}

// MultipleCleanupsAllRun registers three cleanups on one object and reports
// whether every one of them ran.
//
// The flags are atomic because cleanups run on their own goroutine: plain bools
// written there and read here would be a data race, and the race detector says
// so. That is not a detail of this experiment — it is the same fact the
// experiment is measuring.
func MultipleCleanupsAllRun() bool {
	var a, b, c atomic.Bool
	set := func(p *atomic.Bool) { p.Store(true) }

	func() {
		o := &Obj{}
		runtime.AddCleanup(o, set, &a)
		runtime.AddCleanup(o, set, &b)
		runtime.AddCleanup(o, set, &c)
	}()

	allRan := func() bool { return a.Load() && b.Load() && c.Load() }
	for range maxGCs {
		runtime.GC()
		time.Sleep(settle)
		if allRan() {
			return true
		}
	}
	return allRan()
}
