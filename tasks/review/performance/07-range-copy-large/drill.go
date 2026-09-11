// Package drill — C7/07 range-copy-large.
package drill

// Record is a large value type.
type Record struct {
	ID     int
	Name   string
	Data   [256]byte
	Labels [16]string
}

// CountActive counts the records for which active reports true.
func CountActive(records []Record, active func(Record) bool) int {
	n := 0
	for _, r := range records {
		if active(r) {
			n++
		}
	}
	return n
}
