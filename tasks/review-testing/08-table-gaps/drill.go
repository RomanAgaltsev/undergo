// Package drill — C10/08 table-gaps.
package drill

// Classify buckets a score. ok is false for an invalid (negative) score.
func Classify(score int) (bucket string, ok bool) {
	switch {
	case score < 0:
		return "invalid", false
	case score < 50:
		return "low", true
	case score < 80:
		return "mid", true
	default:
		return "high", true
	}
}
