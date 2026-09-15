# Go 1.0 – 1.9

Source: https://go.dev/doc/devel/release

Era E4 is harvested rather than swept. Most of this era has been superseded, and
ten thin files would misrepresent the yield. What follows is the durable
surprises — findings whose behaviour is still observable today — each naming the
release that introduced it. A finding that was later changed belongs in the
*later* release's file, and is not repeated here.

## Goroutine stacks became contiguous, and then became small
→ candidate | alloc
Go 1.3: "Go 1.3 has changed the implementation of goroutine stacks away from the
old, 'segmented' model to a contiguous model. When a goroutine needs more stack
than is available, its stack is transferred to a larger single block of memory.
The overhead of this transfer operation amortizes well and eliminates the old
'hot spot' problem when a calculation repeatedly steps across a segment
boundary."

Go 1.4 then collected the dividend: "The use of contiguous stacks means that
stacks can start smaller without triggering performance issues, so the default
starting size for a goroutine's stack in 1.4 has been reduced from 8192 bytes to
2048 bytes."

Two of the numbers every Go programmer quotes — "a goroutine costs 2 KB" and "a
stack grows by copying" — were both set here, and both are still true. The hot
spot is the part worth a task: under segmented stacks a function call that
straddled a boundary allocated and freed a segment *on every iteration*, so the
cost of a loop depended on where its frame happened to land. Copying makes growth
expensive once instead of cheap repeatedly, which is the same amortisation
argument `types/01-cap-growth` makes about `append`.

## The runtime decided that pointer-typed means pointer, and nothing else does
→ candidate | edges
Go 1.3: "Starting with Go 1.3, the runtime assumes that values with pointer type
contain pointers and other values do not. This assumption is fundamental to the
precise behavior of both stack expansion and garbage collection. Programs that
use package `unsafe` to store integers in pointer-typed values are illegal and
will crash if the runtime detects the behavior. Programs that use package
`unsafe` to store pointers in integer-typed values are also illegal but more
difficult to diagnose during execution. Because the pointers are hidden from the
runtime, a stack expansion or garbage collection may reclaim the memory they
point at, creating dangling pointers."

Go 1.4 completed it: "This rewrite allows the garbage collector in 1.4 to be
fully precise, meaning that it is aware of the location of all active pointers in
the program."

`edges/05-unsafe-pointer-rules` grades six legal patterns, and this paragraph is
where all six come from. The asymmetry it states is the thing to understand: an
integer in a pointer slot *crashes* — the collector will try to follow it —
while a pointer in an integer slot is silent and fatal later, because nothing
looks at it at all. One is a bug you find, the other is a bug that finds you.
Note also that the previous entry's stack copying is the second reason: a moved
stack rewrites every pointer it can identify, and a hidden one keeps the old
address.

## The collector became concurrent
→ candidate | gc
Go 1.5: "The garbage collector is now concurrent and provides dramatically lower
pause times by running, when possible, in parallel with other goroutines." And:
"The 'stop the world' phase of the collector will almost always be under 10
milliseconds and usually much less."

Every fact the `gc` track teaches dates from this release. Before it, a
collection was a pause proportional to the heap; after it, marking runs alongside
the mutator and the pause is proportional to almost nothing. That is the change
that makes `gc/05-gctrace`'s output shape — several phases, two tiny
stop-the-world brackets around a long concurrent middle — and it is why the
cumulative-percentage figure in a `gctrace` line answers a throughput question
and not a latency one.

Ten milliseconds was the claim in 1.5. The next entry is what it became.

## Stop-the-world stack rescanning was eliminated
→ candidate | gc
Go 1.8: "Garbage collection pauses should be significantly shorter than they were
in Go 1.7, usually under 100 microseconds and often as low as 10 microseconds."

From 10 milliseconds to 10 microseconds in three releases, and the mechanism is
named in the linked proposal: the hybrid write barrier, which removed the need to
re-scan every goroutine stack with the world stopped. Before it, the pause grew
with the number of goroutines — a program with a hundred thousand of them paid
for all hundred thousand stacks at the end of every cycle.

This is the entry that makes `weak` and `gc` tasks possible at all: a barrier
that catches writes *during* marking is what lets an object's reachability be
decided without freezing the program, and it is why 1.12's more precise liveness
analysis could make a finalizer run a whole collection sooner.

## GOMAXPROCS started at the number of cores
→ candidate | sched
Go 1.5: "By default, Go programs run with `GOMAXPROCS` set to the number of cores
available; in prior releases it defaulted to 1."

Every task in the `sched` track pins `GOMAXPROCS` explicitly, and this is why
that is necessary rather than fussy: the default is a property of the machine.
The three-year period where it defaulted to 1 is also the reason so much old Go
advice about goroutines is quietly wrong — concurrency without parallelism has
different failure modes, and `memmodel/01-store-buffering` cannot even reproduce
its subject at `GOMAXPROCS=1`.

