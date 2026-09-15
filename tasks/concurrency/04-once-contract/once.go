// Package oncecontract holds four small experiments with sync.Once and the
// Go 1.21 once-wrappers, each one asking what happens after f panics.
package oncecontract

import "sync"

// DidPanic reports whether f panicked, swallowing the panic.
func DidPanic(f func()) (panicked bool) {
	defer func() { panicked = recover() != nil }()
	f()
	return
}

// DoResult records what a sync.Once did across two calls, where the first f
// panics and the second does not.
type DoResult struct {
	FirstPanicked  bool
	SecondPanicked bool
	SecondRanF     bool
}

// RunDo calls once.Do twice: first with a panicking f, then with a plain one.
func RunDo() DoResult {
	var once sync.Once
	calls := 0

	var r DoResult
	r.FirstPanicked = DidPanic(func() {
		once.Do(func() { calls++; panic("boom") })
	})
	before := calls
	r.SecondPanicked = DidPanic(func() {
		once.Do(func() { calls++ })
	})
	r.SecondRanF = calls > before
	return r
}

// OnceFuncSecondCallPanics calls a sync.OnceFunc whose f panics, twice, and
// reports whether the SECOND call panicked.
func OnceFuncSecondCallPanics() bool {
	fn := sync.OnceFunc(func() { panic("boom") })
	_ = DidPanic(fn)
	return DidPanic(fn)
}

// OnceValueFCalls calls a sync.OnceValue whose f panics, twice, and reports how
// many times f itself ran.
func OnceValueFCalls() int {
	calls := 0
	vf := sync.OnceValue(func() int { calls++; panic("boom") })
	_ = DidPanic(func() { _ = vf() })
	_ = DidPanic(func() { _ = vf() })
	return calls
}
