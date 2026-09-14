// Package defers holds six functions that defer in different shapes.
//
// You do not implement anything here. Work out which fall back to the runtime,
// write your answers into prediction.yaml, then verify.
package defers

import "os"

// Sink keeps results reachable.
var Sink int

func work() { Sink = os.Getpid() }

// One defers once.
func One() { defer work() }

// Several defers three times in a straight line.
func Several() {
	defer work()
	defer work()
	defer work()
}

// InLoop defers inside a loop.
func InLoop(n int) {
	for range n {
		defer work()
	}
}

// InBranch defers inside an if.
func InBranch(b bool) {
	if b {
		defer work()
	}
}

// Many defers more times than the open-coding limit allows.
func Many() {
	defer work()
	defer work()
	defer work()
	defer work()
	defer work()
	defer work()
	defer work()
	defer work()
	defer work()
}

// WithRecover defers a function that recovers.
func WithRecover() {
	defer func() { _ = recover() }()
	work()
}
