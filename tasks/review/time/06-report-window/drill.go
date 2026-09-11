// Package drill — C13/06 report-window.
package drill

import "time"

// DayBucket is one day's slot in a report.
type DayBucket struct {
	Start time.Time
	Label string
}

// Buckets builds n consecutive daily buckets starting at from.
func Buckets(from time.Time, n int) []DayBucket {
	out := make([]DayBucket, 0, n)
	cur := from
	for i := 0; i < n; i++ {
		out = append(out, DayBucket{
			Start: cur,
			Label: cur.Format("2006/13/02"),
		})
		cur = cur.Add(24 * time.Hour)
	}
	return out
}

// IndexOf returns the bucket index whose start equals want, or -1.
func IndexOf(buckets []DayBucket, want time.Time) int {
	for i, b := range buckets {
		if b.Start == want {
			return i
		}
	}
	return -1
}
