package heapbase

import (
	"strings"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

func TestPredictions(t *testing.T) {
	local, err := Address(Local)
	if err != nil {
		t.Fatalf("local: %v", err)
	}
	old, err := Address("go1.25.7")
	if err != nil {
		t.Fatalf("go1.25.7: %v", err)
	}

	predict.Check(t, map[string]any{
		"base_prefix_is_c000_on_local":     strings.HasPrefix(local, BasePrefix),
		"base_prefix_is_c000_under_go1_25": strings.HasPrefix(old, BasePrefix),
	})
}
