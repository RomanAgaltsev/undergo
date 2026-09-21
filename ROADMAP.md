# Roadmap

## Tracks

Shipped: all sixteen internals tracks, **every one of them at five tasks or more**.
Memory — `layout` (5 tasks), `alloc` (5),
`types` (5). Language surface — `generics` (5), `iter` (5), `reflect` (5). Lifetime
and collection — `weak` (5), `gc` (5). The machine — `iface` (5), `compiler` (5),
`asm` (5), `edges` (10). Scheduling and memory — `sched` (5), `memmodel` (5),
`concurrency` (8). Versions — `versions` (13). Alongside the 15 `review/*`
categories (135 drills, imported
from loupe) and `design` (36 katas, imported from keystone). 267 tasks in all, 96 of them machine-graded.

These counts are **not** checked by anything. `undergo validate` compares the
catalogue against `README.md` only, so this paragraph drifted from M16 until
somebody happened to read it. Treat it with suspicion, or give it a check.

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
| `compiler` | predict/optimize | inlining budget, bounds-check elimination, PGO, loop lowering |
| `asm` | build/optimize | Plan9 syntax, register ABI, `//go:noescape`, SIMD |
| `edges` | mixed | cgo cost, `defer` tiers, panic/recover, `unsafe.Pointer` rules |
| `review/*` | review | 15 categories × 3 tiers × 3 drills — concurrency, nil-safety, error-handling, context, resource-leaks, api-design, performance, security, correctness, testing, generics, json, time, http-client, typed-nil |
| `design` | design | 36 system-design katas across 8 tracks |
| `concurrency` | predict/build | channel handoff, abandoned results, `Cond`, the `Once` contract, `Pool` clearing, `RWMutex` admission, cancellability, close cascades — derived from `go-concurrency` |
| `versions` | predict | the `go.mod` line as a behaviour switch: GODEBUG defaults, loop variables, language legality, and the edge where a removed switch stops answering; and its second half, where the toolchain changes instead and reaches what no go line can |

## Candidate pool

Each entry records where the idea came from.

Every row below records a source: a `→ candidate` entry in `radar/versions/`, or
a named reference. A row with no source behind it is drift, and gets deleted
rather than kept out of politeness. Ideas may be borrowed; code may not.

The historical sweep (M12) added 158 rows at once. That is a pool, not a
backlog: a row is a lead with a citation, and most will never become tasks.

193 of these rows come from the radar, against its 194 candidates. The one
difference is deliberate: "slice backing stores are stack-allocated in more
cases" is a `→ candidate` in both 1.25 and 1.26 because the change was extended,
and it is one lead, so it gets one row. Any other gap between the radar and this
table is a defect in one direction or the other.

The remaining three rows are **kata-sourced**, not radar-sourced — M9 triaged
the fourteen `go-concurrency` katas and these three could not be built without
adding a module dependency. The pool's rule is "records a source", and a named
kata is one; but the count above is about the radar, so it does not move.

The row count is now stable, and the **`Built` column** is the thing that moves.
M14 claimed three rows without adding any: a lead may be built by more than one
task, and `versions/01` and `versions/05` each join a task that was already
there — one grading the current behaviour, the other grading the version it
changed in. A row is retired by being built, never by being deleted.

### What the pool cannot supply

The row count is not a backlog and it is not a budget. M15 found that most rows
cannot become tasks as written, and M16 measured how many: **97 rows were
unbuilt, and they yielded five tasks.** Three rules do most of that cutting, and
none of them was written down before now.

**The Go 1.21 floor.** `GOTOOLCHAIN` is itself a Go 1.21 feature, so
`internal/toolchain` cannot name anything older — `minToolchain` enforces it. A
version *pair* needs a "before" as well as an "after", so the oldest usable
"before" is 1.21 and the oldest gradeable *change* is one that landed in
**1.22**. Measured at v0.19.0, every unbuilt row cited exactly one version:

| band | unbuilt rows |
|---|---|
| ≤ Go 1.20 | **68** |
| Go 1.21 | 6 — reachable only through a GODEBUG knob |
| ≥ Go 1.22 | 23 |

Seventy percent of the unbuilt pool is out of reach of a toolchain pair, and no
amount of effort changes that. M14b learned it the expensive way: the spec named
the Go 1.13 escape-analysis rewrite as its best candidate, and the candidate was
unreachable.

**The answer form.** A duration, a rate or a ratio cannot be graded — it is a
fact about the machine that ran it. This kills more leads than anything except
the floor, because the release notes are full of performance claims: the two
Green Tea rows, cgo overhead, size-specialized allocation, `sync.Map`'s
contention curve. Some survive by being re-asked: "stop-the-world pauses split
in two" is a duration, but *which metric names exist* is an exact string, and
that became `versions/11`.

**The platform rule, which is not what it looks like.** Gate 2 proves every task
on ubuntu and macos, so it is tempting to say a task must answer the same on
both. That is wrong, and `versions/13` proved it: `0xc000` is the arena hint
from the `default` arm of `runtime/malloc.go`, and `GOARCH == "arm64"` takes its
own, so the task failed on macos and passed on ubuntu.

The rule this repo actually keeps is **one correct answer on every platform a
task claims**. Eleven tasks already pin `requires.arch` or `requires.os` — the
whole `asm` track among them — and `Gradeable` skips them off-platform by
design. So a platform-specific lead is **pinned, not struck**. What disqualifies
a lead is having no single correct answer on a platform it claims, not having a
different answer on a platform it never claimed.

### Three genres, and only one of them downloads anything

The `versions` track is not "the toolchain-pair track". Choosing the wrong genre
makes a task more expensive and less portable than it needs to be.

| genre | tasks | mechanism | downloads? |
|---|---|---|---|
| GODEBUG / `-lang` gate | `versions/01`–`05`, `09` | `internal/goline`, the go.mod `go` line | **no** |
| toolchain pair | `versions/06`–`08`, `10`–`13` | `internal/toolchain` | yes |
| GOEXPERIMENT A/B | none yet | one toolchain, two environments | no |

Since Go 1.21 the go line is the compatibility switch, so genre 1 reaches any
change shipped with a GODEBUG knob and any *language legality* change at any
version — including changes far below the 1.21 floor, because `-lang` has
accepted old language versions for years. Reach for it first.

Genre 3 belongs to a lead's home track rather than to `versions`: a
`GOEXPERIMENT` is a property of one toolchain, not a difference between two.

The cost of a pair task is **per distinct toolchain, not per task**. M16 added
five tasks for two new downloads, because two of them reused toolchains already
required, and gate 2 grew by 22s on ubuntu and 31s on macos — about what M14b's
three tasks cost. A candidate that reuses a toolchain already named by some task
is materially cheaper than an equally good one that does not.

### Counting rows overstates the leads

At least one lead appears twice, under `alloc` and under `layout`. Both are
struck below for the same reason, and they had to be struck together. Any future
count should say whether it is counting rows or leads.

