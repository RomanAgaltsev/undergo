# Review rubric

Walk this top to bottom on every review. It is the same lens whether you are grading a
loupe drill or self-reviewing your own PR before you request review. For each dimension,
ask "does this change get it right, and what would break it?"

## 1. Correctness & edge cases
Does it do what the PR claims? Check boundaries (empty, one, many, max), zero values, integer
overflow/conversion, off-by-one, and the unhappy paths. → taxonomy **C9**.

## 2. Concurrency safety
Shared state guarded? Any data race, unsynchronized map, or goroutine that never exits? Channel
ownership and close discipline clear? → **C1**.

## 3. Error handling & propagation
Every error checked, wrapped with `%w` when the chain matters, and either handled once or
returned — never both. `errors.Is/As` used instead of `==` on wrapped errors. → **C3**.

## 4. Resource lifecycle
Every opened thing closed on every path (bodies, rows, files, tickers, contexts). No `defer`
inside a loop. Cancels actually called. → **C5**.

## 5. Context & cancellation
`context.Context` accepted, propagated, and honored. Deadlines/timeouts where I/O happens. No
work continues after cancel. → **C4**.

## 6. API & interface design
Small, honest interfaces defined at the consumer. Right receiver (pointer vs value). No leaky
abstraction or accidental breaking change to an exported signature. → **C6**.

## 7. Performance
No N+1 queries, needless allocations in hot paths, unbounded buffers, or missing
preallocation. Complexity appropriate to the input size. → **C7**.

## 8. Security
Untrusted input validated at the boundary. No SQL/command injection, weak crypto, or secrets
in logs. → **C8**.

## 9. Readability & naming
Names say what they mean; functions do one thing; comments explain *why*, not *what*. Could a
newcomer follow it? → cross-cutting.

## 10. Test adequacy
Tests assert real behavior, cover the changed branch and its error path, and aren't
flaky/time-dependent. A test that would still pass with the bug present is worthless. → **C10**.

## 11. Change hygiene
One logical change, small diff, atomic commits, a description that states intent and how it
was verified. → see `reviewable-code.md`.
