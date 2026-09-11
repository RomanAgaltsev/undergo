# Drill: 09-combine

- **Category:** C2 — Nil & memory safety (see `../../../../rubric/bug-taxonomy.md`)
- **Tier:** multi-bug

## PR description

> `Combine` merges two int slices, runs a completion hook, and returns the first merged element.
> A `safeFirst` helper is included. Several nil/memory problems are hiding here — find as many as
> you can.

## Files to review

- `drill.go`

## Your task

Review against `../../../../rubric/review-rubric.md`; write findings in the
`../../../../rubric/submission-format.md` format (one line per finding). **Do not run `undergo reveal review/nil-safety/09-combine` until you submit.**
