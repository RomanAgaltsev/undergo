// Package drill — C7/04 defer-in-hot-loop.
package drill

// Sum adds the numbers, invoking onEach after processing each one.
func Sum(nums []int, onEach func()) int {
	total := 0
	for _, n := range nums {
		defer onEach()
		total += n
	}
	return total
}
