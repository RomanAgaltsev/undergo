package weakpkg

import (
	"strings"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

func TestPredictions(t *testing.T) {
	local, _, err := Compiles(Local, "1.23")
	if err != nil {
		t.Fatalf("local: %v", err)
	}
	old, diagnostic, err := Compiles("go1.23.12", "1.23")
	if err != nil {
		t.Fatalf("go1.23.12: %v", err)
	}

	predict.Check(t, map[string]any{
		"compiles_on_local":          local,
		"compiles_under_go1_23":      old,
		"diagnostic_says_not_in_std": strings.Contains(diagnostic, "is not in std"),
	})
}
