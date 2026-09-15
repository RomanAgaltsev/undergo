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
		ticks = TickCount(time.Second, 10500*time.Millisecond)
		year = Year()
	})
	realCost := time.Since(wall)

	// The time slot compares the two clocks against each other rather than
	// against a fixed budget. An earlier version asked whether the real cost
	// was under one second; it measured 0-22ms locally and still failed about
	// half of all gate-2 runs, because gate 2's load is process creation,
	// compilation and antivirus scanning rather than CPU contention. A
	// threshold against a constant is a tolerance. This is an invariant: the
	// bubble slept 90 minutes, so any real cost short of that proves the clock
	// was fake, with a margin no scheduler can close.
	predict.Check(t, map[string]any{
		"slept_minutes":     int(slept / time.Minute),
		"ticks":             ticks,
		"year":              year,
		"fake_exceeds_real": slept > realCost,
	})
}
