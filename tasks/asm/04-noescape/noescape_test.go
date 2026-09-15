//go:build amd64

package noescape

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

var sink int64

// TestPredictions grades the allocation counts and confirms the two functions
// really do compute the same thing.
func TestPredictions(t *testing.T) {
	predict.Check(t, map[string]any{
		"call_free_allocs": int(testing.AllocsPerRun(1000, func() { sink = CallFree() })),
		"call_kept_allocs": int(testing.AllocsPerRun(1000, func() { sink = CallKept() })),
		"same_result":      CallFree() == CallKept(),
	})
}
