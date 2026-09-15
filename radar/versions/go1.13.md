# Go 1.13

Source: https://go.dev/doc/go1.13

## sync.Pool stopped being emptied at every collection
→ candidate | alloc
"`Pool` no longer needs to be completely repopulated after every GC. It now
retains some objects across GCs, as opposed to releasing all objects, reducing
load spikes for heavy users of `Pool`." And separately: "Large `Pool` no longer
increase stop-the-world pause times."

This is the victim cache. Before it, every collection emptied every pool
completely, so a program that pooled expensive objects paid to rebuild all of
them immediately afterwards — the "load spike" the note names. After it, the
current generation is demoted to a victim generation and only dropped at the
*next* collection, so an object survives one GC.

Measurable exactly: put an object in a pool, force one collection, and ask
whether `Get` returns it. Then force a second. The answer changes between the two
and it is a small integer, not a rate.

## A new escape analysis, with an opt-out that no longer exists
→ candidate | alloc
"The compiler has a new implementation of escape analysis that is more precise.
For most Go code should be an improvement (more variables and expressions
allocated on stack instead of heap) ... This new escape analysis may break invalid
code that violates `unsafe.Pointer` safety rules ... the old escape analysis can
be re-enabled with `go build -gcflags=all=-newescape=false`. The option will be
removed in a future release."

`alloc/03-what-escapes` grades which values reach the heap on a current
toolchain. This is the release where the answers changed wholesale, and the note
is unusually candid about the hazard: code that was already violating the
`unsafe.Pointer` rules could *depend* on a variable being heap-allocated, and a
more precise analyser takes that away. The flag is long gone, which makes this a
version-diff that needs two toolchains rather than one flag.

## Defer got 30% faster, one release before it got free
→ candidate | edges
"This release improves the performance of most uses of `defer` by 30%."

Then 1.14: "almost zero overhead compared to calling the deferred function
directly." Two consecutive releases, two different mechanisms — this one made the
existing heap-allocated defer record cheaper, and the next removed it entirely
for the common case. `edges/04-defer-tiers` grades the end state; the arc is what
makes "which tier is this defer in" a question with history.

## Lock fast paths became inlinable
→ candidate | compiler
"The fast paths of `Mutex.Lock`, `Mutex.Unlock`, `RWMutex.Lock`, `RWMutex.RUnlock`
and `Once.Do` are now inlined ... For uncontended cases on amd64, `Once.Do` is
twice as fast and `Mutex` and `RWMutex` methods are up to 10% faster."

`compiler/01-inlining-budget` grades what fits. This is the standard library
being *restructured* to fit: the fast path was split out so that the part which
usually runs is small enough to inline and the slow path stays behind a call.
That is a design technique, not an optimisation the compiler found, and it is
visible in the source.

`memmodel/05-publication` grades `sync.Once` as a publication strategy — this is
where its uncontended cost halved.

## Memory started coming back promptly, and RSS still did not move
→ candidate | gc
"The runtime is now more aggressive at returning memory to the operating system
... Previously, the runtime could retain memory for 5+ minutes after a heap size
spike, it now begins returning it promptly after the heap shrinks." And the
caveat that matters: "on many OSes, including Linux, the OS itself reclaims
memory lazily, so process RSS will not decrease until the system is under memory
pressure."

Pair this with the 1.16 entry, which switched the default from `MADV_FREE` to
`MADV_DONTNEED` precisely so RSS would start moving. Together they are one
question asked twice: what does it mean for a program to "return" memory, and
which of the several numbers reporting it actually changed?

## math/bits gained a constant-time guarantee
→ candidate | asm
"The execution time of `Add`, `Sub`, `Mul`, `RotateLeft` and `ReverseBytes` is now
guaranteed to be independent of the inputs."

`asm/05-carry-chain` is built on `bits.Add64` and asks why the intrinsic exists.
Part of the answer is here: it is not only that Go cannot express the carry flag,
but that the function carries a *timing* promise a hand-written Go fallback with
a branch would break. A documented guarantee about execution time is rare enough
in Go to be worth a task on its own.

## Out-of-range panics started naming the index and the length
→ candidate | compiler
"`s[3]` where `s` is a slice of length 1 produces `runtime error: index out of
range [3] with length 1`."

`compiler/02-bounds-checks` grades which checks survive. This is the other end of
the same machinery: the check that survived now has to carry the index and the
length to the panic, which is information the compiler must keep live to the
point of failure. Predicting what a panic message can contain is a question about
what the bounds check costs.

## reflect.Value.IsZero arrived
→ candidate | reflect
Added here, and then changed in 1.22 to "return true for a floating-point or
complex negative zero" so that it agrees with `==`. Three releases apart, one
function, and the second change is a bug fix to a definition that looked obvious
when it was written. A good pair for "what does zero mean".

## -trimpath removes file system paths from the binary
→ candidate | compiler
"`-trimpath` ... removes all file system paths from the compiled executable, to
improve build reproducibility."

This repository has already met the consequence of *not* using it: M8 measured
that building identical source in two temporary directories of different name
lengths produces binaries of different sizes, because source paths are embedded.
That finding is why the binary-size task was withdrawn, and this is the flag that
would have removed the variable — except that M8 also measured `-trimpath`
shifting sizes by a couple of hundred bytes in both directions.

## Signed shift counts
→ no action
"Go 1.13 removes the restriction that a shift count must be unsigned. This change
eliminates the need for many artificial `uint` conversions." A real ergonomic
improvement and a language change, but nothing underneath it to measure.

## Number literals gained prefixes and separators
→ no action
`0b1011`, `0o660`, `0x1.0p-1021`, `1_000_000`. Syntax, plus matching updates
across `strconv`, `fmt`, `go/scanner` and `text/scanner`.

## Error wrapping arrived
→ no action
`%w`, `errors.Is`, `errors.As`, `errors.Unwrap`. The foundation of modern Go
error handling and worth knowing thoroughly, but it is library design; 1.20's
multiple-error wrapping is the same story continued.

## TLS 1.3 on by default
→ no action
With `GODEBUG=tls13=0` to opt out, removed in 1.14. A good example of the
GODEBUG-gated rollout that 1.21 later formalised, but the protocol itself is not
a machine-level question.