| Candidate | Track | Source | Built |
|---|---|---|---|
| Size-specialized allocation routines make small allocations (<80 bytes) up to 30% cheaper — A/B it with `GOEXPERIMENT=nosizespecializedmalloc` | `alloc` | [Go 1.27](https://go.dev/doc/go1.27) | — |
| Slice backing stores are stack-allocated in more cases — predict which `make` calls reach the heap, with `-d=variablemakehash=n` as the control | `alloc` | [Go 1.25](https://go.dev/doc/go1.25), extended in [1.26](https://go.dev/doc/go1.26) | `alloc/03-what-escapes`, `types/01-cap-growth` |
| Heap metadata moved next to the object and allocation alignment fell from 16 bytes to 8 — which size classes changed, and by how much. **Probed thin 2026-09-21:** effective bytes per allocation are identical under 1.21 and 1.22 at every size from 1 to 128, so the answer as asked is "none of them". Not visible through `TotalAlloc`. Same lead as the `layout` row below; strike both or neither | `alloc` | [Go 1.22](https://go.dev/doc/go1.22) | — |
| A goroutine's starting stack is sized from the average of its predecessors, not from a constant — so a function's first allocation depends on program history | `alloc` | [Go 1.19](https://go.dev/doc/go1.19) | — |
| `strings.Trim` and friends became allocation-free for the no-op case — predict `AllocsPerRun` before and after | `alloc` | [Go 1.18](https://go.dev/doc/go1.18) | — |
| `os.File.WriteString` stopped copying to a `[]byte` — one allocation removed by a compiler-level conversion rule | `alloc` | [Go 1.17](https://go.dev/doc/go1.17) | — |
| Converting a small integer to an interface stopped allocating — where is the boundary, and what is the staticuint64s table | `alloc` | [Go 1.15](https://go.dev/doc/go1.15) | — |
| Small-object allocation stopped degrading at high core counts — an answer that changes shape with `GOMAXPROCS` | `alloc` | [Go 1.15](https://go.dev/doc/go1.15) | — |
| The page allocator stopped contending: large parallel allocations cross the 32 KiB boundary into a different allocator | `alloc` | [Go 1.14](https://go.dev/doc/go1.14) | `alloc/05-size-classes` measures the boundary |
| `sync.Pool` retains objects across one GC via a victim cache — put, force one collection, `Get`; then force two | `alloc` | [Go 1.13](https://go.dev/doc/go1.13) | `concurrency/05-pool-clearing` |
| The 1.13 escape-analysis rewrite changed which values reach the heap wholesale — a genuine two-toolchain version diff | `alloc` | [Go 1.13](https://go.dev/doc/go1.13) | — |
| Memory profiles stopped overcounting large heap allocations — every pre-1.12 profile of an above-size-class allocation was wrong | `alloc` | [Go 1.12](https://go.dev/doc/go1.12) | — |
| The `allocs` profile: bytes allocated ever against bytes live now, and why both are true of one program | `alloc` | [Go 1.11](https://go.dev/doc/go1.11) | `alloc/05-size-classes` uses the quantity |
| The heap became sparse and lost its 512 GiB ceiling, fixing address-space conflicts under `-race` and cgo | `alloc` | [Go 1.11](https://go.dev/doc/go1.11) | — |
| `strings.Builder` avoids `bytes.Buffer`'s copy in `String()` because its API forbids writing afterwards — the restriction *is* the optimisation | `alloc` | [Go 1.10](https://go.dev/doc/go1.10) | `alloc/01-zero-alloc-join` |
| Contiguous growable stacks replaced segmented ones, and the starting size fell to 2048 bytes — where the "2 KB goroutine" number comes from | `alloc` | [Go 1.3](https://go.dev/doc/go1.3), [1.4](https://go.dev/doc/go1.4) | — |
| The experimental portable `simd` package against a pure-Go baseline | `asm` | [Go 1.27](https://go.dev/doc/go1.27) | — |
| Frameless `NOSPLIT` assembly is no longer automatically `NOFRAME` on amd64 — what a frame costs and what it buys | `asm` | [Go 1.21](https://go.dev/doc/go1.21) | `asm/03-register-abi` |
| The register ABI reached arm64 and ppc64 — the same function, two calling conventions, one `-S` dump | `asm` | [Go 1.18](https://go.dev/doc/go1.18) | `asm/03-register-abi` |
| Stack traces mark register-passed arguments as possibly inaccurate — the debugging cost of the register ABI, stated | `asm` | [Go 1.18](https://go.dev/doc/go1.18) | — |
| arm64 keeps frame pointers in every function where amd64 does not — the price of `perf`-compatible unwinding, paid everywhere | `asm` | [Go 1.17](https://go.dev/doc/go1.17) | — |
| Stack traces print arguments, and print `?` where a register-passed value is unknowable — the register ABI's debugging debt | `asm` | [Go 1.17](https://go.dev/doc/go1.17) | — |
| The register calling convention on amd64, and the ABI0 adapters the transition needed | `asm` | [Go 1.17](https://go.dev/doc/go1.17) | `asm/03-register-abi` |
| `go vet` learned that assembly must preserve BP — a register the Go ABI treats as callee-save and hand-written code forgets | `asm` | [Go 1.16](https://go.dev/doc/go1.16) | — |
| `math/bits` carries a documented constant-time guarantee for `Add`, `Sub`, `Mul`, `RotateLeft` and `ReverseBytes` | `asm` | [Go 1.13](https://go.dev/doc/go1.13) | `asm/05-carry-chain` |
| Frame pointers on linux/arm64 cost about 3% — a stated price for a debugging affordance | `asm` | [Go 1.12](https://go.dev/doc/go1.12) | — |
| The assembler stopped rewriting `MOVL $0, AX` as `XORL`, because the peephole clobbered the condition flags | `asm` | [Go 1.10](https://go.dev/doc/go1.10) | `asm/05-carry-chain` is built on flags surviving |
| `math/bits` functions are intrinsics: the portable Go body ships, compiles, and almost never runs — which ones, on which GOARCH | `asm` | [Go 1.9](https://go.dev/doc/go1.9) | `asm/05-carry-chain` |
| PGO build overhead collapsed, and PGO now aligns hot loop blocks for 1–1.5% | `compiler` | [Go 1.23](https://go.dev/doc/go1.23) | `compiler/04-pgo-devirtualization` |
| The compiler overlaps stack slots of locals with disjoint live ranges — predict a frame size, then move one line | `compiler` | [Go 1.23](https://go.dev/doc/go1.23) | `compiler/03-stack-slots` — which measured that it does NOT happen for address-taken locals |
| An inliner that weighs the call site and not just the callee — a panic path discourages inlining where 1.11 encouraged it. **Genre 3, not a version task:** probed 2026-09-21, it does not flip across a 1.21/1.22 pair because `newinliner` is still an opt-in `GOEXPERIMENT` at Go 1.27. Under it the compiler prints exact inlining scores — the panic path `-9`, a plain call `4`, `fmt.Println` `72` — and `main` becomes inlinable. Gradeable as a GOEXPERIMENT A/B here in `compiler`; `toolchain.Run` already accepts env vars | `compiler` | [Go 1.22](https://go.dev/doc/go1.22) | — |
| PGO devirtualizes more, and devirtualization now interleaves with inlining rather than running before it | `compiler` | [Go 1.22](https://go.dev/doc/go1.22) | `compiler/04-pgo-devirtualization` |
| Linker symbol prefixes changed from `go.` to `go:` — every tool that parsed `-S` output by prefix broke | `compiler` | [Go 1.20](https://go.dev/doc/go1.20) | — |
| PGO arrived as a preview that only inlined — the baseline any PGO measurement is against | `compiler` | [Go 1.20](https://go.dev/doc/go1.20) | `compiler/04-pgo-devirtualization` |
| Large value switches became jump tables — predict the crossover count, and what makes a switch ineligible | `compiler` | [Go 1.19](https://go.dev/doc/go1.19) | — |
| Functions containing range loops became inlinable — which loop shapes still block inlining | `compiler` | [Go 1.18](https://go.dev/doc/go1.18) | `compiler/01-inlining-budget` |
| The inliner learned three new shapes at once in 1.16 — a budget task that is really an eligibility task | `compiler` | [Go 1.16](https://go.dev/doc/go1.16) | `compiler/01-inlining-budget` |
| Bounds-check elimination learned about slice creation and about index types narrower than `int` | `compiler` | [Go 1.14](https://go.dev/doc/go1.14) | `compiler/02-bounds-checks` |
| The compiler can emit inlining, escape, BCE and nil-check decisions as JSON (`-json`) — a structured alternative to parsing `-m` | `compiler` | [Go 1.14](https://go.dev/doc/go1.14) | — |
| Mutex, RWMutex and Once fast paths were *restructured* to fit the inlining budget — a design technique visible in the source | `compiler` | [Go 1.13](https://go.dev/doc/go1.13) | `compiler/01-inlining-budget`, `memmodel/05-publication` |
| Out-of-range panics name the index and the length — information the surviving bounds check must keep live to the point of failure | `compiler` | [Go 1.13](https://go.dev/doc/go1.13) | `compiler/02-bounds-checks` |
| `-trimpath` removes file system paths from the binary — the flag M8 needed when identical source in differently-named directories produced different binary sizes | `compiler` | [Go 1.13](https://go.dev/doc/go1.13) | — |
| Functions that do nothing but call another function became inlinable — and `runtime.Callers` stopped matching the source as a result | `compiler` | [Go 1.12](https://go.dev/doc/go1.12) | — |
| Bounds-check elimination learned transitivity, arithmetic offsets, induction and shift narrowing — four abilities in one release | `compiler` | [Go 1.11](https://go.dev/doc/go1.11) | `compiler/02-bounds-checks` |
| The map-clearing idiom `for k := range m { delete(m, k) }` became one runtime call — and why it cannot be `memclr` | `compiler` | [Go 1.11](https://go.dev/doc/go1.11) | `compiler/05-loop-lowering` grades the slice analogue |
| Functions that call `panic` became inlinable — the small argument-checking wrapper stopped being penalised | `compiler` | [Go 1.11](https://go.dev/doc/go1.11) | — |
| The SSA back end (amd64 in 1.7, every architecture in 1.8) — the platform every later `compiler` finding is written against | `compiler` | [Go 1.7](https://go.dev/doc/go1.7), [1.8](https://go.dev/doc/go1.8) | — |
| `errgroup`'s first error wins and its siblings' errors are silently discarded — the mechanism is an internal `errOnce sync.Once` | `concurrency` | go-concurrency kata 09 | — (costs `golang.org/x/sync`) |
| `singleflight` deletes the key on success **and on error** before `Do` returns — and `DoChan` allocates five times what `Do` does | `concurrency` | go-concurrency kata 12 | — (costs `golang.org/x/sync`) |
| Three rate limiters behind one interface: which can burst, and why the `time.Ticker` one strictly cannot | `concurrency` | go-concurrency kata 14 | — (costs `golang.org/x/time`) |
| Closures may now share a code pointer, so comparing function pointers misleads in more cases | `edges` | [Go 1.27](https://go.dev/doc/go1.27) | `versions/08-closure-code-pointers` |
| A struct literal key may be any valid field selector, not just a top-level field name | `edges` | [Go 1.27](https://go.dev/doc/go1.27) | `versions/09-promoted-field-keys` |
| cgo call overhead is down about 30% — measure the boundary | `edges` | [Go 1.26](https://go.dev/doc/go1.26) | — |
| The heap base address is randomized on 64-bit: which observations stop being reproducible, and which never were | `edges` | [Go 1.26](https://go.dev/doc/go1.26) | `versions/13-heap-base` |
| `GOAMD64=v3` fused multiply-add changes the exact floating-point values a program produces. **Pin candidate:** amd64-only, which is not a strike — see the platform rule above. `requires.arch: [amd64]` is how `versions/13` handles the same problem | `edges` | [Go 1.25](https://go.dev/doc/go1.25) | — |
| A Go 1.21 compiler bug delayed nil checks; the same program panics on 1.25 — a toolchain-pair task | `edges` | [Go 1.25](https://go.dev/doc/go1.25) | — |
| `//go:linkname` to unmarked standard-library symbols is now refused | `edges` | [Go 1.23](https://go.dev/doc/go1.23) | `edges/01-linkname` |
| Loop variables are created anew each iteration — gated on the `go.mod` language version, so one toolchain gives both answers | `edges` | [Go 1.22](https://go.dev/doc/go1.22) | `versions/02-loop-variables` |
| `GODEBUG` and the `go` line became the compatibility mechanism — how a behaviour change ships without breaking the Go 1 promise | `edges` | [Go 1.21](https://go.dev/doc/go1.21) | — |
| `panic(nil)` stopped being nil: `recover()` now returns a `*runtime.PanicNilError`, gated on the go line | `edges` | [Go 1.21](https://go.dev/doc/go1.21) | `edges/03-panic-vs-fatal`, `versions/01-panic-nil` |
| Package initialisation order became a specified algorithm rather than an implementation detail | `edges` | [Go 1.21](https://go.dev/doc/go1.21) | `edges/07-init-order` |
| The cgo call boundary got an order of magnitude cheaper — the number people quote for "cgo is slow" has a date on it | `edges` | [Go 1.21](https://go.dev/doc/go1.21) | — |
| `runtime.Pinner` — pinning an object so C may hold it, and what the collector gives up to allow it | `edges` | [Go 1.21](https://go.dev/doc/go1.21) | — |
| `unsafe.SliceData`, `unsafe.String` and `unsafe.StringData` completed the set — the supported spelling of the zero-copy conversion | `edges` | [Go 1.20](https://go.dev/doc/go1.20) | `types/03-zero-copy-strings` |
| cgo is disabled by default when no C toolchain is present — so whether a program builds depends on the machine, not the source | `edges` | [Go 1.20](https://go.dev/doc/go1.20) | — |
| Fatal error tracebacks got shorter unless you ask (`GOTRACEBACK`) — what a crash is allowed to hide | `edges` | [Go 1.19](https://go.dev/doc/go1.19) | — |
| Importing `os` raises the process's file descriptor limit to the hard maximum — an import with a side effect on the kernel | `edges` | [Go 1.19](https://go.dev/doc/go1.19) | — |
| `GOAMD64` selects a microarchitecture level — one source, four instruction sets, and sometimes different floating-point results | `edges` | [Go 1.18](https://go.dev/doc/go1.18) | — |
| Functions containing closures became inlinable, and one function's code pointers multiplied as a result | `edges` | [Go 1.17](https://go.dev/doc/go1.17) | — |
| `unsafe.Add` and `unsafe.Slice` — the two operations that previously required a `uintptr` round trip | `edges` | [Go 1.17](https://go.dev/doc/go1.17) | `edges/05-unsafe-pointer-rules` |
| `GODEBUG=inittrace=1` prints what every package `init` cost in time and bytes | `edges` | [Go 1.16](https://go.dev/doc/go1.16) | — |
| `go test` now fails a test that calls `os.Exit(0)` mid-run — a passing suite that ran nothing | `edges` | [Go 1.16](https://go.dev/doc/go1.16) | `edges/08-exit-during-test` |
| Chained `unsafe.Pointer`-to-`uintptr` conversions became illegal — the loophole in pattern 3 closed | `edges` | [Go 1.15](https://go.dev/doc/go1.15) | `edges/05-unsafe-pointer-rules` |
| `-race` and `-msan` began implying `-d=checkptr` on every platform — why the race gate finds pointer bugs that vet does not | `edges` | [Go 1.15](https://go.dev/doc/go1.15) | `edges/05-unsafe-pointer-rules` |
| `panic` prints derived types rather than bare addresses — what a panic message is allowed to know about its value | `edges` | [Go 1.15](https://go.dev/doc/go1.15) | `edges/10-what-a-panic-prints` |
| `checkptr`'s two rules, stated exactly: alignment on conversion, and same-object arithmetic | `edges` | [Go 1.14](https://go.dev/doc/go1.14) | `edges/05-unsafe-pointer-rules` |
| `defer` became almost free — for *most* uses; which ones still fall back to the heap is the whole question | `edges` | [Go 1.14](https://go.dev/doc/go1.14) | `edges/04-defer-tiers` |
| `math.FMA(x, y, z)` and `x*y + z` can produce different `float64` values — predict which inputs expose the gap | `edges` | [Go 1.14](https://go.dev/doc/go1.14) | — |
| `runtime.Goexit` can no longer be aborted by a recursive `panic`/`recover` — a third exit path that is neither | `edges` | [Go 1.14](https://go.dev/doc/go1.14) | `edges/03-panic-vs-fatal` |
| `hash/maphash` is consistent within a process and different across them — the guarantee people accidentally rely on | `edges` | [Go 1.14](https://go.dev/doc/go1.14) | `edges/09-maphash-seeds` |
| `defer` got 30% faster in 1.13 and nearly free in 1.14 — two mechanisms, two consecutive releases | `edges` | [Go 1.13](https://go.dev/doc/go1.13) | `edges/04-defer-tiers` |
| `Sin`, `Cos`, `Tan` stopped being bit-for-bit reproducible across releases — Go promises accuracy, not reproducibility | `edges` | [Go 1.12](https://go.dev/doc/go1.12) | — |
| Converting a nil `unsafe.Pointer` to `uintptr` and back with arithmetic is invalid — the obvious-looking loophole in pattern 3 | `edges` | [Go 1.12](https://go.dev/doc/go1.12) | `edges/05-unsafe-pointer-rules` |
| The compiler rejects an unused type switch guard — a case where gc was the lenient implementation and gccgo was not | `edges` | [Go 1.11](https://go.dev/doc/go1.11) | — |
| Stack traces stopped including `<autogenerated>` wrappers, so a `runtime.Caller` skip count finally matched the source | `edges` | [Go 1.10](https://go.dev/doc/go1.10) | — |
| Test results are cached, and `-count=1` is the documented escape — what the cache keys on, and when it is wrong | `edges` | [Go 1.10](https://go.dev/doc/go1.10) | the verifier depends on it |
| `go test` runs a high-confidence subset of `go vet` first — which is not the same set `go vet` runs | `edges` | [Go 1.10](https://go.dev/doc/go1.10) | `edges/05-unsafe-pointer-rules` keeps snippets in testdata because of it |
| A `time.Time` carries two clocks: which operations use the monotonic reading, which strip it, and why `==` is treacherous | `edges` | [Go 1.9](https://go.dev/doc/go1.9) | `edges/06-two-clocks` |
| Go pointers passed to C: the sharing rules, enforced at run time, and `GODEBUG=cgocheck` | `edges` | [Go 1.6](https://go.dev/doc/go1.6) | — |
| The runtime assumes pointer-typed means pointer: an integer in a pointer slot crashes, a pointer in an integer slot is silent and fatal later | `edges` | [Go 1.3](https://go.dev/doc/go1.3), [1.4](https://go.dev/doc/go1.4) | `edges/05-unsafe-pointer-rules` |
| Green Tea cuts GC overhead 10–40%, and `GOEXPERIMENT=nogreenteagc` makes it an A/B within one toolchain | `gc` | [Go 1.26](https://go.dev/doc/go1.26) | — |
| Green Tea's further ~10% on Ice Lake / Zen 4 and newer — an answer that depends on the CPU | `gc` | [Go 1.26](https://go.dev/doc/go1.26) | — |
| Stop-the-world pauses split into a stopping phase and a total — two numbers where the trace used to report one | `gc` | [Go 1.22](https://go.dev/doc/go1.22) | `versions/11-four-pause-metrics` |
| Transparent huge pages are managed explicitly on Linux — the runtime overriding a kernel policy, with a measurable result | `gc` | [Go 1.21](https://go.dev/doc/go1.21) | — |
| `GOGC` and `GOMEMLIMIT` became readable as `runtime/metrics` values — configuration a program can inspect about itself | `gc` | [Go 1.21](https://go.dev/doc/go1.21) | — |
| The collector's own internal structures got 2% cheaper — the cost of the collector's bookkeeping, separate from its work | `gc` | [Go 1.20](https://go.dev/doc/go1.20) | — |
| `GOMEMLIMIT` is a soft limit that holds even with `GOGC=off` — the two knobs are not alternatives | `gc` | [Go 1.19](https://go.dev/doc/go1.19) | — |
| The GC CPU limiter caps collection at 50% of CPU — what a program does when it cannot collect fast enough | `gc` | [Go 1.19](https://go.dev/doc/go1.19) | — |
| The pacer started counting stack scanning and globals — the heap goal stopped being a function of the heap alone | `gc` | [Go 1.18](https://go.dev/doc/go1.18) | `gc/05-gctrace` |
| `MADV_DONTNEED` returned as the default so RSS reflects physical memory again — the end of a four-year arc | `gc` | [Go 1.16](https://go.dev/doc/go1.16) | — |
| `runtime/metrics` arrived: a supported, versioned alternative to parsing `gctrace` lines | `gc` | [Go 1.16](https://go.dev/doc/go1.16) | `gc/05-gctrace` parses the unsupported one |
| `ReadMemStats` stopped stopping the world — the instrument stopped perturbing the measurement | `gc` | [Go 1.15](https://go.dev/doc/go1.15) | — |
| The runtime returns memory promptly after a heap spike, and RSS still does not move — what "returned" means | `gc` | [Go 1.13](https://go.dev/doc/go1.13) | — |
| `MADV_FREE` became the default and RSS stopped being the truth for four releases — which measurement changes at each step | `gc` | [Go 1.12](https://go.dev/doc/go1.12) | — |
| Sweeping got faster when most of the heap survives — the cost is paid by the *allocator*, so it shows up as allocation latency | `gc` | [Go 1.12](https://go.dev/doc/go1.12) | — |
| The collector's CPU fraction fell while its duration rose and the total stayed put — three quantities, two of them moved | `gc` | [Go 1.10](https://go.dev/doc/go1.10) | `gc/05-gctrace` |
| Stop-the-world stack rescanning was eliminated by the hybrid write barrier: pauses fell from milliseconds to microseconds | `gc` | [Go 1.8](https://go.dev/doc/go1.8) | — |
| The collector became concurrent — the release every fact in the `gc` track dates from | `gc` | [Go 1.5](https://go.dev/doc/go1.5) | — |
| Generic methods: how many instantiations does GC-shape stenciling emit, and how does it differ from a generic function? | `generics` | [Go 1.27](https://go.dev/doc/go1.27) | `generics/01-instantiation-count`, `generics/02-generic-methods` |
| Generic type aliases are fully supported — does a parameterized alias add an instantiation, or share one? | `generics` | [Go 1.24](https://go.dev/doc/go1.24) | `versions/12-generic-alias` |
| `new` accepts an expression, so `new(f(x))` compiles — a version-diff task a pre-1.26 model gets wrong | `generics` | [Go 1.26](https://go.dev/doc/go1.26) | `generics/04-new-and-self-reference` |
| A generic type may refer to itself in its own type parameter list (`type Adder[A Adder[A]]`) | `generics` | [Go 1.26](https://go.dev/doc/go1.26) | `generics/04-new-and-self-reference` |
| Function type inference generalized to assignment and conversion contexts — not gated by the go.mod language version, so a true A/B needs an older toolchain | `generics` | [Go 1.27](https://go.dev/doc/go1.27) | `generics/05-inference-limits` — the assignment shape only; no 1.26-vs-1.27 A/B, see the radar note |
| Type inference gained four specific new abilities in one release — predict which of a set of calls compile | `generics` | [Go 1.21](https://go.dev/doc/go1.21) | `generics/05-inference-limits` |
| `comparable` may be satisfied by types that are not strictly comparable — a constraint whose name stopped meaning what it says | `generics` | [Go 1.20](https://go.dev/doc/go1.20) | — |
| Generics arrived with six limitations stated as a list — which of them are still true | `generics` | [Go 1.18](https://go.dev/doc/go1.18) | `generics/03-constraint-satisfaction` |
| `context.Background()` and `context.TODO()` can now compare equal — two sentinels that stopped being distinguishable | `iface` | [Go 1.21](https://go.dev/doc/go1.21) | `iface/05-interface-comparison` |
| The spec spelled out how structs and arrays are compared, including the freedom to stop early | `iface` | [Go 1.20](https://go.dev/doc/go1.20) | `iface/05-interface-comparison` |
| `netip.Addr` is comparable and `net.IP` is not — a deliberate redesign around the comparability rule | `iface` | [Go 1.18](https://go.dev/doc/go1.18) | `iface/05-interface-comparison` |
| `go vet` learned to spot impossible interface assertions — same method name, different signature | `iface` | [Go 1.15](https://go.dev/doc/go1.15) | — |
| Overlapping interface embedding became legal, but only for *identical* signatures — predict which declarations compile | `iface` | [Go 1.14](https://go.dev/doc/go1.14) | — |
| Range-over-function iterators and the `iter` package — what `break`, `return` and `goto` do to the yield contract | `iter` | [Go 1.23](https://go.dev/doc/go1.23) | `iter/02-yield-contract`, `iter/01-adapters`, `iter/04-goto-and-labels` |
| Heap metadata moved next to the object and alignment fell from 16 to 8 — what that does to a struct's real footprint. **Probed thin 2026-09-21**, with the `alloc` row above: the same lead, and the same measurement kills both | `layout` | [Go 1.22](https://go.dev/doc/go1.22) | — |
| `atomic.Int64` aligns itself even on 32-bit, where a bare `int64` field does not — the fix for a decade of alignment panics | `layout` | [Go 1.19](https://go.dev/doc/go1.19) | `layout/05-alignment-64bit` |
| Functions became 32-byte aligned to dodge a CPU erratum — code alignment as a correctness measure, not a performance one | `layout` | [Go 1.15](https://go.dev/doc/go1.15) | — |
| `int` is implementation-defined and became 64-bit on 64-bit platforms in 1.1 — so `unsafe.Sizeof` is a property of the build | `layout` | [Go 1.1](https://go.dev/doc/go1.1) | `layout/01-struct-padding` |
| `sync.OnceFunc`, `OnceValue` and `OnceValues` — the publication pattern given three standard spellings | `memmodel` | [Go 1.21](https://go.dev/doc/go1.21) | `memmodel/05-publication` |
| `sync.Map` gained `Swap`, `CompareAndSwap` and `CompareAndDelete` — atomic update on a concurrent map | `memmodel` | [Go 1.20](https://go.dev/doc/go1.20) | — |
| The memory model was revised to pin Go to sequential consistency for `sync/atomic` — the document, not the implementation, changed | `memmodel` | [Go 1.19](https://go.dev/doc/go1.19) | `memmodel/01-store-buffering` |
| The race detector got faster and lost its 8192-goroutine ceiling — what the detector can and cannot see | `memmodel` | [Go 1.19](https://go.dev/doc/go1.19) | `memmodel/02-benign-race` |
| `sync.Mutex.TryLock` and `RWMutex.TryLock` — a lock operation with no happens-before edge on failure | `memmodel` | [Go 1.18](https://go.dev/doc/go1.18) | — |
| `atomic.Value` gained `Swap` and `CompareAndSwap` — read-modify-write on a publication slot | `memmodel` | [Go 1.17](https://go.dev/doc/go1.17) | `memmodel/05-publication` |
| The race detector started catching races it used to miss — a "clean" pre-1.16 run proves less than it seemed to | `memmodel` | [Go 1.16](https://go.dev/doc/go1.16) | `memmodel/02-benign-race` |
| `sync.Map`: a read-only snapshot published through an atomic plus a dirty map behind a mutex, promoted on a miss count | `memmodel` | [Go 1.9](https://go.dev/doc/go1.9) | `memmodel/05-publication` |
| `reflect.TypeFor[T]()` — a type descriptor without a value, and without `(*T)(nil)` | `reflect` | [Go 1.22](https://go.dev/doc/go1.22) | — |
| `reflect.ValueOf` stopped forcing its argument to the heap — reflection stopped being an escape hatch in the literal sense | `reflect` | [Go 1.21](https://go.dev/doc/go1.21) | `alloc/03-what-escapes` |
| The linker deletes dead global map variables — a `map` initialised in `init` and never read costs nothing | `reflect` | [Go 1.21](https://go.dev/doc/go1.21) | — |
| `reflect.Value.Comparable` and `Equal` — asking whether a value *can* be compared before comparing it | `reflect` | [Go 1.20](https://go.dev/doc/go1.20) | `iface/05-interface-comparison` |
| `reflect.Len` and `Cap` accept a pointer to an array — the one place reflection auto-dereferences | `reflect` | [Go 1.19](https://go.dev/doc/go1.19) | — |
| `reflect.MapIter.Reset` makes reflective map iteration allocation-free — the third of three releases to get there | `reflect` | [Go 1.18](https://go.dev/doc/go1.18) | — |
| `reflect.ConvertibleTo` stopped being a guarantee: a convertible pair whose conversion panics | `reflect` | [Go 1.17](https://go.dev/doc/go1.17) | `types/05-compiler-free-conversions` |
| The linker prunes reflection-reachable symbols more aggressively — what `reflect.MethodByName` costs a binary | `reflect` | [Go 1.16](https://go.dev/doc/go1.16) | — |
| `reflect.Zero` stopped allocating, and two zero values stopped comparing the way they used to | `reflect` | [Go 1.16](https://go.dev/doc/go1.16) | — |
| Binary size fell 5% by dropping type metadata the runtime no longer needed — reflection's cost, stated as a number | `reflect` | [Go 1.15](https://go.dev/doc/go1.15) | — |
| `reflect` closed a hole in unexported field access — a rule the package had been failing to enforce | `reflect` | [Go 1.15](https://go.dev/doc/go1.15) | `reflect/03-settability` |
| `reflect.StructOf` gained unexported fields via `PkgPath` — the export boundary as something you construct | `reflect` | [Go 1.14](https://go.dev/doc/go1.14) | `reflect/02-structof` |
| `reflect.Value.IsZero` arrived in 1.13 and was corrected in 1.22 for negative zero — one function, two definitions of zero | `reflect` | [Go 1.13](https://go.dev/doc/go1.13) | — |
| `reflect.MapIter` — reflective map iteration that does not allocate a slice of every key | `reflect` | [Go 1.12](https://go.dev/doc/go1.12) | — |
| An embedded pointer to an unexported struct type used to punch a hole through the export check — `CanSet` was wrong for years | `reflect` | [Go 1.10](https://go.dev/doc/go1.10) | `reflect/03-settability` |
| The `goroutineleak` profile is generally available — plant a leak of each shape and make the profile name them | `sched` | [Go 1.27](https://go.dev/doc/go1.27) | `sched/05-goroutine-leak-profile` |
| Timer channels are always unbuffered now that `asynctimerchan` is gone | `sched` | [Go 1.27](https://go.dev/doc/go1.27) | `sched/03-timer-channels`, `versions/05-go-line-limits`, `versions/06-timer-channel-buffer` |
| `GOMAXPROCS` is container-aware and updates itself as the cgroup quota changes. **Pin candidate:** needs Linux cgroups, which is not a strike — `requires.os: [linux]`. Note gate 2 would then prove it on one runner only | `sched` | [Go 1.25](https://go.dev/doc/go1.25) | — |
| `testing/synctest` is GA: virtualized time in a bubble — a task, and the harness that makes `sched` and `memmodel` gradeable at all | `sched` | [Go 1.25](https://go.dev/doc/go1.25) | `sched/02-synctest-bubble` |
| `runtime/trace.FlightRecorder` as the data source for a trace-analyzer task | `sched` | [Go 1.25](https://go.dev/doc/go1.25) | — |
| A new runtime-internal mutex (`GOEXPERIMENT=nospinbitmutex`) — only worth a task if the spin/park boundary can be made visible | `sched` | [Go 1.24](https://go.dev/doc/go1.24) | — |
| Timers and tickers are collected without `Stop`, and the behaviour is gated on the `go.mod` go line rather than the toolchain | `sched` | [Go 1.23](https://go.dev/doc/go1.23) | `sched/03-timer-channels` |
| Windows timer resolution went from 15.6ms to 0.5ms — an answer that depends on the OS | `sched` | [Go 1.23](https://go.dev/doc/go1.23) | — |
| Mutex profiles now scale by the number of blocked goroutines — the third correction the radar found to Go's contention profiling | `sched` | [Go 1.22](https://go.dev/doc/go1.22) | — |
| The execution tracer was rewritten and can be streamed — traces stopped being a whole-program stop | `sched` | [Go 1.22](https://go.dev/doc/go1.22) | — |
| Stack traces name the goroutine that created each goroutine — the 1.11 opt-in ancestor made unconditional | `sched` | [Go 1.21](https://go.dev/doc/go1.21) | `sched/05-goroutine-leak-profile` |
| New metrics for `GOMAXPROCS`, cgo calls and mutex wait — runtime state a program can read about itself | `sched` | [Go 1.20](https://go.dev/doc/go1.20) | — |
| The CPU profiler moved to per-thread timers on Linux — a profiler that had been under-sampling busy threads | `sched` | [Go 1.18](https://go.dev/doc/go1.18) | — |
| Goroutine scheduling latency became a `runtime/metrics` distribution — how long a runnable goroutine waits for a P | `sched` | [Go 1.17](https://go.dev/doc/go1.17) | `sched/01-runnext-ordering` |
| Block profiles stopped favouring rare long events — a debiasing fix to an instrument, not to the program | `sched` | [Go 1.17](https://go.dev/doc/go1.17) | — |
| Non-blocking receives on closed channels got faster — the fast path of the idiom every cancellation uses | `sched` | [Go 1.15](https://go.dev/doc/go1.15) | — |
| `os` and `net` started retrying on `EINTR` — async preemption's signals were interrupting syscalls | `sched` | [Go 1.15](https://go.dev/doc/go1.15) | `sched/04-async-preemption` |
| Goroutines became asynchronously preemptible — on every platform *except* four, where the old behaviour still holds | `sched` | [Go 1.14](https://go.dev/doc/go1.14) | `sched/04-async-preemption` |
| Unlocking a contended mutex hands the CPU directly to the next waiter — which changes acquisition *order*, not just speed | `sched` | [Go 1.14](https://go.dev/doc/go1.14) | `sched/01-runnext-ordering` |
| Timers became cheaper "with no user visible changes" — a claim that later releases falsified | `sched` | [Go 1.14](https://go.dev/doc/go1.14) | `sched/03-timer-channels` |
| Timer and deadline code started scaling with CPU count — step one of a rewrite that spans six releases | `sched` | [Go 1.12](https://go.dev/doc/go1.12) | — |
| Non-blocking descriptors use the runtime poller instead of a thread — the difference between costing a thread and costing nothing | `sched` | [Go 1.11](https://go.dev/doc/go1.11) | — |
| `GODEBUG=tracebackancestors=N` extends a traceback with the stacks that created each goroutine | `sched` | [Go 1.11](https://go.dev/doc/go1.11) | `sched/05-goroutine-leak-profile` |
| The mutex profile learned about reader/writer contention — a starving writer used to produce an empty profile | `sched` | [Go 1.11](https://go.dev/doc/go1.11) | — |
| `LockOSThread` calls nest, and a locked thread is retired rather than reused — lock and unlock in a loop and count threads | `sched` | [Go 1.10](https://go.dev/doc/go1.10) | — |
| The 1024 ceiling on `GOMAXPROCS` was removed | `sched` | [Go 1.10](https://go.dev/doc/go1.10) | — |
| `GOMAXPROCS` defaults to the number of cores — which is why every `sched` task must pin it | `sched` | [Go 1.5](https://go.dev/doc/go1.5) | — |
| Preemption at (non-inlined) function entry: for twelve years, whether a loop could be preempted depended on an inlining decision | `sched` | [Go 1.2](https://go.dev/doc/go1.2) | `sched/04-async-preemption`, `compiler/01-inlining-budget` |
| The builtin map is a Swiss table — predict real memory for a key count, A/B with `GOEXPERIMENT=noswissmap` | `types` | [Go 1.24](https://go.dev/doc/go1.24) | `types/04-map-memory` — memory half only; the `noswissmap` A/B is still open |
| `sync.Map` is a hash-trie — the contention curve, and what the old implementation was warming up | `types` | [Go 1.24](https://go.dev/doc/go1.24) | — |
| Shrinking a slice now zeroes the tail beyond the new length — `s = s[:0]` stopped keeping its elements alive | `types` | [Go 1.22](https://go.dev/doc/go1.22) | `types/02-aliasing` |
| `reflect.Value.IsZero` agrees with `==` for negative zero — a correction to a definition that looked obvious | `types` | [Go 1.22](https://go.dev/doc/go1.22) | `versions/10-negative-zero` |
| Slice-to-array conversion (not to array *pointer*) — which lengths panic, and at what point | `types` | [Go 1.20](https://go.dev/doc/go1.20) | — |
| `append`'s growth formula changed: the 1024-element doubling threshold became a smooth 1.25× ramp | `types` | [Go 1.18](https://go.dev/doc/go1.18) | `types/01-cap-growth` |
| `strings.Clone` exists specifically to stop sharing memory — when a substring keeps a megabyte alive | `types` | [Go 1.18](https://go.dev/doc/go1.18) | `types/02-aliasing` |
| A type conversion that panics — the shape that made `ConvertibleTo` stop being a guarantee | `types` | [Go 1.17](https://go.dev/doc/go1.17) | `types/05-compiler-free-conversions` |
| `fmt` prints maps in key-sorted order while `range` stays randomised — both true at once, with stated ordering rules for `NaN` keys | `types` | [Go 1.12](https://go.dev/doc/go1.12) | — |
| `append(s, make([]T, n)...)` became one allocation at the final size rather than two | `types` | [Go 1.11](https://go.dev/doc/go1.11) | `types/01-cap-growth` |
| `bytes.Split` and `Fields` clip returned subslices to capacity so appending cannot overwrite the input | `types` | [Go 1.10](https://go.dev/doc/go1.10) | `types/02-aliasing` |
| Type aliases: `reflect.TypeOf` cannot tell `T1 = T2` apart, and `case byte:` with `case uint8:` is a compile error | `types` | [Go 1.9](https://go.dev/doc/go1.9) | — |
| Three-index slicing `a[2:4:7]` — the only way to hand out a window without handing out the room behind it | `types` | [Go 1.2](https://go.dev/doc/go1.2) | `types/02-aliasing` |
| `runtime.AddCleanup` vs `SetFinalizer` — four enumerated differences, four predictions | `weak` | [Go 1.24](https://go.dev/doc/go1.24) | `weak/02-cleanup-vs-finalizer` |
| The `weak` package — the floor for the whole track | `weak` | [Go 1.24](https://go.dev/doc/go1.24) | `weak/01-weak-cache` |
| `AddCleanup` runs concurrently and `unique` handles reclaim in a single GC cycle — every pre-1.25 measurement is stale | `weak` | [Go 1.25](https://go.dev/doc/go1.25) | `weak/02-cleanup-vs-finalizer` + `weak/05-single-cycle-reclaim` |
| `GODEBUG=checkfinalizers=1` names the classic finalizer mistakes — plant each one | `weak` | [Go 1.25](https://go.dev/doc/go1.25) | `weak/03-lifetime-traps` |
| The `unique` package: interning where handle comparison reduces to a pointer compare | `weak` | [Go 1.23](https://go.dev/doc/go1.23) | `weak/04-unique-interning` |
| More precise liveness analysis made finalizers run sooner — and made "still in scope" stop meaning "still alive" | `weak` | [Go 1.12](https://go.dev/doc/go1.12) | `weak/05-single-cycle-reclaim` |

## Invalidated by a release

| Task | Release | What changed | Status |
|---|---|---|---|
| `layout/01-struct-padding` | [Go 1.23](https://go.dev/doc/go1.23) | `structs.HostLayout` exists because "struct layout order is not guaranteed by the language spec". The sealed hint said "Go lays a struct out in declaration order", stating a gc implementation detail as a language guarantee. The five measured answers were unaffected. | **fixed 2026-09-11** — hint corrected, explanation gained a section on spec guarantees vs implementation, task re-sealed |

## Release radar

Every Go release gets a triage pass answering two questions: what does it
create, and what does it invalidate. The radar also runs backwards over Go's
whole history, because a behaviour that *changed* is the best kind of task —
"this was true in 1.21; is it still?"

Sweep order, most actionable first. **The sweep is complete: every Go release
from 1.0 to 1.27 is triaged**, 271 entries across 19 files, 194 candidates and
1 invalidation.

| Era | Versions | Status |
|---|---|---|
| E1 | 1.23 – 1.27 | **done** — 60 entries across 5 releases: 36 candidates, 1 invalidation. Swept twice: the first pass (2026-09-11) read the runtime, toolchain and library sections and caught three language changes but missed six; the second (2026-09-13) read every release's "Changes to the language" section on its own. |
| E2 | 1.18 – 1.22 | **done** — 80 entries across 5 releases: 59 candidates, 0 invalidations. |
| E3 | 1.10 – 1.17 | **done** — 114 entries across 8 releases: 85 candidates, 0 invalidations. The densest era: the SSA optimisations, the escape-analysis rewrite, async preemption and `checkptr` all land here. |
| E4 | 1.0 – 1.9 | **done** — harvested rather than swept, into one file: 17 entries, 14 candidates. Ten thin files would have misrepresented the yield. |

That the historical sweep produced **zero invalidations** is itself a finding.
The one invalidation on record came from E1 — the newest era, where the
catalogue's answers are youngest and least settled. Old behaviour that survived
to 1.27 is old behaviour that has stopped moving.

Output lands in `radar/versions/go1.NN.md`, each entry tagged `→ candidate`,
`→ invalidates <task-id>`, or `→ no action`. `undergo radar-check` runs in CI and
fails if an invalidation names a task that does not exist, so the dataset cannot
rot as the catalogue changes. Era E4 is the one exception to one file per
release, which is why the check counts *radar files*.

## Standing decisions

**Authoring comes before solving, deliberately.** 241 tasks ship and none have been
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

**Binary size is not measurable reliably enough to grade, and the binary-size
golf task was withdrawn.** (It never reached a tag, so the number it had been
given was later reused by `compiler/04-pgo-devirtualization`.) Spec §8 T12's
"binary-size golf" was built, passed locally on Windows
and in a Linux container, and then failed intermittently on the ubuntu CI runner
across three attempts — including after every slot was moved behind an 8 KiB
threshold. Measured evidence for why: the same snippet is +24 bytes on Linux and
exactly 0 on Windows because PE section padding rounds it away; building the same
source in two temporary directories of different name lengths changes the size,
because source paths are embedded; and `-trimpath` shifts sizes by a couple of
hundred bytes in both directions rather than removing the dependence. A binary's
size has many contributors that have nothing to do with the question being asked.

The idea stays in the pool. A future attempt should measure something with a
defined meaning — a section size from `go tool nm`, or the count of symbols
retained from a named package — rather than the size of a file on disk.

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

**Every track carries at least five tasks, and as of v0.13.0 every one does.**
Ten where the material supports it. This is discipline, not machinery — there is
deliberately no CI gate. What keeps it honest is that the ceilings are written
down, so a track sitting below ten is legibly finished rather than neglected:

| Verdict | Tracks |
|---|---|
| Ten comfortably | `sched`, `alloc`, `types`, `layout`, `generics`, `compiler`, `asm` |
| Ten if a blocked item unblocks | `gc` (needs the Green Tea invariant), `edges` (cgo needs a C toolchain; the 1.21 nil-check task needs two toolchains), `memmodel` (IRIW is legality-only, being unobservable on x86) |
| About eight | `iface`, `reflect`, `iter` |
| About seven, and will not reach ten | `weak` — the `weak` and `unique` packages have a genuinely small surface |

**`sched` and `memmodel` are gradeable because of answer form, not new
machinery.** M10 added ten tasks and no harness code beyond a failure log. Every
slot in both tracks is an ordering, a predicate, a legality set or a
compile-time constant; not one grades a duration or a count the scheduler can
perturb. Three levers made it possible: `testing/synctest`'s fake clock, the
`goroutineleak` profile, and the race gate, which turns the detector from a
guard into a grader.

Two slots were written, measured, and replaced before shipping, both for the
same reason. `sched/02` asked whether the test finished in under a second — a
tolerance against a constant, which measured 0–22ms locally and failed half of
all gate-2 runs, because gate 2's load is process creation and compilation
rather than CPU contention. It now compares the two clocks against each other.
The same task's ticker count was a coin flip, 9 in 28 runs and 10 in 32: the
tenth tick and the deadline landed on the same instant of fake time, and
`select` chooses among ready cases at random. **A fake clock removes
non-determinism that comes from scheduling; it cannot remove non-determinism
that is specified.**

**The race detector is not a neutral observer, and now in three ways.**
`requires.default_build` exists for tasks whose answers hold only on an
unmodified build, and each milestone has found a new reason for it. M7: a
compiler optimisation that `-race` disables (`types/01-cap-growth`). M10: a task
whose subject *is* a data race, which is most of `memmodel`. Also M10, and
caught by the gate rather than by design: **a task measuring scheduling order.**
`sched/01-runnext-ordering` measures `4 1 2 3` unmodified and `4 2 1 3` under
`-race`, in a program with no data race at all — instrumentation inserts work
between a goroutine being scheduled and reaching the mutex, so the `runnext`
slot survives and the queue order behind it does not. Reach for the detector to
answer *is there a race*, never *in what order do these run*.

**Gate 1 builds only what the host can build, and says what it skipped.**
`requires.arch` gated gate 2 but not gate 1, so an amd64-only `.s` failed to
build on the arm64 runner — and every assembly task carried a portable Go
fallback to compensate, which meant a green arm64 run proved the fallback worked
and said nothing about the assembly. Gate 1 now consults `manifest.Buildable`,
deliberately narrower than `Gradeable`: only arch and OS decide whether a package
compiles, because a task needing an extra toolchain still builds with the default
one and a task needing a newer Go should fail loudly rather than vanish.

This is safe **only because CI runs both architectures** — `ubuntu-latest` is
amd64 and `macos-latest` is arm64. A task pinning a platform no runner has would
be built nowhere and nobody would be told. New assembly tasks need no fallback;
`asm/02-write-plan9` keeps its own because its README teaches the pattern.

**An answer that depends on the machine is not an answer, and it takes CI to
find out.** Three M10 slots shipped green locally and failed in CI, each for a
different reason, and the pattern is worth naming because it recurs:

- A **wall-clock budget** (`sched/02`: "did this finish in under a second")
  measured 0–22ms locally and failed half of all gate-2 runs. Gate 2's load is
  process creation, compilation and antivirus scanning, not CPU contention, so a
  CPU-load probe "proved" it safe while measuring the wrong thing. A threshold
  against a constant is a tolerance; compare two measured quantities instead.
- A **struct width** (`memmodel/02`: tearing) was two words, which is two `MOVQ`s
  on amd64 and a single `STP` on arm64 — so it did not tear there at all. Widened
  to three words rather than pinned, because the claim being taught is
  architecture-independent and only the demonstration was accidentally amd64.
- A **race being won** (`memmodel/05`: duplicate lazy initialisation) was true on
  twelve cores and false on a CI runner. The constructor now yields part-way
  through, which makes the window reliable rather than lucky.

**The verifier says what a failing predict task got wrong.** Gate 2 captures a
failing task's output rather than printing it, so a reference solution never
reaches a public log — which left "frozen tests did not pass" as the only
diagnostic, and two CI failures in a row were investigated blind. The output now
goes to `.undergo/ci-failures/<id>.log` (gitignored) and only the path is
printed. In addition, for a **predict** task the failing slot names are printed:
`predict.Check` emits only `slot "name": correct` and never a measured value, and
the names are already public in `task.yaml`. Build and optimize tasks print
nothing, because their output can contain the overlaid reference source.

## Not yet

Contests — scored challenges, leaderboards and seasons — are deliberately out of
scope until the catalogue is large enough to be worth competing over.
