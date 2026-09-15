package classes

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions grades the bytes charged across the whole curve.
func TestPredictions(t *testing.T) {
	predict.Check(t, map[string]any{
		"n1":          BytesPerAlloc(1),
		"n8":          BytesPerAlloc(8),
		"n9":          BytesPerAlloc(9),
		"n17":         BytesPerAlloc(17),
		"n33":         BytesPerAlloc(33),
		"n600":        BytesPerAlloc(600),
		"n1000":       BytesPerAlloc(1000),
		"n32768":      BytesPerAlloc(32768),
		"n32769":      BytesPerAlloc(32769),
		"one_pointer": PointerBytesPerAlloc(),
	})
}
