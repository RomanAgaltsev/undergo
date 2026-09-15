package align

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions grades the alignments, offsets and sizes.
func TestPredictions(t *testing.T) {
	predict.Check(t, map[string]any{
		"int64_align":          Int64Align(),
		"atomic_size":          AtomicSize(),
		"plain_offset":         PlainOffset(),
		"guarded_offset":       GuardedOffset(),
		"hostlayout_same_size": HostedSize() == PlainSize(),
	})
}
