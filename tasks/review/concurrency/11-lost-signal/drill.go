// Package drill — C1/11 lost-signal.
package drill

import "sync"

// Gate holds the server back until its configuration has loaded.
type Gate struct {
	mu   sync.Mutex
	cond *sync.Cond
}

// NewGate returns a closed Gate.
func NewGate() *Gate {
	g := &Gate{}
	g.cond = sync.NewCond(&g.mu)
	return g
}

// Wait blocks until Open is called.
func (g *Gate) Wait() {
	g.mu.Lock()
	g.cond.Wait()
	g.mu.Unlock()
}

// Open lets the waiter through.
func (g *Gate) Open() {
	g.cond.Signal()
}
