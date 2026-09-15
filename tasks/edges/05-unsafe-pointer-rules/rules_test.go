package rules

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions grades, for each snippet, whether vet rejects it.
func TestPredictions(t *testing.T) {
	measured := map[string]any{}
	for _, name := range Snippets {
		rejected, err := VetRejects(name)
		if err != nil {
			t.Fatalf("%v", err)
		}
		measured["vet_rejects_"+name] = rejected
	}
	predict.Check(t, measured)
}
