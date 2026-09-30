// Package drill — C7/13 sharded-unpadded.
package drill

import "sync"

const shards = 16

// shard is one lock and the counts it guards.
type shard struct {
	mu     sync.Mutex
	counts map[uint64]int
}

// Counter counts events per user. Users are spread over sixteen shards by id,
// each with its own lock, so goroutines counting different users never wait
// for one another.
type Counter struct {
	shards [shards]shard
}

// New returns an empty Counter.
func New() *Counter {
	c := &Counter{}
	for i := range c.shards {
		c.shards[i].counts = make(map[uint64]int)
	}
	return c
}

// Inc counts one event for user id.
func (c *Counter) Inc(id uint64) {
	s := &c.shards[id%shards]
	s.mu.Lock()
	s.counts[id]++
	s.mu.Unlock()
}

// Count returns how many events user id has had.
func (c *Counter) Count(id uint64) int {
	s := &c.shards[id%shards]
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.counts[id]
}
