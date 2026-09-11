// Package drill — C5/06 unlock-on-panic.
package drill

import "sync"

// Registry tracks integer counters by key.
type Registry struct {
	mu    sync.Mutex
	items map[string]int
}

// Top returns the key (from keys) with the highest counter.
func (r *Registry) Top(keys []string) string {
	r.mu.Lock()
	best := keys[0]
	for _, k := range keys {
		if r.items[k] > r.items[best] {
			best = k
		}
	}
	r.mu.Unlock()
	return best
}
