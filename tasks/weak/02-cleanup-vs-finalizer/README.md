# 02 — AddCleanup against SetFinalizer

Go has two ways to run code when an object is collected. `runtime.SetFinalizer`
has been there since the beginning; `runtime.AddCleanup` arrived in Go 1.24 to
replace it. Six experiments in `lifetime.go` run the same shape against each.

Read the code, then predict.

| slot | question |
|---|---|
| `cleanup_gcs_needed` | collections before an `AddCleanup` runs |
| `finalizer_gcs_needed` | collections before a `SetFinalizer` runs |
| `cleanup_self_reference_panics` | registering a cleanup whose *argument* is the object — does it panic? |
| `finalizer_runs_on_self_reference` | a finalizer whose closure *captures* the object — does it ever run? |
| `finalizer_can_resurrect` | can a finalizer make its object reachable again? |
| `multiple_cleanups_all_run` | three cleanups on one object — do all three run? |

Rows three and four are the same mistake made two ways. They do not have the
same answer, and the difference is most of why `AddCleanup` exists.

```
undergo verify weak/02-cleanup-vs-finalizer
```

The first two have a widely-repeated answer. Check it.

## Questions to answer in writing

1. Enumerate the differences between `AddCleanup` and `SetFinalizer`. Put the one
   that most often causes a real bug first, and say what the bug looks like in
   production.
2. A finalizer that keeps its own object reachable never runs, and the object is
   never freed — a leak that no tool will point at. `AddCleanup` refuses the
   equivalent. Say what it is about the *signature* of `AddCleanup` that lets the
   runtime catch this, and why `SetFinalizer` cannot.
3. `AddCleanup` takes the cleanup argument separately from the object rather than
   handing the cleanup function the object itself. Name two things that
   separation buys beyond catching the self-reference.
4. Both mechanisms run at a time nobody controls. If you need a file descriptor
   closed, name what you should use instead — and then say why these exist at all,
   given that answer.
5. Since Go 1.25 cleanups run concurrently. Every function in `lifetime.go` waits
   on a channel rather than reading a flag after `runtime.GC()`. Say what would
   go wrong with the flag version, and how often.
