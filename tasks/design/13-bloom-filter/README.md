# Kata k13 — Bloom filter

**Track:** B — Storage & indexing · **Tier:** T1

## Problem
Design a probabilistic set-membership structure that answers "definitely not present" or
"probably present" for a target false-positive rate, in a small fixed amount of memory.

## Requirements
**Functional**
- Add a key.
- Test membership.
- Size the filter from an expected element count n and a target false-positive rate.

**Non-functional**
- No false negatives — a key that was added always tests present.
- The measured false-positive rate stays near the target.
- Memory is fixed by (n, target FPR).

## Given
- Up to ~10M keys; target false-positive rate ≈ 1%.
- Single process.

## Your core must
Expose a `Filter` with:
- `Add(key []byte)` — record a key.
- `Test(key []byte) bool` — false ⇒ definitely absent; true ⇒ probably present.

Use a bit array of m bits and k hash functions, with m and k derived from n and the
target FPR. **Tests must show:** no false negatives (every added key tests true); the
measured FPR over random non-members stays within the configured bound; the sizing math
(m, k) matches the standard formulas.

## Design should also address (in DESIGN.md)
- The sizing formulas: m = −n·ln p / (ln 2)² and k = (m/n)·ln 2.
- Double hashing (Kirsch–Mitzenmacher) to derive k hashes from two base hashes.
- Counting or scalable bloom filters for deletion and growth, and their cost.
- Where a bloom filter sits on an LSM read path to skip SSTables (tie-in to k14).

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
