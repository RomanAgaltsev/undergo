package optimize

import (
	"fmt"
	"strings"
	"testing"
)

var _ TB = (*testing.T)(nil)

// recorder stands in for *testing.T so a grading can be inspected rather than
// obeyed. Fatal does not end the goroutine here because every Check path that
// reaches it returns immediately afterwards.
type recorder struct {
	errs  []string
	logs  []string
	fatal string
}

func (r *recorder) Helper() {}

func (r *recorder) Errorf(format string, args ...any) {
	r.errs = append(r.errs, fmt.Sprintf(format, args...))
}

func (r *recorder) Logf(format string, args ...any) {
	r.logs = append(r.logs, fmt.Sprintf(format, args...))
}

func (r *recorder) Fatal(args ...any) { r.fatal = fmt.Sprint(args...) }

func (r *recorder) failed() bool { return len(r.errs) > 0 || r.fatal != "" }

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

// Check has to fail a candidate that misses its target, and gate 2 can never
// show that: it runs only reference solutions, which meet their targets by
// construction. The metric is allocations because that is the one this grader
// can judge deterministically — a time ratio is a property of the machine, and
// asserting one here would be the same defect the answer-form rule keeps out of
// the catalogue.
func TestCheckFailsACandidateThatMissesItsTarget(t *testing.T) {
	tests := []struct {
		name       string
		metric     Metric
		target     float64
		baseline   func(*testing.B)
		candidate  func(*testing.B)
		wantFailed bool
	}{
		{
			name:   "candidate allocates where none is allowed",
			metric: MetricAllocs, target: 0,
			baseline: free, candidate: allocating,
			wantFailed: true,
		},
		{
			name:   "candidate is allocation-free as required",
			metric: MetricAllocs, target: 0,
			baseline: allocating, candidate: free,
		},
		{
			name:   "candidate allocates more bytes than allowed",
			metric: MetricBytes, target: 0,
			baseline: free, candidate: allocating,
			wantFailed: true,
		},
		{
			name:   "unknown metric is fatal rather than a silent pass",
			metric: Metric("binary_size"), target: 1,
			baseline: free, candidate: free,
			wantFailed: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := &recorder{}
			Check(r, tc.metric, tc.target, tc.baseline, tc.candidate)
			if got := r.failed(); got != tc.wantFailed {
				t.Errorf("failed = %v, want %v (errs=%v fatal=%q)",
					got, tc.wantFailed, r.errs, r.fatal)
			}
		})
	}
}
