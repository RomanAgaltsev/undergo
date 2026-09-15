# Go 1.21

Source: https://go.dev/doc/go1.21

## GODEBUG and the go line became the compatibility mechanism
→ candidate | edges
"Go 1.21 formalizes Go's use of the GODEBUG environment variable to control the
default behavior for changes that are non-breaking according to the compatibility
policy but nonetheless may cause existing programs to break ... it now chooses
between the old and new behavior based on the `go` line in the workspace's
`go.work` file or else the main module's `go.mod` file."

This is the machinery behind every "gated on the `go.mod` line" answer in the
catalogue — `sched/03-timer-channels` grades two such behaviours without naming
where the mechanism came from. The task is to predict which of a set of
behaviours follow the module and which follow the toolchain, and to say why the
compatibility promise forced the distinction.

## panic(nil) stopped being nil
→ candidate | edges
"Go 1.21 now defines that if a goroutine is panicking and recover was called
directly by a deferred function, the return value of recover is guaranteed not
to be nil. To ensure this, calling panic with a nil interface value (or an
untyped nil) causes a run-time panic of type `*runtime.PanicNilError`."

Before 1.21, `panic(nil)` was recovered with `recover()` returning `nil`, which
is indistinguishable from not panicking at all — so `if r := recover(); r != nil`
silently swallowed it. Gated: `GODEBUG=panicnil=1`, set automatically for a
module declaring `go 1.20` or earlier. This extends `edges/03-panic-vs-fatal`
naturally, which currently asks what recover cannot catch and not what recover
cannot distinguish.

## reflect.ValueOf stopped forcing its argument to the heap
→ candidate | reflect
"In Go 1.21, `ValueOf` no longer forces its argument to be allocated on the heap,
allowing a `Value`'s content to be allocated on the stack. Most operations on a
`Value` also allow the underlying value to be stack allocated."

Directly measurable with `testing.AllocsPerRun`: the same reflection code
allocates differently across this boundary. The gradeable question is which
`reflect.Value` operations still force the heap and why — the answer is about
escape analysis meeting an API that erases types.

## Package initialisation order became an algorithm
→ candidate | edges
The notes give it as a procedure: "Sort all packages by import path. Repeat until
the list of packages is empty: find the first package in the list for which all
imports are already initialized, initialize that package and remove it." And the
consequence: "This may change the behavior of some programs that rely on a
specific initialization ordering that was not expressed by explicit imports."

Observable with `init` functions that print. The mechanism is the answer: given a
dependency graph, produce the order, and say which orderings were possible before
and are not now.

## Type inference gained four specific new abilities
→ candidate | generics
The notes enumerate them rather than summarising, which is what makes them
gradeable: generic functions may now be passed as arguments to generic functions
with inference of the callee's type arguments; "type inference now considers
methods when a value is assigned to an interface"; methods of type arguments and
constraints are matched; and untyped constant arguments of different kinds are
resolved "using the same approach as an operator with untyped constant operands"
rather than erroring.

Also a tightening: "Type inference is now precise when matching corresponding
types in assignments: component types ... must be identical (given suitable type
arguments) to match, otherwise inference fails."
`generics/05-inference-limits` grades which calls need an explicit type argument;
this release is where several of its answers changed.

## Frameless nosplit assembly functions are no longer automatically NOFRAME
→ candidate | asm
"On amd64, frameless nosplit assembly functions are no longer automatically
marked as `NOFRAME`. Instead, the `NOFRAME` attribute must be explicitly
specified if desired, which is already the behavior on other architectures
supporting frame pointers. With this, the runtime now maintains the frame
pointers for stack transitions."

Every assembly task in this repository writes `NOSPLIT` without thinking about
it. The question is what the frame pointer buys — profiling and unwinding — and
what an explicit `NOFRAME` gives up.

## The linker deletes dead global map variables
→ candidate | reflect
"the linker (with help from the compiler) is now capable of deleting dead
(unreferenced) global map variables, if the number of entries in the variable
initializer is sufficiently large, and if the initializer expressions are
side-effect free."

Two conditions, both checkable: *sufficiently large* and *side-effect free*.
`reflect/05-methodbyname-linker` counts retained symbols; this is the same
question asked of data rather than code, and the size threshold makes it a
genuine predict-then-measure.

