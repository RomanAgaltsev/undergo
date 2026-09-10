# Drill: 04-closed-channel-send

- **Category:** C1 — Concurrency & races (see `../../../rubric/bug-taxonomy.md`)
- **Tier:** obvious

## PR description

> `Broadcast` streams each value to a channel, closes it when done, and for an empty input
> sends a `-1` sentinel so the consumer knows there was nothing to process.

## Files to review

- `drill.go`

## Your task

Review against `../../../rubric/review-rubric.md`; write findings in the
`../../../rubric/submission-format.md` format. **Do not run `undergo reveal review-concurrency/04-closed-channel-send` until you submit.**
