// Package tree is a binary search tree you can walk, stop walking, and resume.
//
// Two things to build. All is a push sequence — it hands values to a loop. Walk
// is a pull sequence — the caller asks for the next value when it wants one.
package tree

import "iter"

// Tree is an unbalanced binary search tree.
type Tree struct {
	Value       int
	Left, Right *Tree
}

// Insert adds v to the tree and returns the (possibly new) root.
func Insert(t *Tree, v int) *Tree {
	if t == nil {
		return &Tree{Value: v}
	}
	if v < t.Value {
		t.Left = Insert(t.Left, v)
	} else if v > t.Value {
		t.Right = Insert(t.Right, v)
	}
	return t
}

// All returns an in-order sequence over the tree.
//
// It must stop descending as soon as the loop body leaves. A recursive walk
// that ignores what yield returned will keep going through the rest of the
// tree, which the tests check for.
func All(t *Tree) iter.Seq[int] {
	panic("implement All")
}

// Walk turns All into a pull sequence: next returns the following value and
// whether there was one, and stop releases whatever next is holding.
//
// After stop, next must report false. Calling stop more than once must be
// safe.
func Walk(t *Tree) (next func() (int, bool), stop func()) {
	panic("implement Walk")
}
