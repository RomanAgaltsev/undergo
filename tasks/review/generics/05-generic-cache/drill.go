// Package drill — C11/05 generic-cache.
package drill

// Cache is a simple in-memory key/value cache, generic over K and V.
type Cache[K comparable, V any] struct {
	m map[K]any
}

// NewCache builds an empty cache.
func NewCache[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{m: make(map[K]any)}
}

// Set stores v under k.
func (c *Cache[K, V]) Set(k K, v V) {
	c.m[k] = v
}

// Get returns the value stored under k.
func (c *Cache[K, V]) Get(k K) V {
	return c.m[k].(V)
}

// Keys returns all keys currently in the cache.
func (c *Cache[K, V]) Keys() []any {
	keys := make([]any, 0, len(c.m))
	for k := range c.m {
		keys = append(keys, k)
	}
	return keys
}
