// Command park parks goroutines and reports what they cost. It runs the one
// phase named by its argument and prints one number.
package main

import (
	"fmt"
	"os"
	"runtime"
)

// depth is how far a deep goroutine recurses, through 1 KiB frames.
const depth = 32

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: park parked|deep|returned|unreachable")
		os.Exit(2)
	}
	switch os.Args[1] {
	case "parked":
		// Stack bytes per goroutine waiting on a channel.
		fmt.Println(perGoroutine(10_000, 0, func(ready chan<- struct{}, stop <-chan struct{}) {
			ready <- struct{}{}
			<-stop
		}))
	case "deep":
		// The same, for a goroutine waiting at the bottom of the recursion.
		fmt.Println(perGoroutine(1_000, 0, func(ready chan<- struct{}, stop <-chan struct{}) {
			descend(depth, func() {
				ready <- struct{}{}
				<-stop
			})
		}))
	case "returned":
		// A goroutine that recursed, came back up, and then waits — measured
		// after six more collections than the other phases.
		fmt.Println(perGoroutine(1_000, 6, func(ready chan<- struct{}, stop <-chan struct{}) {
			descend(depth, func() {})
			ready <- struct{}{}
			<-stop
		}))
	case "unreachable":
		// How many of 10 000 goroutines were reclaimed once nothing but they
		// themselves could reach the channel they wait on.
		fmt.Println(reclaimed(10_000))
	default:
		fmt.Fprintf(os.Stderr, "park: unknown phase %q\n", os.Args[1])
		os.Exit(2)
	}
}

// stackInuse collects twice and returns the bytes held in stack spans.
func stackInuse() uint64 {
	runtime.GC()
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.StackInuse
}

// descend recurses d times through 1 KiB frames and calls bottom at the end.
//
//go:noinline
func descend(d int, bottom func()) byte {
	var frame [1024]byte
	frame[d%len(frame)] = byte(d)
	if d == 0 {
		bottom()
		return frame[0]
	}
	return descend(d-1, bottom) + frame[1]
}

// perGoroutine starts n goroutines running body, waits until every one has
// signalled ready, collects extraGCs more times, and returns the growth in
// stack memory divided by n. It releases the goroutines before returning.
func perGoroutine(n, extraGCs int, body func(ready chan<- struct{}, stop <-chan struct{})) uint64 {
	before := stackInuse()
	ready := make(chan struct{})
	stop := make(chan struct{})
	for range n {
		go body(ready, stop)
	}
	for range n {
		<-ready
	}
	for range extraGCs {
		runtime.GC()
	}
	after := stackInuse()
	close(stop)
	if after < before {
		return 0
	}
	return (after - before) / uint64(n)
}

// reclaimed parks n goroutines on a channel that only they can reach once this
// function's inner frame has returned, collects twice, and returns how many of
// them the runtime got rid of.
func reclaimed(n int) int {
	ready := make(chan struct{})
	func() {
		lost := make(chan struct{})
		for range n {
			go func() {
				ready <- struct{}{}
				<-lost
			}()
		}
	}()
	for range n {
		<-ready
	}
	runtime.GC()
	runtime.GC()
	return n - (runtime.NumGoroutine() - 1)
}
