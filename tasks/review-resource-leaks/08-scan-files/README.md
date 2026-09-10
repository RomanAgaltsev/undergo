# Drill: 08-scan-files

- **Category:** C5 — Resource leaks (see `../../../rubric/bug-taxonomy.md`)
- **Tier:** multi-bug

## PR description

> `CountLines` opens each path, scans it line by line, guards against pathologically long lines,
> and sums the line counts. Several resource/handling problems are hiding here — find as many as
> you can.

## Files to review

- `drill.go`

## Your task

Review against `../../../rubric/review-rubric.md`; write findings in the
`../../../rubric/submission-format.md` format (one line per finding). **Do not run `undergo reveal review-resource-leaks/08-scan-files` until you submit.**
