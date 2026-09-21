// Package optimize grades an optimize-mode task by running the reference
// baseline and the solver's candidate in the same process, on the same machine,
// in the same run.
//
// Absolute numbers are meaningless between machines; that is why a time target
// is expressed as a ratio to the baseline rather than a nanosecond count.
package optimize

import (
	"fmt"
	"testing"
)

// Metric is what an optimize task is scored on.
type Metric string

// The metrics supported in v1. Allocation and byte targets are absolute
// per-op counts; the time target is a ratio to the baseline.
const (
	MetricNs     Metric = "ns"
	MetricAllocs Metric = "allocs"
	MetricBytes  Metric = "bytes"
)

// Result is one graded comparison.
type Result struct {
	Baseline  testing.BenchmarkResult
	Candidate testing.BenchmarkResult
	Value     float64
	Target    float64
	Pass      bool
}

// Measure runs both benchmarks and scores the candidate.
func Measure(m Metric, target float64, baseline, candidate func(*testing.B)) (Result, error) {
	base := testing.Benchmark(baseline)
	cand := testing.Benchmark(candidate)

	r := Result{Baseline: base, Candidate: cand, Target: target}
	switch m {
	case MetricAllocs:
		r.Value = float64(cand.AllocsPerOp())
	case MetricBytes:
		r.Value = float64(cand.AllocedBytesPerOp())
	case MetricNs:
		// Not NsPerOp: it truncates to whole nanoseconds, so a sub-nanosecond
		// operation on either side would collapse the ratio to 0 or divide by 0.
		baseNs, err := nsPerOp(base)
		if err != nil {
			return r, fmt.Errorf("optimize: baseline: %w", err)
		}
		candNs, err := nsPerOp(cand)
		if err != nil {
			return r, fmt.Errorf("optimize: candidate: %w", err)
		}
		r.Value = candNs / baseNs
	default:
		return r, fmt.Errorf("optimize: unknown metric %q", m)
	}
	r.Pass = r.Value <= target
	return r, nil
}

// TB is the part of *testing.T that Check uses.
//
// As in internal/predict, it is here so that Check can be tested: gate 2 proves
// only that a reference solution passes, so a Check that never failed would be
// invisible to every gate. *testing.T satisfies it, so a task's optimize.Check(t,
// ...) is unchanged.
type TB interface {
	Helper()
	Errorf(format string, args ...any)
	Fatal(args ...any)
	Logf(format string, args ...any)
}

// Check runs Measure and fails the test when the target is missed.
func Check(t TB, m Metric, target float64, baseline, candidate func(*testing.B)) {
	t.Helper()

	r, err := Measure(m, target, baseline, candidate)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("baseline: %d ns/op, %d B/op, %d allocs/op",
		r.Baseline.NsPerOp(), r.Baseline.AllocedBytesPerOp(), r.Baseline.AllocsPerOp())
	t.Logf("candidate: %d ns/op, %d B/op, %d allocs/op",
		r.Candidate.NsPerOp(), r.Candidate.AllocedBytesPerOp(), r.Candidate.AllocsPerOp())
	if !r.Pass {
		t.Errorf("%s = %.3f, want <= %.3f", m, r.Value, r.Target)
	}
}

// nsPerOp is testing.BenchmarkResult.NsPerOp in floating point, so that costs
// below one nanosecond stay measurable.
func nsPerOp(r testing.BenchmarkResult) (float64, error) {
	if r.N <= 0 || r.T <= 0 {
		return 0, fmt.Errorf("benchmark measured nothing (N=%d, T=%v)", r.N, r.T)
	}
	return float64(r.T.Nanoseconds()) / float64(r.N), nil
}
