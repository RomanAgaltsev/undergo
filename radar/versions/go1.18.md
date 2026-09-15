# Go 1.18

Source: https://go.dev/doc/go1.18

## Generics, and six limitations stated as a list
→ candidate | generics
The release enumerates what the compiler could not yet do, which is what makes it
gradeable rather than a feature announcement. Among them: "The Go compiler cannot
handle type declarations inside generic functions or methods"; "only supports
calling a method `m` on a value `x` of type parameter type `P` if `m` is
explicitly declared by `P`'s constraint interface"; "does not support accessing a
struct field `x.f` where `x` is of type parameter type even if all types in the
type parameter's type set have a field `f`"; and "Embedding a type parameter, or
a pointer to a type parameter, as an unnamed field in a struct type is not
permitted."

Several were later lifted — 1.20's front-end rewrite "enables type declarations
within generic functions and methods" — and the field-access one still stands.
"Which of these compiles, and on which release" is a version-diff task with a
list of answers rather than a judgement.

## append's growth formula changed
→ candidate | types
"The built-in function `append` now uses a slightly different formula when
deciding how much to grow a slice when it must allocate a new underlying array.
The new formula is less prone to sudden transitions in allocation behavior."

`types/01-cap-growth` grades the exact capacity sequence, measured on a current
toolchain. This release is where that sequence changed shape: the old formula
doubled until a fixed threshold and then switched abruptly, and the new one eases
between the two regimes. The task is the diff — the same `append` loop produces
two different sequences, and the second is smoother in a way the release note
describes but does not quantify.

## The register ABI reached arm64 and ppc64
→ candidate | asm
"Go 1.17 implemented a new way of passing function arguments and results using
registers instead of the stack on 64-bit x86 architecture on selected operating
systems. Go 1.18 expands the supported platforms to include 64-bit ARM
(`GOARCH=arm64`), big- and little-endian 64-bit PowerPC ... On 64-bit ARM and
64-bit PowerPC systems, benchmarking shows typical performance improvements of
10% or more."

`asm/03-register-abi` is pinned to amd64 and grades the nine-register integer
sequence. arm64 has a different sequence and a different count, so the same task
asked on the other architecture has different answers — which is exactly why that
task is pinned rather than written portably.

## Stack traces mark register-passed arguments as possibly wrong
→ candidate | asm
"Go 1.17 generally improved the formatting of arguments in stack traces, but
could print inaccurate values for arguments passed in registers. This is improved
in Go 1.18 by printing a question mark (`?`) after each value that may be
inaccurate."

A direct, visible consequence of the register ABI: once an argument lives in a
register, it may have been overwritten by the time the traceback is produced, so
the runtime can no longer promise the value it prints. Predict which arguments in
a traceback carry the `?`.

## GOAMD64 selects a microarchitecture level
→ candidate | edges
"the new `GOAMD64` environment variable, which selects at compile time a minimum
target version of the AMD64 architecture ... `v1`, `v2`, `v3`, or `v4`. Each
higher level requires, and takes advantage of, additional processor features. The
`GOAMD64` environment variable defaults to `v1`."

This is the switch behind the open pool row about fused multiply-add changing
floating-point results: at `v3` the compiler may use FMA, which computes
`a*b+c` with one rounding instead of two, so the same expression produces a
different `float64`. A compile-time flag that changes arithmetic results is about
as good a predict-then-measure as this catalogue gets.

## Functions containing range loops became inlinable
→ candidate | compiler
"The compiler now can inline functions that contain range loops or labeled for
loops."

Before this, a `range` anywhere in a function made it non-inlinable outright,
regardless of cost. `compiler/01-inlining-budget` grades what fits the budget;
this release is where the rule changed from a structural veto to a cost. The
version-diff question is which of a set of small functions were inlinable before
and after, and why a structural veto existed at all.

