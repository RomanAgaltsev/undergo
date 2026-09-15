package runnext

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions grades four positions of the observed order. It never prints
// the order itself.
func TestPredictions(t *testing.T) {
	order := SpawnOrder(4)
	if len(order) != 4 {
		t.Fatalf("expected 4 goroutines to report, got %d", len(order))
	}

	predict.Check(t, map[string]any{
		"first":  order[0],
		"second": order[1],
		"third":  order[2],
		"fourth": order[3],
	})
}
