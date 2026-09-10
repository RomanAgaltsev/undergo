package padding

import (
	"testing"
	"unsafe"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions compares your prediction.yaml against what the compiler
// actually did. It never prints the real values.
func TestPredictions(t *testing.T) {
	predict.Check(t, map[string]any{
		"sizeof_header":    unsafe.Sizeof(Header{}),
		"sizeof_reordered": unsafe.Sizeof(Reordered{}),
		"sizeof_trailing":  unsafe.Sizeof(Trailing{}),
		"offset_kind":      unsafe.Offsetof(Header{}.Kind),
		"align_header":     unsafe.Alignof(Header{}),
	})
}
