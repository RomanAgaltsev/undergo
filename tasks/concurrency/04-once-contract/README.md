# 04 — What Once promises when f panics

`sync.Once` guarantees `f` runs once. Everybody knows that. Nobody reads what
happens when that one run panics.

There are two families here and they do not agree with each other.

| slot | the situation |
|---|---|
| `do_panic_reaches_caller` | `once.Do(f)` where `f` panics. Does the panic reach the caller of `Do`? |
| `do_second_call_runs_f` | a second `once.Do(g)` afterwards, with a different `g`. Does `g` run? |
| `do_second_call_panics` | does that second `Do` panic? |
| `oncefunc_second_call_panics` | `fn := sync.OnceFunc(f)` where `f` panics. `fn()` panicked once. Does the **second** `fn()` panic? |
| `oncevalue_calls_f_twice` | `vf := sync.OnceValue(f)` where `f` panics. After two calls to `vf()`, did `f` run twice? |

Three of these five follow from "f runs once". Two do not, and they are in the
same package as the three that do.

```
undergo verify concurrency/04-once-contract
```

## Questions to answer in writing

1. State the `Do` contract for a panicking `f` in one sentence, as the standard
   library states it. Note the word it uses instead of "completed".
2. `sync.Once` is not just a bool. Say what it holds, why the fast path is an
   atomic load rather than a mutex acquisition, and what `Do` does with the
   mutex on the slow path.
3. `OnceFunc` and `OnceValue` behave differently from `Do` on the second call.
   Give the reason, in terms of what each one owes its caller.
4. A lazy initialiser calls a remote service and panics on failure. Say what
   your program does on the second request under `Do`, and under `OnceValue`,
   and which of the two you would rather debug at 3am.
5. `sync.Once` may not be copied after first use. Say what goes wrong if it is,
   and name the tool that catches it.
