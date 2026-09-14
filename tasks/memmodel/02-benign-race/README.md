# 02 — Three ways the harmless race is not

Everybody has heard these three defences, and usually said one of them:

> *"It's only a bool."*
> *"A 64-bit write is atomic anyway."*
> *"Worst case we read a stale value once."*

`benign.go` checks what actually happens.

| slot | question |
|---|---|
| `torn_uint64_seen` | did a reader ever see a `uint64` neither writer wrote? |
| `torn_struct_seen` | did a reader ever see a two-word struct assembled from two writes? |
| `spin_loop_exits` | does a loop spinning on a plain `bool` ever notice it changed? |

One of the three defences turns out to be true on this hardware. Work out which
before you run it — and be ready to say whether "true on this hardware" is the
same as "true".

```
undergo verify memmodel/02-benign-race
```

## Questions to answer in writing

1. One tearing experiment finds nothing and the other finds tearing every time.
   Say what makes the difference, and name a type or a platform where the
   *first* one would tear too.
2. `spin_loop_exits` may not be the answer you predicted. The compiler is
   **permitted** to load `done` once and reuse it forever — the program has a
   race, so the model promises nothing. Build `testdata/hoist` with
   `go build -gcflags=-S` and find out whether it did. Then say what your
   program's correctness is resting on.
3. Following from that: if the compiler's current choice is what makes your
   racy code work, what happens at the next toolchain upgrade? Say why "it
   passes our tests" is not evidence here.
4. "Worst case we read a stale value once" assumes staleness is bounded. Say
   whether the Go memory model bounds it, and quote the sentence that settles
   it.
5. Fix the spin loop with a one-word change. Say precisely what that word buys
   you — ordering, atomicity, or both — and which of the three defences it
   makes true.
