// Package adapters builds five sequence adapters on top of iter.Seq.
//
// Every one of them is four or five lines. The difficulty is not the
// transformation — it is the yield contract: an adapter is both a consumer of
// the sequence it wraps and a producer for the one consuming it, and it has to
// honour the contract in both directions.
package adapters

import "iter"

// Map returns a sequence of f applied to each element of seq.
//
// It must stop pulling from seq as soon as its own consumer stops.
func Map[T, R any](seq iter.Seq[T], f func(T) R) iter.Seq[R] {
	panic("undergo: implement Map")
}

// Filter returns the elements of seq for which keep reports true.
func Filter[T any](seq iter.Seq[T], keep func(T) bool) iter.Seq[T] {
	panic("undergo: implement Filter")
}

// Take returns at most the first n elements of seq.
//
// It must stop the source once it has n — not filter afterwards. A Take over an
// infinite sequence has to terminate.
//
// n <= 0 yields nothing and must not touch the source at all.
func Take[T any](seq iter.Seq[T], n int) iter.Seq[T] {
	panic("undergo: implement Take")
}

// Zip pairs elements of a and b positionally, stopping at the shorter of the
// two. The longer sequence must not be drained past the pair that ended it.
func Zip[A, B any](a iter.Seq[A], b iter.Seq[B]) iter.Seq2[A, B] {
	panic("undergo: implement Zip")
}

// Chunk groups seq into slices of at most size elements. The final chunk may be
// short. A size <= 0 yields nothing.
//
// Each chunk handed to the consumer is theirs to keep: chunks must not share a
// backing array with one another.
func Chunk[T any](seq iter.Seq[T], size int) iter.Seq[[]T] {
	panic("undergo: implement Chunk")
}
