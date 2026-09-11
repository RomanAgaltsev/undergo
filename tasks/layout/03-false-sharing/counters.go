// Package counters holds two shapes of per-worker counter: the obvious one and
// the padded one.
//
// You do not implement anything here. Work out the answers, write them into
// prediction.yaml, then verify — and run the benchmarks to see what the layout
// costs.
package counters

// CacheLine is the line size this task assumes: 64 bytes, which is what amd64
// and arm64 use. A machine with a different line size is excluded by
// requires.arch in task.yaml.
const CacheLine = 64

// Naive is one counter per worker, packed as tightly as the compiler likes.
type Naive struct {
	N uint64
}

// Padded pads each counter out so that no two of them share a line.
type Padded struct {
	N   uint64
	pad [CacheLine - 8]byte //nolint:unused // padding that is never read is the entire point
}

// Workers is how many counters the benchmark contends over.
const Workers = 8

// NaiveCounters and PaddedCounters are the arrays the benchmark hammers.
var (
	NaiveCounters  [Workers]Naive
	PaddedCounters [Workers]Padded
)
