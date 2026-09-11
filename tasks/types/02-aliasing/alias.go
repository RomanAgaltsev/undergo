// Package alias performs a fixed sequence of slice operations on one backing
// array. Predict what each view holds at the end.
//
// You do not implement anything here. Work out the answers, write them into
// prediction.yaml, then verify.
package alias

// Run performs the sequence below and returns the four views it produced.
//
//	base       := []int{0, 1, 2, 3, 4, 5, 6, 7}
//	twoIndex   := base[2:4]        // capacity runs to the end of base
//	threeIndex := base[2:4:4]      // capacity capped at 4
//	twoIndex    = append(twoIndex, 99)
//	grown      := append(threeIndex, 77)
//	copy(base[1:], base[:3])
func Run() (base, twoIndex, threeIndex, grown []int) {
	base = []int{0, 1, 2, 3, 4, 5, 6, 7}

	twoIndex = base[2:4]
	threeIndex = base[2:4:4]

	twoIndex = append(twoIndex, 99)
	//nolint:gocritic // appendAssign is the exercise: this append detaches from base.
	grown = append(threeIndex, 77)

	copy(base[1:], base[:3])

	return base, twoIndex, threeIndex, grown
}