1.25's container-aware default — still an open pool row — is this same setting
learning about cgroup limits, which is the third answer to "how many cores are
available" in the language's life.

## Preemption arrived at function entry, twelve years before it arrived anywhere else
→ candidate | sched
Go 1.2: "In prior releases, a goroutine that was looping forever could starve out
other goroutines on the same thread, a serious problem when GOMAXPROCS provided
only one user thread. In Go 1.2, this is partially addressed: The scheduler is
invoked occasionally upon entry to a function. This means that any loop that
includes a (non-inlined) function call can be pre-empted, allowing other
goroutines to run on the same thread."

`sched/04-async-preemption` grades the 1.14 mechanism, and this is the one it
replaced. Read the two together and the interesting word is **non-inlined**: for
twelve years, whether a loop could be preempted depended on whether the compiler
had decided to inline the call inside it. An inlining budget change could hang a
program. That is the strongest argument in the language's history for why
`compiler/01-inlining-budget` is not a trivia track.

## Slicing learned to specify a capacity
→ candidate | types
Go 1.2: "Go 1.2 adds new syntax to allow a slicing operation to specify the
capacity as well as the length. A second colon introduces the capacity value,
which must be less than or equal to the capacity of the source slice or array,
adjusted for the origin." The example given is `array[2:4:7]`, of which the notes
say: "It is impossible to use this new slice value to access the last three
elements of the original array."

`types/02-aliasing` grades exactly this, and 1.10's capacity-clipping of
`bytes.Split` results is the standard library adopting it defensively. The
question worth asking is why the form is *needed*: a two-index slice inherits the
source's capacity, so handing one to a caller hands them write access to
everything past its end. The third index is the only way to give away a window
without giving away the room it looks into.

## The SSA back end
→ candidate | compiler
Go 1.7, for 64-bit x86: "The new back end, based on SSA, generates more compact,
more efficient code and provides a better platform for optimizations such as
bounds check elimination. The new back end reduces the CPU time required by our
benchmark programs by 5-35%." Go 1.8 extended it: "In Go 1.8, that back end has
been developed further and is now used for all architectures ... reduces the CPU
time required by our benchmark programs by 20-30% on 32-bit ARM systems."

Half the `compiler` track measures things that only exist because of this.
Bounds-check elimination is named in the note itself, and every later entry the
radar found about it — transitive inequalities in 1.11, slice-creation facts in
1.14 — is an analysis written against SSA. The `-d=ssa/check_bce` flag those
tasks use is a debug hook into this back end, and `GOSSAFUNC` dumps its passes.

## math/bits, and functions the compiler recognises as instructions
→ candidate | asm
Go 1.9: "Go 1.9 includes a new package, `math/bits`, with optimized
implementations for manipulating bits. On most architectures, functions in this
package are additionally recognized by the compiler and treated as intrinsics for
additional performance."

`asm/05-carry-chain` rests on `bits.Add64`, and this is its origin. The
gradeable idea is the word *intrinsic*: the package ships a portable Go
implementation that is real, compilable, and almost never executed, because the
compiler substitutes a single instruction for the call. So the source you can
read is not the code that runs, and `go tool compile -S` is the only way to find
out which of the package's functions are intrinsics on your architecture — the
answer differs by GOARCH, and the portable fallback silently costs ten times
more where it does not.

1.13 later added a timing guarantee on top of the same functions.

## A time.Time carries two clocks
→ candidate | edges
Go 1.9: "The `time` package now transparently tracks monotonic time in each
`Time` value, making computing durations between two `Time` values a safe
operation in the presence of wall clock adjustments." And: "If a `Time` value has
a monotonic clock reading, its string representation (as returned by `String`)
now includes a final field `m=±value`."

One struct, two independent readings, and which one is used depends on the
operation: `Sub` and `Before` prefer the monotonic reading when both operands
have one, `Format` never uses it, and operations like `Round`, `Truncate`, `UTC`
and `AddDate` **strip** it. That makes `t.Round(0)` the documented way to drop it
and makes `==` genuinely treacherous, because it compares the monotonic reading
too — two `Time` values naming the same instant can compare unequal, and the same
value round-tripped through `UTC()` no longer equals itself.

`sched/02` uses a fake clock and `sched/03-timer-channels` grades timer
behaviour; this is the other half of the subject, and it is a comparison question
with an exact answer rather than a duration.

## sync.Map, and a map you can read without a lock
→ candidate | memmodel
Go 1.9: "The new `Map` type in the `sync` package is a concurrent map with
amortized-constant-time loads, stores, and deletes. It is safe for multiple
goroutines to call a `Map`'s methods concurrently."

`memmodel/05-publication` grades publication strategies, and this is the standard
library's most elaborate one: a read-only snapshot published through an atomic,
plus a dirty map behind a mutex, plus a miss counter that decides when to promote
one to the other. A `Load` that hits the snapshot takes no lock at all — the
happens-before edge comes from the atomic load of the snapshot pointer, exactly
the pattern that task measures.

