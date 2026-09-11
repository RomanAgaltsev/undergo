// Package drill — C11/02 predicate-any.
package drill

// Filter returns the elements of s for which keep reports true.
func Filter[T any](s []T, keep func(any) bool) []T {
	out := make([]T, 0, len(s))
	for _, v := range s {
		if keep(v) {
			out = append(out, v)
		}
	}
	return out
}
