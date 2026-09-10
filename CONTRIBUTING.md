# Contributing

## Adding a task

```
go run ./cmd/undergo new --id <track>/<NN>-<slug> --mode <mode> --title "..." --difficulty N
```

That scaffolds `tasks/<track>/<NN>-<slug>/` with a manifest, a README skeleton
and an `_solution/` directory. Then:

1. **Write the stub.** It must compile and vet clean; its tests must fail.
2. **Write the frozen tests.** These are the specification. For `predict`, the
   test computes the truth and hands it to `predict.Check`. For `optimize`, it
   calls `optimize.Check` with a baseline and a candidate benchmark.
3. **Write `_solution/`** — `HINT.md` (one nudge, naming the mechanism, not the
   answer), `EXPLANATION.md` (why the answer is what it is), `solution/` (files
   that overlay the task directory), and for a `predict` task a reference
   `prediction.yaml`.
4. **Seal it:** `go run ./cmd/undergo seal <id>`. This produces
   `solution.sealed` and deletes the plaintext. `_solution/` is gitignored, so
   it cannot be committed by accident.
5. **Prove it:** `task ci`.

## Rules

- **Ideas may be borrowed; code may not.** If a task was inspired by a talk, a
  release note, an issue or another repository, record it in `inspired_by:`.
  Never copy code from a source whose licence you have not checked.
- **Task IDs are permanent.** Renaming one invalidates every user's progress
  file, so a retired task gets `deprecated: true` rather than a new name.
- **Declare what your answer depends on.** If a `predict` answer differs by
  architecture or Go version, say so in `requires`. The harness will skip rather
  than mark a correct solver wrong — that is far better than a wrong grade.
- **Never weaken a frozen test** to make a solution pass.
- Conventional Commits. `task ci` must be green.

## What makes a good task

The test is whether a solver can be *vaguely right*. "Does this allocate?" is a
coin flip; "how many allocations, and which line causes each" is a task. Prefer
questions where the mechanism is the answer.

## Self-graded tasks

A `review` or `design` task has no machine grader, and the repo says so rather
than faking one. Its sealed blob carries `HINT.md` and `EXPLANATION.md` but no
`solution/`, CI gate 2 skips it and reports how many it skipped, and `go vet` is
never run over a review drill — a vet diagnostic on planted-defect code *is* the
answer, and CI logs are public.
