package cancellability

import (
	"testing"
	"testing/synctest"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

func TestPredictions(t *testing.T) {
	d := Detach()
	got := map[string]any{
		"without_cancel_done_is_nil":             d.DoneIsNil,
		"without_cancel_keeps_values":            d.KeepsValues,
		"without_cancel_err_after_parent_cancel": d.ErrStillNilAfterCancel,
		"detached_child_is_uncancellable":        d.ChildIsUncancellable,
		"cause_of_plain_cancel_is_canceled":      CauseOfPlainCancel(),
		"cause_of_live_context_is_canceled":      CauseOfLiveContext(),
		"err_when_cause_differs":                 ErrAndCauseDiffer(),
	}

	synctest.Test(t, func(t *testing.T) {
		got["afterfunc_runs_on_already_cancelled_ctx"] = AfterFuncRunsOnDeadContext()
	})

	predict.Check(t, got)
}
