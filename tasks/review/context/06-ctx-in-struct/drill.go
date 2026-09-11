// Package drill — C4/06 ctx-in-struct.
package drill

import "context"

// Worker performs jobs.
type Worker struct {
	ctx  context.Context
	name string
}

// NewWorker stores the context for later use.
func NewWorker(ctx context.Context, name string) *Worker {
	return &Worker{ctx: ctx, name: name}
}

// Run executes the worker's job using the stored context.
func (w *Worker) Run() error {
	return w.doWork(w.ctx)
}

func (w *Worker) doWork(ctx context.Context) error {
	_ = ctx
	return nil
}
