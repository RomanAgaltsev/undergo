package shapes

import "testing"

var sink float64

// BenchmarkProfileWorkload is the workload the profile is collected from: Hot,
// called overwhelmingly with a Rect. It is a benchmark rather than a test so
// that it runs only when asked for by name.
func BenchmarkProfileWorkload(b *testing.B) {
	shapes := []Shape{Rect{2, 3}, Rect{4, 5}, Rect{6, 7}, Rect{8, 9}}
	for b.Loop() {
		for _, s := range shapes {
			sink += Hot(s)
		}
	}
}
