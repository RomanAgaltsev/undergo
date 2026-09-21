package negzero

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

func TestPredictions(t *testing.T) {
	local, err := Read(Local)
	if err != nil {
		t.Fatalf("local: %v", err)
	}
	old, err := Read("go1.21.13")
	if err != nil {
		t.Fatalf("go1.21.13: %v", err)
	}

	predict.Check(t, map[string]any{
		"scalar_is_zero_on_local":           local.Scalar == "true",
		"scalar_is_zero_under_go1_21":       old.Scalar == "true",
		"small_struct_is_zero_on_local":     local.SmallStruct == "true",
		"small_struct_is_zero_under_go1_21": old.SmallStruct == "true",
		"big_struct_is_zero_on_local":       local.BigStruct == "true",
		"big_struct_is_zero_under_go1_21":   old.BigStruct == "true",
		"array_is_zero_on_local":            local.Array == "true",
		"array_is_zero_under_go1_21":        old.Array == "true",
		"field_is_zero_on_local":            local.Field == "true",
		"field_is_zero_under_go1_21":        old.Field == "true",
	})
}
