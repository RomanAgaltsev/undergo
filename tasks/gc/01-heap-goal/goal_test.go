package goal

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions reads the real heap goal and compares it with yours.
// It never prints a ratio.
func TestPredictions(t *testing.T) {
	predict.Check(t, map[string]any{
		"formula_gogc_100":           Formula(100),
		"large_heap_ratio_gogc_100":  LargeHeapRatio(100),
		"large_heap_ratio_gogc_400":  LargeHeapRatio(400),
		"small_heap_exceeds_formula": SmallHeapExceedsFormula(),
		"large_heap_matches_formula": LargeHeapMatchesFormula(),
		"gap_shrinks_as_heap_grows":  GapShrinksAsHeapGrows(),
	})
}