## The pacer started counting stack scanning
→ candidate | gc
"The garbage collector now includes non-heap sources of garbage collector work
(e.g., stack scanning) when determining how frequently to run. As a result,
garbage collector overhead is more predictable when these sources are
significant." And the honest caveat: "some Go applications may now use less
memory and spend more time on garbage collection, or vice versa, than before.
The intended workaround is to tweak `GOGC`."

`gc/01-heap-goal` computes the goal from the live heap. This is the release where
that input stopped being the whole story: a program with many deep goroutine
stacks now collects more often at the same heap size, which is measurable by
holding the heap constant and varying stack depth.

## netip.Addr is comparable and net.IP is not
→ candidate | iface
"Compared to the existing `net.IP` type, the `netip.Addr` type takes less memory,
is immutable, and is comparable so it supports `==` and can be used as a map
key."

`iface/05-interface-comparison` grades which comparisons panic, and this is the
standard library's own worked example: `net.IP` is a `[]byte`, so two of them in
interfaces panic on comparison and cannot key a map, while `netip.Addr` is an
array-backed struct and can. The same information, two types, one design decision
— and the memory difference is measurable with `unsafe.Sizeof`.

## Trim became allocation free
→ candidate | alloc
"`Trim`, `TrimLeft`, and `TrimRight` are now allocation free and, especially for
small ASCII cutsets, up to 10 times faster."

Measurable with `testing.AllocsPerRun`, and the mechanism is the interesting
half: trimming returns a subslice of its input, so the allocation that used to
happen was in building the cutset, not in the result. `alloc/02-allocs-per-op`
asks exactly this kind of question.

## strings.Clone exists to stop sharing memory
→ candidate | types
"`Clone` copies the input `string` without the returned cloned `string`
referencing the input string's memory."

The whole point is retention: a small substring of a large string keeps the large
one alive, because they share a backing array. `types/03-zero-copy-strings`
grades the conversions that share memory deliberately; this is the function that
exists to stop sharing it, and the gradeable question is when you need it — which
is a `weak`/`gc` question about what keeps what alive.

## reflect.MapIter.Reset enables allocation-free map iteration
→ candidate | reflect
"`MapIter.Reset` changes its receiver to iterate over a different map. The use of
`MapIter.Reset` allows allocation-free iteration over many maps." One iterator
reused across many maps rather than one per map — measurable, and a nice example
of an API shaped entirely by allocation concerns.

## sync gained TryLock
→ candidate | memmodel
"`Mutex.TryLock`, `RWMutex.TryLock`, and `RWMutex.TryRLock` will acquire the lock
if it is not currently held." The interesting question is what happens-before
edge a *failed* `TryLock` establishes — the answer is none, which is a genuinely
easy thing to get wrong when writing lock-free fast paths.

## The CPU profiler moved to per-thread timers on Linux
→ candidate | sched
"This increases the maximum CPU usage that a profile can observe, and reduces
some forms of bias." Before this, a process-wide timer meant a program using many
cores could not be profiled accurately above a ceiling. A profiler whose own
resolution is the limit is a good thing to have measured rather than assumed.

## declared but not used inside a function literal
→ no action
"The Go 1.18 compiler now correctly reports `declared but not used` errors for
variables that are set inside a function literal but are never used. Before Go
1.18, the compiler did not report an error in such cases." A real behaviour
change and it breaks compilation of previously valid programs, but it is a
compile error rather than a mechanism to reason about.

## Compile speed regressed by roughly 15%
→ no action
"Because of changes in the compiler related to supporting generics, the Go 1.18
compile speed can be roughly 15% slower than the Go 1.17 compile speed. The
execution time of the compiled code is not affected." Worth knowing as history —
1.20 recovered it — but it is a fact about the toolchain.

## Fuzzing
→ no action
A significant addition, and the warning is practical ("may occupy a large amount
of storage (possibly several GBs)"), but a testing facility rather than something
underneath Go.

## The linker emits far fewer relocations
→ no action
"most codebases will link faster, require less memory to link, and generate
smaller binaries." Binary size is exactly the quantity this repository has
already found unmeasurable to a gradeable standard.
