package promoted

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

func TestPredictions(t *testing.T) {
	once26, _, err := Compiles("1.26", ProgramOnce)
	if err != nil {
		t.Fatalf("once at go 1.26: %v", err)
	}
	twice26, _, err := Compiles("1.26", ProgramTwice)
	if err != nil {
		t.Fatalf("twice at go 1.26: %v", err)
	}
	once27, _, err := Compiles("1.27", ProgramOnce)
	if err != nil {
		t.Fatalf("once at go 1.27: %v", err)
	}
	mid, deep, err := Shadowed("1.27")
	if err != nil {
		t.Fatalf("shadowed: %v", err)
	}
	if mid != "7" && deep != "7" {
		t.Fatalf("neither field holds the value: mid=%q deep=%q", mid, deep)
	}

	depth := 2
	if mid == "7" {
		depth = 1
	}

	predict.Check(t, map[string]any{
		"promoted_once_compiles_under_go1_26":  once26,
		"promoted_twice_compiles_under_go1_26": twice26,
		"promoted_once_compiles_under_go1_27":  once27,
		"shadowed_key_sets_depth":              depth,
	})
}
