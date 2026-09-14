package cycles

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions compares your ordering against the one this machine produced.
// It never prints a cycle count.
func TestPredictions(t *testing.T) {
	predict.Check(t, map[string]any{
		"ordering":             Ordering(48),
		"gogc50_beats_gogc800": Cycles(50, 48) > Cycles(800, 48),
	})
}
