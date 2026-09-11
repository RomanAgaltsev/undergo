// Package drill — C9/09 schedule.
package drill

import "time"

// SameInstant reports whether a and b are the same moment.
func SameInstant(a, b time.Time) bool {
	return a == b
}

// DaysBetween returns the number of days between start and end.
func DaysBetween(start, end time.Time) int {
	return int(end.Sub(start).Hours()) / 24
}

// Weekday returns the weekday at the given offset from Sunday.
func Weekday(offset int) time.Weekday {
	return time.Weekday(offset % 7)
}
