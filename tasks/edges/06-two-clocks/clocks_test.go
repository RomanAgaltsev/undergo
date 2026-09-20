package clocks

import (
	"testing"
	"time"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

func TestPredictions(t *testing.T) {
	ref := Reference()
	stripped := Stripped(ref)

	round, err := JSONRoundTrip(ref)
	if err != nil {
		t.Fatalf("json round trip: %v", err)
	}

	predict.Check(t, map[string]any{
		"now_carries_monotonic":           Monotonic(ref),
		"add_keeps_monotonic":             Monotonic(ref.Add(time.Second)),
		"utc_keeps_monotonic":             Monotonic(ref.UTC()),
		"round_zero_keeps_monotonic":      Monotonic(stripped),
		"json_round_trip_keeps_monotonic": Monotonic(round),
		"equal_method_says_same":          ref.Equal(stripped),
		// The == is the subject: this task grades the difference between
		// comparing instants and comparing memory. Equal is the right answer
		// everywhere else, which is why the linter is right to say so.
		"double_equals_says_same":      ref == stripped, //nolint:staticcheck // QF1009 is the lesson
		"sub_of_stripped_pair_matches": ref.Sub(stripped) == 0,
	})
}
