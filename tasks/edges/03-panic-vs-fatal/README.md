# 03 — What recover cannot catch

Every program in `testdata/` has the same shape: a deferred `recover` that prints
`RECOVERED` if it catches something. Then each one goes wrong differently.

Predict which are recovered.

| slot | what the program does |
|---|---|
| `ordinary_recovered` | `panic("ordinary")` |
| `child_goroutine_recovered` | panics inside a goroutine, recovers in `main` |
| `concurrent_map_recovered` | reads and writes one map from two goroutines |
| `nil_map_write_recovered` | writes to a nil map |
| `deadlock_recovered` | receives from a channel nothing will send on |
| `concurrent_map_is_fatal` | does the runtime say `fatal error` for the map case? |

```
undergo verify edges/03-panic-vs-fatal
```

Two are recovered. The pair worth staring at is the two map failures, which are
both "misusing a map" and do not end the same way at all.

Each program runs as its own process, because two of these take the process down
and a test binary cannot host that.

## Questions to answer in writing

1. State the rule dividing recoverable panics from fatal errors. Your rule should
   predict both map cases without special-casing them.
2. A panic in a child goroutine is not caught by a `recover` in `main`, even
   though `main` is still running. Say where `recover` looks, and why that is the
   only sensible design.
3. Concurrent map access is fatal rather than a panic. Say what the runtime is
   protecting and why letting the program continue would be worse than killing it.
4. Name the standard-library type that exists because of question 3, and say when
   it is the wrong answer.
5. A deadlock is detected and reported. Say what the runtime checked to know, and
   name a deadlock it cannot detect.
