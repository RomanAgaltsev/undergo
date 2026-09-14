# Roadmap

## Tracks

Shipped: all fourteen internals tracks. Memory — `layout` (3 tasks), `alloc` (4),
`types` (4). Language surface — `generics` (5), `iter` (3), `reflect` (3). Lifetime
and collection — `weak` (4), `gc` (4). The machine — `iface` (4), `compiler` (4),
`asm` (2), `edges` (4). Alongside the 15 `review/*` categories (135 drills, imported
from loupe) and `design` (36 katas, imported from keystone). 215 tasks in all, 44 of
them machine-graded.

Planned, in rough order:

| Track | Mode | What it covers |
|---|---|---|
| `layout` | predict | size, alignment, offsets, padding, false sharing |
| `alloc` | predict | allocation counts, escape analysis, interface boxing, stack growth |
| `types` | predict/build | slice growth, aliasing, `unsafe.String`/`Slice`, map memory |
| `iface` | predict | eface vs iface, nil-interface traps, devirtualization |
| `reflect` | build | tag decoding, settability, `StructOf`, reflect vs generics vs codegen |
| `generics` | build/predict | type sets, inference limits, GC-shape stenciling, dictionaries, generic methods |
| `iter` | build | `iter.Seq`/`Seq2` adapters, `iter.Pull`, the yield contract |
| `weak` | build/predict | `weak.Pointer`, `AddCleanup` vs `SetFinalizer`, `unique.Handle` |
| `gc` | predict | GC cycles, GOGC × GOMEMLIMIT, write barriers, Green Tea |
| `sched` | predict/build | interleaving, async preemption, `LockOSThread`, timers, traces |
| `memmodel` | predict | happens-before, store buffering, seqlocks, the benign-race myth |
| `compiler` | predict/optimize | inlining budget, bounds-check elimination, PGO, binary size |
| `asm` | build/optimize | Plan9 syntax, register ABI, `//go:noescape`, SIMD |
| `edges` | mixed | cgo cost, `defer` tiers, panic/recover, `unsafe.Pointer` rules |
| `review/*` | review | 15 categories × 3 tiers × 3 drills — concurrency, nil-safety, error-handling, context, resource-leaks, api-design, performance, security, correctness, testing, generics, json, time, http-client, typed-nil |
| `design` | design | 36 system-design katas across 8 tracks |
| `concurrency` | build | imported from go-concurrency |

## Candidate pool

Each entry records where the idea came from. Ideas may be borrowed; code may not.

Every row below is derived from a `→ candidate` entry in `radar/versions/`. A row
with no entry behind it is drift, and gets deleted rather than kept out of
politeness.

