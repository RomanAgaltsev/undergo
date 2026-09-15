package convert

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions grades the allocation count of each conversion.
func TestPredictions(t *testing.T) {
	b := []byte("hello")

	predict.Check(t, map[string]any{
		"map_index":   int(MapIndex(b)),
		"comparison":  int(Comparison(b)),
		"range_over":  int(RangeOver(b)),
		"type_switch": int(TypeSwitch(b)),
		"assignment":  int(Assignment(b)),
	})
}
