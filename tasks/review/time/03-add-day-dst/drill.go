// Package drill — C13/03 add-day-dst.
package drill

import "time"

// Tomorrow returns the same wall-clock time on the next calendar day.
func Tomorrow(t time.Time) time.Time {
	return t.Add(24 * time.Hour)
}
