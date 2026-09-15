package directsend

import (
	"testing"
	"testing/synctest"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// Every case follows the same three beats: start the goroutine, use
// synctest.Wait to make "it is parked" a fact rather than a hope, then act and
// measure. Outside a bubble the middle beat would be a sleep.
func TestPredictions(t *testing.T) {
	got := map[string]any{}

	// A receiver is already waiting when the send happens.
	synctest.Test(t, func(t *testing.T) {
		ch := make(chan int, 4)
		done := make(chan struct{})
		go func() { <-ch; close(done) }()
		synctest.Wait()
		ch <- 1
		got["len_after_send_with_parked_receiver"] = len(ch)
		<-done
	})

	// The same send with nobody waiting.
	synctest.Test(t, func(t *testing.T) {
		ch := make(chan int, 4)
		ch <- 1
		got["len_after_send_with_no_receiver"] = len(ch)
	})

	// Can a channel hold a full buffer and a parked receiver at once? The
	// receiver is started against a full buffer; if it parked, len is still 2.
	synctest.Test(t, func(t *testing.T) {
		ch := make(chan int, 2)
		ch <- 1
		ch <- 2
		go func() { <-ch }()
		synctest.Wait()
		got["full_buffer_and_parked_receiver_possible"] = len(ch) == 2
		<-ch
	})

	// A full buffer and a parked sender. Which value does a receive return?
	synctest.Test(t, func(t *testing.T) {
		ch := make(chan int, 2)
		ch <- 10
		ch <- 20
		go func() { ch <- 30 }()
		synctest.Wait()
		got["value_received_when_buffer_full_and_sender_parked"] = <-ch
		<-ch
		<-ch
	})

	predict.Check(t, got)
}
