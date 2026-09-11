# Go 1.25

Source: https://go.dev/doc/go1.25

## GOMAXPROCS is container-aware and updates itself
→ candidate | sched
On Linux "the runtime considers the CPU bandwidth limit of the cgroup containing
the process", and on all OSes "the runtime periodically updates `GOMAXPROCS` if
the number of logical CPUs available or the cgroup CPU bandwidth limit change".
Predict `GOMAXPROCS` for a container with a fractional CPU quota, and what
happens to a running program when the quota changes under it. The escape hatches
`containermaxprocs=0` and `updatemaxprocs=0` make it an A/B within one toolchain,
and `runtime.SetDefaultGOMAXPROCS` is the programmatic view of the same thing.

## testing/synctest is generally available
→ candidate | sched
"Within the bubble, time is virtualized: `time` package functions operate on a
fake clock and the clock moves forward instantaneously if all goroutines in the
bubble are blocked", with `Wait()` blocking until every goroutine in the bubble
is blocked. This is a task in its own right — predict what a bubbled program
observes — and it is also the harness that makes deterministic `sched` and
`memmodel` tasks gradeable at all. Read it before writing either track.

## The compiler stack-allocates slice backing stores in more cases
→ candidate | alloc
This is where the mechanism landed; 1.26 extends it. "The compiler can now
allocate the backing store for slices on the stack in more situations", with
`-gcflags=all=-d=variablemakehash=n` to turn it off. 1.25 carries the warning
that makes it an `unsafe` task as well as an allocation one: the change "has the
potential to amplify the effects of incorrect `unsafe.Pointer` usage."

## AddCleanup runs concurrently, and unique reclaims in a single cycle
→ candidate | weak
"Cleanup functions scheduled by `AddCleanup` are now executed concurrently and in
parallel", and values containing `unique.Handle`s that "previously required
multiple garbage collection cycles to collect ... are collected promptly in a
single cycle". The whole `weak` track depends on when things are reclaimed, so
this dates every answer in it: a measurement taken on 1.24 is wrong here. Any
such task must pin `requires.go` to at least 1.25.

## GODEBUG=checkfinalizers=1 diagnoses finalizer and cleanup mistakes
→ candidate | weak
It "runs diagnostics on each garbage collection cycle, and will also regularly
report the finalizer and cleanup queue lengths to stderr". A build task: plant
each classic finalizer mistake — the resurrection, the self-reference, the
cleanup that captures its own object — and make the diagnostic name it.

## runtime/trace.FlightRecorder
→ candidate | sched
"A lightweight way to capture a runtime execution trace by continuously
recording the trace into an in-memory ring buffer." The planned `sched` task
"write an analyzer over a `runtime/trace`" gets a better data source: capture
only the window around an event instead of tracing everything.

## GOAMD64=v3 enables fused multiply-add
→ candidate | edges
"The compiler will now use fused multiply-add instructions ... This may change
the exact floating-point values that a program generates." One program, two
`GOAMD64` levels, different numbers — with the mechanism (one rounding step
instead of two) as the answer. Needs `requires.arch: [amd64]` and must state the
`GOAMD64` level it assumes.

## A Go 1.21 compiler bug delayed nil pointer checks
→ candidate | edges
Programs that "previously executed incorrectly will now correctly panic with a
nil-pointer exception". This is the version-diff genre in its purest form — the
same program behaves differently on 1.24 and 1.25 — so it belongs to the E2
sweep's toolchain-pair tasks rather than to a single-toolchain task.

## The Green Tea collector arrives as an experiment
→ no action
`GOEXPERIMENT=greenteagc` here, default in 1.26, where it is already a candidate.

## testing.AllocsPerRun now panics under parallel tests
→ no action
Nothing in this repo calls it yet, so there is nothing to invalidate — but it
constrains every future allocation task: a predict task that measures with
`AllocsPerRun` must not call `t.Parallel()`. Worth remembering when the `alloc`
track grows.

## Unhandled panic output now reads "[recovered, repanicked]"
→ no action
Observable, but a solver could only be wrong about the wording.

## DWARF5 debug info, linker -funcalign, VMA names on Linux
→ no action
Toolchain and observability changes with no behaviour a task could grade.

## sync.WaitGroup.Go, testing T.Attr/T.Output, crypto.MessageSigner
→ no action
API conveniences.
