package typednil

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions evaluates each comparison and compares it with yours.
func TestPredictions(t *testing.T) {
	predict.Check(t, map[string]any{
		"true_nil_is_nil":                    TrueNilIsNil(),
		"nil_pointer_in_interface_is_nil":    NilPointerInInterfaceIsNil(),
		"returned_nil_pointer_is_nil":        ReturnedNilPointerIsNil(),
		"comparing_to_typed_nil_works":       ComparingToTypedNilWorks(),
		"reflect_says_nil":                   ReflectSaysNil(),
		"errors_is_nil_agrees_with_equality": ErrorsIsNilAgreesWithEquality(),
	})
}
