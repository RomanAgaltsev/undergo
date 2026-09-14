package edges

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// Runs is how many times each program is executed. A guaranteed edge holds on
// every single run; observing an unguaranteed one holding proves nothing
// either way, which is the lesson rather than a flaw in the measurement.
const Runs = 500

func alwaysHolds(f func() bool) bool {
	for range Runs {
		if !f() {
			return false
		}
	}
	return true
}

// TestPredictions grades which programs held on every run.
func TestPredictions(t *testing.T) {
	predict.Check(t, map[string]any{
		"unbuffered_send_always": alwaysHolds(UnbufferedSend),
		"buffered_full_always":   alwaysHolds(BufferedSendFull),
		"close_always":           alwaysHolds(CloseChannel),
		"mutex_always":           alwaysHolds(MutexUnlockLock),
		"once_always":            alwaysHolds(OnceDo),
		"waitgroup_always":       alwaysHolds(WaitGroupWait),
		"atomic_always":          alwaysHolds(AtomicStoreLoad),
		"sleep_always":           alwaysHolds(SleepOnly),
		"no_sync_always":         alwaysHolds(NoSync),
	})
}
