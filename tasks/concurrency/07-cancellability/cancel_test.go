package cancellability

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

type key struct{}

var errBecause = errors.New("because")

func TestPredictions(t *testing.T) {
	got := map[string]any{}

	// A cancellable context carrying a value, detached with WithoutCancel,
	// then cancelled.
	parent, cancel := context.WithCancel(context.WithValue(context.Background(), key{}, "v"))
	detached := context.WithoutCancel(parent)
	cancel()

	got["without_cancel_done_is_nil"] = detached.Done() == nil
	got["without_cancel_keeps_values"] = detached.Value(key{}) == "v"
	got["without_cancel_err_after_parent_cancel"] = detached.Err() == nil

	// A cancellable context derived FROM the detached one.
	sub, cancelSub := context.WithCancel(detached)
	cancelSub()
	got["detached_child_is_uncancellable"] = sub.Err() == nil

	// Cause on a context cancelled the ordinary way, and on one still live.
	plain, cancelPlain := context.WithCancel(context.Background())
	cancelPlain()
	got["cause_of_plain_cancel_is_canceled"] = errors.Is(context.Cause(plain), context.Canceled)

	live, cancelLive := context.WithCancel(context.Background())
	defer cancelLive()
	got["cause_of_live_context_is_canceled"] = errors.Is(context.Cause(live), context.Canceled)

	// Cancelled with an explicit cause: Err and Cause disagree on purpose.
	caused, cancelCause := context.WithCancelCause(context.Background())
	cancelCause(errBecause)
	got["err_when_cause_differs"] = errors.Is(caused.Err(), context.Canceled) &&
		errors.Is(context.Cause(caused), errBecause)

	// AfterFunc registered on a context that is ALREADY cancelled.
	synctest.Test(t, func(t *testing.T) {
		dead, cancelDead := context.WithCancel(context.Background())
		cancelDead()
		ran := make(chan struct{})
		context.AfterFunc(dead, func() { close(ran) })
		synctest.Wait()
		select {
		case <-ran:
			got["afterfunc_runs_on_already_cancelled_ctx"] = true
		default:
			got["afterfunc_runs_on_already_cancelled_ctx"] = false
		}
	})

	predict.Check(t, got)
}
