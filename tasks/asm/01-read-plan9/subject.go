// Package subject holds four small functions whose compiled form you are asked
// to read.
//
// You do not implement anything here. Work out the answers, write them into
// prediction.yaml, then verify.
package subject

// AddInts returns a + b.
func AddInts(a, b int) int { return a + b }

// LenOf returns the length of a string.
func LenOf(s string) int { return len(s) }

// FirstOf returns the first element of a slice, or zero.
func FirstOf(xs []int) int {
	if len(xs) == 0 {
		return 0
	}
	return xs[0]
}

// DivMod returns the quotient and remainder.
func DivMod(a, b int) (int, int) { return a / b, a % b }
