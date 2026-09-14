package cycles

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions compares each pair against your answers. It never prints a
// cycle count.
func TestPredictions(t *testing.T) {
	const garbage = 48

	predict.Check(t, map[string]any{
		"gogc50_beats_gogc400":  Beats(50, 400, garbage),
		"gogc100_beats_gogc800": Beats(100, 800, garbage),
		"gogc400_beats_gogc100": Beats(400, 100, garbage),
		"gogc800_beats_gogc50":  Beats(800, 50, garbage),
		"gogc50_beats_gogc800":  Beats(50, 800, garbage),
	})
}
