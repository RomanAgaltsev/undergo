package compare

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions grades which comparisons panic and which merely differ.
func TestPredictions(t *testing.T) {
	predict.Check(t, map[string]any{
		"slices_panic":          Panics([]int{1}, []int{1}),
		"maps_panic":            Panics(map[int]int{}, map[int]int{}),
		"funcs_panic":           Panics(func() {}, func() {}),
		"struct_slice_panics":   Panics(WithSlice{}, WithSlice{}),
		"array_of_slice_panics": Panics([1][]int{{1}}, [1][]int{{1}}),
		"plain_struct_panics":   Panics(Plain{1, 2}, Plain{1, 2}),
		"array_panics":          Panics([2]int{1, 2}, [2]int{1, 2}),
		"different_types_panic": Panics([]int{1}, "hello"),
		"plain_struct_equal":    Equal(Plain{1, 2}, Plain{1, 2}),
	})
}
