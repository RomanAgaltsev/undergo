package counters

import (
	"sync"
	"testing"
	"unsafe"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// lineOf returns which cache line an address falls on, counting from the start
// of the array it belongs to.
func lineOf(base, addr uintptr) uintptr { return (addr - base) / CacheLine }

func sameLine[T any](arr *[Workers]T, i, j int) bool {
	base := uintptr(unsafe.Pointer(arr))
	return lineOf(base, uintptr(unsafe.Pointer(&arr[i]))) ==
		lineOf(base, uintptr(unsafe.Pointer(&arr[j])))
}

// TestPredictions grades the layout, which is deterministic. The throughput
// cliff it causes is the subject of the benchmarks below and of the first
// written question, but it is not graded: a ratio measured on an unknown number
// of cores is not something you can be exactly right about.
func TestPredictions(t *testing.T) {
	predict.Check(t, map[string]any{
		"sizeof_naive":            unsafe.Sizeof(Naive{}),
		"sizeof_padded":           unsafe.Sizeof(Padded{}),
		"naive_per_line":          CacheLine / unsafe.Sizeof(Naive{}),
		"naive_neighbours_share":  sameLine(&NaiveCounters, 0, 1),
		"padded_neighbours_share": sameLine(&PaddedCounters, 0, 1),
	})
}

// BenchmarkNaive and BenchmarkPadded exist to be run, not graded: they show the
// cliff your prediction describes. Run them with:
//
//	go test -bench . -cpu 8
//
// They increment without synchronisation on purpose — this measures hardware
// coherency traffic, not correct concurrent code. Never ship either shape.
func BenchmarkNaive(b *testing.B) {
	benchmark(b, func(i int) *uint64 { return &NaiveCounters[i].N })
}

func BenchmarkPadded(b *testing.B) {
	benchmark(b, func(i int) *uint64 { return &PaddedCounters[i].N })
}

func benchmark(b *testing.B, at func(int) *uint64) {
	var wg sync.WaitGroup
	for w := range Workers {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			p := at(w)
			for range b.N / Workers {
				*p++
			}
		}(w)
	}
	wg.Wait()
}
