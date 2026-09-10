# Roadmap

## Tracks

Shipped: `layout`, `alloc`, `reflect` (one exemplar task each).

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
| `review` | review | imported from loupe |
| `design` | design | imported from keystone |
| `concurrency` | build | imported from go-concurrency |

## Candidate pool

Each entry records where the idea came from. Ideas may be borrowed; code may not.

| Candidate | Track | Source |
|---|---|---|
| Small allocations got ~30% cheaper via size-specialized routines — measure it | `alloc` | Go 1.27 release notes |
| Generic methods: what does the compiler emit? | `generics` | Go 1.27 release notes |
| `goroutineleak` profile: plant a leak and find it | `sched` | Go 1.27 release notes |
| Green Tea GC vs the old collector on the same workload | `gc` | Go 1.26 release notes |
| cgo call cost before and after the 1.26 improvement | `edges` | Go 1.26 release notes |
| `simd/archsimd` vs hand-written Plan9 vs pure Go | `asm` | Go 1.26 release notes |
| Heap base randomization: what stops being reproducible? | `edges` | Go 1.26 release notes |

## Release radar

Every Go release gets a triage pass answering two questions: what does it
create, and what does it invalidate. The radar also runs backwards over Go's
whole history, because a behaviour that *changed* is the best kind of task —
"this was true in 1.21; is it still?"

Sweep order, most actionable first:

| Era | Versions |
|---|---|
| E1 | 1.23 → current |
| E2 | 1.18 – 1.22 |
| E3 | 1.10 – 1.17 |
| E4 | 1.0 – 1.9 |

Output lands in `radar/versions/go1.NN.md`, each entry tagged `→ candidate`,
`→ invalidates <task-id>`, or `→ no action`.

## Not yet

Contests — scored challenges, leaderboards and seasons — are deliberately out of
scope until the catalogue is large enough to be worth competing over.
