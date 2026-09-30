package earlyreturn

import (
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions races the naive Once and sync.Once and grades three slots.
func TestPredictions(t *testing.T) {
	start := runtime.NumGoroutine()
	naive := Race(&Naive{})
	std := Race(&sync.Once{})
	settle(t, start)

	predict.Check(t, map[string]any{
		"naive_runs_f_once":          naive.Calls == 1,
		"naive_second_returns_early": naive.ReturnedEarly,
		"once_second_returns_early":  std.ReturnedEarly,
	})
}

// settle fails the test if goroutines Race started are still running a while
// after it returned. Race waits for all of them, so this only catches a
// change to Race that stops waiting.
func settle(t *testing.T, start int) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for runtime.NumGoroutine() > start {
		if time.Now().After(deadline) {
			t.Fatalf("Race left goroutines running")
		}
		time.Sleep(time.Millisecond)
	}
}
