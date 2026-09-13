// Package pull walks a binary search tree and merges two sorted sequences.
//
// The merge is the reason this task exists: it needs one element from each
// sequence in hand at the same time, and two push sequences cannot both be in
// control. Inverting one of them is what iter.Pull is for — and what it hands
// back alongside the next function is the part people forget.
package pull

import "iter"

// Node is a binary search tree node. A nil *Node is an empty tree.
type Node struct {
	Value       int
	Left, Right *Node
}

// InOrder yields the tree's values in sorted order.
//
// It must stop walking as soon as its consumer stops — including partway down
// the left spine, where the obvious recursive version has no way to say so.
func (n *Node) InOrder() iter.Seq[int] {
	panic("undergo: implement InOrder")
}

// Merge yields the elements of two sorted sequences in sorted order.
//
// The contract the frozen tests enforce:
//
//   - the output is sorted, and every element of both inputs appears;
//   - an empty side is handled — the other side comes through whole;
//   - when the consumer stops early, BOTH sources are stopped;
//   - no goroutine outlives the merge.
//
// That last pair is one requirement wearing two hats. Whatever you use to pull
// from a sequence has to be released however this function exits.
func Merge(a, b iter.Seq[int]) iter.Seq[int] {
	panic("undergo: implement Merge")
}
