package reclaim

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions grades how many forced collections each mechanism needed.
//
// The cleanup slot is a predicate rather than a count, and deliberately. The
// exact number measured 1 once and 2 five times over six fresh processes —
// running a cleanup is asynchronous with respect to the collection that queued
// it, so the count is a race. An invariant with room is worth more than a
// number that is right most of the time.
func TestPredictions(t *testing.T) {
	cleanup := CleanupGCs()

	predict.Check(t, map[string]any{
		"weak_gcs":           WeakGCs(),
		"finalizer_gcs":      FinalizerGCs(),
		"cleanup_within_two": cleanup >= 1 && cleanup <= 2,
		"pinned_survives":    PinnedSurvives(),
	})
}
