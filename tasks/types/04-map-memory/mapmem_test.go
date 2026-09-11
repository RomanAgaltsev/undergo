package mapmem

import (
	"math"
	"runtime"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// heapBytes reports the live heap after a collection, so the measurement either
// side of Build sees settled memory rather than garbage in flight.
func heapBytes() uint64 {
	runtime.GC()
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return ms.HeapAlloc
}

// bytesPerEntry measures the heap a full map occupies, divided by its entries.
func bytesPerEntry(build func() map[int64]int64) float64 {
	before := heapBytes()
	m := build()
	after := heapBytes()
	runtime.KeepAlive(m)
	return float64(after-before) / float64(N)
}

// TestPredictions compares your prediction.yaml against what the map really
// costs. It never prints the measured numbers.
func TestPredictions(t *testing.T) {
	perEntry := bytesPerEntry(Build)

	predict.Check(t, map[string]any{
		// Rounded to the nearest 8 bytes: the question is what a map costs per
		// entry, not whether you can match a fraction.
		"bytes_per_entry_rounded": int(math.Round(perEntry/8) * 8),
		"payload_bytes_per_entry": 16,
	})
}
