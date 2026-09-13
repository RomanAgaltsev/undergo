package settable

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// panics reports whether f panicked, and never reports what it said.
func panics(f func()) (did bool) {
	defer func() {
		if recover() != nil {
			did = true
		}
	}()
	f()
	return false
}

// TestPredictions runs each operation and compares the verdicts with yours.
func TestPredictions(t *testing.T) {
	predict.Check(t, map[string]any{
		"set_through_value":      panics(SetThroughValue),
		"set_through_pointer":    panics(SetThroughPointer),
		"set_unexported":         panics(SetUnexported),
		"set_wrong_type":         panics(SetWrongType),
		"set_slice_element":      panics(SetSliceElement),
		"set_map_value_in_place": panics(SetMapValueInPlace),
		"elem_of_nil_pointer":    panics(ElemOfNilPointer),
		"append_and_assign":      panics(AppendAndAssign),
	})
}
