// Package growth appends to a slice one element at a time and records every
// capacity the runtime chose along the way.
//
// You do not implement anything here. Work out the answers, write them into
// prediction.yaml, then verify.
package growth

// Sink makes the slice escape.
//
// This is not decoration. Without it the compiler may build the backing array on
// the stack, where the growth pattern is different — so the answer to this task
// would depend on an optimisation rather than on append. Question 3 in the
// README asks you to remove it and watch that happen.
var Sink []int

// Sequence appends n elements one at a time and returns every distinct capacity
// in the order it appeared.
func Sequence(n int) []int {
	var s []int
	var caps []int
	last := -1
	for i := range n {
		s = append(s, i)
		if cap(s) != last {
			caps = append(caps, cap(s))
			last = cap(s)
		}
	}
	Sink = s
	return caps
}