| Candidate | Track | Source | Built |
|---|---|---|---|
| Size-specialized allocation routines make small allocations (<80 bytes) up to 30% cheaper — A/B it with `GOEXPERIMENT=nosizespecializedmalloc` | `alloc` | [Go 1.27](https://go.dev/doc/go1.27) | — |
| Slice backing stores are stack-allocated in more cases — predict which `make` calls reach the heap, with `-d=variablemakehash=n` as the control | `alloc` | [Go 1.25](https://go.dev/doc/go1.25), extended in [1.26](https://go.dev/doc/go1.26) | `alloc/03-what-escapes`, `types/01-cap-growth` |
| The experimental portable `simd` package against a pure-Go baseline | `asm` | [Go 1.27](https://go.dev/doc/go1.27) | — |
| PGO build overhead collapsed, and PGO now aligns hot loop blocks for 1–1.5% | `compiler` | [Go 1.23](https://go.dev/doc/go1.23) | — |
| The compiler overlaps stack slots of locals with disjoint live ranges — predict a frame size, then move one line | `compiler` | [Go 1.23](https://go.dev/doc/go1.23) | `compiler/03-stack-slots` — which measured that it does NOT happen for address-taken locals |
| Closures may now share a code pointer, so comparing function pointers misleads in more cases | `edges` | [Go 1.27](https://go.dev/doc/go1.27) | — |
| cgo call overhead is down about 30% — measure the boundary | `edges` | [Go 1.26](https://go.dev/doc/go1.26) | — |
| The heap base address is randomized on 64-bit: which observations stop being reproducible, and which never were | `edges` | [Go 1.26](https://go.dev/doc/go1.26) | — |
| `GOAMD64=v3` fused multiply-add changes the exact floating-point values a program produces | `edges` | [Go 1.25](https://go.dev/doc/go1.25) | — |
| A Go 1.21 compiler bug delayed nil checks; the same program panics on 1.25 — a toolchain-pair task | `edges` | [Go 1.25](https://go.dev/doc/go1.25) | — |
| `//go:linkname` to unmarked standard-library symbols is now refused | `edges` | [Go 1.23](https://go.dev/doc/go1.23) | `edges/01-linkname` |
| A struct literal key may be any valid field selector, not just a top-level field name | `edges` | [Go 1.27](https://go.dev/doc/go1.27) | — |
| Green Tea cuts GC overhead 10–40%, and `GOEXPERIMENT=nogreenteagc` makes it an A/B within one toolchain | `gc` | [Go 1.26](https://go.dev/doc/go1.26) | — |
| Green Tea's further ~10% on Ice Lake / Zen 4 and newer — an answer that depends on the CPU | `gc` | [Go 1.26](https://go.dev/doc/go1.26) | — |
| Generic methods: how many instantiations does GC-shape stenciling emit, and how does it differ from a generic function? | `generics` | [Go 1.27](https://go.dev/doc/go1.27) | `generics/01-instantiation-count`, `generics/02-generic-methods` |
| Generic type aliases are fully supported — does a parameterized alias add an instantiation, or share one? | `generics` | [Go 1.24](https://go.dev/doc/go1.24) | — |
| `new` accepts an expression, so `new(f(x))` compiles — a version-diff task a pre-1.26 model gets wrong | `generics` | [Go 1.26](https://go.dev/doc/go1.26) | `generics/04-new-and-self-reference` |
| A generic type may refer to itself in its own type parameter list (`type Adder[A Adder[A]]`) | `generics` | [Go 1.26](https://go.dev/doc/go1.26) | `generics/04-new-and-self-reference` |
| Function type inference generalized to assignment and conversion contexts — not gated by the go.mod language version, so a true A/B needs an older toolchain | `generics` | [Go 1.27](https://go.dev/doc/go1.27) | `generics/05-inference-limits` — the assignment shape only; no 1.26-vs-1.27 A/B, see the radar note |
| Range-over-function iterators and the `iter` package — what `break`, `return` and `goto` do to the yield contract | `iter` | [Go 1.23](https://go.dev/doc/go1.23) | `iter/02-yield-contract` and `iter/01-adapters` — `goto` out of a range body is still uncovered |
| The `goroutineleak` profile is generally available — plant a leak of each shape and make the profile name them | `sched` | [Go 1.27](https://go.dev/doc/go1.27) | — |
| Timer channels are always unbuffered now that `asynctimerchan` is gone | `sched` | [Go 1.27](https://go.dev/doc/go1.27) | — |
| `GOMAXPROCS` is container-aware and updates itself as the cgroup quota changes | `sched` | [Go 1.25](https://go.dev/doc/go1.25) | — |
| `testing/synctest` is GA: virtualized time in a bubble — a task, and the harness that makes `sched` and `memmodel` gradeable at all | `sched` | [Go 1.25](https://go.dev/doc/go1.25) | — |
| `runtime/trace.FlightRecorder` as the data source for a trace-analyzer task | `sched` | [Go 1.25](https://go.dev/doc/go1.25) | — |
| A new runtime-internal mutex (`GOEXPERIMENT=nospinbitmutex`) — only worth a task if the spin/park boundary can be made visible | `sched` | [Go 1.24](https://go.dev/doc/go1.24) | — |
| Timers and tickers are collected without `Stop`, and the behaviour is gated on the `go.mod` go line rather than the toolchain | `sched` | [Go 1.23](https://go.dev/doc/go1.23) | — |
| Windows timer resolution went from 15.6ms to 0.5ms — an answer that depends on the OS | `sched` | [Go 1.23](https://go.dev/doc/go1.23) | — |
| The builtin map is a Swiss table — predict real memory for a key count, A/B with `GOEXPERIMENT=noswissmap` | `types` | [Go 1.24](https://go.dev/doc/go1.24) | `types/04-map-memory` — memory half only; the `noswissmap` A/B is still open |
| `sync.Map` is a hash-trie — the contention curve, and what the old implementation was warming up | `types` | [Go 1.24](https://go.dev/doc/go1.24) | — |
| `runtime.AddCleanup` vs `SetFinalizer` — four enumerated differences, four predictions | `weak` | [Go 1.24](https://go.dev/doc/go1.24) | `weak/02-cleanup-vs-finalizer` |
| The `weak` package — the floor for the whole track | `weak` | [Go 1.24](https://go.dev/doc/go1.24) | `weak/01-weak-cache` |
| `AddCleanup` runs concurrently and `unique` handles reclaim in a single GC cycle — every pre-1.25 measurement is stale | `weak` | [Go 1.25](https://go.dev/doc/go1.25) | `weak/02-cleanup-vs-finalizer` — the concurrency half; the single-cycle reclaim claim is untested |
| `GODEBUG=checkfinalizers=1` names the classic finalizer mistakes — plant each one | `weak` | [Go 1.25](https://go.dev/doc/go1.25) | `weak/03-lifetime-traps` |
| The `unique` package: interning where handle comparison reduces to a pointer compare | `weak` | [Go 1.23](https://go.dev/doc/go1.23) | `weak/04-unique-interning` |

## Invalidated by a release

| Task | Release | What changed | Status |
|---|---|---|---|
| `layout/01-struct-padding` | [Go 1.23](https://go.dev/doc/go1.23) | `structs.HostLayout` exists because "struct layout order is not guaranteed by the language spec". The sealed hint said "Go lays a struct out in declaration order", stating a gc implementation detail as a language guarantee. The five measured answers were unaffected. | **fixed 2026-09-11** — hint corrected, explanation gained a section on spec guarantees vs implementation, task re-sealed |

## Release radar

Every Go release gets a triage pass answering two questions: what does it
create, and what does it invalidate. The radar also runs backwards over Go's
whole history, because a behaviour that *changed* is the best kind of task —
"this was true in 1.21; is it still?"

Sweep order, most actionable first:

| Era | Versions | Status |
|---|---|---|
| E1 | 1.23 – 1.27 | **done** — 60 entries across 5 releases: 36 candidates, 1 invalidation. Swept twice: the first pass (2026-09-11) read the runtime, toolchain and library sections and caught three language changes but missed six; the second (2026-09-13) read every release's "Changes to the language" section on its own. |
| E2 | 1.18 – 1.22 | later |
| E3 | 1.10 – 1.17 | later |
| E4 | 1.0 – 1.9 | later |

Output lands in `radar/versions/go1.NN.md`, each entry tagged `→ candidate`,
`→ invalidates <task-id>`, or `→ no action`. `undergo radar-check` runs in CI and
fails if an invalidation names a task that does not exist, so the dataset cannot
rot as the catalogue changes.

## Standing decisions

**Authoring comes before solving, deliberately.** 215 tasks ship and none have been
solved here. That is a choice, not a backlog: the catalogue is being built out
first, and the solver's experience — whether a hint is one rung or two, whether a
question is answerable without its explanation — is unvalidated until someone
works through a track. Gate 2 proves every reference solution compiles and passes;
it proves nothing about whether a task teaches. Revisit when a track is picked up
in a solver's clone.

**cgo tasks cannot ship while the authoring machine has no C toolchain.** Spec §8
T14's "cgo boundary cost" is deferred, and not merely for want of a compiler: the
`cgo` build tag is satisfied by `CGO_ENABLED=1`, which is set here, rather than by
a compiler existing. So a `//go:build cgo` file with a `//go:build !cgo` fallback
does **not** degrade gracefully — the cgo file is selected and the build fails with
`C compiler "gcc" not found`, breaking `task ci` locally rather than only in CI.
Revisit when the machine has a toolchain, or build the task as a `testdata`
snippet compiled with `CGO_ENABLED=0` explicitly.

**Assembly tasks ship a portable fallback and pin `requires.arch`.** Gate 1 builds
every task on every platform and CI includes an arm64 macOS runner, where an
amd64-only `.s` fails with "missing function body" — `requires.arch` gates gate 2,
not gate 1. The `//go:build amd64` declaration plus a `//go:build !amd64` Go
implementation builds and vets on windows/amd64, linux/arm64 and darwin/arm64.
The cost is that a green run on arm64 exercised the fallback, so an asm task
should make which implementation ran visible in its output.

**An invariant is not automatically robust — measure its margin.** `gc/02` shipped
a four-way GOGC ordering in M7 that measured 8/8 stable locally, 8/8 under the race
detector, and green on both CI platforms. It then failed once during a heavily
loaded build, because a full ordering is a conjunction of adjacent comparisons and
its weakest link was 3 collections against 1. It now asks about pairs an order of
magnitude apart. Reaching for an invariant instead of a tolerance was right; not
checking how much room the invariant had was not.

**Noisy measurements get an invariant, not a tolerance.** GC cycle counts vary 39–68
across identical runs, and `GOMEMLIMIT`-clamped heap goals give four distinct values in
eight. Rather than widen the answer with a range, `gc/02` and `gc/03` ask for an
ordering and a predicate, which measured 6/6 and 8/8 stable and need no harness change.
`sched` and `memmodel` face the same problem in M10 and should reach for the same tool
before proposing tolerance slots.

The two Green Tea candidates stay open for this reason: a 10–40% overhead reduction is
a distribution, not a value, and nobody has yet found the invariant underneath it.

**The race gate covers the solutions, not just the machinery.** `task race` runs
the race detector over `./cmd/... ./internal/...` *and* over every reference
solution via `ci-verify --race`. It is a required check. It is deliberately not
part of `task ci`, because the detector needs cgo and a Windows checkout without a
C toolchain cannot run it — `undergo doctor` reports which of the native and
container routes this machine has, and `task race:docker` is the way round a
missing compiler.

## Not yet

Contests — scored challenges, leaderboards and seasons — are deliberately out of
scope until the catalogue is large enough to be worth competing over.
