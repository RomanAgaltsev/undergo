// Package drill — C6/04 error-first-return.
package drill

import "strconv"

// Parse converts s to an int.
func Parse(s string) (error, int) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return err, 0
	}
	return nil, n
}

// ParseOK converts s to an int.
func ParseOK(s string) (int, error) {
	return strconv.Atoi(s)
}
