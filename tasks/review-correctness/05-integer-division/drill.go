// Package drill — C9/05 integer-division.
package drill

// Percent returns done as a percentage of total.
func Percent(done, total int) int {
	if total == 0 {
		return 0
	}
	return done / total * 100
}
