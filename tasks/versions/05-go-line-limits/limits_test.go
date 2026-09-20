package limits

import (
	"strings"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

func TestPredictions(t *testing.T) {
	cap119, err := AfterCap("1.19")
	if err != nil {
		t.Fatalf("go 1.19: %v", err)
	}
	cap127, err := AfterCap("1.27")
	if err != nil {
		t.Fatalf("go 1.27: %v", err)
	}

	// Reach for a GODEBUG that no longer exists.
	out, err := Startup("1.27", "GODEBUG=asynctimerchan=1")
	if err != nil {
		t.Fatalf("startup probe: %v", err)
	}

	predict.Check(t, map[string]any{
		"after_cap_at_go119":            cap119,
		"after_cap_at_go127":            cap127,
		"init_ran_with_removed_godebug": strings.Contains(out, "init"),
		"output_says_fatal_error":       strings.Contains(out, "fatal error"),
	})
}
