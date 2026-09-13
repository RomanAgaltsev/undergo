// Package contract counts how many times an iterator's yield is called when the
// consumer stops in five different ways.
//
// You do not implement anything here. Work out the answers, write them into
// prediction.yaml, then verify.
package contract

// Counted yields 1..10 and records how many times it called yield.
//
// It honours the contract: when yield reports false it stops. Every consumer
// below shares this one iterator, so the only thing that varies between the
// slots is how the consumer leaves the loop.
type Counted struct {
	Yields int
}

// Seq is the sequence itself.
func (c *Counted) Seq(yield func(int) bool) {
	for i := 1; i <= 10; i++ {
		c.Yields++
		if !yield(i) {
			return
		}
	}
}

// Plain ranges to the end.
func Plain() int {
	c := &Counted{}
	for v := range c.Seq {
		_ = v
	}
	return c.Yields
}

// Break leaves the loop with break, after the third element.
func Break() int {
	c := &Counted{}
	for v := range c.Seq {
		if v == 3 {
			break
		}
	}
	return c.Yields
}

// Return leaves the loop — and the function — with return, after the third.
func Return() int {
	c := &Counted{}
	count := func() {
		for v := range c.Seq {
			if v == 3 {
				return
			}
		}
	}
	count()
	return c.Yields
}

// Labelled runs the sequence twice and breaks out of the enclosing labelled
// loop from inside the range body, on the second pass.
//
// One Counted serves both passes, so its total spans them.
func Labelled() int {
	c := &Counted{}
outer:
	for pass := range 2 {
		for v := range c.Seq {
			if pass == 1 && v == 3 {
				break outer
			}
		}
	}
	return c.Yields
}

// Panicking panics inside the loop body at the fifth element, and recovers
// outside the loop.
func Panicking() (yields int) {
	c := &Counted{}
	defer func() {
		_ = recover()
		yields = c.Yields
	}()

	for v := range c.Seq {
		if v == 5 {
			panic("stop")
		}
	}
	return c.Yields
}

// Misbehaving is an iterator that ignores the contract: it keeps calling yield
// after yield has already returned false.
//
// The question is what happens to it — not to the consumer that wrote a
// perfectly ordinary break.
func Misbehaving(yield func(int) bool) {
	for i := 1; i <= 10; i++ {
		yield(i) // deliberately ignoring the result
	}
}

// MisbehavingPanics reports whether ranging over Misbehaving and breaking out of
// it panics.
func MisbehavingPanics() (panicked bool) {
	defer func() {
		if recover() != nil {
			panicked = true
		}
	}()

	for v := range Misbehaving {
		if v == 3 {
			break
		}
	}
	return false
}
