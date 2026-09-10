// Package drill — C11/09 typed-pool.
package drill

// Pool is a simple free-list of reusable T values.
type Pool[T any] struct {
	New   func() any
	items []any
}

// Get returns a value from the pool, creating one if empty.
func (p *Pool[T]) Get() any {
	if len(p.items) == 0 {
		return p.New()
	}
	last := p.items[len(p.items)-1]
	p.items = p.items[:len(p.items)-1]
	return last.(T)
}

// Put returns x to the pool for reuse.
func (p *Pool[T]) Put(x any) {
	p.items = append(p.items, x)
}
