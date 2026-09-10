// Package drill — C13/08 truncate-day.
package drill

import "time"

// StartOfDay returns local midnight at the start of t's day.
func StartOfDay(t time.Time) time.Time {
	return t.Truncate(24 * time.Hour)
}
