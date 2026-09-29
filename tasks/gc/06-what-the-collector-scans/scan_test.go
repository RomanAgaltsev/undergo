package scan

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions grades five survivals and one vet verdict.
func TestPredictions(t *testing.T) {
	vet, err := VetFlags()
	if err != nil {
		t.Fatalf("running go vet: %v", err)
	}

	predict.Check(t, map[string]any{
		"in_bytes_survives":     InBytes(),
		"through_view_survives": ThroughView(),
		"in_pointers_survives":  InPointers(),
		"in_words_survives":     InWords(),
		"in_anys_survives":      InAnys(),
		"vet_flags_a_store":     vet,
	})
}
