# Go 1.10

Source: https://go.dev/doc/go1.10

## The assembler stopped silently clobbering the condition flags
→ candidate | asm
"The assembler also no longer implements `MOVL $0, AX` as an `XORL` instruction,
to avoid clearing the condition flags unexpectedly."

`asm/05-carry-chain` is built entirely on the carry flag surviving between
`ADDQ` and `ADCQ`. This is the assembler itself having been a hazard to that: it
rewrote a zeroing move into the shorter `XORL`, which is a well-known peephole
optimisation and which **destroys the flags**. Assembly that looked flag-safe in
source was not in the object file.

A task can ask which instructions clobber flags, and the deeper question is what
an assembler is allowed to change about code you wrote by hand — the answer being
that Plan 9 assembly is not a literal encoding, and that is easy to forget.

## strings.Builder, and why a restricted API can avoid a copy
→ candidate | alloc
"A new type `Builder` is a replacement for `bytes.Buffer` for the use case of
accumulating text into a `string` result. The `Builder`'s API is a restricted
subset of `bytes.Buffer`'s that allows it to safely avoid making a duplicate copy
of the data during the `String` method."

`alloc/01-zero-alloc-join` asks for a join with no allocations, and `Builder` is
the idiomatic answer. The gradeable part is the *reason*: `bytes.Buffer.String()`
must copy because the buffer may be written again afterwards, which would mutate
a string. `Builder` forbids copying itself and has no `Read`, so the bytes are
provably final and the string can alias them.

Same trick, same constraint as `types/03-zero-copy-strings` and
`types/05-compiler-free-conversions`: a `[]byte` can become a `string` for free
exactly when nothing can write to it afterwards. Here the guarantee comes from the
API's shape rather than from the compiler or from `unsafe`.

## Split and Fields started returning capacity-clipped subslices
→ candidate | types
"The `Fields`, `FieldsFunc`, `Split`, and `SplitAfter` functions have always
returned subslices of their inputs. Go 1.10 changes each returned subslice to
have capacity equal to its length, so that appending to one cannot overwrite
adjacent data in the original input."

`types/02-aliasing` grades who still shares a backing array after append, copy
and three-index slicing. This is the standard library adopting the three-index
form defensively across a whole package, and the bug it fixes is exactly that
task's subject: `append` to a subslice with spare capacity writes into the next
field of the original string.

Before 1.10, `bytes.Split(data, sep)[0] = append(...)` could silently corrupt
`[1]`. A clean predict-then-measure: the length is unchanged, the capacity is not.

## LockOSThread calls began to nest
→ candidate | sched
"Previously, calling `LockOSThread` more than once in a row was equivalent to
calling it once, and a single `UnlockOSThread` always unlocked the thread. Now,
the calls nest: if `LockOSThread` is called multiple times, `UnlockOSThread` must
be called the same number of times in order to unlock the thread." And: "the
runtime now treats locked threads as unsuitable for reuse or for creating new
threads."

Spec §8 T10 names `LockOSThread` as `sched` material and it is still unbuilt. Two
observable facts here: the nesting count, and the fact that a locked thread is
**retired** rather than returned to the pool — so a program that locks and
unlocks in a loop accumulates threads. Measurable through the thread count.

## The collector's CPU moved rather than shrank
→ candidate | gc
"The garbage collector has been modified to reduce its impact on allocation
latency. It now uses a smaller fraction of the overall CPU when running, but it
may run more of the time. The total CPU consumed by the garbage collector has not
changed significantly."

Three quantities, and only two of them moved: the instantaneous fraction fell,
the duration rose, the total stayed put. That is the shape of every honest
latency-against-throughput trade, and it is why a single "GC overhead" number
answers almost nothing — which is the point `gc/05-gctrace` makes about the
cumulative percentage being useless for anything time-sensitive.

## reflect stopped letting you set through an unexported embedded pointer
→ candidate | reflect
"In structs, embedded pointers to unexported struct types were previously
incorrectly reported with an empty `PkgPath` ... with the result that for those
fields, `Value.CanSet` incorrectly returned true and `Value.Set` incorrectly
succeeded. The underlying metadata has been corrected; for those fields, `CanSet`
now correctly returns false and `Set` now correctly panics."

`reflect/03-settability` grades which reflect calls panic. This is a case where
the answer was *wrong* for years: embedding a pointer to an unexported type
punched a hole through the export check, and `encoding/json` could unmarshal into
fields it had no business touching. The note says so plainly — "this may affect
reflection-based unmarshalers that could previously unmarshal into such fields
but no longer can."

## Stack traces stopped including compiler-generated wrappers
→ candidate | edges
"Stack traces no longer include implicit wrapper functions (previously marked
`<autogenerated>`), unless a fault or panic happens in the wrapper itself. As a
result, skip counts passed to functions like `Caller` should now always match the
structure of the code as written, rather than depending on optimization decisions
and implementation details."

A skip count that depended on whether the optimiser had inserted a wrapper. 1.12
then removed initialisation functions from tracebacks and 1.17 made inlining
produce multiple frames per PC — three releases of the same struggle, which is
that `runtime.Caller(n)` asks a question about source structure and gets an
answer about generated code.

## The GOMAXPROCS ceiling was removed
→ candidate | sched
"There is no longer a limit on the `GOMAXPROCS` setting. (In Go 1.9 the limit was
1024.)" Worth an entry because every `sched` task pins `GOMAXPROCS` and because
1.25's container-aware default — still an open pool row — is the same setting
gaining a mind of its own.

## Test results became cacheable, and -count=1 became the escape hatch
→ candidate | edges
"`go test` will print the previous test output, replacing the elapsed time with
the string '(cached)' ... The idiomatic way to bypass test caching is to use
`-count=1`."

This repository's own verifier runs `go test -count=1` for exactly this reason, so
the machinery under a task's grade depends on this release. The gradeable question
is what the cache keys on — "the test executable and command line ... and the
files and environment variables consulted by that run" — because a task that
reads an environment variable the cache does not know about would be cached
wrongly.

## go test began running go vet
→ candidate | edges
"The `go test` command now automatically runs `go vet` on the package being
tested, to identify significant problems before running the test. Any such
problems are treated like build errors and prevent execution of the test. Only a
high-confidence subset of the available `go vet` checks are enabled."

`edges/05-unsafe-pointer-rules` keeps its failing snippets in `testdata` precisely
because a package that fails vet cannot be tested. That constraint starts here,
and the "high-confidence subset" is the detail worth knowing: the checks `go test`
runs are not the checks `go vet` runs.

## The build cache arrived
→ no action
"a cache of recently built packages ... The old advice to add the `-i` flag for
speed ... is no longer necessary", and out-of-date detection moved to content
hashing: "Modification times are no longer consulted or relevant." A foundational
change to how Go builds, and entirely invisible to a running program.

## The x86 assembler gained 359 instructions
→ no action
"the full AVX, AVX2, BMI, BMI2, F16C, FMA3, SSE2, SSE3, SSSE3, SSE4.1, and SSE4.2
extension sets." Background for the `GOAMD64` and `simd` pool rows — FMA3 is what
1.14's `math.FMA` and `GOAMD64=v3` later reach for — but accepting an instruction
is not itself observable.

## DWARF debug information improved
→ no action
"constant values are now recorded; line number information is more accurate ...
each package is now presented as its own DWARF compilation unit." Debugger
quality, not program behaviour.
