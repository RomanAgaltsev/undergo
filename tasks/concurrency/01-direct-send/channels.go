// Package directsend holds four small channel experiments.
//
// Every function here must be called from inside a synctest bubble: each one
// uses synctest.Wait to turn "the other goroutine is parked" from a hope into a
// fact before it measures anything.
package directsend

import "testing/synctest"

// BufferAfterSendToParkedReceiver parks a receiver on a channel of capacity 4,
// sends one value, and reports len(ch) immediately afterwards.
func BufferAfterSendToParkedReceiver() int {
	ch := make(chan int, 4)
	done := make(chan struct{})
	go func() { <-ch; close(done) }()
	synctest.Wait()
	ch <- 1
	n := len(ch)
	<-done
	return n
}

// BufferAfterSendWithNoReceiver does the same send with nobody waiting.
func BufferAfterSendWithNoReceiver() int {
	ch := make(chan int, 4)
	ch <- 1
	return len(ch)
}

// FullBufferKeepsParkedReceiver fills a channel of capacity 2, starts a
// receiver, and reports whether the buffer is still full once everything has
// settled — that is, whether the receiver is parked rather than served.
func FullBufferKeepsParkedReceiver() bool {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	go func() { <-ch }()
	synctest.Wait()
	full := len(ch) == 2
	<-ch
	return full
}

// ReceiveWithFullBufferAndParkedSender fills a channel of capacity 2 with 10 and
// 20, parks a sender holding 30, and reports which value a receive returns.
func ReceiveWithFullBufferAndParkedSender() int {
	ch := make(chan int, 2)
	ch <- 10
	ch <- 20
	go func() { ch <- 30 }()
	synctest.Wait()
	v := <-ch
	<-ch
	<-ch
	return v
}
