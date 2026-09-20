# versions — the go.mod line as a behaviour switch

Every task here compiles the same source twice with the **same toolchain** and
gets two different answers. The only thing that changes is the `go` line in
`go.mod`.

That is possible because Go 1.21 made the go line the compatibility mechanism:
GODEBUG defaults, several compiler semantics, and language legality are all
resolved from the main module's `go` directive. It is why Go can change what a
program does without breaking the Go 1 promise — your program's behaviour
changes when *you* edit `go.mod`, and not before.

Three kinds of change are reachable this way:

| kind | example |
|---|---|
| a GODEBUG default keyed to the go line | `panic(nil)` — task 01 |
| compiler semantics keyed to the go line | loop variables — task 02 |
| language legality keyed to the go line | `for range 3` — task 03 |

And one kind is **not**: a GODEBUG that has since been removed. Task 05 is about
that edge, and it is the most useful thing in the track.

## The one hard rule

A task here may only name go lines **at or below** the toolchain it runs on.
A higher line is not a version diff — it is a toolchain download waiting to
happen. `internal/goline` enforces this and sets `GOTOOLCHAIN=local`, so the
attempt fails loudly instead of quietly reaching for the network.

Nothing in this track needs the network. Every program is a few lines of
standard library.

## The second half: toolchain pairs

Tasks 06 and up change the **toolchain** instead of the go line. `GOTOOLCHAIN`
names a version and the go command fetches it, so behaviour that no go line can
select — a GODEBUG that has been removed, a package that did not exist, a
compiler that made a different choice — is reachable after all.

These tasks declare `requires.toolchains` and need the network **once** per
toolchain; afterwards it is in the module cache. Where it cannot be fetched they
skip with an explanation, and `undergo doctor` says which ones resolve. CI runs
`ci-verify --require-toolchains`, so there they are proven or the gate is red.

The two halves answer different questions. The go line asks *what does my module
declare*; the toolchain asks *what does my compiler have*.
