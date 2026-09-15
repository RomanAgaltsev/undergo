// Package cascade builds cancellable fan-out/fan-in pipelines.
package cascade

import "context"

// Gen emits vs on the returned channel, then closes it.
//
// It stops early if ctx is cancelled, and closes the channel in every case.
func Gen(ctx context.Context, vs ...int) <-chan int {
	panic("TODO: implement Gen")
}

// Square reads in until it is closed or ctx is cancelled, emitting v*v for each
// value. It closes the returned channel when it stops, for either reason.
func Square(ctx context.Context, in <-chan int) <-chan int {
	panic("TODO: implement Square")
}

// FanOut returns n channels, each fed from in. A value from in goes to exactly
// one of them — whichever is ready first. Every returned channel is closed when
// in is closed or ctx is cancelled.
//
// It panics if n is not positive.
func FanOut(ctx context.Context, in <-chan int, n int) []<-chan int {
	panic("TODO: implement FanOut")
}

// FanIn merges ins into one channel, closed once every input is drained or ctx
// is cancelled. Order across inputs is not specified.
func FanIn(ctx context.Context, ins ...<-chan int) <-chan int {
	panic("TODO: implement FanIn")
}
