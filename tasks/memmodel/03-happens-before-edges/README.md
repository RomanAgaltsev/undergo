# 03 — Which of these assertions can fail

`edges.go` holds nine two-goroutine programs. Every one of them does the same
thing: a goroutine writes `payload = 42`, and the main goroutine reads it back.
They differ in exactly one respect — **the act that separates the write from
the read**.

The test runs each program 500 times and asks whether the read saw 42 on *every*
run.

| slot | the separating act |
|---|---|
| `unbuffered_send_always` | send on an unbuffered channel |
| `buffered_full_always` | send on a capacity-1 channel that is already full |
| `close_always` | `close` |
| `mutex_always` | `Unlock` then `Lock` |
| `once_always` | both goroutines call `sync.Once.Do` |
| `waitgroup_always` | `Done` then return from `Wait` |
| `atomic_always` | `atomic.Bool` store then load |
| `sleep_always` | nothing — the reader sleeps a millisecond |
| `no_sync_always` | nothing — the reader looks immediately |

The last two have **identical synchronisation**: none. They differ only in how
long the reader waits. Predict both, and be ready to say what that pair proves.

```
undergo verify memmodel/03-happens-before-edges
```

## Questions to answer in writing

1. Seven of these programs are *guaranteed* by the Go memory model. Two are
   not. Name the two, and say which of them nevertheless held on all 500 runs.
2. Following from that: this task grades "held on every run". Say why that is a
   strictly weaker claim than "guaranteed", and give a program that would be
   graded correct here while still being broken.
3. `buffered_full_always` fills the channel before starting the writer, and the
   reader receives twice. Say exactly which edge that arrangement creates, and
   what a send on a *non-full* buffered channel would have guaranteed instead.
4. `close` and a send both release a waiting receiver. Say whether they
   establish the same edge, and cite the memory model for each.
5. `sleep_always` held 500 times out of 500. Say what would have to be true of
   the machine, the scheduler, or the compiler for it to fail — and say why a
   million passing runs would still not make it correct.
