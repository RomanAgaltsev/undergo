// Package reclaim counts garbage collections rather than waiting for them.
//
// Every measurement here forces collections with runtime.GC and asks after each
// one whether something has happened yet. That is deterministic in a way that
// counting the collector's own cycles is not — gc/02 in this repository exists
// because cycle counts under a workload vary from run to run. A forced
// collection is an event you caused, so counting those is counting your own
// actions.
package reclaim

import (
	"runtime"
	"weak"
)

// Big is large enough not to share a block with anything else.
type Big struct{ b [512]byte }

// gcsUntil forces collections until check reports true, and returns how many
// were needed. It returns -1 if the limit is reached first.
func gcsUntil(limit int, check func() bool) int {
	for i := 1; i <= limit; i++ {
		runtime.GC()
		if check() {
			return i
		}
	}
	return -1
}

func fired(ch chan struct{}) func() bool {
	return func() bool {
		select {
		case <-ch:
			return true
		default:
			return false
		}
	}
}

// WeakGCs reports how many collections until a dropped weak.Pointer's target is
// gone.
func WeakGCs() int {
	var w weak.Pointer[Big]
	func() {
		b := &Big{}
		w = weak.Make(b)
		runtime.KeepAlive(b)
	}()
	return gcsUntil(8, func() bool { return w.Value() == nil })
}

// FinalizerGCs reports how many collections until a SetFinalizer function runs.
func FinalizerGCs() int {
	ch := make(chan struct{})
	func() {
		b := &Big{}
		runtime.SetFinalizer(b, func(*Big) { close(ch) })
	}()
	return gcsUntil(8, fired(ch))
}

// CleanupGCs reports how many collections until an AddCleanup function runs.
func CleanupGCs() int {
	ch := make(chan struct{})
	func() {
		b := &Big{}
		runtime.AddCleanup(b, func(c chan struct{}) { close(c) }, ch)
	}()
	return gcsUntil(8, fired(ch))
}

// PinnedSurvives reports whether a weak pointer whose target is still
// referenced survives a collection.
func PinnedSurvives() bool {
	b := &Big{}
	w := weak.Make(b)
	runtime.GC()
	alive := w.Value() != nil
	runtime.KeepAlive(b)
	return alive
}
