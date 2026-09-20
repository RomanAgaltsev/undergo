package seed

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

func TestPredictions(t *testing.T) {
	at123, err := SeedHonoured("1.23")
	if err != nil {
		t.Fatalf("go 1.23: %v", err)
	}
	at124, err := SeedHonoured("1.24")
	if err != nil {
		t.Fatalf("go 1.24: %v", err)
	}
	restored, err := SeedHonoured("1.24", "GODEBUG=randseednop=0")
	if err != nil {
		t.Fatalf("go 1.24 with randseednop=0: %v", err)
	}

	predict.Check(t, map[string]any{
		"seed_honoured_at_go123":                    at123,
		"seed_honoured_at_go124":                    at124,
		"seed_honoured_at_go124_with_randseednop_0": restored,
	})
}
