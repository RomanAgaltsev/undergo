# 07 — When a context stops being cancellable

Every `context.Context` has a `Done()` channel — except the ones that do not.
Go 1.21 added a way to make more of them, and two functions that make `Err()`
and `Cause()` deliberately disagree.

| slot | question |
|---|---|
| `without_cancel_done_is_nil` | is `context.WithoutCancel(parent).Done()` nil? |
| `without_cancel_keeps_values` | does the detached context still see a value set on the parent? |
| `without_cancel_err_after_parent_cancel` | the parent is cancelled. Is the detached context's `Err()` still nil? |
| `detached_child_is_uncancellable` | `WithCancel` derived **from** the detached context, then cancelled. Is its `Err()` nil? |
| `cause_of_plain_cancel_is_canceled` | after an ordinary `cancel()`, does `context.Cause(ctx)` match `context.Canceled`? |
| `cause_of_live_context_is_canceled` | on a context nobody has cancelled, does `Cause(ctx)` match `context.Canceled`? |
| `err_when_cause_differs` | after `WithCancelCause(...)` and `cancel(errBecause)`, is `Err()` still `Canceled` while `Cause()` is `errBecause`? |
| `afterfunc_runs_on_already_cancelled_ctx` | `AfterFunc` registered on a context that was cancelled *before* the call. Does it run? |

Six of the eight are `true`. Finding the two that are not is the task, and both
are places where a reasonable mental model gives the wrong answer.

```
undergo verify concurrency/07-cancellability
```

## Questions to answer in writing

1. State, in one sentence, what it means for a context to be cancellable — in
   terms of what its `Done()` returns. List the three kinds of context in the
   standard library for which that is false.
2. A receive on a nil channel blocks forever. Given that, explain why
   `select { case <-ctx.Done(): ... }` on a `WithoutCancel` context is not a bug.
3. `WithoutCancel` detaches from *the parent*. Say precisely what it does not
   detach, and use that to explain slot four.
4. `Err()` and `Cause()` can return different errors for the same context, by
   design. Give the reason each exists and say which one belongs in a
   `errors.Is` check in library code.
5. Name the real situation `WithoutCancel` was added for — the shape where you
   must keep working after the caller has gone — and say what people wrote
   before Go 1.21 to get the same effect, and what that version got wrong.
6. Slot eight has a practical consequence for cleanup code. State it as a rule
   about when it is safe to register an `AfterFunc`.
