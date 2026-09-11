// Package drill — C13/04 parse-location.
package drill

import "time"

// ParseLocalStamp parses a "2006-01-02 15:04:05" timestamp that the operator
// entered in the server's local time zone.
func ParseLocalStamp(s string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", s)
}
