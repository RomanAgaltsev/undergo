# 05 — How many collections a pooled object survives

`sync.Pool`'s documentation is unusually blunt about its own guarantees:

> *Any item stored in the Pool may be removed automatically at any time without
> notification.*

"At any time" is the contract. It is not the behaviour. The behaviour is
specific, it changed once, and it is measurable.

`pool.go` puts one object into a pool and then forces collections, counting
how many times the pool has to call `New`. `GOMAXPROCS` is pinned to 1, because
a pool's storage is per-P and the question is otherwise not well posed.

| slot | question |
|---|---|
| `new_calls_before_any_gc` | how many times `New` ran to serve the first `Get` from an empty pool |
| `survives_one_gc` | after `Put`, one `runtime.GC()`, then `Get` — did the pool avoid calling `New`? |
| `survives_two_gc` | after `Put`, **two** collections, then `Get` — did it? |
| `new_calls_after_two_gc` | how many times `New` ran to serve that last `Get` |

Two of the four are the same question asked once and twice. The difference
between the answers is the whole task.

```
undergo verify concurrency/05-pool-clearing
```

## Questions to answer in writing

1. Name the mechanism that makes `survives_one_gc` and `survives_two_gc`
   disagree, and the release that introduced it.
2. Describe what the pool looked like before that release, and say concretely
   what a service that pooled expensive objects paid for on every collection.
3. A `Pool` has per-P storage with a private slot and a shared queue. Say why
   this test pins `GOMAXPROCS`, and what the measurement would mean without it.
4. Given the documented contract, is a design that *depends* on an object
   surviving one collection correct? Answer for a buffer pool and for a pool of
   objects whose construction requires a network call, and say why the answer
   differs.
5. Name the two ways to misuse `Put`. One is about what you keep, the other is
   about what you hand over.
6. `sync.Pool` is not a cache and the documentation says so. State the one
   property a cache has that `Pool` deliberately does not.
