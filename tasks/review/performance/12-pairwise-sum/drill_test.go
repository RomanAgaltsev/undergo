package drill

import (
	"strconv"
	"testing"
)

func TestTotal(t *testing.T) {
	cases := []struct {
		xs   []int
		want int
	}{
		{nil, 0},
		{[]int{1, 2}, 3},
		{[]int{1, 2, 3, 4}, 10},
	}
	for _, c := range cases {
		if got := Total(c.xs); got != c.want {
			t.Errorf("Total(%v) = %d, want %d", c.xs, got, c.want)
		}
	}
}

var data = make([]int, 1<<16)

func BenchmarkTotal(b *testing.B) {
	for b.Loop() {
		Total(data)
	}
}

// BenchmarkReportBaseline is the old concatenating renderer, kept to show
// what Report saves.
func BenchmarkReportBaseline(b *testing.B) {
	s := ""
	for i := 0; i < b.N; i++ {
		if i > 0 {
			s += ","
		}
		s += strconv.Itoa(i)
	}
}

func BenchmarkReport(b *testing.B) {
	totals := make([]int, 1000)
	for b.Loop() {
		Report(totals)
	}
}
