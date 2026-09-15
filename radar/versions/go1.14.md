# Go 1.14

Source: https://go.dev/doc/go1.14

## Goroutines became asynchronously preemptible
→ candidate | sched
"Goroutines are now asynchronously preemptible. As a result, loops without
function calls no longer potentially deadlock the scheduler or significantly
delay garbage collection. This is supported on all platforms except
`windows/arm`, `darwin/arm`, `js/wasm`, and `plan9/*`."

`sched/04-async-preemption` grades exactly this, using
`GODEBUG=asyncpreemptoff=1` as the A/B. The release note adds the part that task
does not state: the mechanism is **not universal**. On the excluded platforms the
pre-1.14 behaviour is still what happens, so a tight loop still freezes the
world. Neither of undergo's CI runners is on that list, but the exclusion is the
honest answer to "is this fixed" — it is fixed where signals can do the job.

## Defer became almost free
→ candidate | edges
"This release improves the performance of most uses of `defer` to incur almost
zero overhead compared to calling the deferred function directly. As a result,
`defer` can now be used in performance-critical code without overhead concerns."

`edges/04-defer-tiers` grades which defers are open-coded and which fall back to
the heap. The words that matter here are **most uses**: the optimisation applies
to defers whose count is known at compile time and which are not in a loop, so
the old cost is still there for the cases that fail those conditions. A release
note saying "almost zero" and a task asking "which ones" are the same fact at two
resolutions.

## checkptr arrived, with its two rules stated
→ candidate | edges
"`-d=checkptr` as a compile-time option for adding instrumentation to check that
Go code is following `unsafe.Pointer` safety rules dynamically", enabled by
default with `-race` or `-msan` except on Windows, and it checks exactly two
things:

1. "When converting `unsafe.Pointer` to `*T`, the resulting pointer must be
   aligned appropriately for `T`."
2. "If the result of pointer arithmetic points into a Go heap object, one of the
   `unsafe.Pointer`-typed operands must point into the same object."

`edges/05-unsafe-pointer-rules` measures that `go vet` finds two of three
violations and misses the out-of-bounds one, and this is why: the second rule is
a *value* property that only exists at run time. The checks were specified here,
five releases before that task's `-race` requirement became universal in 1.15.

## Bounds-check elimination learned about slice creation
→ candidate | compiler
"Bounds check elimination now uses information from slice creation and can
eliminate checks for indexes with types smaller than `int`."

`compiler/02-bounds-checks` grades which checks survive. Two distinct new
sources of information here — the length a slice was created with, and the range
implied by a narrow index type — and the second is the more surprising: indexing
with a `uint8` cannot exceed 255, so a slice known to be longer needs no check at
all.

## math.FMA, and one rounding instead of two
→ candidate | edges
"The new `FMA` function computes `x*y+z` in floating point with no intermediate
rounding of the `x*y` computation. Several architectures implement this
computation using dedicated hardware instructions for additional performance."

This is the explicit spelling of what `GOAMD64=v3` later lets the compiler do
implicitly, which is the open pool row about floating-point results changing.
`math.FMA(x, y, z)` and `x*y + z` can produce **different `float64` values** for
the same inputs, and predicting which inputs expose the difference is a real
question with an exact answer.

## Overlapping interface embedding became legal
→ candidate | iface
"methods from an embedded interface may have the same names and identical
signatures as methods already present in the (embedding) interface. This solves
problems that typically (but not exclusively) occur with diamond-shaped embedding
graphs. Explicitly declared methods in an interface must remain unique, as
before."

Two rules that look alike and are not: duplication through embedding is fine,
duplication by explicit declaration is not. And *identical* signatures — a
same-named method with a different signature is still an error, which is the
condition 1.15's "impossible interface assertion" vet check keys on. Predict
which of a set of interface declarations compile.

## The page allocator stopped contending
→ candidate | alloc
"The page allocator is more efficient and incurs significantly less lock
contention at high values of `GOMAXPROCS`. This is most noticeable as lower
latency and higher throughput for large allocations being done in parallel and at
a high rate."

Note *large* allocations — the ones above the size-class ceiling that
`alloc/05-size-classes` measures, which come from the page allocator directly
rather than from a per-P cache. The 32 KiB boundary that task grades is exactly
the line between the two allocators, and this is the release where crossing it
stopped being a scalability cliff.

## Unlocking a contended mutex yields directly to the next waiter
→ candidate | sched
"Unlocking a highly contended `Mutex` now directly yields the CPU to the next
goroutine waiting for that `Mutex`. This significantly improves the performance
of highly contended mutexes on high CPU count machines."

A direct handoff rather than a wake-and-race. That changes the *order* goroutines
acquire a lock in, not just the speed — which is the kind of observable
`sched/01-runnext-ordering` is built on, and it interacts with `runnext`: the
woken goroutine goes into that slot.

## Timers became cheaper and stopped switching contexts
→ candidate | sched
"Internal timers, used by `time.After`, `time.Tick`, `net.Conn.SetDeadline`, and
friends, are more efficient, with less lock contention and fewer context
switches. This is a performance improvement that should not cause any user
visible changes."

The last sentence is the claim to test. `sched/03-timer-channels` grades timer
behaviour that *did* become user-visible in 1.23 and 1.27, so this is the first
step of a rewrite whose later steps changed semantics. Predicting where an
invisible change becomes a visible one is the version-diff genre at its best.

## runtime.Goexit can no longer be aborted
→ candidate | edges
"`runtime.Goexit` can no longer be aborted by a recursive `panic`/`recover`."

`edges/03-panic-vs-fatal` grades what `recover` can and cannot catch. `Goexit` is
a third thing — neither a panic nor a normal return — and until this release a
sufficiently determined `recover` could cancel it, leaving a goroutine running
that the runtime believed was gone.

## reflect.StructOf gained unexported fields
→ candidate | reflect
"`StructOf` now supports creating struct types with unexported fields, by setting
the `PkgPath` field in a `StructField` element." `reflect/02-structof` builds a
type at run time; this is where the export boundary became something you could
construct rather than only observe.

## hash/maphash is per-process
→ candidate | edges
"The hash value of a given byte sequence is consistent within a single process,
but will be different in different processes."

The same determinism property Go's map iteration order has, stated as an API
guarantee. Worth an entry because the reason is the same in both cases — hash
flooding — and because "consistent within a process" is exactly the guarantee
people accidentally rely on across processes.

## The compiler can emit optimization decisions as JSON
→ candidate | compiler
"machine-readable logs of key optimizations using the `-json` flag, including
inlining, escape analysis, bounds-check elimination, and nil-check elimination",
and "detailed escape analysis diagnostics (`-m=2`) now work again."

Several tasks in this repository parse `-m` or `-S` text output. A structured
form of the same decisions is the more robust instrument, and knowing it exists
is worth as much as the parsing.

## Go binaries on Windows have DEP enabled
→ no action
Data Execution Prevention. A platform hardening default with no observable
behaviour from Go code.

## js.Value can no longer be compared with ==
→ no action
"`js.Value` values can no longer be compared using the `==` operator, and instead
must be compared using their `Equal` method." An interesting echo of
`iface/05-interface-comparison` — a type that stopped being comparable — but it
is confined to `js/wasm`, which this repository does not target.

## Vendoring became the default when a vendor directory is present
→ no action
`-mod=vendor` by default for modules declaring `go 1.14` with a top-level
`vendor` directory. Build tooling.
