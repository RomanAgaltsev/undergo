package bubble

import (
	"testing"
	"testing/synctest"
	"time"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions runs the same three calls inside a bubble and grades what the
// fake clock reported. It never prints a duration.
func TestPredictions(t *testing.T) {
	var (
		slept time.Duration
		ticks int
		year  int
	)

	wall := time.Now()
	synctest.Test(t, func(t *testing.T) {
		slept = SleepTwice(time.Hour, 30*time.Minute)
		ticks = TickCount(time.Second, 10*time.Second)
		year = Year()
	})
	realCost := time.Since(wall)

	predict.Check(t, map[string]any{
		"slept_minutes":     int(slept / time.Minute),
		"ticks":             ticks,
		"year":              year,
		"real_under_second": realCost < time.Second,
	})
}
