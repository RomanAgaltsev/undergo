// Package drill — C9/02 midpoint-overflow.
package drill

// Search returns the index of target in the ascending-sorted slice, or -1.
func Search(sorted []int, target int) int {
	lo, hi := 0, len(sorted)-1
	for lo <= hi {
		mid := (lo + hi) / 2
		switch {
		case sorted[mid] == target:
			return mid
		case sorted[mid] < target:
			lo = mid + 1
		default:
			hi = mid - 1
		}
	}
	return -1
}
