# Drill: 05-insecure-tls

- **Category:** C8 — Security (see `../../../rubric/bug-taxonomy.md`)
- **Tier:** obvious

## PR description

> Adds `newClient`, an HTTP client with a sensible 10-second timeout and a custom TLS transport
> for the upstream API.

## Files to review

- `drill.go`

## Your task

Review against `../../../rubric/review-rubric.md`; write findings in the
`../../../rubric/submission-format.md` format. **Do not run `undergo reveal review-security/05-insecure-tls` until you submit.**
