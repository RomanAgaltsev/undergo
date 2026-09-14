# 05 — Publishing a half-built object

`publish.go` holds five ways of lazily building one shared `Config` and handing
it to every goroutine that asks:

| strategy | how it publishes |
|---|---|
| `PlainPointer` | `if s.p == nil { s.p = New() }` — no synchronisation |
| `AtomicPointer` | `atomic.Pointer` with a compare-and-swap |
| `OnceGuarded` | `sync.Once` |
| `MutexGuarded` | mutex held across the check **and** the build |
| `DoubleChecked` | unlocked read first, then the lock, then check again |

Eight goroutines hammer each one from a cold start, four hundred times over. The
test asks a single question of each: **did callers ever receive more than one
distinct `Config`?** — that is, did the object get built twice?

| slot | strategy |
|---|---|
| `plain_pointer_duplicates` | `PlainPointer` |
| `atomic_pointer_duplicates` | `AtomicPointer` |
| `once_duplicates` | `OnceGuarded` |
| `mutex_duplicates` | `MutexGuarded` |
| `double_checked_duplicates` | `DoubleChecked` |

Two of the five strategies are incorrect. Only one of them is caught here.
Predict all five, then read question 3.

```
undergo verify memmodel/05-publication
```

## Questions to answer in writing

1. Name the two incorrect strategies. For each, describe the interleaving that
   breaks it — two goroutines, step by step.
2. `PlainPointer` builds the object twice and every caller still gets a
   *complete* one. Say why the "partially constructed object" everybody warns
   about does not appear in this code, and what would have to be different for
   it to.
3. One incorrect strategy produced no wrong answer at all in four hundred
   attempts. Run the package under `-race` and find out what the detector says
   about it. Then say what that tells you about testing this class of bug.
4. `DoubleChecked` reads the pointer without the lock and then re-checks under
   it. Say exactly which of the two reads is the problem, and why adding a lock
   afterwards does not repair it.
5. `sync.Once` and `atomic.Pointer` are both correct here. Say when you would
   choose each, on grounds other than correctness — think about what happens
   when the constructor is expensive, and when it can fail.
