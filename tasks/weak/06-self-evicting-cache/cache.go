// Package cache holds values without keeping them alive, and forgets them on
// its own once they are gone.
//
// weak/01 built a cache whose dead entries lingered until somebody called
// Reap. This one has no Reap: every entry removes itself some time after its
// value is collected. That one change brings a second goroutine into a type
// that had none, and gives an old value a way to destroy a newer entry.
package cache

// Cache maps keys to values it does not own. It is safe for concurrent use.
type Cache[K comparable, V any] struct {
	// entries, and whatever guards them, are yours to declare.
}

// New returns an empty Cache.
func New[K comparable, V any]() *Cache[K, V] {
	panic("undergo: implement New")
}

// Put stores a weak reference to v under k, replacing any entry already there.
//
// It must not keep v alive. Once v is collected the entry removes itself — but
// only if it is still the entry this Put created. Put(k, nil) removes k.
func (c *Cache[K, V]) Put(k K, v *V) {
	panic("undergo: implement Put")
}

// Get returns the value stored under k while it is still reachable elsewhere,
// and (nil, false) otherwise.
func (c *Cache[K, V]) Get(k K) (*V, bool) {
	panic("undergo: implement Get")
}

// Len reports how many entries the cache holds.
//
// An entry whose value has been collected still counts until it has removed
// itself, which happens some time after the collection, never during it.
func (c *Cache[K, V]) Len() int {
	panic("undergo: implement Len")
}
