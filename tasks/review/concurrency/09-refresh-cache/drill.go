// Package drill — C1/09 refresh-cache.
package drill

import (
	"sync"
	"time"
)

// Cache is a concurrent map with a background refresher.
type Cache struct {
	mu   sync.RWMutex
	data map[string]int
}

// NewCache returns a cache that stamps a refresh marker every interval.
func NewCache(interval time.Duration) *Cache {
	c := &Cache{data: make(map[string]int)}
	go func() {
		for range time.Tick(interval) {
			c.mu.RLock()
			c.data["_refreshed"] = int(time.Now().Unix())
			c.mu.RUnlock()
		}
	}()
	return c
}

// Get returns the value stored for key.
func (c *Cache) Get(key string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data[key]
}

// Set stores v for key.
func (c *Cache) Set(key string, v int) {
	c.mu.Lock()
	c.data[key] = v
	c.mu.RUnlock()
}
