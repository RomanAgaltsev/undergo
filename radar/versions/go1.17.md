# Go 1.17

Source: https://go.dev/doc/go1.17

## The register calling convention, and the adapters it needed
→ candidate | asm
"Benchmarks for a representative set of Go packages and programs show performance
improvements of about 5%" and "a typical reduction in binary size of about 2%",
on `linux/amd64`, `darwin/amd64` and `windows/amd64`.

`asm/03-register-abi` grades where the arguments go. The half that release note
adds is what happens at the boundary: "the compiler generates adapter functions
that convert between the register-based calling convention and the stack-based
one", so "taking the address of a Go function in assembly or via
`reflect.ValueOf(fn).Pointer()` or `unsafe.Pointer` returns the address of the
adapter", and there is overhead "calling an assembly function indirectly via a
`func` value" and "calling Go functions from assembly".

That is the same class of surprise as the open 1.27 pool row about closures
sharing a code pointer: a function's address is not the thing you think it is.

## Functions containing closures became inlinable, and their code pointers multiplied
→ candidate | edges
"Functions containing closures can now be inlined ... a function with a closure
may produce a distinct closure code pointer for each inline location ... this may
reveal bugs in code using `reflect` or `unsafe.Pointer` to compare functions by
code pointer."

This is the origin of the open 1.27 pool row. One closure literal in the source
becomes several code pointers in the binary, one per inline site — so comparing
`reflect.ValueOf(f).Pointer()` for two values derived from the same source
closure can report either equal or unequal depending on inlining. Count the
distinct pointers and explain the count.

## A type conversion that panics
→ candidate | types
Slice to array pointer: "The conversion panics if `len(s)` is less than `N`",
and the notes flag what is novel about it — this is the **first case where a type
conversion can panic at run time**.

Conversions had always been either compile-time errors or total. Predict which of
a set of conversions panic, then say what that does to the mental model of
`T(x)`: the syntax no longer tells you whether control can leave.

## reflect.ConvertibleTo stopped being a guarantee
→ candidate | reflect
The direct consequence of the above: "`Type.ConvertibleTo` is no longer a
sufficient guarantee that `Value.Convert` will not panic", and the new
`Value.CanConvert` method exists to answer the question properly.

This is the same shape as `iface/05-interface-comparison`'s
`reflect.Value.Comparable`: a static predicate that was total became partial, and
the library grew a value-aware check beside it. Predict which of a set of
`ConvertibleTo`/`CanConvert` pairs disagree.

## unsafe.Add and unsafe.Slice
→ candidate | edges
`unsafe.Add(ptr, len)` "adds `len` to `ptr`" and `unsafe.Slice(ptr, len)`
"returns a slice whose underlying array starts at `ptr` and whose length and
capacity are `len`". `edges/05-unsafe-pointer-rules` uses both.

The sentence worth grading is the disclaimer: "The rules remain unchanged. In
particular, existing programs that correctly use `unsafe.Pointer` remain valid,
and new programs must still follow the rules when using `unsafe.Add` or
`unsafe.Slice`." A friendlier spelling of pointer arithmetic is not a safer one,
and `edges/05` already measures that `go vet` cannot tell the difference.

## Stack traces print arguments, and admit when they cannot
→ candidate | asm
Arguments are "printed separately, separated by commas", aggregates "delimited by
curly braces", return values are no longer printed "as they were usually
inaccurate" — and the caveat that 1.18 later turned into a `?` marker: "the value
of an argument that only lives in a register and is not stored to memory may be
inaccurate."

A traceback stopped being able to tell the truth about arguments in the same
release the arguments moved into registers. Predict which arguments a traceback
can still report accurately, and why.

## Goroutine scheduling latency became a distribution
→ candidate | sched
`runtime/metrics` gained "a distribution of goroutine scheduling latencies" along
with counters for total bytes and objects allocated and freed.

Scheduling latency is the quantity `sched/01-runnext-ordering` reasons about
qualitatively — how long a runnable goroutine waits for a P — and this is where it
became measurable. A histogram also makes the `GOMAXPROCS=1` case gradeable
without timing anything by hand.

## Block profiles stopped favouring rare long events
→ candidate | sched
"Block profiles are no longer biased toward infrequent long events over frequent
short events."

The same correction the mutex profile got in 1.22, applied to blocking five
releases earlier. Both are sampling-bias fixes, and the pair makes a good
question: what does a profiler have to do to report *total* blocked time rather
than *sampled* blocked time, and what was the old number actually measuring?

## os.File.WriteString stopped copying
→ candidate | alloc
"`File.WriteString` has been optimized to not make a copy of the input string."
Previously the string was converted to `[]byte` to reach the underlying write,
which allocated for every call. Measurable with `AllocsPerRun`, and it is the
same conversion question `types/05-compiler-free-conversions` asks — here answered
inside the standard library rather than by the compiler.

## atomic.Value gained Swap and CompareAndSwap
→ candidate | memmodel
`memmodel/05-publication` grades five lazy-initialisation strategies and names
`atomic.Value` as one of the correct ones. Compare-and-swap turns it from
publish-once into a primitive that can also *replace* a published value safely,
which is a different problem — and the interesting question is what happens to
readers holding the old value.

## arm64 keeps frame pointers everywhere
→ candidate | asm
"Go programs now maintain stack frame pointers on 64-bit ARM on all operating
systems. Previously, this was only done on Linux, macOS, and iOS." Pairs with
1.21's change to `NOFRAME` on amd64: the frame pointer costs a register and an
instruction or two per call, and buys profilers and debuggers the ability to
unwind without DWARF.

## strconv adopted Ryū
→ no action
"more than 99% faster on worst-case inputs." A genuinely dramatic number and a
famous algorithm, but the observable — float formatting produces the same string
faster — offers nothing to predict.

## //go:build replaced // +build
→ no action
The `go` command "prefers them over `// +build` lines", `gofmt` synchronises the
two, and `vet` warns when they disagree. Every file in this repository uses the
new form, but it is build syntax.

## Pruned module graphs
→ no action
A significant change to how the `go` command loads dependencies, with real
consequences for build times and `go.mod` contents, but nothing underneath Go
that a task could measure.

## URL queries stopped accepting semicolons
→ no action
"`example?a=1;b=2&c=3`" now parses as "`map[c:[3]]`". A visible behaviour change
with a security rationale, and a good cautionary tale about parser
disagreements, but it is library semantics rather than machine behaviour.
