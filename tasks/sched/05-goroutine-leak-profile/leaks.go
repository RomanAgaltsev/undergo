// Package leaks plants one goroutine leak of each classic shape and asks the
// Go 1.27 goroutineleak profile which of them it can name.
//
// None of these functions returns anything. Each starts a goroutine that will
// never finish, and drops the only thing that could have released it.
package leaks

import (
	"bytes"
	"context"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
)

// BlockedSend sends on an unbuffered channel nobody receives from.
func BlockedSend() {
	ch := make(chan int)
	go func() { ch <- 1 }()
}

// NilRecv receives from a nil channel.
func NilRecv() {
	var ch chan int
	go func() { <-ch }()
}

// WaitGroupNeverDone waits on a counter nothing decrements.
func WaitGroupNeverDone() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { wg.Wait() }()
}

// MutexNeverUnlocked takes a lock the caller is still holding.
func MutexNeverUnlocked() {
	var mu sync.Mutex
	mu.Lock()
	go func() { mu.Lock() }()
}

// ContextNeverDone waits on a context that has no cancellation.
//
// context.Background().Done() returns a nil channel, so this is the nil-receive
// shape wearing a context's clothes — which is the point.
func ContextNeverDone() {
	ctx := context.Background()
	go func() { <-ctx.Done() }()
}

// PlantAll plants one of every shape above.
func PlantAll() {
	BlockedSend()
	NilRecv()
	WaitGroupNeverDone()
	MutexNeverUnlocked()
	ContextNeverDone()
}

// Report is what one profiling pass saw.
type Report struct {
	CountBefore int
	CountAfter  int
	Named       map[string]bool
}

// Survey plants every shape, lets the leaked goroutines reach their blocking
// points, then profiles.
//
// It reads Count both before and after WriteTo, because those turn out not to
// be the same question.
//
// GOMAXPROCS is pinned to 1 and the surveying goroutine yields repeatedly
// first. Go's own corpus of these patterns, in
// runtime/testdata/testgoroutineleakprofile, does exactly this and says why:
// the tests are not flaky if and only if the profiler runs after the leaky
// goroutines have reached their blocking points.
func Survey() Report {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(1))
	PlantAll()

	for range 10 {
		runtime.Gosched()
	}
	runtime.GC()

	p := pprof.Lookup("goroutineleak")
	before := p.Count()

	var buf bytes.Buffer
	_ = p.WriteTo(&buf, 1)
	after := p.Count()

	text := buf.String()
	named := map[string]bool{}
	for _, fn := range []string{
		"BlockedSend", "NilRecv", "WaitGroupNeverDone",
		"MutexNeverUnlocked", "ContextNeverDone",
	} {
		// A profile names a function by its FULL IMPORT PATH, not by package
		// name: the frame reads
		// ".../tasks/sched/05-goroutine-leak-profile.BlockedSend.func1", never
		// "leaks.BlockedSend". Matching on the package name finds nothing and
		// reports every shape as undetected, which looks like a real answer.
		named[fn] = strings.Contains(text, "."+fn+".func")
	}
	return Report{CountBefore: before, CountAfter: after, Named: named}
}
