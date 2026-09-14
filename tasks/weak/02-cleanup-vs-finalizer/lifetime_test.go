package lifetime

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions runs all six experiments and compares them with yours.
// It never prints a result.
func TestPredictions(t *testing.T) {
	predict.Check(t, map[string]any{
		"cleanup_gcs_needed":               CleanupGCsNeeded(),
		"finalizer_gcs_needed":             FinalizerGCsNeeded(),
		"cleanup_self_reference_panics":    CleanupSelfReferencePanics(),
		"finalizer_runs_on_self_reference": FinalizerRunsOnSelfReference(),
		"finalizer_can_resurrect":          FinalizerCanResurrect(),
		"multiple_cleanups_all_run":        MultipleCleanupsAllRun(),
	})
}
