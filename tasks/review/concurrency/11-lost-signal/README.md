# Drill: 11-lost-signal

- **Category:** C1 — Concurrency & races (see `../../../../rubric/bug-taxonomy.md`)
- **Tier:** subtle

## PR description

> Startup ordering: the server calls `gate.Wait()` before it accepts
> connections, and the config loader calls `gate.Open()` once the configuration
> is in. A `sync.Cond` lets the server sleep instead of polling a flag.

## Files to review

- `drill.go`

## Your task

Review against `../../../../rubric/review-rubric.md`; write findings in the
`../../../../rubric/submission-format.md` format. **Do not run `undergo reveal review/concurrency/11-lost-signal` until you submit.**
