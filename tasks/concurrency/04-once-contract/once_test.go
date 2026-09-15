package oncecontract

import (
	"sync"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

func didPanic(f func()) (panicked bool) {
	defer func() { panicked = recover() != nil }()
	f()
	return
}

func TestPredictions(t *testing.T) {
	got := map[string]any{}

	// sync.Once.Do, where the first f panics.
	var once sync.Once
	calls := 0
	got["do_panic_reaches_caller"] = didPanic(func() {
		once.Do(func() { calls++; panic("boom") })
	})
	before := calls
	got["do_second_call_panics"] = didPanic(func() {
		once.Do(func() { calls++ })
	})
	got["do_second_call_runs_f"] = calls > before

	// sync.OnceFunc, where f panics. Called twice.
	fn := sync.OnceFunc(func() { panic("boom") })
	_ = didPanic(fn)
	got["oncefunc_second_call_panics"] = didPanic(fn)

	// sync.OnceValue, where f panics. Called twice; how often did f run?
	vCalls := 0
	vf := sync.OnceValue(func() int { vCalls++; panic("boom") })
	_ = didPanic(func() { _ = vf() })
	_ = didPanic(func() { _ = vf() })
	got["oncevalue_calls_f_twice"] = vCalls == 2

	predict.Check(t, got)
}
