# Drill: 03-ticker-rows

- **Category:** C5 — Resource leaks (see `../../../../rubric/bug-taxonomy.md`)
- **Tier:** multi-bug

## PR description

> Adds a `Poller` that counts rows on a ticker in a background goroutine and hands back a
> stop function. There's a DB query per tick and a prepared-statement helper. Several
> resources aren't managed correctly — find as many as you can.

## Files to review

- `drill.go`

## Your task

Review against `../../../../rubric/review-rubric.md`; write findings in the
`../../../../rubric/submission-format.md` format (one line per finding). **Do not run `undergo reveal review/resource-leaks/03-ticker-rows` until you submit.**
