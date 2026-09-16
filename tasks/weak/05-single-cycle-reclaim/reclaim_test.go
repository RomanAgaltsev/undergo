package reclaim

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions grades how many forced collections each mechanism needed.
//
// Two of the four slots are predicates rather than counts, and deliberately:
// both AddCleanup and SetFinalizer run their function asynchronously with
// respect to the collection that queued it, so the exact count is a race.
//
// The cleanup slot was made a predicate when it measured 1 once and 2 five
// times over six fresh processes. The finalizer slot measured 1 all six times
// and was left as an exact count — which was wrong: over twelve processes it
// measures 1 seven times and 2 five times, and it failed gate 2 on ubuntu while
// passing on macOS. Six samples of a near-even coin flip is not evidence.
//
// An invariant with room is worth more than a number that is right most of the
// time.
func TestPredictions(t *testing.T) {
	finalizer := FinalizerGCs()
	cleanup := CleanupGCs()

	predict.Check(t, map[string]any{
		"weak_gcs":             WeakGCs(),
		"finalizer_within_two": finalizer >= 1 && finalizer <= 2,
		"cleanup_within_two":   cleanup >= 1 && cleanup <= 2,
		"pinned_survives":      PinnedSurvives(),
	})
}
