# Go 1.19

Source: https://go.dev/doc/go1.19

## The memory model was revised, and pinned Go to sequential consistency
→ candidate | memmodel
"The Go memory model has been revised to align Go with the memory model used by
C, C++, Java, JavaScript, Rust, and Swift" — and the sentence the whole
`memmodel` track rests on: "Go only provides sequentially consistent atomics, not
any of the more relaxed forms found in other languages."

`memmodel/01-store-buffering` grades exactly this: `(0,0)` is legal with plain
variables and illegal with atomics, because Go has no relaxed ordering to reach
for. Before 1.19 the document did not say so in those terms. The version-diff
question is not about behaviour — the implementation already did this — but about
when a guarantee you were relying on became one.

## atomic.Int64 aligns itself, even on 32-bit
→ candidate | layout
The new `atomic.Bool`, `Int32`, `Int64`, `Uint32`, `Uint64`, `Uintptr` and
`Pointer[T]` types, with the guarantee that matters: "`Int64` and `Uint64` are
automatically aligned to 64-bit boundaries in structs and allocated data, even on
32-bit systems."

This is the origin of `layout/05-alignment-64bit`. Before these types, a 64-bit
atomic field had to be placed by hand — first in the struct, or padded — and the
bug was invisible on the developer's 64-bit machine. "These types hide the
underlying values so that all accesses are forced to use the atomic APIs" is the
second half: the type makes the misuse unrepresentable rather than merely
documented.

## GOMEMLIMIT, and it holds even with GOGC=off
→ candidate | gc
"The limit works in conjunction with `runtime/debug.SetGCPercent` / `GOGC`, and
will be respected even if `GOGC=off`, allowing Go programs to always make maximal
use of their memory limit."

`gc/03-memory-limit` grades which knob sets the heap goal; this is where the
second knob arrived and where the `GOGC=off` case — the one that looks like it
should disable collection entirely and does not — comes from. Note also what the
limit excludes: "mappings of the binary itself, memory managed in other
languages, and memory held by the operating system on behalf of the Go program."

## The GC CPU limiter caps collection at 50%
→ candidate | gc
"In order to limit the effects of GC thrashing when the program's live heap size
approaches the soft memory limit, the Go runtime also attempts to limit total GC
CPU utilization to 50%, excluding idle time, choosing to use more memory over
preventing application progress ... the new runtime metric
`/gc/limiter/last-enabled:gc-cycle` reports when this last occurred."

A hard number, a stated policy — *exceed the limit rather than stop making
progress* — and a metric that says whether it fired. That triple is what makes it
gradeable: drive a program into the limiter and predict which of the memory limit
and the progress guarantee gives way.

## Goroutine stacks are sized from history
→ candidate | alloc
"The runtime will now allocate initial goroutine stacks based on the historic
average stack usage of goroutines. This avoids some of the early stack growth and
copying needed in the average case in exchange for at most 2x wasted space on
below-average goroutines."

A stated tradeoff with a bound. The gradeable version: a program that spawns
many shallow goroutines after some deep ones pays a different per-goroutine cost
than one that does not, and the direction is measurable through
`runtime/metrics` stack figures. Before this, every goroutine started at the same
fixed size.

## Large switches became jump tables
→ candidate | compiler
"The compiler now uses a jump table to implement large integer and string switch
statements. Performance improvements for the switch statement vary but can be on
the order of 20% faster. (`GOARCH=amd64` and `GOARCH=arm64` only.)"

The threshold is the question: *large* is a specific number of cases, and below
it the compiler emits a comparison chain or a binary search instead. Readable
from `-gcflags=-S`, which makes it a predict-then-look task in the shape
`compiler/05-loop-lowering` already uses. Relevant to `iface/04-assertion-cost`,
where a type switch is the same decision on a different key.

## Fatal error tracebacks got shorter unless you ask
→ candidate | edges
"Unrecoverable fatal errors (such as concurrent map writes, or unlock of unlocked
mutexes) now print a simpler traceback excluding runtime metadata (equivalent to
a fatal panic) unless `GOTRACEBACK=system` or `crash`."

`edges/03-panic-vs-fatal` grades which failures are recoverable and names the
concurrent map write as the fatal case. This is the same boundary seen from the
output side: what the runtime chooses to print is a function of whether it
considers the failure the program's fault or its own.

## Importing os changes the process's file descriptor limit
→ candidate | edges
"On Unix operating systems, Go programs that import package `os` now
automatically increase the open file limit (`RLIMIT_NOFILE`) to the maximum
allowed value; that is, they change the soft limit to match the hard limit."

An `import` with an observable side effect on process state, visible from a
subprocess before and after. That is a genuinely surprising thing for an import
to do, and the reason — "artificially low limits set on some systems for
compatibility with very old C programs using the select system call" — is the
kind of history that makes a good explanation.

## The race detector got faster and lost its goroutine ceiling
→ candidate | memmodel
Thread sanitizer v3 is "typically 1.5x to 2x faster, uses half as much memory,
and it supports an unlimited number of goroutines." The ceiling is the
interesting half: v2 had a fixed maximum, so a program with enough goroutines
could not be raced at all. Anything in this repository that reasons about what
`-race` can and cannot observe — the whole `memmodel` track, and `sched/01`'s
`default_build` pin — sits on top of what the detector is capable of.

## reflect.Len and Cap accept a pointer to an array
→ candidate | reflect
"`Value.Len` and `Value.Cap` now successfully operate on a pointer to an array
and return the length of that array, to match what the builtin `len` and `cap`
functions do." A small change that exposes a real asymmetry: `len` has always
auto-dereferenced an array pointer, and `reflect` did not. Predict which of a set
of `reflect` calls now succeed.

## Doc comments became a parsed format
→ no action
Links, lists and headings, `gofmt` reformatting them, and a new `go/doc/comment`
package. Real and visible, but it is a documentation format rather than program
behaviour.

## os/exec stopped resolving PATH entries relative to the current directory
→ no action
"`Command` and `LookPath` no longer allow results from a PATH search to be found
relative to the current directory. This removes a common source of security
problems but may also break existing programs." Important, and worth knowing for
any task that shells out — several in this repository do — but it is a security
policy rather than a mechanism to measure.

## sort was rewritten to pattern-defeating quicksort
→ no action
"faster for several common scenarios." No ordering guarantee changed —
`sort.Sort` was never stable and still is not — so there is nothing a solver
could predict that was not already true.

## The linker emits compressed DWARF in the standard format
→ no action
`SHF_COMPRESSED` instead of the legacy `.zdebug`. Affects tools reading the
binary, not the binary's behaviour.
