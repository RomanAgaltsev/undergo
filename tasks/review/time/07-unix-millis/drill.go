// Package drill — C13/07 unix-millis.
package drill

import "time"

// At converts a unix-milliseconds timestamp to a time.Time.
func At(msTimestamp int64) time.Time {
	return time.Unix(msTimestamp, 0)
}
