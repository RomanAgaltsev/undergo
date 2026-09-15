# Go 1.20

Source: https://go.dev/doc/go1.20

## The spec spelled out how structs and arrays are compared
→ candidate | iface
"Struct values are compared one field at a time, considering fields in the order
they appear in the struct type definition, and stopping at the first mismatch.
Array values are compared one element at a time, in increasing index order." And
the reason it matters: this "affects whether certain comparisons must panic".

The obvious reading — that a mismatch in an earlier field short-circuits before
reaching an uncomparable one — is **wrong for interface comparison**, and it was
measured to be sure. Comparing two `any` values holding
`struct{ A int; S []int }` panics whether or not the `A` fields differ, because
the runtime resolves the dynamic type's equality function *before* comparing any
field, and a type containing a slice has none.

So there are two distinct mechanisms with the same syntax: field-at-a-time
short-circuiting applies to types the compiler already accepts as comparable,
while an interface holding an uncomparable type fails earlier and unconditionally.
`iface/05-interface-comparison` grades the second; this wording is what makes the
distinction precise, and a task could ask for both.

## Linker symbol prefixes changed from go. to go:
→ candidate | compiler
"The linker now uses `go:` and `type:` prefixes for compiler-generated symbols
rather than `go.` and `type.` ... This avoids confusion for user packages whose
name starts with `go.`."

Every symbol-counting task in this repository depends on the spelling of a
symbol name: `generics/01-instantiation-count` and `reflect/05-methodbyname-linker`
both match against `go tool nm` output. A matcher written against a pre-1.20
binary finds nothing in a post-1.20 one, silently — the same failure mode this
repository has now met twice with import-path-qualified names. `debug/gosym`
handles both spellings; a hand-written matcher does not.

## unsafe.SliceData, String and StringData completed the set
→ candidate | edges
Together with 1.17's `unsafe.Slice`, these "provide the complete ability to
construct and deconstruct slice and string values, without depending on their
exact representation." This is what deprecated `reflect.SliceHeader` and
`reflect.StringHeader`, and it is the supported spelling that
`edges/05-unsafe-pointer-rules` and `types/03-zero-copy-strings` both use. The
gradeable question is what "without depending on their exact representation"
buys — the old header structs encoded the layout in the type system, and these
do not.

## Slice to array conversion
→ candidate | types
"`[4]byte(x)` now works" where 1.17 allowed only the pointer form. The
interesting half is the failure condition: the conversion panics if the slice is
shorter than the array, which makes it a run-time check in the middle of what
looks like a static conversion. Predict which of a set of conversions panic.

## comparable may be satisfied by types that are not strictly comparable
→ candidate | generics
"Comparable types (such as ordinary interfaces) may now satisfy `comparable`
constraints, even if the type arguments are not strictly comparable (comparison
may panic at runtime)."

This is the exact seam between the two mechanisms above. Before 1.20 the
constraint meant "comparison cannot panic"; after, it means "comparison is
permitted, and may panic". So a generic map keyed by a type parameter can now be
instantiated with an interface and blow up at run time on insert — which is
`iface/05`'s failure mode reached through the type system rather than through
`any`.

## cgo is disabled by default when there is no C toolchain
→ candidate | edges
"when the `CGO_ENABLED` environment variable is unset, the `CC` environment
variable is unset, and the default C compiler ... is not found in the path,
`CGO_ENABLED` defaults to `0`."

This repository's standing decision records that cgo tasks cannot ship because
`CGO_ENABLED=1` is set on the authoring machine, so a `//go:build cgo` file is
selected and then fails to build. That decision and this default are the same
mechanism seen from two sides, and the gradeable question is what actually
selects the `cgo` build tag — the answer is the variable, not the compiler.

## reflect gained Comparable and Equal
→ candidate | reflect
"`Value.Comparable` and `Value.Equal` can be used to compare two `Value`s for
equality ... `Comparable` reports whether `Equal` is a valid operation for a
given `Value` receiver."

This is the run-time comparability check that `iface/05` names as one of the two
ways to make an `any`-keyed cache safe, and `reflect/04-deepequal-cycles`
implements a hand-rolled version of the comparison. Asking a solver to predict
`Comparable` for a set of values is the same question as "which of these
comparisons panic", asked without the panic.

## PGO arrived as a preview, inlining only
→ candidate | compiler
"Go 1.20 uses PGO to more aggressively inline functions at hot call sites ...
enabling profile-guided inlining optimization improves performance about 3–4%."
Devirtualization came in 1.21 and the interleaving in 1.22, so this is the first
step of a three-release arc. `compiler/04-pgo-devirtualization` sits at the end
of it; the version-diff question is which optimisation a given profile enables at
each point along the way.

## sync.Map gained atomic update operations
→ candidate | memmodel
"The new `Swap`, `CompareAndSwap`, and `CompareAndDelete` methods allow existing
map entries to be updated atomically." Compare-and-swap on a map entry is a
different synchronisation primitive from the load-then-store that `sync.Map`
previously forced, and `memmodel/05-publication` grades exactly the difference
between those two shapes.

## The garbage collector's internal structures got 2% cheaper
→ candidate | gc
"reduces memory overheads and improves overall CPU performance by up to 2%", and
the collector "behaves less erratically with respect to goroutine assists in some
circumstances". The assist behaviour is the gradeable half: an assist is work
charged to an allocating goroutine, so "less erratic" means the charge became
more predictable, which is measurable through `runtime/metrics`.

## New metrics for GOMAXPROCS, cgo calls and mutex wait
→ candidate | sched
`/sched/gomaxprocs:threads`, `/cgo/go-to-c-calls:calls` and
`/sync/mutex/wait/total:seconds`. The cgo counter is the one worth a task: it
counts crossings of a boundary whose cost 1.21 then reduced tenfold, so the two
together make the boundary measurable from inside Go.

## Time histograms became less precise to save memory
→ no action
"Time-based histogram metrics are now less precise, but take up much less
memory." A real tradeoff, and worth knowing before building anything on those
histograms, but the change is in the metric's own resolution rather than in
program behaviour.

## math/rand seeds itself
→ no action
"automatically seeds the global random number generator ... with a random value",
with `GODEBUG=randautoseed=0` to restore the old behaviour. The 1.22 entry covers
where this ended up; this is the first step and adds nothing a solver could
measure that 1.22 does not.

## errors.Join and multiple %w
→ no action
An error may now "wrap more than one error by providing an `Unwrap` method that
returns a `[]error`", and `fmt.Errorf` accepts several `%w` verbs. Genuinely
useful and genuinely a mechanism — `errors.Is` now walks a tree rather than a
chain — but it is library behaviour with no machine-level question underneath.

## Build speed improved by up to 10%
→ no action
"Go 1.18 and 1.19 saw regressions in build speed, largely due to the addition of
support for generics ... Go 1.20 improves build speeds by up to 10%, bringing it
back in line with Go 1.17." A fact about the toolchain, not about any program it
produces.

## The race detector on macOS no longer needs cgo
→ no action
"race-detector-enabled programs can be built and run without Xcode", while on
Linux, Unix and Windows "a host C toolchain is required to use the race
detector." That asymmetry is exactly why this repository's race gate runs in a
container, but it is a toolchain property rather than something a task can ask
about.
