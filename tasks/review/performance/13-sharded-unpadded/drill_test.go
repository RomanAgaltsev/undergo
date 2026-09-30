package drill

import (
	"sync/atomic"
	"testing"
)

// BenchmarkInc has every goroutine count events for a user of its own, so no
// two goroutines ever want the same shard.
func BenchmarkInc(b *testing.B) {
	c := New()
	var next atomic.Uint64
	b.RunParallel(func(pb *testing.PB) {
		id := next.Add(1)
		for pb.Next() {
			c.Inc(id)
		}
	})
}
