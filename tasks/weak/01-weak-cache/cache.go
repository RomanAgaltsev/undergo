// Package cache holds values without keeping them alive.
//
// An ordinary map is a strong reference: anything you put in it lives as long
// as the map does. This one uses weak.Pointer, so the garbage collector is free
// to take a value back the moment nothing else wants it.
package cache

// Cache maps keys to values it does not own.
//
// Hold weak references and nothing else. A cache that also keeps the value in
// an ordinary field is a strong cache wearing a weak one's name, and every test
// in this task is built to catch that.
type Cache[K comparable, V any] struct {
	// entries is yours to declare.
}

// New returns an empty Cache.
func New[K comparable, V any]() *Cache[K, V] {
	panic("undergo: implement New")
}

// Put stores a weak reference to v under k.
//
// It must not keep v alive. Replacing an existing key replaces its entry.
func (c *Cache[K, V]) Put(k K, v *V) {
	panic("undergo: implement Put")
}

// Get returns the value stored under k while it is still reachable elsewhere.
//
// Once nothing else holds the value, Get reports (nil, false) — the same answer
// it gives for a key that was never stored.
func (c *Cache[K, V]) Get(k K) (*V, bool) {
	panic("undergo: implement Get")
}

// Len reports how many entries the cache holds.
//
// This counts entries, not live values: an entry whose value has been collected
// is still an entry until Reap removes it. The difference is the point of the
// task, so do not be tempted to make Len filter.
func (c *Cache[K, V]) Len() int {
	panic("undergo: implement Len")
}

// Reap removes every entry whose value has been collected.
//
// Without it the map grows for the lifetime of the process even though the
// values are long gone.
func (c *Cache[K, V]) Reap() {
	panic("undergo: implement Reap")
}
