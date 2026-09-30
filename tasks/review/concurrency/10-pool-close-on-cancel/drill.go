// Package drill — C1/10 pool-close-on-cancel.
package drill

import (
	"context"
	"errors"
	"sync"
)

// ErrClosed is returned by Submit once the pool has shut down.
var ErrClosed = errors.New("pool closed")

// Pool runs submitted jobs on a fixed set of workers until its context is
// cancelled.
type Pool struct {
	jobs   chan func()
	closed chan struct{}
	wg     sync.WaitGroup
}

// New starts n workers. When ctx is cancelled the pool shuts down: the workers
// finish what is queued and exit, and Submit starts returning ErrClosed.
func New(ctx context.Context, n int) *Pool {
	p := &Pool{
		jobs:   make(chan func(), n),
		closed: make(chan struct{}),
	}
	for range n {
		p.wg.Add(1)
		go p.work()
	}
	go func() {
		<-ctx.Done()
		close(p.closed)
		close(p.jobs)
	}()
	return p
}

func (p *Pool) work() {
	defer p.wg.Done()
	for job := range p.jobs {
		job()
	}
}

// Submit queues job, or returns ErrClosed if the pool has shut down.
func (p *Pool) Submit(job func()) error {
	select {
	case <-p.closed:
		return ErrClosed
	default:
	}
	p.jobs <- job
	return nil
}

// Wait blocks until every worker has exited.
func (p *Pool) Wait() { p.wg.Wait() }
