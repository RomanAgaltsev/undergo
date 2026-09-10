package join

import (
	"strings"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/optimize"
)

var corpus = []string{"alpha", "beta", "gamma", "delta", "epsilon", "zeta"}

func TestJoinIntoMatchesStringsJoin(t *testing.T) {
	tests := [][]string{
		nil,
		{},
		{"only"},
		corpus,
		{"", "", ""},
	}
	for _, parts := range tests {
		want := strings.Join(parts, ", ")
		got := string(JoinInto(nil, parts, ", "))
		if got != want {
			t.Errorf("JoinInto(%q) = %q, want %q", parts, got, want)
		}
	}
}

func TestJoinIntoAppendsToExistingContent(t *testing.T) {
	dst := []byte("prefix: ")
	got := string(JoinInto(dst, []string{"a", "b"}, "-"))
	if got != "prefix: a-b" {
		t.Errorf("JoinInto = %q, want %q", got, "prefix: a-b")
	}
}

// BenchmarkBaseline is the reference point named in task.yaml.
func BenchmarkBaseline(b *testing.B) {
	for b.Loop() {
		_ = joinBaseline(corpus, ", ")
	}
}

// BenchmarkCandidate exercises your JoinInto with a pre-sized destination.
func BenchmarkCandidate(b *testing.B) {
	buf := make([]byte, 0, 256)
	for b.Loop() {
		buf = JoinInto(buf[:0], corpus, ", ")
	}
	_ = buf
}

// TestOptimizeTarget is the graded benchmark: with a pre-sized destination,
// JoinInto must allocate nothing per operation.
func TestOptimizeTarget(t *testing.T) {
	optimize.Check(t, optimize.MetricAllocs, 0, BenchmarkBaseline, BenchmarkCandidate)
}
