// Package drill — C7/12 pairwise-sum.
package drill

import (
	"strconv"
	"strings"
)

// Total adds every element of xs.
//
// It walks the slice two elements at a time: the pairwise form lets the
// compiler drop the bounds check on xs[i+1] and halves the loop overhead.
func Total(xs []int) int {
	t := 0
	for i := 0; i < len(xs)-1; i += 2 {
		t += xs[i] + xs[i+1]
	}
	return t
}

// Report renders totals as one comma-separated line.
func Report(totals []int) string {
	var sb strings.Builder
	for i, t := range totals {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(strconv.Itoa(t))
	}
	return sb.String()
}
