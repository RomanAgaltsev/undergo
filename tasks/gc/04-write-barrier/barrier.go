// Package barrier holds five one-line stores. Some of them compile to a write
// barrier and some do not.
//
// You do not implement anything here. Work out which, write your answers into
// prediction.yaml, then verify.
package barrier

// Node is a heap object with one pointer field and one scalar field.
type Node struct {
	Next *Node
	N    int
}

// Sink keeps values reachable so nothing here is optimised away.
var Sink *Node

// PointerStore writes a pointer into a heap object.
func PointerStore(n *Node, next *Node) { n.Next = next }

// ScalarStore writes an int into a heap object.
func ScalarStore(n *Node, v int) { n.N = v }

// NilStore writes nil into a heap object's pointer field.
func NilStore(n *Node) { n.Next = nil }

// LocalStore writes a pointer into a Node that never leaves this function.
func LocalStore(next *Node) int {
	var n Node
	n.Next = next
	return n.N
}

// SliceStore writes a pointer into a slice element.
func SliceStore(ns []*Node, i int, v *Node) { ns[i] = v }
