# Go 1.15

Source: https://go.dev/doc/go1.15

## Converting a small integer to an interface stopped allocating
→ candidate | alloc
"Converting a small integer value into an interface value no longer causes
allocation."

This is the origin of `alloc/04-boxing-cache`, which grades why boxing 255 is
free and boxing 256 is not. The release note says only "a small integer value"
and leaves the boundary unstated, which is exactly the gradeable part: the
runtime keeps a static array of the first 256 byte-sized values, so boxing any of
them stores a pointer into that array instead of allocating. Above it there is
nothing to point at.

The version-diff question is what the same loop cost before 1.15 — one allocation
per conversion, every time — and why an interface needs a pointer at all for a
value that fits in a word.

## Small-object allocation stopped degrading at high core counts
→ candidate | alloc
"Allocation of small objects now performs much better at high core counts, and
has lower worst-case latency."

The mechanism is the per-P allocation caches and how they are refilled from the
central lists; contention on that refill is what scaled badly. A task can measure
allocation throughput against `GOMAXPROCS` and ask where the curve bends, which
is a question about the allocator's structure rather than about a benchmark
number.

## Chained unsafe.Pointer to uintptr conversions became illegal
→ candidate | edges
"Previously, in some cases, the compiler allowed multiple chained conversions
(for example, `syscall.Syscall(…, uintptr(uintptr(ptr)), …)`). The compiler now
requires exactly one conversion."

`edges/05-unsafe-pointer-rules` grades which snippets `go vet -unsafeptr`
rejects, and this is the compiler enforcing a rule the documentation always
stated. The doubled conversion looks harmless and is not: each `uintptr` is a
point at which the pointer is no longer a reference, so nesting them widens the
window in which the object is invisible to the collector.

## -race and -msan began implying -d=checkptr everywhere
→ candidate | edges
"The `-race` and `-msan` flags now always enable `-d=checkptr`, which checks uses
of `unsafe.Pointer`. This was previously the case on all OSes except Windows."

`edges/05`'s explanation rests on `checkptr` being available under `-race`. This
is where that became unconditional, and the Windows exception is the interesting
history: for a while the same program was checked on Linux and not on Windows,
so a clean run meant different things on different machines.

## Non-blocking receives on closed channels got faster
→ candidate | sched
"Non-blocking receives on closed channels now perform as well as non-blocking
receives on open channels."

A `select` with a `default` over a closed channel used to take a slower path.
That is measurable, and the mechanism is worth asking about: a closed channel has
no sender to coordinate with, so the fast path had to learn a case it previously
punted to the general one. Relevant to `memmodel/03-happens-before-edges`, which
grades what a receive from a closed channel guarantees.

## Functions became 32-byte aligned to dodge a CPU erratum
→ candidate | layout
"The toolchain now mitigates Intel CPU erratum SKX102 on `GOARCH=amd64` by
aligning functions to 32 byte boundaries and padding jump instructions. While
this padding increases binary sizes, this is more than made up for by the binary
size improvements mentioned above." A `-spectre` flag was added to the compiler
and assembler in the same release, "provided mainly as a 'defense in depth'
mechanism".

Function alignment is directly observable — take the address of a set of
functions and look at the low bits. A hardware bug reaching up through the
toolchain into the layout of compiled code is a good story, and the padding is
measurable in both directions.

## Binary size fell 5% by dropping type metadata
→ candidate | reflect
"Go 1.15 reduces typical binary sizes by around 5% compared to Go 1.14 by
eliminating certain types of GC metadata and more aggressively eliminating unused
type metadata."

*Unused type metadata* is the phrase that matters. `reflect/05-methodbyname-linker`
grades what the linker retains when reflection makes reachability undecidable;
this is the complementary case, where the linker proved metadata unreachable and
dropped it. The two together are one question: what does the linker have to
believe about reflection to remove anything at all?

## vet learned to spot impossible interface assertions
→ candidate | iface
"The vet tool now warns about type assertions from one interface type to another
interface type when the type assertion will always fail. This will happen if both
interface types implement a method with the same name but with a different type
signature."

Statically impossible, and until this release silently compiled. It pairs with
`iface/05-interface-comparison`: both are cases where the interface's type
erasure lets the compiler accept something that cannot succeed, and the question
is what the compiler would have to know to reject more.

## panic prints derived types, not just their addresses
→ candidate | edges
"If `panic` is invoked with a value whose type is derived from any of `bool`,
`complex64`, ... `string`, `uint` ... then the value will be printed, instead of
just its address. Previously, this was only true for values of exactly these
types."

`edges/03-panic-vs-fatal` grades what the runtime does with different failures.
This is the same machinery choosing what it can print: a named type with an
underlying `string` used to come out as a pointer, because the runtime matched on
the exact type rather than its kind.

## os and net started retrying EINTR
→ candidate | sched
"Packages `os` and `net` now automatically retry system calls that fail with
`EINTR`. Previously this led to spurious failures, which became more common in Go
1.14 with the addition of asynchronous preemption. Now this is handled
transparently."

A direct, named consequence of `sched/04-async-preemption`'s subject: preempting
goroutines with signals means more signals, which means more interrupted system
calls. One release added the mechanism and the next had to absorb its fallout —
which is the most instructive kind of pair the radar can find.

## ReadMemStats stopped stopping the world
→ candidate | gc
"Several functions, including `ReadMemStats` and `GoroutineProfile`, no longer
block if a garbage collection is in progress." A measurement tool that perturbs
what it measures is a recurring theme in this catalogue — `sched/01` is pinned to
the default build for exactly that reason — and this is the standard library
fixing an instance of it.

## reflect closed a hole in unexported field access
→ candidate | reflect
"Package `reflect` now disallows accessing methods of all non-exported fields,
whereas previously it allowed accessing those of non-exported, embedded fields."
Embedding used to promote a method past the export check. Predict which of a set
of `reflect` accesses now panic.

## The linker got 20% faster and used 30% less memory
→ no action
"linking is 20% faster and requires 30% less memory on average, for ELF-based
OSes ... on amd64", from "a newly redesigned object file format, and a revamping
of internal phases to increase concurrency". The first half of the two-release
linker modernisation that 1.16 completed. Toolchain performance.

## vet warns about string(int)
→ no action
"`string(9786)` does not evaluate to the string `"9786"`; it evaluates to the
string `"\xe2\x98\xba"`, or `"☺"`." A genuine and common confusion, and the note
says the language change was being considered — but the behaviour itself is
specified and unchanged, so there is nothing to predict that reading the spec
would not settle.

## time/tzdata can be embedded
→ no action
"Either approach increases the size of the program by about 800 KB." A concrete
number attached to an import, which is almost interesting, but it is binary size.
