# 05 — How many collections until it is gone

Four questions about *when* the collector gets round to something.

Every measurement forces collections with `runtime.GC()` and checks after each
one. That is deliberate: counting the collector's own cycles under a workload is
not reproducible — `gc/02` exists because of it — but counting collections you
caused yourself is.

| slot | question |
|---|---|
| `weak_gcs` | collections until a dropped `weak.Pointer`'s target is gone |
| `finalizer_within_two` | did a `SetFinalizer` function run within two collections? |
| `cleanup_within_two` | did an `AddCleanup` function run within two collections? |
| `pinned_survives` | does a weak pointer survive a collection while its target is still referenced? |

Notice that three of these are counts and one is a **predicate**. That is not
an accident, and question 3 is about why.

```
undergo verify weak/05-single-cycle-reclaim
```

## Questions to answer in writing

1. The folklore says finalizers need "two GCs". Your measurement disagrees.
   Say what the two collections in that story actually are, and what each one
   does — then say which of them this task is counting.
2. `AddCleanup` was added in Go 1.24 to replace `SetFinalizer` for most uses.
   List three things it does differently, and say which of them is the reason
   `SetFinalizer` is now discouraged.
3. The cleanup slot is a predicate because the exact count measured 1 once and
   2 five times across six fresh processes. Say what is asynchronous here, and
   why that makes an exact count the wrong thing to ask for.
4. This task forces collections rather than waiting. Say what that buys, and —
   more importantly — what it **hides** about how these mechanisms behave in a
   program that never calls `runtime.GC`.
5. A `weak.Pointer` whose target is also referenced by a live `unique.Handle`:
   say what you expect, and how you would check it.
