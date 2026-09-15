// Command hoist spins until another goroutine sets a plain bool.
//
// The program has a data race, so the compiler is entitled to assume the
// variable cannot change underneath the loop, load it once, and reuse the
// value. It prints "terminated" if the loop exits.
//
// Its caller kills it if it does not.
package main

import (
	"fmt"
	"runtime"
	"time"
)

var done bool
var spins uint64

func main() {
	runtime.GOMAXPROCS(2)

	go func() {
		time.Sleep(50 * time.Millisecond)
		done = true
	}()

	for !done {
		spins++
	}
	fmt.Println("terminated", spins > 0)
}
