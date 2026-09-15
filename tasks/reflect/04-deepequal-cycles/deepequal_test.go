package deepequal

import (
	"reflect"
	"testing"
	"time"
)

type node struct {
	Name string
	Next *node
}

type holder struct {
	Items []int
	Meta  map[string]int
	Child *holder
}

// cases covers everything except the cyclic values, which need their own
// construction and are tested separately.
func cases() []struct {
	name string
	a, b any
} {
	return []struct {
		name string
		a, b any
	}{
		{"equal ints", 1, 1},
		{"unequal ints", 1, 2},
		{"different types", 1, int64(1)},
		{"equal strings", "x", "x"},
		{"nil and nil", nil, nil},
		{"nil and value", nil, 1},
		{"equal slices, different arrays", []int{1, 2}, []int{1, 2}},
		{"unequal slices", []int{1, 2}, []int{1, 3}},
		{"nil slice and empty slice", []int(nil), []int{}},
		{"maps, insertion order differs", map[string]int{"a": 1, "b": 2}, map[string]int{"b": 2, "a": 1}},
		{"unequal maps", map[string]int{"a": 1}, map[string]int{"a": 2}},
		{"nil map and empty map", map[string]int(nil), map[string]int{}},
		{"equal structs", holder{Items: []int{1}}, holder{Items: []int{1}}},
		{"unequal structs", holder{Items: []int{1}}, holder{Items: []int{2}}},
		{"pointers to equal values", &node{Name: "a"}, &node{Name: "a"}},
		{"pointers to unequal values", &node{Name: "a"}, &node{Name: "b"}},
		{"nil pointer and pointer", (*node)(nil), &node{}},
		{"nested structs", holder{Child: &holder{Items: []int{1}}}, holder{Child: &holder{Items: []int{1}}}},
	}
}

// TestAgainstDeepEqual checks every acyclic case against the standard library,
// so the semantics are pinned rather than invented.
func TestAgainstDeepEqual(t *testing.T) {
	for _, c := range cases() {
		// Compared inline rather than through a variable. The stub's Equal
		// panics unconditionally, so staticcheck can prove a following
		// comparison unreachable and reports the variable as never read.
		if Equal(c.a, c.b) != reflect.DeepEqual(c.a, c.b) {
			t.Errorf("%s: Equal disagrees with reflect.DeepEqual", c.name)
		}
	}
}

// selfReferential builds a one-element ring: n.Next points back at n.
func selfReferential(name string) *node {
	n := &node{Name: name}
	n.Next = n
	return n
}

// TestTerminatesOnSelfReference is the whole point of the task. It runs Equal
// in a goroutine so a non-terminating implementation fails rather than hanging
// the suite until the timeout.
func TestTerminatesOnSelfReference(t *testing.T) {
	done := make(chan bool, 1)
	go func() { done <- Equal(selfReferential("a"), selfReferential("a")) }()

	select {
	case got := <-done:
		if !got {
			t.Error("two structurally identical rings should be equal")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Equal did not terminate on a self-referential value")
	}
}

// TestCyclicAgainstAcyclic checks that a cycle is not mistaken for equality
// with a value that merely looks the same for one step.
func TestCyclicAgainstAcyclic(t *testing.T) {
	cyclic := selfReferential("a")
	acyclic := &node{Name: "a", Next: &node{Name: "a"}}

	done := make(chan bool, 1)
	go func() { done <- Equal(cyclic, acyclic) }()

	select {
	case got := <-done:
		if got {
			t.Error("a ring and a two-node chain are not structurally equal")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Equal did not terminate when one side was cyclic")
	}
}

// TestLongerCycle uses a two-element ring, which a visited set keyed on a
// single pointer rather than on the pair will get wrong.
func TestLongerCycle(t *testing.T) {
	ring := func() *node {
		a := &node{Name: "a"}
		b := &node{Name: "b"}
		a.Next, b.Next = b, a
		return a
	}

	done := make(chan bool, 1)
	go func() { done <- Equal(ring(), ring()) }()

	select {
	case got := <-done:
		if !got {
			t.Error("two structurally identical two-element rings should be equal")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Equal did not terminate on a two-element cycle")
	}
}
