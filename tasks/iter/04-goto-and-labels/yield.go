// Package yield watches what each way out of a range-over-function body does
// to the sequence that is feeding it.
//
// Every function here ranges over the same sequence, which yields 1 to 5 and
// records what yield returned on each call. They differ only in how the loop
// body leaves.
package yield

import "iter"

// Trace is what one sequence saw: how many times it called yield, and what
// each call returned.
type Trace struct {
	Calls   int
	Returns []bool
}

// counted returns a sequence over 1..5 that records its own yield results into
// tr, and stops as soon as yield says to.
func counted(tr *Trace) iter.Seq[int] {
	return func(y func(int) bool) {
		for i := 1; i <= 5; i++ {
			ok := y(i)
			tr.Calls++
			tr.Returns = append(tr.Returns, ok)
			if !ok {
				return
			}
		}
	}
}

// LastReturn reports what yield returned on its final call, or true if it was
// never called.
func (t Trace) LastReturn() bool {
	if len(t.Returns) == 0 {
		return true
	}
	return t.Returns[len(t.Returns)-1]
}

// Complete runs the loop to the end without leaving early.
func Complete() Trace {
	var tr Trace
	sum := 0
	for v := range counted(&tr) {
		sum += v
	}
	_ = sum
	return tr
}

// Break leaves with a plain break on the third value.
func Break() Trace {
	var tr Trace
	n := 0
	for range counted(&tr) {
		n++
		if n == 3 {
			break
		}
	}
	return tr
}

// Goto leaves with a goto to a label outside the loop, on the second value.
func Goto() Trace {
	var tr Trace
	n := 0
	for range counted(&tr) {
		n++
		if n == 2 {
			goto done
		}
	}
done:
	return tr
}

// LabelledBreak leaves an inner loop and the range loop together, on the
// fourth inner iteration.
func LabelledBreak() Trace {
	var tr Trace
	n := 0
outer:
	for range counted(&tr) {
		for range 2 {
			n++
			if n == 4 {
				break outer
			}
		}
	}
	return tr
}

// Return leaves the enclosing function from inside the loop body.
func Return() Trace {
	var tr Trace
	run := func(tr *Trace) {
		for range counted(tr) {
			return
		}
	}
	run(&tr)
	return tr
}

// Panicked leaves by panicking, recovered outside the loop.
func Panicked() (tr Trace) {
	defer func() { _ = recover() }()
	n := 0
	for range counted(&tr) {
		n++
		if n == 2 {
			panic("out")
		}
	}
	return tr
}
