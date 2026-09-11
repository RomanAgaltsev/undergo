package growth

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// after returns the capacity that follows v in the sequence, or -1 if v never
// appears or is last.
func after(caps []int, v int) int {
	for i, c := range caps {
		if c == v && i+1 < len(caps) {
			return caps[i+1]
		}
	}
	return -1
}

// TestPredictions compares your prediction.yaml against the capacities append
// actually chose. It never prints them.
func TestPredictions(t *testing.T) {
	caps := Sequence(600)

	at := func(i int) int {
		if i < len(caps) {
			return caps[i]
		}
		return -1
	}

	predict.Check(t, map[string]any{
		"first_cap":     at(0),
		"third_cap":     at(2),
		"cap_after_256": after(caps, 256),
		"cap_after_512": after(caps, 512),
		"distinct_caps": len(caps),
	})
}
