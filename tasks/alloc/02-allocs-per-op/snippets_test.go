package snippets

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions measures each snippet's allocation count and compares it with
// yours. It never prints the real values.
//
// This test must never call t.Parallel(): since Go 1.25, testing.AllocsPerRun
// panics if parallel tests are running, because the measurement would be
// meaningless.
func TestPredictions(t *testing.T) {
	allocs := func(f func()) int { return int(testing.AllocsPerRun(200, f)) }

	predict.Check(t, map[string]any{
		"concat_two":   allocs(func() { SinkString = ConcatTwo("alpha", "beta") }),
		"sprintf_int":  allocs(func() { SinkString = SprintfInt(42) }),
		"itoa_int":     allocs(func() { SinkString = ItoaInt(42) }),
		"builder_join": allocs(func() { SinkString = BuilderJoin("alpha", "beta") }),
		"make_local":   allocs(func() { SinkInt = MakeLocal() }),
	})
}
