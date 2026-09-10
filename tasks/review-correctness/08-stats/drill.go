// Package drill — C9/08 stats.
package drill

// Summarize returns the average-as-percentage and the total of xs.
func Summarize(xs []int) (avgPct, total int) {
	var sum int32
	for i := 0; i <= len(xs); i++ {
		sum += int32(xs[i])
	}
	total = int(sum)
	avgPct = total / len(xs) * 100
	return avgPct, total
}
