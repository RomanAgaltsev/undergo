# 06 — Who gets in while a writer is waiting

`RWMutex` lets any number of readers hold the lock at once. A writer needs it
alone.

So: a reader **R1** holds the read lock. A writer **W** arrives and blocks. Now a
second reader **R2** arrives.

R2 wants a read lock. Read locks are shared. R1 is holding one. Does R2 get in?

| slot | question |
|---|---|
| `second_reader_parked` | does R2 block, or does it acquire immediately? |
| `second_reader_admitted_while_writer_pending` | did R2 record its acquisition before R1 released? |
| `acquisition_order` | the three names in the order they acquired, comma-separated, e.g. `R1,W,R2` |
| `writer_waits_for_first_reader` | did W have to wait for R1 at all? |

Whichever way you answer, the other way has an obvious failure mode. Work out
what both of them are before you run this — that pair is the real question, and
it is what the standard library had to choose between.

```
undergo verify concurrency/06-rwmutex-admission
```

## Questions to answer in writing

1. Suppose R2 *were* admitted. Describe the traffic pattern — realistic, not
   contrived — under which W is never scheduled at all, and name it.
2. Suppose R2 is not admitted. State what the lock is now costing readers, and
   in what sense a read lock has stopped being purely shared.
3. Go's `RWMutex` documentation contains a sentence warning about recursive read
   locking. Quote it, and connect it to your answer: construct the deadlock in
   three steps.
4. Name a workload where a plain `Mutex` is the correct choice even at 99%
   reads. Give the property of the critical section that decides it.
5. This test cannot use a `synctest` bubble. Say why — the answer is about what
   `synctest` counts as durably blocked, and it is worth knowing before you
   reach for a bubble elsewhere.
