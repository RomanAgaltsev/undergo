package poolclearing

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

func TestPredictions(t *testing.T) {
	s := Measure()

	predict.Check(t, map[string]any{
		"new_calls_before_any_gc": s.NewCallsBeforeAnyGC,
		"survives_one_gc":         s.SurvivesOneGC,
		"survives_two_gc":         s.SurvivesTwoGC,
		"new_calls_after_two_gc":  s.NewCallsAfterTwoGC,
	})
}
