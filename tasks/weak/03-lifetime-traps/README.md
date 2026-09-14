# 03 — What checkfinalizers will and will not tell you

Go 1.25 added `GODEBUG=checkfinalizers=1`, which inspects finalizers and cleanups
during collection and reports the ones that can never work. Four programs, each
run under that setting. Predict whether the runtime complains about each.

| slot | the program |
|---|---|
| `trap1_diagnosed` | a finalizer whose closure captures the object it finalizes |
| `trap2_diagnosed` | a package-level variable also holds the object, so the finalizer can never run |
| `trap3_diagnosed` | a cleanup whose argument is not the object but points at it |
| `trap4_diagnosed` | a correct cleanup — the control |

Two of these are diagnosed and two are not. The pair that is *not* is the
interesting half: one of them is correct code, and the other is a genuine leak
that runs forever without a word.

```
undergo verify weak/03-lifetime-traps
```

When the runtime does find something it does not print a warning and carry on —
it throws. Each program therefore runs as its own process, which is why the
snippets live in `testdata/` rather than in this package.

## Questions to answer in writing

1. Two programs are diagnosed. State the single property they share, in one
   sentence, precisely enough to predict a fifth case from it.
2. `trap2` leaks exactly as surely as `trap1` — the finalizer never runs and the
   object never dies. Say why the runtime can see one and not the other.
3. From your answer to 1, say what this tool is *for* — and name the class of
   finalizer bug you must still find by other means.
4. `trap3` passes the check `AddCleanup` performs at registration time, and is
   still wrong. Explain how it slipped through, and what a cleanup argument may
   safely contain.
5. The check runs during a collection and throws on failure. Say why it is a
   `GODEBUG` and not on by default, and where you would turn it on.
