// Package drill — C9/07 truncating-conversion.
package drill

import "strconv"

// ParseAmount parses a decimal amount (in cents) from s.
func ParseAmount(s string) (int32, error) {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return int32(n), nil
}
