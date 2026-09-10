# Drill: 02-slice-alias

- **Category:** C2 — Nil & memory safety (see `../../../rubric/bug-taxonomy.md`)
- **Tier:** subtle

## PR description

> Adds framing helpers: `Split` cuts a byte buffer into a 4-byte header and the rest, and
> `Frame` appends a checksum to the header. `CopyHeader` returns a standalone header copy.

## Files to review

- `drill.go`

## Your task

Review against `../../../rubric/review-rubric.md`; write findings in the
`../../../rubric/submission-format.md` format. **Do not run `undergo reveal review-nil-safety/02-slice-alias` until you submit.**
