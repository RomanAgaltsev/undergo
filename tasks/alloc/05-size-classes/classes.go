// Package classes measures what the allocator actually charges for a slice of
// a given length, which is not the length.
//
// Importing testing from a non-test file is deliberate: the measurement is the
// subject of the task, so it belongs in the code you are reading.
package classes

import "testing"

// Sinks keep the allocations from being optimised away. They are exported
// because a package-level variable that is only ever written is flagged as
// unused, and a benchmark sink is exactly that.
var (
	SinkBytes []byte
	SinkPtrs  []*int
)

// BytesPerAlloc reports the bytes charged for one make([]byte, n).
func BytesPerAlloc(n int) int64 {
	r := testing.Benchmark(func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			SinkBytes = make([]byte, n)
		}
	})
	return r.AllocedBytesPerOp()
}

// PointerBytesPerAlloc reports the same for a one-element slice of pointers,
// which is scannable and so cannot share a block with anything else.
func PointerBytesPerAlloc() int64 {
	r := testing.Benchmark(func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			SinkPtrs = make([]*int, 1)
		}
	})
	return r.AllocedBytesPerOp()
}
