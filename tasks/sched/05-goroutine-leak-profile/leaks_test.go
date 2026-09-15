package leaks

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions grades which shapes the profile named, and the two counts. It
// never prints the profile text.
func TestPredictions(t *testing.T) {
	r := Survey()

	predict.Check(t, map[string]any{
		"blocked_send_named":   r.Named["BlockedSend"],
		"nil_recv_named":       r.Named["NilRecv"],
		"waitgroup_named":      r.Named["WaitGroupNeverDone"],
		"mutex_named":          r.Named["MutexNeverUnlocked"],
		"context_named":        r.Named["ContextNeverDone"],
		"count_before_writeto": r.CountBefore,
		"count_after_writeto":  r.CountAfter,
	})
}
