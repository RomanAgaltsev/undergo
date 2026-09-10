# Drill: 03-lookup-chain

- **Category:** C2 — Nil & memory safety (see `../../../rubric/bug-taxonomy.md`)
- **Tier:** multi-bug

## PR description

> Adds a `Registry` of users by id and a `Greeting` that formats a looked-up user's name
> together with a display name pulled from an `any` payload. `Tag` records tag hits; there's
> a `DisplayName` helper too. Several unsafe assumptions crept in — find as many as you can.

## Files to review

- `drill.go`

## Your task

Review against `../../../rubric/review-rubric.md`; write findings in the
`../../../rubric/submission-format.md` format (one line per finding). **Do not run `undergo reveal review-nil-safety/03-lookup-chain` until you submit.**