The amortisation is the honest part of the claim, and it is measurable: the
promotion copies the whole dirty map, so a workload of misses pays O(n)
occasionally rather than O(1) always. 1.24 rewrote the internals wholesale
(HashTrieMap) without changing any of this, which is itself the point.

## Go's pointers stopped being C's pointers
→ candidate | edges
Go 1.6: "Go and C may share memory allocated by Go when a pointer to that memory
is passed to C as part of a `cgo` call, provided that the memory itself contains
no pointers to Go-allocated memory, and provided that C does not retain the
pointer after the call returns. These rules are checked by the runtime during
program execution: if the runtime detects a violation, it prints a diagnosis and
crashes the program."

A second set of pointer rules, enforced dynamically like `checkptr` and for the
same reason: a collector that moves stacks and tracks reachability cannot see
what C is holding. The note's closing advice is unusually blunt — the checks can
be turned off, "but note that the vast majority of code identified by the checks
is subtly incompatible with garbage collection in one way or another. Disabling
the checks will typically only lead to more mysterious failure modes."

Worth an entry beside `edges/05-unsafe-pointer-rules` because the shape is
identical: a rule the type system cannot express, a checker that finds most
violations at run time, and a flag people reach for instead of fixing the code.

## int became 64 bits, and the language still does not say so
→ candidate | layout
Go 1.1: "The language allows the implementation to choose whether the `int` type
and `uint` types are 32 or 64 bits. Previous Go implementations made `int` and
`uint` 32 bits on all systems. Both the gc and gccgo implementations now make
`int` and `uint` 64 bits on 64-bit platforms such as AMD64/x86-64."

The first clause is the finding. `int` is *implementation-defined*, which means
`unsafe.Sizeof(int(0))` is a property of the build and not of the program —
which is why `layout/01-struct-padding` and `layout/05-alignment-64bit` are
written against fixed-width types, and why a struct's size can differ between
`GOARCH=amd64` and `GOARCH=386` without any field changing. The same paragraph
also explains why a slice header is three words rather than three fixed
quantities.

## Type aliases
→ candidate | types
Go 1.9: "a type alias declaration has the form `type T1 = T2`. This declaration
introduces an alias name `T1` — an alternate spelling — for the type denoted by
`T2`; that is, both `T1` and `T2` denote the same type."

*The same type*, as against a defined type, which is a new one sharing a
representation. The distinction is invisible in a struct's layout and completely
visible everywhere else: an alias shares the original's method set and cannot
take methods of its own, `reflect.TypeOf` cannot tell the two names apart, and a
type switch cannot have a case for each. `byte`/`uint8` and `rune`/`int32` have
always been aliases, which is why `case byte:` and `case uint8:` in one switch is
a compile error — a fact that predates the syntax that explains it.

## The Go 1 compatibility promise
→ no action
Go 1.0: "It is intended that programs written to the Go 1 specification will
continue to compile and run correctly, unchanged, over the lifetime of that
specification."

Nothing to measure, and the reason this radar exists. The promise is what makes
"was true in 1.21 — is it still?" a question with a meaningful answer, and its
stated exceptions are the map of where this catalogue lives: "packages that
import `unsafe` may depend on internal properties of the Go implementation. We
reserve the right to make changes to the implementation that may break such
programs"; a program depending on buggy behaviour "may break if the bug is
fixed"; and "the Go toolchain (compilers, linkers, build tools, and so on) is
under active development and may change behavior."

Every track here except `types` and `iface` studies something the promise
explicitly does not cover. That is not a defect in the catalogue — it is the
definition of internals, and it is why `GODEBUG` settings and `go.mod`-gated
behaviour had to be invented later to make even the *unpromised* parts change
predictably.

## The toolchain stopped being C
→ no action
Go 1.4 translated most of the runtime — "the garbage collector, concurrency
support, interface management, maps, slices, strings" — from C to Go, and Go 1.5
finished the job: "The compiler and runtime are now written entirely in Go (with
a little assembler). C is no longer involved in the implementation, and so the C
compiler that was once necessary for building the distribution is gone."

Foundational, and the precondition for the precise collector two entries above,
since a Go runtime is a runtime the collector can scan. But nothing a program can
observe about itself — which is exactly why it is listed here rather than as a
candidate.

## The race detector arrived
→ no action
Go 1.1: "A major addition to the tools is a race detector, a way to find bugs in
programs caused by concurrent access of the same variable, where at least one of
the accesses is a write ... To enable it, set the `-race` flag when building or
testing your program."

This repository's sixth CI gate is one `go test -race` invocation, and seven
tasks carry `requires.default_build` because their subject cannot survive it. But
the detector is an instrument, not a behaviour: what it *finds* is the memory
model, which `memmodel/02-benign-race` already grades directly.
