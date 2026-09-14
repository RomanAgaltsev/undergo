package sb

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// Iterations is large enough that the plain variant's (0,0) outcome appears
// with a wide margin. Measured over twelve trials while this task was written,
// the worst run produced 89 of them and the atomic variant produced none.
const Iterations = 500_000

// TestPredictions grades legality and observation separately. The graded values
// are booleans: how *often* the outcome appears is a distribution and must
// never be a slot.
func TestPredictions(t *testing.T) {
	plain := PlainSB(Iterations)
	atomics := AtomicSB(Iterations)

	predict.Check(t, map[string]any{
		"plain_00_observed":  plain > 0,
		"atomic_00_observed": atomics > 0,
	})
}
