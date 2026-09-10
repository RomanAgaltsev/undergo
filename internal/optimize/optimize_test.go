package optimize

import (
	"strings"
	"testing"
)

// sink is package-level so the allocation below escapes and reaches the heap.
// Assigning to a blank identifier instead would let it stay on the stack, and
// the "allocating" benchmark would report zero allocations.
var sink []byte

func allocating(b *testing.B) {
	for b.Loop() {
		sink = make([]byte, 64)
	}
}

func free(b *testing.B) {
	var n int
	for b.Loop() {
		n++
	}
	_ = n
}

func TestMeasureAllocsIsAbsolute(t *testing.T) {
	got, err := Measure(MetricAllocs, 0, allocating, free)
	if err != nil {
		t.Fatal(err)
	}
	if !got.Pass {
		t.Errorf("an allocation-free candidate failed a target of 0 allocs (value %v)", got.Value)
	}

	got, err = Measure(MetricAllocs, 0, free, allocating)
	if err != nil {
		t.Fatal(err)
	}
	if got.Pass {
		t.Error("an allocating candidate passed a target of 0 allocs")
	}
}

func TestMeasureNsIsARatio(t *testing.T) {
	got, err := Measure(MetricNs, 1.0, allocating, free)
	if err != nil {
		t.Fatal(err)
	}
	if got.Value <= 0 {
		t.Fatalf("ratio = %v, want > 0", got.Value)
	}
	if !got.Pass {
		t.Errorf("a faster candidate missed a ratio target of 1.0 (value %v)", got.Value)
	}
}

// A candidate whose per-op cost is under a nanosecond must still produce a
// usable ratio: testing.BenchmarkResult.NsPerOp truncates to whole nanoseconds
// and would report 0 for both sides.
func TestMeasureNsSurvivesSubNanosecondOps(t *testing.T) {
	got, err := Measure(MetricNs, 1.0, free, free)
	if err != nil {
		t.Fatal(err)
	}
	if got.Value <= 0 {
		t.Fatalf("ratio = %v, want > 0 even below 1 ns/op", got.Value)
	}
}

func TestMeasureRejectsAnUnknownMetric(t *testing.T) {
	_, err := Measure(Metric("binary_size"), 1, free, free)
	if err == nil || !strings.Contains(err.Error(), "binary_size") {
		t.Fatalf("err = %v, want it to name the unknown metric", err)
	}
}
