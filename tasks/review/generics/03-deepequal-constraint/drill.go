// Package drill — C11/03 deepequal-constraint.
package drill

import "reflect"

// IndexOf returns the index of the first element equal to target, or -1.
func IndexOf[T any](s []T, target T) int {
	for i, v := range s {
		if reflect.DeepEqual(v, target) {
			return i
		}
	}
	return -1
}
