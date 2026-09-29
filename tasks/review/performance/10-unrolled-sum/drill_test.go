package drill

import "testing"

var input = make([]int, 1<<16)

func BenchmarkSum(b *testing.B) {
	for b.Loop() {
		Sum(input)
	}
}

func BenchmarkSumUnrolled(b *testing.B) {
	for b.Loop() {
		SumUnrolled(input)
	}
}
