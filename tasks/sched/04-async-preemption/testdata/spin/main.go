// Command spin runs a tight loop with no function calls, no allocation and no
// channel operations, while another goroutine asks for a garbage collection.
//
// It prints "true" if the collection finished while the loop was still running,
// and "false" if it did not. It is run twice: once normally, once with
// GODEBUG=asyncpreemptoff=1.
//
// There is no timeout and no select here, deliberately. Under
// GODEBUG=asyncpreemptoff=1 with one P, the spinning goroutine holds the only
// processor and *nothing else runs at all* — not the collector, and not a
// goroutine waiting on a timer. A select over a timer would therefore not
// measure the collector; it would be evaluated only after the loop ended, with
// both cases already ready, and select chooses among ready cases at random.
//
// Sampling an atomic from inside the loop asks the question directly and has no
// race to lose. An atomic load compiles to a plain move with no call, so the
// loop still contains no safe point.
package main

import (
	"fmt"
	"runtime"
	"sync/atomic"
)

// spinIterations is large enough that a collection has ample time to finish
// while the loop runs, and small enough that the program ends promptly even
// when the loop cannot be preempted.
const spinIterations = 1 << 30

var sink atomic.Uint64

func main() {
	runtime.GOMAXPROCS(1)

	var collected atomic.Bool
	go func() {
		runtime.GC()
		collected.Store(true)
	}()

	var x uint64
	var duringSpin bool
	for i := uint64(0); i < spinIterations; i++ {
		x += i
		if i == spinIterations-1 {
			duringSpin = collected.Load()
		}
	}
	sink.Store(x)

	fmt.Println(duringSpin)
}
