# Drill: 06-unlock-on-panic

- **Category:** C5 — Resource leaks (see `../../../rubric/bug-taxonomy.md`)
- **Tier:** subtle

## PR description

> `Registry.Top` returns the key with the highest counter, holding the registry mutex while it
> scans.

## Files to review

- `drill.go`

## Your task

Review against `../../../rubric/review-rubric.md`; write findings in the
`../../../rubric/submission-format.md` format. **Do not run `undergo reveal review-resource-leaks/06-unlock-on-panic` until you submit.**
