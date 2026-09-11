package boxing

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// allocs reports how many heap allocations boxing v performs.
//
// This test must never call t.Parallel(): testing.AllocsPerRun panics if
// parallel tests are running.
func allocs(v int) int {
	return int(testing.AllocsPerRun(200, func() { Sink = Box(v) }))
}

// TestPredictions measures the cost of boxing at several values and finds the
// exact point where it changes. It never prints what it found.
func TestPredictions(t *testing.T) {
	// The threshold is discovered rather than assumed: the smallest value whose
	// boxing allocates. Grading it forces an exact answer instead of "small
	// integers are cheap".
	threshold := -1
	for v := 0; v <= 4096; v++ {
		if allocs(v) > 0 {
			threshold = v
			break
		}
	}

	predict.Check(t, map[string]any{
		"allocs_0":         allocs(0),
		"allocs_255":       allocs(255),
		"allocs_256":       allocs(256),
		"allocs_100000":    allocs(100000),
		"first_allocating": threshold,
	})
}
