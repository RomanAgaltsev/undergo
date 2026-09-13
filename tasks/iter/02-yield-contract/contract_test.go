package contract

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions runs each consumer and compares the yield counts with yours.
// It never prints a count.
func TestPredictions(t *testing.T) {
	predict.Check(t, map[string]any{
		"plain_yields":       Plain(),
		"break_yields":       Break(),
		"return_yields":      Return(),
		"labelled_yields":    Labelled(),
		"panicking_yields":   Panicking(),
		"misbehaving_panics": MisbehavingPanics(),
	})
}
