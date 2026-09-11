# Roadmap

## Tracks

Shipped: `layout`, `alloc`, `reflect` (one exemplar each), the 15 `review/*` categories
(135 drills, imported from loupe), and `design` (36 katas, imported from keystone).

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

| Candidate | Track | Source |
|---|---|---|
| Size-specialized allocation routines make small allocations (<80 bytes) up to 30% cheaper — A/B it with `GOEXPERIMENT=nosizespecializedmalloc` | `alloc` | [Go 1.27](https://go.dev/doc/go1.27) |
| Slice backing stores are stack-allocated in more cases — predict which `make` calls reach the heap, with `-d=variablemakehash=n` as the control | `alloc` | [Go 1.25](https://go.dev/doc/go1.25), extended in [1.26](https://go.dev/doc/go1.26) |
| The experimental portable `simd` package against a pure-Go baseline | `asm` | [Go 1.27](https://go.dev/doc/go1.27) |
| PGO build overhead collapsed, and PGO now aligns hot loop blocks for 1–1.5% | `compiler` | [Go 1.23](https://go.dev/doc/go1.23) |
| The compiler overlaps stack slots of locals with disjoint live ranges — predict a frame size, then move one line | `compiler` | [Go 1.23](https://go.dev/doc/go1.23) |
| Closures may now share a code pointer, so comparing function pointers misleads in more cases | `edges` | [Go 1.27](https://go.dev/doc/go1.27) |
| cgo call overhead is down about 30% — measure the boundary | `edges` | [Go 1.26](https://go.dev/doc/go1.26) |
| The heap base address is randomized on 64-bit: which observations stop being reproducible, and which never were | `edges` | [Go 1.26](https://go.dev/doc/go1.26) |
| `GOAMD64=v3` fused multiply-add changes the exact floating-point values a program produces | `edges` | [Go 1.25](https://go.dev/doc/go1.25) |
| A Go 1.21 compiler bug delayed nil checks; the same program panics on 1.25 — a toolchain-pair task | `edges` | [Go 1.25](https://go.dev/doc/go1.25) |
| `//go:linkname` to unmarked standard-library symbols is now refused | `edges` | [Go 1.23](https://go.dev/doc/go1.23) |
| Green Tea cuts GC overhead 10–40%, and `GOEXPERIMENT=nogreenteagc` makes it an A/B within one toolchain | `gc` | [Go 1.26](https://go.dev/doc/go1.26) |
| Green Tea's further ~10% on Ice Lake / Zen 4 and newer — an answer that depends on the CPU | `gc` | [Go 1.26](https://go.dev/doc/go1.26) |
| Generic methods: how many instantiations does GC-shape stenciling emit, and how does it differ from a generic function? | `generics` | [Go 1.27](https://go.dev/doc/go1.27) |
| Range-over-function iterators and the `iter` package — what `break`, `return` and `goto` do to the yield contract | `iter` | [Go 1.23](https://go.dev/doc/go1.23) |
| The `goroutineleak` profile is generally available — plant a leak of each shape and make the profile name them | `sched` | [Go 1.27](https://go.dev/doc/go1.27) |
| Timer channels are always unbuffered now that `asynctimerchan` is gone | `sched` | [Go 1.27](https://go.dev/doc/go1.27) |
| `GOMAXPROCS` is container-aware and updates itself as the cgroup quota changes | `sched` | [Go 1.25](https://go.dev/doc/go1.25) |
| `testing/synctest` is GA: virtualized time in a bubble — a task, and the harness that makes `sched` and `memmodel` gradeable at all | `sched` | [Go 1.25](https://go.dev/doc/go1.25) |
| `runtime/trace.FlightRecorder` as the data source for a trace-analyzer task | `sched` | [Go 1.25](https://go.dev/doc/go1.25) |
| A new runtime-internal mutex (`GOEXPERIMENT=nospinbitmutex`) — only worth a task if the spin/park boundary can be made visible | `sched` | [Go 1.24](https://go.dev/doc/go1.24) |
| Timers and tickers are collected without `Stop`, and the behaviour is gated on the `go.mod` go line rather than the toolchain | `sched` | [Go 1.23](https://go.dev/doc/go1.23) |
| Windows timer resolution went from 15.6ms to 0.5ms — an answer that depends on the OS | `sched` | [Go 1.23](https://go.dev/doc/go1.23) |
| The builtin map is a Swiss table — predict real memory for a key count, A/B with `GOEXPERIMENT=noswissmap` | `types` | [Go 1.24](https://go.dev/doc/go1.24) |
| `sync.Map` is a hash-trie — the contention curve, and what the old implementation was warming up | `types` | [Go 1.24](https://go.dev/doc/go1.24) |
| `runtime.AddCleanup` vs `SetFinalizer` — four enumerated differences, four predictions | `weak` | [Go 1.24](https://go.dev/doc/go1.24) |
| The `weak` package — the floor for the whole track | `weak` | [Go 1.24](https://go.dev/doc/go1.24) |
| `AddCleanup` runs concurrently and `unique` handles reclaim in a single GC cycle — every pre-1.25 measurement is stale | `weak` | [Go 1.25](https://go.dev/doc/go1.25) |
| `GODEBUG=checkfinalizers=1` names the classic finalizer mistakes — plant each one | `weak` | [Go 1.25](https://go.dev/doc/go1.25) |
| The `unique` package: interning where handle comparison reduces to a pointer compare | `weak` | [Go 1.23](https://go.dev/doc/go1.23) |

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
| E1 | 1.23 – 1.27 | **done** — 54 entries across 5 releases: 31 candidates, 1 invalidation |
| E2 | 1.18 – 1.22 | later |
| E3 | 1.10 – 1.17 | later |
| E4 | 1.0 – 1.9 | later |

Output lands in `radar/versions/go1.NN.md`, each entry tagged `→ candidate`,
`→ invalidates <task-id>`, or `→ no action`. `undergo radar-check` runs in CI and
fails if an invalidation names a task that does not exist, so the dataset cannot
rot as the catalogue changes.

## Not yet

Contests — scored challenges, leaderboards and seasons — are deliberately out of
scope until the catalogue is large enough to be worth competing over.
