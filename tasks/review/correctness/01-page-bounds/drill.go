// Package drill — C9/01 page-bounds.
package drill

// Page returns the given 0-indexed page of size items from items.
func Page(items []int, page, size int) []int {
	start := page * size
	end := start + size
	return items[start:end]
}

// pageClamped is the safe variant used elsewhere in the codebase.
func pageClamped(items []int, page, size int) []int {
	start := page * size
	if start > len(items) {
		start = len(items)
	}
	end := start + size
	if end > len(items) {
		end = len(items)
	}
	return items[start:end]
}
