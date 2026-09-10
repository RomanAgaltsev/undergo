// Package drill — C10/02 happy-path-only.
package drill

import "errors"

// Divide returns a/b, or an error when b is zero.
func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("divide by zero")
	}
	return a / b, nil
}
