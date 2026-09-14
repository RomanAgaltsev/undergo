package shapes

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions measures the shape and cost of interface values and compares
// them with yours. It never prints a size.
func TestPredictions(t *testing.T) {
	predict.Check(t, map[string]any{
		"empty_interface_words":   EmptyInterfaceWords(),
		"method_interface_words":  MethodInterfaceWords(),
		"slice_words":             SliceWords(),
		"boxing_small_int_allocs": BoxingSmallIntAllocates(),
		"boxing_large_int_allocs": BoxingLargeIntAllocates(),
		"boxing_pointer_allocs":   BoxingPointerAllocates(),
		"boxing_struct_allocs":    BoxingStructAllocates(),
	})
}
