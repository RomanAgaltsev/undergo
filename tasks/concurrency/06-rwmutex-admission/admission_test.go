package rwadmission

import (
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

type recorder struct {
	mu    sync.Mutex
	order []string
}

func (r *recorder) note(s string) { r.mu.Lock(); r.order = append(r.order, s); r.mu.Unlock() }

func (r *recorder) snapshot() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]string(nil), r.order...)
}

func pendingWriter(mu *sync.RWMutex, r *recorder) { mu.Lock(); r.note("W"); mu.Unlock() }
func secondReader(mu *sync.RWMutex, r *recorder)  { mu.RLock(); r.note("R2"); mu.RUnlock() }

// waitParked polls the goroutine profile until some goroutine is inside fn, or
// the deadline passes. A synctest bubble cannot be used here: a goroutine
// blocked on a mutex is not "durably blocked" in the bubble's sense, so
// synctest.Wait would never return. The profile names functions by full import
// path, hence the "."+fn match.
func waitParked(t *testing.T, fn string) bool {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		var sb strings.Builder
		if err := pprof.Lookup("goroutine").WriteTo(&sb, 1); err != nil {
			t.Fatal(err)
		}
		if strings.Contains(sb.String(), "."+fn) {
			return true
		}
		runtime.Gosched()
		time.Sleep(time.Millisecond)
	}
	return false
}

func TestPredictions(t *testing.T) {
	got := map[string]any{}

	var mu sync.RWMutex
	r := &recorder{}

	// R1 takes the read lock and holds it for the whole experiment.
	mu.RLock()
	r.note("R1")

	var wg sync.WaitGroup
	wg.Go(func() { pendingWriter(&mu, r) })
	if !waitParked(t, "pendingWriter") {
		t.Fatal("the writer never blocked; the experiment did not set up")
	}

	// R2 now arrives at a lock that is held for reading, with a writer pending.
	wg.Go(func() { secondReader(&mu, r) })
	got["second_reader_parked"] = waitParked(t, "secondReader")

	before := r.snapshot()
	got["second_reader_admitted_while_writer_pending"] = len(before) > 1 && before[1] == "R2"

	mu.RUnlock()
	wg.Wait()

	final := r.snapshot()
	got["acquisition_order"] = strings.Join(final, ",")
	got["writer_waits_for_first_reader"] = final[0] == "R1"

	predict.Check(t, got)
}
