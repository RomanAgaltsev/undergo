# Drill: 08-decode

- **Category:** C2 — Nil & memory safety (see `../../../rubric/bug-taxonomy.md`)
- **Tier:** multi-bug

## PR description

> `Decode` converts a loosely-typed `map[string]any` payload into a `map[string]string`, pulling
> out an id and a nested user. A `safeString` helper is included. Several nil/type-safety
> problems are hiding here — find as many as you can.

## Files to review

- `drill.go`

## Your task

Review against `../../../rubric/review-rubric.md`; write findings in the
`../../../rubric/submission-format.md` format (one line per finding). **Do not run `undergo reveal review-nil-safety/08-decode` until you submit.**
