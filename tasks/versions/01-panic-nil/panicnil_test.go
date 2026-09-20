package panicnil

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

func TestPredictions(t *testing.T) {
	got := map[string]any{}

	for _, c := range []struct {
		slot string
		line string
		env  []string
	}{
		{"type_at_go120", "1.20", nil},
		{"type_at_go121", "1.21", nil},
		{"type_at_go121_with_panicnil_1", "1.21", []string{"GODEBUG=panicnil=1"}},
		{"type_at_go120_with_panicnil_0", "1.20", []string{"GODEBUG=panicnil=0"}},
	} {
		out, err := RecoveredType(c.line, c.env...)
		if err != nil {
			t.Fatalf("%s: %v", c.slot, err)
		}
		got[c.slot] = out
	}

	predict.Check(t, got)
}
