// Package drill — C10/04 only-error-nil.
package drill

import "strconv"

// Parse converts s to an int.
func Parse(s string) (int, error) {
	return strconv.Atoi(s)
}
