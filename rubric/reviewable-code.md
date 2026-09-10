# Writing reviewable code

The other half of code review: making *your* changes easy to review well. Apply this to your
own PRs on every project (strider, keystone, streamcache, …). A reviewer who can hold your
change in their head finds more real problems and invents fewer imaginary ones.

## Before you open the PR

- **One logical change per PR.** If you can describe it with an "and", split it.
- **Keep the diff small.** Big diffs get rubber-stamped; small diffs get read.
- **Atomic commits, Conventional Commits messages.** Each commit builds and tells one step of
  the story.
- **Self-review first.** Read your own diff as a hostile stranger would — run this repo's
  `review-rubric.md` over it. Most of your reviewer's future comments, you can catch here.
- **Tests that fail before the fix and pass after.** A test that passes with the bug still in
  is not a test.

## In the PR description

- **Intent** — what problem this solves (link the issue).
- **Approach** — the shape of the change and any alternative you rejected.
- **How verified** — the exact commands you ran and what you observed (not "it works").

## During review

- Answer the *why*, not just the *what*. If a reviewer asks, the code or a comment should
  already have said it — fix the code, don't just reply in the thread.
- Treat every comment as signal, even the wrong ones: if a reviewer misread it, the code was
  probably unclear.
