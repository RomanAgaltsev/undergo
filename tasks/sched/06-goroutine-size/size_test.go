package goroutinesize

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// bracket names the decade, in KiB, that a size in bytes falls in.
func bracket(bytes int) string {
	switch kib := float64(bytes) / 1024; {
	case kib < 1:
		return "<1"
	case kib < 10:
		return "1-10"
	case kib < 100:
		return "10-100"
	default:
		return ">100"
	}
}

// TestPredictions runs each phase in its own process and grades four slots.
// It never prints a size.
func TestPredictions(t *testing.T) {
	measure := func(phase string) int {
		t.Helper()
		v, err := Measure(phase)
		if err != nil {
			t.Fatalf("measuring: %v", err)
		}
		return v
	}
	parked := measure("parked")
	deep := measure("deep")
	returned := measure("returned")
	gone := measure("unreachable")

	predict.Check(t, map[string]any{
		"stack_kib_per_goroutine": bracket(parked),
		"grows_with_stack_depth":  deep > 2*parked,
		"shrinks_after_return":    returned < deep/2,
		"parked_reclaimed_by_gc":  gone > 0,
	})
}
