package limit

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions asks the runtime which knob set the goal, and compares the
// answers with yours. It never prints a goal.
func TestPredictions(t *testing.T) {
	predict.Check(t, map[string]any{
		"clamped_gogc800_limit16":    ClampedByLimit(800, 16),
		"clamped_gogc100_limit512":   ClampedByLimit(100, 512),
		"clamped_gogc_off_limit32":   ClampedByLimit(-1, 32),
		"clamped_gogc100_limit_huge": ClampedByLimit(100, 1<<20),
		"limit_can_be_exceeded":      LimitCanBeExceeded(),
	})
}
