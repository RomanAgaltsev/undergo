# Submission format & scoring

## How to submit a review

Review the drill cold — read its `README.md` and code, **do not open `ANSWERS.md`** — then
write your findings, **one per line**, in this exact format:

```
path:line · <category> · <severity> · <one sentence: why it's a defect + the fix>
```

- **path:line** — the file and line the defect anchors to (e.g. `drill.go:42`).
- **category** — one of the taxonomy IDs `C1`–`C10` (see `bug-taxonomy.md`).
- **severity** — one of:
  - **blocker** — must fix before merge (data race, security hole, data loss).
  - **major** — should fix (a real bug, wrong behavior on a path).
  - **minor** — nice to fix (a latent issue, poor edge handling).
  - **nit** — style/preference, non-blocking.
- **why + fix** — one sentence: what's wrong and what you'd change.

Example:
```
drill.go:31 · C1 · blocker · shared counts map is written from goroutines with no mutex — guard with sync.Mutex or use sync/atomic
drill.go:58 · C5 · major · resp.Body is never closed — add `defer resp.Body.Close()` after the error check
```

Your **submission** is the full set of findings for one drill. Submit it, then ask for grading.

## Scoring

Grading (by Claude, against the drill's sealed `ANSWERS.md` + this rubric) reports:

- **Recall** = planted defects you found / total planted. *(Did you catch the bugs?)*
- **Precision** = true findings / all findings you submitted. *(Did you over-flag?)*
- **Severity accuracy** = fraction of found defects whose severity you rated correctly.
- **Rubric coverage** = which of the 11 rubric dimensions your review visibly considered.
- **Decoys** = how many planted decoys (things that look wrong but aren't) you correctly left alone.

Record the numbers in `log/REVIEW-LOG.md`. The point is the **trend**: recall and precision
climbing over sessions is the measurable proof your review skill is improving.

## Scoring a clean drill

A **clean drill** (`drills/clean/…`) has **zero planted defects** — it is correct code that
looks suspicious. Scoring inverts:

- **Recall** is N/A (nothing to find).
- **Precision is the whole game.** Every real-defect finding you submit against a clean drill
  is a **false positive** and drives precision toward 0.
- **Ideal submission:** *"No defects."* Nits (style/preference) are acceptable and don't count
  against precision; a **major/blocker** finding on clean code is the failure this tier trains
  out. For full credit, name each red herring and say *why it is correct* — recognizing why
  suspicious code is fine is the skill.

**Blind mode (true precision test):** the `clean/` folder name telegraphs the answer. For an
honest test, have the grader hand you a drill *by number from either the category tiers or the
clean tier* without telling you which — review it before you learn the tier, then check. Log
precision as usual; a clean drill you correctly leave alone is precision `1.0`.
