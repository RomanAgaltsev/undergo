# Go 1.11

Source: https://go.dev/doc/go1.11

## Bounds-check elimination learned to chain inequalities
→ candidate | compiler
"The compiler now performs significantly more aggressive bounds-check and branch
elimination. Notably, it now recognizes transitive relations, so if `i<j` and
`j<len(s)`, it can use these facts to eliminate the bounds check for `s[i]`. It
also understands simple arithmetic such as `s[i-10]` and can recognize more
inductive cases in loops. Furthermore, the compiler now uses bounds information
to more aggressively optimize shift operations."

`compiler/02-bounds-checks` grades which checks survive, and M8 measured one of
its answers directly: `xs[i+1]` keeps **no** bounds check because the compiler
chains the inequality. That capability starts here. Four distinct new abilities
in one paragraph — transitivity, arithmetic offsets, induction, and shift
narrowing — and each is separately testable with `-d=ssa/check_bce`.

## The map-clearing idiom became one call
→ candidate | compiler
"The compiler now optimizes map clearing operations of the form
`for k := range m { delete(m, k) }`."

`compiler/05-loop-lowering` grades the slice version of exactly this — a
hand-written zeroing loop recognised and replaced with `memclr`. This is the map
analogue, and it is the reason `clear()` on a map in 1.21 could be defined
without a performance argument: the idiom it replaces had already been a single
runtime call for three years.

The gradeable pair is what the two idioms have in common — a loop whose shape the
compiler matches — and what they do not: a map clear cannot be `memclr`, because
the buckets hold pointers.

## Slice extension by append got its own optimisation
→ candidate | types
"The compiler now optimizes slice extension of the form `append(s, make([]T, n)...)`."

The idiom for growing a slice by `n` zero values. Without the optimisation it
allocates a temporary slice, zeroes it, and copies it in; with it, the
destination is grown and zeroed directly. Measurable with `AllocsPerRun`, and it
connects `types/01-cap-growth` to `alloc/05-size-classes`: the growth is one
allocation at the final size rather than two at different ones.

## The heap became sparse and lost its 512 GiB ceiling
→ candidate | alloc
"The runtime now uses a sparse heap layout so there is no longer a limit to the
size of the Go heap (previously, the limit was 512GiB). This also fixes rare
'address space conflict' failures in mixed Go/C binaries or binaries compiled
with `-race`."

Before this the runtime reserved one contiguous arena, which is why the limit
existed and why a C library or the race detector mapping something in the way
could break the process outright. The open `edges` pool row about the heap base
address being randomised in 1.26 is the same subject eight releases later: what
the runtime is allowed to assume about where its memory lives.

## Non-blocking descriptors started using the poller instead of a thread
→ candidate | sched
"When a non-blocking descriptor is passed to `NewFile`, the resulting `*File`
will be kept in non-blocking mode. This means that I/O for that `*File` will use
the runtime poller rather than a separate thread, and that the `SetDeadline`
methods will work."

Two mechanisms for the same operation, chosen by a property of the file
descriptor. A blocking descriptor parks an OS thread for the duration of the
read; a non-blocking one registers with netpoll and parks only the goroutine.
That is the difference between costing a thread and costing nothing, and it is
observable through the thread count.

## Tracebacks can name the goroutine that created a goroutine
→ candidate | sched
"Setting the `GODEBUG=tracebackancestors=N` environment variable now extends
tracebacks with the stacks at which goroutines were created, where N limits the
number of ancestor goroutines to report."

1.21 made a reduced form of this unconditional — creator IDs in every traceback —
so this is the opt-in ancestor that became a default. For a leaked goroutine the
creation stack is the only thing that says *which call site* leaked, which is the
question `sched/05-goroutine-leak-profile` asks of the profile.

## Functions that call panic became inlinable
→ candidate | compiler
"More functions are now eligible for inlining by default, including functions
that call `panic`."

A panic path used to make a whole function ineligible, which penalised exactly
the small argument-checking wrappers that should be inlined — the function whose
body is `if x == nil { panic(...) }; return x.f`. 1.21's `newinliner` then went
the other way and *discouraged* inlining at call sites on panic paths, which is
the same observation used twice: a panic is cold, so the code around it should be
cheap to reach and not worth expanding.

## The mutex profile learned about reader/writer contention
→ candidate | sched
"The mutex profile now includes reader/writer contention for `RWMutex`.
Writer/writer contention was already included in the mutex profile."

An `RWMutex` whose readers were starving a writer produced an empty profile
before this. With 1.17's block-profile debiasing and 1.22's contention scaling,
this is the third correction the radar has found to Go's contention profiling —
each one a case where the tool reported less than was happening.

## The allocs profile arrived
→ candidate | alloc
"a new 'allocs' profile type that profiles total number of bytes allocated since
the program began (including garbage-collected bytes) ... Now
`go test -memprofile=...` reports an 'allocs' profile instead of 'heap' profile."

`alloc/05-size-classes` measures exactly this quantity through
`AllocedBytesPerOp`, which is `TotalAlloc` differenced. The distinction the note
draws — bytes *allocated ever* against bytes *live now* — is the one that makes
"my program allocates 4 GB" and "my program uses 40 MB" both true.

## The compiler rejects an unused type switch guard
→ candidate | edges
"The compiler now rejects unused variables declared in a type switch guard, such
as `x` in `switch x := v.(type) {}`. This was already rejected by both `gccgo`
and go/types."

A place where the language's three implementations disagreed, and gc was the
lenient one. Worth an entry because it is a reminder that "the compiler accepts
it" and "it is Go" are different claims — the same distinction `layout/05` makes
about struct field order.

## DWARF became compressed by default
→ no action
"DWARF sections are now compressed by default because of the expanded and more
accurate debug information produced by the compiler ... To disable DWARF
compression, pass `-ldflags=-compressdwarf=false`." Binary size, with a flag.

## Modules arrived, experimentally
→ no action
"preliminary support for a new concept called 'modules,' an alternative to GOPATH
with integrated support for versioning and package distribution." The beginning
of the story that 1.16 finished by making it the default. Nothing underneath Go.

## WebAssembly
→ no action
An experimental `js/wasm` port whose binaries are "at minimum around 2 MB, or 500
KB compressed". A port this repository does not target, and one where several
answers — async preemption, for instance — differ.

## The assembler accepts AVX512
→ no action
Relevant background for the open `asm` pool row about the experimental `simd`
package, but accepting an instruction is not a behaviour a task can measure.
