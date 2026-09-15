package publish

import (
	"sync"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// Attempts, Goroutines and Rounds are sized so a broken strategy is caught
// reliably. Each attempt starts a fresh Getter, because first use is the only
// moment a publication bug can happen.
const (
	Attempts   = 400
	Goroutines = 8
	Rounds     = 200
)

// duplicateInit reports whether callers ever received more than one distinct
// Config from the same Getter — that is, whether the object was built twice.
//
// It also fails the test outright if any caller received an object that was not
// fully built. That is a different bug and no strategy here should exhibit it;
// it is checked rather than graded so that a surprise cannot pass unnoticed.
func duplicateInit(t *testing.T, newGetter func() Getter) bool {
	t.Helper()

	for range Attempts {
		s := newGetter()
		var mu sync.Mutex
		seen := map[*Config]bool{}

		var wg sync.WaitGroup
		wg.Add(Goroutines)
		for g := range Goroutines {
			go func(marker int) {
				defer wg.Done()
				for range Rounds {
					c := s.Get(marker)
					if !c.Complete() {
						t.Error("a caller received a Config that was not fully built")
						return
					}
					mu.Lock()
					seen[c] = true
					mu.Unlock()
				}
			}(g + 1)
		}
		wg.Wait()

		if len(seen) > 1 {
			return true
		}
	}
	return false
}

// TestPredictions grades which strategies ever built the object more than once.
func TestPredictions(t *testing.T) {
	predict.Check(t, map[string]any{
		"plain_pointer_duplicates":  duplicateInit(t, func() Getter { return &PlainPointer{} }),
		"atomic_pointer_duplicates": duplicateInit(t, func() Getter { return &AtomicPointer{} }),
		"once_duplicates":           duplicateInit(t, func() Getter { return &OnceGuarded{} }),
		"mutex_duplicates":          duplicateInit(t, func() Getter { return &MutexGuarded{} }),
		"double_checked_duplicates": duplicateInit(t, func() Getter { return &DoubleChecked{} }),
	})
}
