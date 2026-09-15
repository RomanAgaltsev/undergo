// Package benign takes the defences of the "harmless" data race one at a time
// and checks each against what the machine actually does.
//
// Every race here is deliberate. That is the subject, and it is why this task
// is pinned to the default build.
package benign

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Pair is three words wide, so a write to it is not one instruction on any
// architecture this repository runs on.
//
// Three rather than two is deliberate and was learned the hard way. A 16-byte
// struct compiles to two MOVQs on amd64 — tearable — but to a single STP
// (store-pair) on arm64, which does not tear in practice. At 24 bytes arm64
// needs STP plus MOVD, so the write is genuinely two instructions everywhere.
type Pair struct{ A, B, C uint64 }

// TornUint64Seen has two goroutines write different patterns to one uint64
// while a third reads it, and reports whether the reader ever saw a value
// neither writer wrote.
func TornUint64Seen(d time.Duration) bool {
	const a, b = 0x1111111111111111, 0x2222222222222222
	var v uint64
	torn := false

	stop := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(3)
	go writeLoop(&wg, stop, func() { v = a })
	go writeLoop(&wg, stop, func() { v = b })
	go writeLoop(&wg, stop, func() {
		if got := v; got != 0 && got != a && got != b {
			torn = true
		}
	})

	time.Sleep(d)
	close(stop)
	wg.Wait()
	return torn
}

// TornStructSeen runs the same experiment on a two-word struct. Both writers
// publish a Pair whose fields agree, so a Pair with differing fields was
// assembled from two different writes.
func TornStructSeen(d time.Duration) bool {
	pa, pb := Pair{A: 1, B: 1, C: 1}, Pair{A: 2, B: 2, C: 2}
	var v Pair
	torn := false

	stop := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(3)
	go writeLoop(&wg, stop, func() { v = pa })
	go writeLoop(&wg, stop, func() { v = pb })
	go writeLoop(&wg, stop, func() {
		if got := v; got.A != got.B || got.B != got.C {
			torn = true
		}
	})

	time.Sleep(d)
	close(stop)
	wg.Wait()
	return torn
}

// writeLoop runs body until stop is closed.
func writeLoop(wg *sync.WaitGroup, stop <-chan struct{}, body func()) {
	defer wg.Done()
	for {
		select {
		case <-stop:
			return
		default:
			body()
		}
	}
}

// SpinLoopExits builds and runs testdata/hoist, which spins on a plain bool set
// by another goroutine, and reports whether the loop ever ended.
//
// The compiler is entitled to load that bool once and reuse it, because the
// program has a race and is therefore promised nothing. Whether it does so is a
// different question from whether it may.
//
// The build is deliberately outside the budget. Putting `go run` under a
// timeout measures how fast the machine compiles as much as how the program
// behaves, and a cold or slow runner then fails the task for a reason that has
// nothing to do with the question. Compiled separately, the budget covers only
// the run — which takes about sixty milliseconds against a budget of thirty
// seconds.
func SpinLoopExits(budget time.Duration) (bool, error) {
	dir, err := os.MkdirTemp("", "undergo-hoist-")
	if err != nil {
		return false, err
	}
	defer func() { _ = os.RemoveAll(dir) }()

	bin := filepath.Join(dir, "hoist")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	if out, err := exec.Command("go", "build", "-o", bin, "./testdata/hoist").CombinedOutput(); err != nil {
		return false, fmt.Errorf("building testdata/hoist: %w: %s", err, out)
	}

	cmd := exec.Command(bin)
	done := make(chan struct{})
	var out []byte
	var runErr error
	go func() {
		out, runErr = cmd.Output()
		close(done)
	}()

	select {
	case <-done:
		if runErr != nil {
			return false, runErr
		}
		return strings.HasPrefix(strings.TrimSpace(string(out)), "terminated"), nil
	case <-time.After(budget):
		_ = cmd.Process.Kill()
		<-done
		return false, nil
	}
}
