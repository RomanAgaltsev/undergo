# Drill: 10-pool-close-on-cancel

- **Category:** C1 — Concurrency & races (see `../../../../rubric/bug-taxonomy.md`)
- **Tier:** subtle

## PR description

> `Pool` runs jobs on a fixed set of workers. Cancelling its context shuts it
> down cleanly: queued jobs finish, the workers exit, and from then on `Submit`
> returns `ErrClosed` instead of blocking.

## Files to review

- `drill.go`

## Your task

Review against `../../../../rubric/review-rubric.md`; write findings in the
`../../../../rubric/submission-format.md` format. **Do not run `undergo reveal review/concurrency/10-pool-close-on-cancel` until you submit.**
