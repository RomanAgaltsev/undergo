// Package drill — C13/02 bad-layout.
package drill

import "time"

// Stamp renders t as "YYYY-MM-DD HH:mm:ss" for a log line.
func Stamp(t time.Time) string {
	return t.Format("YYYY-MM-DD HH:mm:ss")
}
