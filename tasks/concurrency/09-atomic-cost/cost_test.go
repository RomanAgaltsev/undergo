package atomiccost

import (
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// margin is how far apart two costs must be before their order counts as a
// measurement and not as noise. Every ordering graded here was measured at
// least 1.8 times apart on every platform the task claims, in both builds.
const margin = 1.5

// TestPredictions grades three orderings. It never prints a cost.
func TestPredictions(t *testing.T) {
	if runtime.GOMAXPROCS(0) < 2 {
		t.Fatalf("this task needs GOMAXPROCS of at least 2: half of it is about cores competing for one value")
	}
	c := Measure()

	pairs := map[string][2]float64{
		"atomic_costs_more_than_plain": {c.Atomic, c.Plain},
		"mutex_costs_more_than_atomic": {c.Locked, c.Atomic},
		"sharing_slows_atomic":         {c.AtomicShared, c.Atomic},
	}
	measured := make(map[string]any, len(pairs))
	for slot, p := range pairs {
		r := p[0] / p[1]
		if r < margin && r > 1/margin {
			t.Fatalf("inconclusive: two of the costs came out too close to call on this run. " +
				"That is noise, not your prediction; run it again on a quieter machine.")
		}
		measured[slot] = r > 1
	}
	predict.Check(t, measured)
}

// showEnv switches TestShowCosts on.
const showEnv = "UNDERGO_SHOW_COSTS"

// TestShowCosts prints the costs the graded test compares. It is for after you
// have committed your predictions, and it only runs when you ask for it.
func TestShowCosts(t *testing.T) {
	if os.Getenv(showEnv) == "" {
		t.Skip("set " + showEnv + "=1 to print the costs")
	}
	c := Measure()
	fmt.Printf("GOMAXPROCS=%d, ns per increment:\n", runtime.GOMAXPROCS(0))
	fmt.Printf("  plain          %8.2f\n", c.Plain)
	fmt.Printf("  atomic         %8.2f  (%.1fx plain)\n", c.Atomic, c.Atomic/c.Plain)
	fmt.Printf("  mutex          %8.2f  (%.1fx atomic)\n", c.Locked, c.Locked/c.Atomic)
	fmt.Printf("  atomic shared  %8.2f  (%.1fx atomic)\n", c.AtomicShared, c.AtomicShared/c.Atomic)
	fmt.Printf("  mutex shared   %8.2f  (%.1fx atomic shared)\n", c.LockedShared, c.LockedShared/c.AtomicShared)
}
