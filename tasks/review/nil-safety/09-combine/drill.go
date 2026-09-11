// Package drill — C2/09 combine.
package drill

// Combine appends extra onto existing, runs a hook, and returns the first
// element of the merged slice.
func Combine(existing, extra []int, hook func()) int {
	merged := append(existing, extra...)
	hook()
	return merged[0]
}

// safeFirst returns the first element or 0 for an empty slice.
func safeFirst(xs []int) int {
	if len(xs) == 0 {
		return 0
	}
	return xs[0]
}
