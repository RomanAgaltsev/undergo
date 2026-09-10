# Kata k35 — Idempotent ledger

**Track:** H — Application systems · **Tier:** T2

## Problem
Design a financial ledger that records money movements safely under retries: every
transaction is double-entry (debits equal credits) and carries an idempotency key, so a
client that retries a post never double-applies it — the core of a payments ledger.

## Requirements
**Functional**
- Post a transaction that moves money between accounts (debit one, credit another).
- Query an account's balance.
- A retried post with the same idempotency key is a no-op (applied at most once).

**Non-functional**
- Per-transaction invariant: debits equal credits (money is conserved).
- No double-apply under at-least-once/retrying callers.
- Overdraft (negative balance) is rejected; balances never silently drift.
- Concurrency-safe under concurrent posts to the same account.

## Given
- A high transaction rate; clients retry on timeout (at-least-once).
- Balances must be exactly reconstructable from the posted transactions.

## Your core must
Expose a `Ledger` with:
- `Post(txn Txn) error` — where `Txn` carries an **idempotency key** and its double-entry lines.
- `Balance(acct string) int64` — an account's current balance.

Enforce double-entry + idempotency-key dedup. **Tests must show:** posting the same
idempotency key twice applies it exactly once (the retry is a no-op); each transaction's
debits equal its credits; a transaction that would overdraw an account is rejected and
leaves all balances unchanged.

## Design should also address (in DESIGN.md)
- The idempotency-key store: where it lives, its retention/TTL, and dedup cost (cf. k19).
- Enforcing the double-entry invariant and atomic (all-or-nothing) multi-line posts.
- Running-balance vs sum-of-entries for `Balance` — the read/write and audit tradeoff.
- Isolation under concurrent posts to one account; where a real system would need a lock/txn.

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
