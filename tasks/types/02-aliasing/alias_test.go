package alias

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions grades individual cells rather than whole slices, so a wrong
// answer names one misunderstanding instead of five.
func TestPredictions(t *testing.T) {
	base, twoIndex, threeIndex, grown := Run()

	predict.Check(t, map[string]any{
		"cap_two_index":   cap(base[2:4]),
		"cap_three_index": cap(threeIndex),
		"base_4":          base[4],
		"base_1":          base[1],
		"base_3":          base[3],
		"grown_0":         grown[0],
		"len_two_index":   len(twoIndex),
	})
}
