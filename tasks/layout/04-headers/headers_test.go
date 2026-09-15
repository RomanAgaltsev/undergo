package headers

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions grades the header sizes, the struct total, and one offset.
func TestPredictions(t *testing.T) {
	s := Sizes()

	predict.Check(t, map[string]any{
		"string_size":    s["string"],
		"slice_size":     s["slice"],
		"map_size":       s["map"],
		"iface_size":     s["iface"],
		"total_size":     s["total"],
		"map_offset":     OffsetOfMap(),
		"big_slice_size": BigSliceSize(),
	})
}
