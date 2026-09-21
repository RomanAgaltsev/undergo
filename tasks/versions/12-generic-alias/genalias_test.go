package genalias

import (
	"strings"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

func TestPredictions(t *testing.T) {
	localOK, _, err := Compiles(Local)
	if err != nil {
		t.Fatalf("local: %v", err)
	}
	oldOK, diagnostic, err := Compiles("go1.23.12")
	if err != nil {
		t.Fatalf("go1.23.12: %v", err)
	}
	identical, typeName, err := Identity(Local)
	if err != nil {
		t.Fatalf("identity on local: %v", err)
	}

	predict.Check(t, map[string]any{
		"compiles_on_local":              localOK,
		"compiles_under_go1_23":          oldOK,
		"diagnostic_names_an_experiment": strings.Contains(diagnostic, "GOEXPERIMENT"),
		"instantiations_are_identical":   identical == "true",
		"alias_type_name":                typeName,
	})
}
