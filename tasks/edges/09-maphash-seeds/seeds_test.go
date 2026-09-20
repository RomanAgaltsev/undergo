package seeds

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

func TestPredictions(t *testing.T) {
	sameSeed, differentSeeds, stringAndBytes := InProcess()

	across, err := AcrossProcesses()
	if err != nil {
		t.Fatalf("running testdata/twice: %v", err)
	}

	predict.Check(t, map[string]any{
		"same_seed_same_key_agrees": sameSeed,
		"different_seeds_agree":     differentSeeds,
		"string_and_bytes_agree":    stringAndBytes,
		"two_processes_agree":       across,
	})
}
