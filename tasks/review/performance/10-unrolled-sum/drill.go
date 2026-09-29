// Package drill — C7/10 unrolled-sum.
package drill

// Sum adds every element of xs.
func Sum(xs []int) int {
	total := 0
	for _, x := range xs {
		total += x
	}
	return total
}

// SumUnrolled adds every element of xs four at a time, so the loop runs a
// quarter as many iterations as Sum's.
func SumUnrolled(xs []int) int {
	total := 0
	for j := 0; j < len(xs)/4; j += 4 {
		total += xs[j] + xs[j+1] + xs[j+2] + xs[j+3]
	}
	return total
}