## Background and TODO can now compare equal
→ candidate | iface
"An optimization means that the results of calling `Background` and `TODO` and
converting them to a shared type can be considered equal. In previous releases
they were always different."

`iface/05-interface-comparison` is built on the two-step comparison of dynamic
type then dynamic value. This is that rule producing a surprising answer in the
standard library: two context values that were reliably unequal became equal
because of a representation change, not a semantic one.

## Huge pages are managed explicitly on Linux
→ candidate | gc
"small heaps should see less memory used (up to 50% in pathological cases) while
large heaps should see fewer broken huge pages for dense parts of the heap,
improving CPU usage and latency by up to 1%." And a warning worth an entry on its
own: "the runtime no longer tries to work around a particular problematic Linux
configuration setting, which may result in higher memory overheads."

Measurable through `runtime/metrics` heap and RSS figures, and the "pathological
case" is the interesting half: what shape of heap was being penalised.

## GOGC and GOMEMLIMIT became metrics
→ candidate | gc
"A few previously-internal GC metrics, such as live heap size, are now available.
`GOGC` and `GOMEMLIMIT` are also now available as metrics." `gc/01-heap-goal` and
`gc/03-memory-limit` both compute against these values; reading them back from
`runtime/metrics` closes the loop and makes a self-checking task possible.

## The cgo call boundary got an order of magnitude cheaper
→ candidate | edges
"On Unix platforms, this setup is now preserved across multiple calls from the
same thread. This significantly reduces the overhead of subsequent C to Go calls
from ~1-3 microseconds per call to ~100-200 nanoseconds per call."

A ten-fold difference between the first call on a thread and the rest, which is a
mechanism question rather than a benchmark: what is being set up, and why does it
survive. Blocked for the same reason as the existing cgo pool row — this machine
has no C toolchain — but recorded so the row exists when one does.

## sync.OnceFunc, OnceValue and OnceValues
→ candidate | memmodel
"capture a common use of `Once` to lazily initialize a value on first use."
`memmodel/05-publication` grades five lazy-initialisation strategies and these
are a sixth, with the interesting wrinkle that `OnceValue` returns a *function*
rather than guarding a field — so the happens-before edge is established by a
different act.

## runtime.Pinner
→ candidate | edges
"Pinners may be used to 'pin' Go memory such that it may be used more freely by
non-Go code. For instance, passing Go values that reference pinned Go memory to
C code is now allowed. Previously, passing any such nested reference was
disallowed by the cgo pointer passing rules."

A direct extension of `edges/05-unsafe-pointer-rules`: pinning changes which
programs are legal, and the question is what the collector promises about a
pinned object that it does not promise about any other.

## Stack traces name the goroutine that created each goroutine
→ candidate | sched
"Textual stack traces produced by Go programs ... now include the IDs of the
goroutines that created each goroutine in the stack trace." The creator edge is
exactly what turns a goroutine leak from "something leaked" into "this call site
leaked", which is the question `sched/05-goroutine-leak-profile` asks of the
profile rather than of a traceback.

## Tail latency fell by up to 40%, with a stated tradeoff
→ no action
"applications may see up to a 40% reduction in application tail latency and a
small decrease in memory use. Some applications may also observe a small loss in
throughput ... the previous release's throughput/memory tradeoff may be recovered
... by increasing `GOGC` and/or `GOMEMLIMIT` slightly." The tradeoff is real and
worth knowing, but every number here is a distribution over unnamed applications
— exactly the shape the `gc` track already refuses to grade.

## min, max and clear
→ no action
`clear` matters to `compiler/05-loop-lowering`, which measures that it compiles
to the same `memclr` call as the hand-written loop — but that task already exists
and the builtins themselves are syntax.

## Build speed improved by up to 6%
→ no action
"largely thanks to building the compiler itself with PGO." A fact about the
toolchain's own build, not observable from any Go program.

## cgocheck=2 moved from GODEBUG to GOEXPERIMENT
→ no action
"this mode has to be selected at build time instead of startup time." A real
distinction — build-time against startup-time selection — but the checker itself
is what would be graded, and that needs a C toolchain.
