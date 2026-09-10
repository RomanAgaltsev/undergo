// Package drill — C13/01 equal-operator.
package drill

import "time"

// SameInstant reports whether a and b represent the same time.
func SameInstant(a, b time.Time) bool {
	return a == b
}
