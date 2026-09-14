package timers

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions grades the three capacities and the collection question.
func TestPredictions(t *testing.T) {
	predict.Check(t, map[string]any{
		"timer_cap":                TimerCap(),
		"ticker_cap":               TickerCap(),
		"after_cap":                AfterCap(),
		"dropped_ticker_collected": DroppedTickerCollected(),
	})
}
