# Go 1.23

Source: https://go.dev/doc/go1.23

## structs.HostLayout, and what the spec does not guarantee about layout
→ invalidates | layout/01-struct-padding
The new `structs` package exists because "without this marker, struct layout
order is not guaranteed by the language spec". The sealed **hint** for
`layout/01-struct-padding` opened with "Go lays a struct out in declaration
order" — true of the gc compiler, not promised by the language. In a repo whose
whole point is citing the mechanism rather than a rule of thumb, that sentence
had to attribute the ordering to the implementation.

Fixed on 2026-09-11 and re-sealed: the hint now says *the gc compiler* and points
at the distinction, and the explanation gained a closing section on what the spec
guarantees versus what merely happens to be true, citing `structs.HostLayout`.
The five measured answers were never affected.

## Range-over-function iterators and the iter package
→ candidate | iter
The range clause now accepts `func(func() bool)`, `func(func(K) bool)` and
`func(func(K, V) bool)`, and `iter` "provides the basic definitions for working
with user-defined iterators". This is the floor for the whole `iter` track:
every task in it needs `requires.go` of at least 1.23. The yield contract — what
`break`, `return` and `goto` do to the iterator function — is the task.

## Timers and tickers are collected without Stop, and their channels are unbuffered
→ candidate | sched
Two guarantees: they "become eligible for garbage collection immediately, even
if their Stop methods have not been called", and the channel "is now unbuffered,
with capacity 0", so "for any call to a Reset or Stop method, no stale values
prepared before that call will be sent or received after the call". The sharpest
detail is the gate: "these new behaviors are only enabled when the main Go
program is in a module with a `go.mod` `go` line using Go 1.23.0 or later" —
behaviour that depends on the go directive, not the toolchain, which is a trap
worth a task of its own. `asynctimerchan=1` reverted it until 1.27 removed the
setting permanently.

## The unique package
→ candidate | weak
"Any value of comparable type may be canonicalized with the new `Make[T]`
function", and comparing two handles reduces "down to a simple pointer
comparison". Note the floor for the `weak` track is split: `unique` is 1.23,
while `weak` and `runtime.AddCleanup` are 1.24, and the single-cycle reclamation
of handles is 1.25. A task that measures interning memory must pin the right one.

## Windows timer resolution improved to 0.5ms from 15.6ms
→ candidate | sched
A timer-skew or sleep-precision answer measured on Windows before 1.23 is off by
a factor of thirty. Any such task must declare `requires.os` and say which
resolution it assumes — this is the clearest small example in the sweep of why
`requires` exists.

## PGO overhead collapsed, and hot blocks are aligned on 386 and amd64
→ candidate | compiler
PGO build overhead went from "100%+ build time increase" to "single digit
percentages", and PGO now aligns hot loop blocks for "an additional 1-1.5% at a
cost of an additional 0.1% text and binary size", disabled with
`-d=alignhot=0`. The planned `compiler` task on measuring a PGO win gets both a
cheaper build and a second, separable effect to attribute.

## The compiler overlaps stack slots of disjoint locals
→ candidate | compiler
"The compiler ... can now overlap the stack frame slots of local variables
accessed in disjoint regions of a function, which reduces stack usage." Predict
the stack frame size of a function whose locals have disjoint live ranges, and
explain why moving one line changes it.

## linkname to unmarked standard-library symbols is now refused
→ candidate | edges
The linker "disallows using a `//go:linkname` directive to refer to internal
symbols in the standard library ... that are not marked with `//go:linkname` on
their definitions", with `-checklinkname=0` to disable. Fits the `edges` task on
what is and is not a legal way to reach into the runtime.

## Profile stack depth raised from 32 to 128 frames
→ no action
It changes what a profile shows, but a solver could only be wrong about a number
with no mechanism behind it.

## Generic type aliases behind GOEXPERIMENT=aliastypeparams
→ no action
Preview only, and not usable across package boundaries here.

## Panic traceback indentation, SetCrashOutput, trace flush on crash
→ no action
Output formatting and crash plumbing.

## pidfd on Linux, macOS 11 minimum, kernel 2.6.32 deprecation
→ no action
Platform support.
