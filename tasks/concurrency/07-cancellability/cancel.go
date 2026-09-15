// Package cancellability asks what remains of a context once it has been
// detached, and where Err and Cause deliberately disagree.
package cancellability

import (
	"context"
	"errors"
	"testing/synctest"
)

type key struct{}

// ErrBecause is the cause handed to WithCancelCause below.
var ErrBecause = errors.New("because")

// Detachment records what WithoutCancel left behind.
type Detachment struct {
	DoneIsNil              bool
	KeepsValues            bool
	ErrStillNilAfterCancel bool
	ChildIsUncancellable   bool
}

// Detach builds a cancellable context carrying a value, detaches it with
// WithoutCancel, cancels the parent, and then derives a cancellable child FROM
// the detached context and cancels that too.
func Detach() Detachment {
	var d Detachment

	parent, cancel := context.WithCancel(context.WithValue(context.Background(), key{}, "v"))
	detached := context.WithoutCancel(parent)
	cancel()

	d.DoneIsNil = detached.Done() == nil
	d.KeepsValues = detached.Value(key{}) == "v"
	d.ErrStillNilAfterCancel = detached.Err() == nil

	sub, cancelSub := context.WithCancel(detached)
	cancelSub()
	d.ChildIsUncancellable = sub.Err() == nil

	return d
}

// CauseOfPlainCancel reports whether Cause matches context.Canceled after an
// ordinary cancel.
func CauseOfPlainCancel() bool {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return errors.Is(context.Cause(ctx), context.Canceled)
}

// CauseOfLiveContext reports whether Cause matches context.Canceled on a
// context nobody has cancelled.
func CauseOfLiveContext() bool {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return errors.Is(context.Cause(ctx), context.Canceled)
}

// ErrAndCauseDiffer reports whether Err stays Canceled while Cause carries the
// supplied reason.
func ErrAndCauseDiffer() bool {
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(ErrBecause)
	return errors.Is(ctx.Err(), context.Canceled) && errors.Is(context.Cause(ctx), ErrBecause)
}

// AfterFuncRunsOnDeadContext registers an AfterFunc on a context that was
// cancelled BEFORE the call, and reports whether it ran.
//
// Must be called from inside a synctest bubble.
func AfterFuncRunsOnDeadContext() bool {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	ran := make(chan struct{})
	context.AfterFunc(ctx, func() { close(ran) })
	synctest.Wait()
	select {
	case <-ran:
		return true
	default:
		return false
	}
}
