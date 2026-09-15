package preempt

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions runs the spinner under both settings and grades whether the
// collector got through. It never prints a duration.
func TestPredictions(t *testing.T) {
	normal, err := RunSpinner(false)
	if err != nil {
		t.Fatalf("running the spinner: %v", err)
	}
	off, err := RunSpinner(true)
	if err != nil {
		t.Fatalf("running the spinner with asyncpreemptoff=1: %v", err)
	}

	predict.Check(t, map[string]any{
		"gc_completes_normally":   normal,
		"gc_completes_preemptoff": off,
	})
}
