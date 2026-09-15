# Go 1.16

Source: https://go.dev/doc/go1.16

## Memory is returned to the OS eagerly again, and RSS started telling the truth
→ candidate | gc
"The runtime now defaults to releasing memory to the operating system promptly
(using `MADV_DONTNEED`), rather than lazily when the operating system is under
memory pressure (using `MADV_FREE`). This means process-level memory statistics
like RSS will more accurately reflect the amount of physical memory being used by
Go processes."

The gradeable shape is a disagreement between two measurements of the same
program: Go's own heap figures are identical under both settings, while the
operating system's RSS differs by however much has been freed but not reclaimed.
`GODEBUG=madvdontneed=0` restores the lazy behaviour, so it is an A/B in one
toolchain. Predict which of the two numbers moves.

This also explains why "my Go service uses more memory than it says it does" was
a real and correct observation before 1.16 and is a misunderstanding after.

## runtime/metrics arrived
→ candidate | gc
"a stable interface for reading implementation-defined metrics from the Go
runtime ... supersedes existing functions like `runtime.ReadMemStats` and
`debug.GCStats` ... significantly more general and efficient."

`gc/01-heap-goal`, `gc/03-memory-limit` and `gc/05-gctrace` all sit on top of
this package, and `gc/05`'s last question is why a supported API beats parsing a
debug format. This is where the supported API came from, and the interesting
question is what "implementation-defined" buys: a metric may disappear between
releases, so a robust reader must handle a name it does not recognise.

## The inliner learned three new shapes
→ candidate | compiler
"The compiler can now inline functions with non-labeled `for` loops, method
values, and type switches ... The inliner can also detect more cases where a
called function is a constant, allowing more inlining of indirect calls."

This is the first step of an arc: 1.16 admits non-labeled `for` loops, 1.18 adds
range loops *and* labeled loops, and 1.21 adds call-site heuristics. Before all
of it, any loop at all made a function non-inlinable outright.

`compiler/01-inlining-budget` grades what fits the budget on a current toolchain.
The version-diff question is which of a set of functions were inlinable at each
point on that arc, and why a structural veto was the starting position — the
answer is that inlining a loop makes the cost model much harder to get right.

## GODEBUG=inittrace=1 prints what every package init cost
→ candidate | edges
"the runtime emits a single line to standard error for each package `init`,
summarizing its execution time and memory allocation."

Directly observable, and it pairs with 1.21's initialisation-order algorithm: one
release made the order deterministic and this one made the cost visible. Predict
which packages dominate a program's startup, then check — the answer is usually
not the one people guess, because a transitively imported package with a large
table costs more than the one they are thinking about.

## The race detector started catching races it used to miss
→ candidate | memmodel
"This release fixes a discrepancy between the race detector and the Go memory
model. The race detector now more precisely follows the channel synchronization
rules of the memory model. As a result, the detector may now report races it
previously missed."

The whole `memmodel` track treats `-race` as an oracle for "is there a race
here". This is a release where the oracle was wrong — it modelled channel
synchronisation more strongly than the memory model actually promises, so some
genuine races looked synchronised. A clean run before 1.16 meant less than a
clean run after.

## reflect.Zero stopped allocating, and the values stopped comparing
→ candidate | reflect
"The `Zero` function is now more efficient, avoiding allocations. Code that
incorrectly compares the returned `Value` to another `Value` using `==` or
`DeepEqual` may now get different results."

Two tasks meet here. `reflect/04-deepequal-cycles` implements `DeepEqual`
semantics, and `iface/05-interface-comparison` grades what `==` does to values
whose representation you cannot see. A `reflect.Value` is a struct, so comparing
two of them compares their internals — which is why the documentation says not to
— and an allocation optimisation was enough to change the answer.

## vet learned that assembly must preserve BP
→ candidate | asm
"The vet tool now warns about amd64 assembly that clobbers the BP register (the
frame pointer) without saving and restoring it, as required by the calling
convention ... An easy fix is to set the frame size to a nonzero value, which
causes the function prologue and epilogue to preserve BP for you."

`asm/04-noescape` records that `go vet` checks a named offset exactly and does not
check the frame size at all. This is the opposite case: a check vet *does*
perform on assembly, and one whose fix is to change the very number vet will not
validate. The pair is a good question about what a linter can and cannot know.

## The linker prunes symbols more aggressively
→ candidate | reflect
"linking is 20-25% faster than 1.15 and requires 5-15% less memory on average for
`linux/amd64` ... Most binaries are also smaller as a result of more aggressive
symbol pruning."

`reflect/05-methodbyname-linker` counts which method symbols survive into a
binary. This is a release where the answer changed in the solver's favour — more
symbols pruned means a direct-call binary retains fewer — and it is the reason
that task's archive and binary counts diverge as sharply as they do.

## go test now fails on os.Exit(0)
→ candidate | edges
"A test binary that calls `os.Exit(0)` during execution of a test function will
now be considered to fail, to avoid the case where a test calls `os.Exit(0)` and
thereby exits without having run all the tests."

A test that passed by exiting successfully now fails, which is a genuinely
surprising inversion. It also matters to this repository: several tasks run
subprocesses and grade their exit codes, and `TestMain` is explicitly exempt.

## GODEBUG=madvdontneed is no longer needed
→ no action
"Systems that currently set `GODEBUG=madvdontneed=1` no longer need that setting."
Covered by the entry above; recorded separately only because a reader searching
for the flag should land somewhere.

## strconv.ParseFloat adopted Eisel-Lemire
→ no action
"improving performance by up to a factor of 2 ... can also speed up decoding
textual formats like `encoding/json`." Faster, identical results. Paired with
1.17's Ryū for formatting, the two together are a nice story about float
conversion being harder than it looks, but neither offers anything to predict.

## embed and io/fs
→ no action
Significant additions — `//go:embed`, the `fs.FS` interface, and the
producer/consumer split across `embed`, `zip`, `os.DirFS`, `http.FS` and
`template.ParseFS`. Library design rather than machine behaviour.

## Modules on by default
→ no action
`GO111MODULE` defaults to `on`, build commands no longer modify `go.mod`
implicitly, and `go install pkg@version` becomes the way to install a tool. The
end of an era, and nothing underneath Go.

## The timezone database shrank by 350 KB
→ no action
The slim format. Binary size, which this repository has already found
unmeasurable to a gradeable standard.
