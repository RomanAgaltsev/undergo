# 04 — Which stores are barriered

Go's collector runs concurrently with your program, which means the program can
move pointers around while the collector is still deciding what is reachable. A
**write barrier** is the small piece of code the compiler inserts around a pointer
store so the collector does not lose track.

Not every store gets one. `barrier.go` has five. Predict which.

| slot | the store |
|---|---|
| `pointer_store_barriered` | `n.Next = next` |
| `scalar_store_barriered` | `n.N = v`, an int |
| `nil_store_barriered` | `n.Next = nil` |
| `local_store_barriered` | the same pointer store into a `Node` that never escapes |
| `slice_store_barriered` | `ns[i] = v`, into a `[]*Node` |

Two of these are false. One of the three that is true surprises most people.

```
undergo verify gc/04-write-barrier
```

The test compiles this package with `-gcflags=-S` and looks for the barrier in
the compiler's own assembly listing, so the answer comes from the compiler rather
than from anyone's reasoning.

## Questions to answer in writing

1. Say what a write barrier is for, in terms of **what the collector would
   otherwise miss**. "Because concurrency" is not an answer; name the specific
   thing that can go wrong.
2. Storing `nil` is still barriered. Explain why, in terms of what the barrier
   records — and note that this rules out one obvious guess about what it is for.
3. The scalar store needs no barrier. State the rule in one sentence.
4. The non-escaping store needs none either. Is that the same reason as the
   scalar, or a different one? Say which, precisely.
5. Barriers are not free, and they are only needed while a collection is running.
   Say what the runtime does about that, and what the cost is the rest of the
   time.
