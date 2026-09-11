// Package drill — C9/06 float-equality.
package drill

// AllocatedFully reports whether the weights sum to exactly 1.0.
func AllocatedFully(weights []float64) bool {
	total := 0.0
	for _, w := range weights {
		total += w
	}
	return total == 1.0
}
