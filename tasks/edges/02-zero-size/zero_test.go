package zero

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions evaluates each comparison and compares it with yours.
func TestPredictions(t *testing.T) {
	predict.Check(t, map[string]any{
		"two_empty_structs_share_an_address": TwoEmptyStructsShareAnAddress(),
		"two_zero_arrays_share_an_address":   TwoZeroArraysShareAnAddress(),
		"trailing_field_adds_padding":        TrailingFieldAddsPadding(),
		"leading_field_adds_padding":         LeadingFieldAddsPadding(),
		"zero_length_slices_share_data":      ZeroLengthSlicesOfOneArrayShareAData(),
		"empty_slice_equals_nil":             EmptySliceAndNilSliceAreEqual(),
	})
}
