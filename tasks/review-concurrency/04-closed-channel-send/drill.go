// Package drill — C1/04 closed-channel-send.
package drill

// Broadcast sends every value to out, closes it, and (for an empty input)
// sends a sentinel -1 to signal "nothing to do".
func Broadcast(values []int, out chan int) {
	for _, v := range values {
		out <- v
	}
	close(out)
	if len(values) == 0 {
		out <- -1
	}
}
