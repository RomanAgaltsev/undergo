package yield

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions grades how many yield calls each exit route completed, and
// what yield returned on the last of them.
func TestPredictions(t *testing.T) {
	complete, brk, gt := Complete(), Break(), Goto()
	lbl, ret, pan := LabelledBreak(), Return(), Panicked()

	predict.Check(t, map[string]any{
		"complete_calls":       complete.Calls,
		"break_calls":          brk.Calls,
		"goto_calls":           gt.Calls,
		"labelled_break_calls": lbl.Calls,
		"return_calls":         ret.Calls,
		"panic_calls":          pan.Calls,
		"break_last_return":    brk.LastReturn(),
		"panic_last_return":    pan.LastReturn(),
	})
}
