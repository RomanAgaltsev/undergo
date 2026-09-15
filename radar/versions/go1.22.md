# Go 1.22

Source: https://go.dev/doc/go1.22

## Loop variables are created anew each iteration
→ candidate | edges
"Previously, the variables declared by a 'for' loop were created once and updated
by each iteration. In Go 1.22, each iteration of the loop creates new variables."
Gated on the module's `go` line, so one program compiled two ways gives two
different printed outputs — a value, not a judgement. The task is not "what does
this print" alone but "which declaration site does the closure capture, and what
does the compiler emit differently". Pairs with the `//go:debug` mechanism from
1.21: the behaviour follows the module, not the toolchain.

## Heap metadata moved next to the object, and alignment fell from 16 to 8
→ candidate | layout
"The runtime now keeps type-based garbage collection metadata nearer to each
heap object, improving the CPU performance ... by 1–3%" and reducing memory "by
approximately 1% by deduplicating redundant metadata". The gradeable consequence
is the one the notes flag as a hazard: "some objects' addresses that were
previously always aligned to a 16 byte (or higher) boundary will now only be
aligned to an 8 byte boundary." `GOEXPERIMENT=noallocheaders` restores the old
layout, so the alignment of a heap object is an A/B within one toolchain. Note
this is *allocation* alignment, not `unsafe.Alignof`, which is a property of the
type — `layout/05-alignment-64bit` measures the second and is unaffected.

## The same change moved the size class boundaries
→ candidate | alloc
"this change adjusts the size class boundaries of the memory allocator, so some
objects may be moved up a size class." `alloc/05-size-classes` grades the cost
curve on 1.27, which already includes this; the task here is the diff — which
sizes changed class, and why moving metadata next to the object costs a class
boundary. Needs the `noallocheaders` A/B to be answerable.

## PGO devirtualizes more, and devirtualization now interleaves with inlining
→ candidate | compiler
"Profile-guided Optimization (PGO) builds can now devirtualize a higher
proportion of calls than previously possible" and "the compiler now interleaves
devirtualization and inlining, so interface method calls are better optimized."
Most programs "see between 2 and 14% improvement". The interleaving is the
gradeable half and the reason `compiler/04-pgo-devirtualization` works at all: a
devirtualized call becomes an inlining candidate in the same pass, so the two
optimisations compound rather than running in sequence.

## An inliner that weighs the call site, not just the callee
→ candidate | compiler
`GOEXPERIMENT=newinliner` enables "heuristics to boost inlinability at call sites
deemed 'important' (for example, in loops) and discourage inlining at call sites
deemed 'unimportant' (for example, on panic paths)." `compiler/01-inlining-budget`
grades a budget charged against the callee alone. This is the same question asked
of the caller, and it is an A/B in one toolchain.

## Shrinking a slice now zeroes the tail
→ candidate | types
"Functions that shrink the size of a slice (`Delete`, `DeleteFunc`, `Compact`,
`CompactFunc`, and `Replace`) now zero the elements between the new length and
the old length." Directly observable through the backing array: reslice past the
new length after a `slices.Delete` and the old elements are gone. The mechanism
is the answer — this exists so that dropped pointers do not keep objects alive,
which makes it a `weak`/`gc` question wearing a `types` hat.

## Mutex profiles now scale by the number of blocked goroutines
→ candidate | sched
"if 100 goroutines are blocked on a mutex for 10 milliseconds, a mutex profile
will now record 1 second of delay instead of 10 milliseconds." A pre-1.22 profile
understated contention by exactly the blocked-goroutine count. Predict the
recorded delay for a given contention shape, then explain which number answers
"how long was this mutex held" and which answers "how much work was waiting".

## Stop-the-world pauses split into stopping and total
→ candidate | gc
Four new histograms: `/sched/pauses/stopping/gc:seconds`,
`/sched/pauses/stopping/other:seconds`, `/sched/pauses/total/gc:seconds`,
`/sched/pauses/total/other:seconds`. "The 'stopping' metrics report the time
taken from deciding to stop the world until all goroutines are stopped." That
difference is the preemption cost, which is exactly what `sched/04`'s
non-preemptible loop inflates — the two metrics together make that measurable
rather than anecdotal.

## The execution tracer was rewritten and can be streamed
→ candidate | sched
Traces are "partitioned regularly on-the-fly and as a result may be processed in
a streamable way", carry "complete durations for all system calls" and
"information about the operating system threads that goroutines executed on", and
no longer "depend on the reliability of the platform's clock". This is the floor
for any trace-analyzer task: a pre-1.22 trace cannot answer the thread question
at all. `GOEXPERIMENT=noexectracer2` restores the old implementation.

## reflect.TypeFor
→ candidate | reflect
"Previously, to get the `reflect.Type` value for a type, one had to use
`reflect.TypeOf((*T)(nil)).Elem()`. This may now be written as
`reflect.TypeFor[T]()`." The interesting version of the question is why the old
spelling needed a nil pointer and an `Elem()` at all, which is a question about
what a `reflect.Type` is obtained *from*.

## IsZero agrees with == for negative zero
→ candidate | types
"`Value.IsZero` will now return true for a floating-point or complex negative
zero, and will return true for a struct value if a blank field (a field named
`_`) somehow has a non-zero value. These changes make `IsZero` consistent with
comparing a value to zero using the language `==` operator." Two separate
surprises in one sentence: that `-0.0 == 0.0` is true, and that a blank field can
hold a value at all.

## The linker's -s now implies -w
→ no action
"`-s` also implies the `-w` flag, which can be negated with `-w=0`." Real and
worth knowing — `reflect/05-methodbyname-linker` asks about `-w` — but it is a
flag convention rather than something a solver can measure a mechanism from.

## darwin/amd64 builds position-independent executables by default
→ no action
Relevant background for the open `edges` pool row on heap address randomisation,
since PIE is what makes load addresses vary, but the change itself is a build
default rather than an observable behaviour.

## Range over integers
→ no action
`for range 10`. Used throughout this repository's own task code, but it is
syntax: there is no mechanism underneath for a solver to be wrong about.

## math/rand's top-level functions are seeded randomly and use ChaCha8
→ no action
"The global generator accessed by top-level functions is unconditionally randomly
seeded" and ChaCha8 is now used by "`math/rand`'s top-level functions (when not
explicitly seeded) and the Go runtime". A genuine determinism change, but the
observable — that output differs between runs — is the documented guarantee
rather than a surprise.
