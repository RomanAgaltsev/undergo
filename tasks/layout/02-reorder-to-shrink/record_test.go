package record

import (
	"math"
	"testing"
	"unsafe"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions compares your prediction.yaml against the layout the
// compiler actually chose. It never prints the real values.
func TestPredictions(t *testing.T) {
	event := unsafe.Sizeof(Event{})
	packed := unsafe.Sizeof(Packed{})
	waste := event - packed

	predict.Check(t, map[string]any{
		"sizeof_event":      event,
		"sizeof_packed":     packed,
		"offset_event_kind": unsafe.Offsetof(Event{}.Kind),
		"waste_per_record":  waste,
		// Rounded to whole mebibytes: the arithmetic has to be right, but you
		// are not asked to match a twelve-digit fraction.
		"waste_mib": int(math.Round(float64(uint64(waste)*Million) / (1 << 20))),
	})
}
