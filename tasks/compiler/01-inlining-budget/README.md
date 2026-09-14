# 01 — What fits the inliner

The inliner has a **cost model**: it walks a function, adds up what it finds, and
inlines if the total stays under a budget. It will tell you both numbers.

Six functions in `budget.go`. Predict.

| slot | question |
|---|---|
| `inline_budget` | the budget, as the compiler reports it |
| `trivial_cost` | the cost of `Trivial` |
| `loop_cost` | the cost of `Loop` |
| `calls_inlinable_cost` | the cost of `CallsInlinable`, which calls `Trivial` three times |
| `sprawling_cost` | the cost of `Sprawling`, which does not fit |
| `with_defer_inlines` | does `WithDefer` inline? |
| `with_go_inlines` | does `WithGo` inline? |

```
undergo verify compiler/01-inlining-budget
```

`CallsInlinable` is the one to think about hardest: three calls, and the question
is what those calls cost when the callee is itself inlinable. Compare your answer
against `loop_cost` before you commit to it.

## Questions to answer in writing

1. What does the cost number count? Name three things that add to it and one that
   does not.
2. `WithDefer` and `WithGo` both refuse, and the compiler gives the same *shape*
   of reason for each. Quote both, and say why these two constructs in particular
   defeat the inliner.
3. Say whether the cost of an inlinable callee counts toward its caller's budget,
   and justify it from `calls_inlinable_cost`. Then say what that means for a
   chain of four small wrappers.
4. `Sprawling` misses by a small margin. Say what you would change to bring it
   under, and whether you should.
5. `//go:noinline` exists, and so does the budget. Name a case where you would
   reach for the pragma deliberately.
