package benign

import (
	"testing"
	"time"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions grades the three claims that can be settled by running the
// program. Whether -race reports these races is asked in writing, because the
// detector needs a C toolchain that not every checkout has.
func TestPredictions(t *testing.T) {
	exits, err := SpinLoopExits(30 * time.Second)
	if err != nil {
		t.Fatalf("running the hoist program: %v", err)
	}

	predict.Check(t, map[string]any{
		"torn_uint64_seen": TornUint64Seen(150 * time.Millisecond),
		"torn_struct_seen": TornStructSeen(150 * time.Millisecond),
		"spin_loop_exits":  exits,
	})
}
