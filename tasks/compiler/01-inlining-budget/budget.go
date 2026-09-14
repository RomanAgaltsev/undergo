// Package budget holds six functions spanning the inliner's cost budget.
//
// You do not implement anything here. Work out the answers, write them into
// prediction.yaml, then verify.
package budget

import "os"

// Trivial is as small as a function gets.
func Trivial(a, b int) int { return a + b }

// Loop does a little arithmetic in a loop.
func Loop(n int) int {
	t := 0
	for i := range n {
		t += i * i
		t -= i
		t += i % 3
	}
	return t
}

// CallsInlinable calls Trivial three times. Whether the inner cost counts
// toward this one's budget is the question.
func CallsInlinable(a, b, c int) int {
	return Trivial(Trivial(a, b), Trivial(b, c))
}

// Sprawling is deliberately long.
func Sprawling(n int) int {
	t := 0
	for i := range n {
		t += i * i
		t -= i
		t += i % 3
		t *= 2
		t /= 3
		t += i * i
		t -= i
		t += i % 5
		t *= 2
		t /= 3
		t += i * i
		t -= i
		t += i % 7
		t *= 2
		t /= 3
		t += i * i
		t -= i
		t += i % 11
	}
	return t
}

// WithDefer is small, and defers.
func WithDefer() int {
	defer func() { _ = os.Getpid() }()
	return 1
}

// WithGo is small, and starts a goroutine.
func WithGo(done chan struct{}) int {
	go func() { close(done) }()
	return 1
}
