// Package drill — C11/01 any-numeric-switch.
package drill

// Sum totals a slice of numbers. It is generic over T so callers can pass
// []int or []float64.
func Sum[T any](vals []T) float64 {
	var total float64
	for _, v := range vals {
		switch x := any(v).(type) {
		case int:
			total += float64(x)
		case int64:
			total += float64(x)
		case float64:
			total += x
		}
	}
	return total
}
