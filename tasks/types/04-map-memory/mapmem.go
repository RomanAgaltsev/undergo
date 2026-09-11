// Package mapmem builds a large map and asks what it really costs.
//
// You do not implement anything here. Work out the answers, write them into
// prediction.yaml, then verify.
package mapmem

// N is how many entries the map holds.
const N = 1 << 20

// Build returns a map[int64]int64 with N entries, built without a size hint.
func Build() map[int64]int64 {
	m := map[int64]int64{}
	for i := range int64(N) {
		m[i] = i
	}
	return m
}

// BuildWithHint returns the same map, built with the size known up front.
func BuildWithHint() map[int64]int64 {
	m := make(map[int64]int64, N)
	for i := range int64(N) {
		m[i] = i
	}
	return m
}
