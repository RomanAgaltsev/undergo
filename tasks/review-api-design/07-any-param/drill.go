// Package drill — C6/07 any-param.
package drill

// Sum adds up the numbers in values, which may be a []int or a []int64.
func Sum(values any) int {
	switch v := values.(type) {
	case []int:
		total := 0
		for _, n := range v {
			total += n
		}
		return total
	case []int64:
		total := 0
		for _, n := range v {
			total += int(n)
		}
		return total
	default:
		return 0
	}
}
