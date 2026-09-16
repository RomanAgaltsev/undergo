// Package rwadmission watches who acquires an RWMutex, and in what order, when
// a writer is waiting behind a reader.
package rwadmission

import (
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"time"
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
// the deadline passes.
//
// A synctest bubble cannot be used here: a goroutine blocked on a mutex is not
// "durably blocked" in the bubble's sense, so synctest.Wait would never return
// and the test would hang rather than fail. The profile names functions by full
// import path, hence the "."+fn match.
func waitParked(fn string) bool {
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		var sb strings.Builder
		if err := pprof.Lookup("goroutine").WriteTo(&sb, 1); err != nil {
			return false
		}
		if strings.Contains(sb.String(), "."+fn) {
			return true
		}
		runtime.Gosched()
		time.Sleep(time.Millisecond)
	}
	return false
}

// Admission records what happened when a second reader arrived at a read-locked
// RWMutex with a writer already waiting.
type Admission struct {
	SetUp                bool
	SecondReaderParked   bool
	SecondReaderAdmitted bool
	Order                string
	WriterWaitedForFirst bool
}

// Measure runs the experiment: R1 holds a read lock, W blocks on Lock, then R2
// asks for a read lock.
func Measure() Admission {
	var a Admission
	var mu sync.RWMutex
	r := &recorder{}

	mu.RLock()
	r.note("R1")

	var wg sync.WaitGroup
	wg.Go(func() { pendingWriter(&mu, r) })
	a.SetUp = waitParked("pendingWriter")

	wg.Go(func() { secondReader(&mu, r) })
	a.SecondReaderParked = waitParked("secondReader")

	before := r.snapshot()
	a.SecondReaderAdmitted = len(before) > 1 && before[1] == "R2"

	mu.RUnlock()
	wg.Wait()

	final := r.snapshot()
	a.Order = strings.Join(final, ",")
	a.WriterWaitedForFirst = final[0] == "R1"
	return a
}
