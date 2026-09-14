package sum

import (
	"runtime"
	"testing"
)

func TestSumsCorrectly(t *testing.T) {
	cases := []struct {
		name string
		in   []int64
		want int64
	}{
		{"empty", nil, 0},
		{"one", []int64{42}, 42},
		{"several", []int64{1, 2, 3, 4, 5}, 15},
		{"negatives", []int64{-5, 10, -20}, -15},
		{"odd length", []int64{1, 2, 3}, 6},
	}
	for _, c := range cases {
		if got := SumInt64(c.in); got != c.want {
			t.Errorf("%s: SumInt64(%v) = %d, want %d", c.name, c.in, got, c.want)
		}
	}
}

func TestHandlesALongSlice(t *testing.T) {
	xs := make([]int64, 1000)
	var want int64
	for i := range xs {
		xs[i] = int64(i)
		want += int64(i)
	}
	if got := SumInt64(xs); got != want {
		t.Errorf("SumInt64(0..999) = %d, want %d", got, want)
	}
}

// TestDoesNotAllocate is what //go:noescape is claiming: the slice does not
// escape into the assembly, so passing it costs nothing.
func TestDoesNotAllocate(t *testing.T) {
	xs := make([]int64, 64)

	avg := testing.AllocsPerRun(100, func() {
		Sink = SumInt64(xs)
	})
	if avg != 0 {
		t.Errorf("SumInt64 allocated %v times, want 0", avg)
	}
}

// TestRunsOnThisArchitecture records which implementation was under test, so a
// pass on arm64 is never mistaken for a pass on the assembly.
func TestRunsOnThisArchitecture(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		t.Skipf("the assembly is amd64-only; this is %s, so the Go fallback is what ran", runtime.GOARCH)
	}
}
