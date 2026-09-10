# Kata k34 — URL shortener

**Track:** H — Application systems · **Tier:** T1

## Problem
Design a URL shortener: map a long URL to a short code and resolve that code back to the
original URL — the core of bit.ly / tinyurl. Codes must be short, unique, and cheap to
resolve, with optional user-chosen custom aliases.

## Requirements
**Functional**
- Shorten a long URL into a short code.
- Resolve a code back to its original URL.
- Support a caller-supplied custom alias.

**Non-functional**
- Codes are short and collision-free.
- Resolve is O(1) and read-heavy-friendly (far more resolves than shortens).
- A custom alias that is already taken is rejected, not silently overwritten.

## Given
- Billions of URLs over the system's life; the keyspace must not run out for short codes.
- Read-heavy: resolves vastly outnumber shortens.

## Your core must
Expose a `Shortener` with:
- `Shorten(url string) (code string, err error)` — mint a short code for a URL.
- `Resolve(code string) (url string, ok bool)` — the original URL for a code.
- A way to shorten with a **custom alias** (rejected on collision).

Generate codes as **base62 of a monotonic counter** (not a hash). **Tests must show:**
shorten→resolve round-trips to the original URL; two different URLs get distinct codes;
a custom alias that collides with an existing code is rejected.

## Design should also address (in DESIGN.md)
- Base62-of-a-counter vs hash-of-URL: keyspace exhaustion, code length, and collision handling.
- Whether identical URLs collapse to one code (dedup) or each shorten mints a new code.
- Sharding the counter across nodes without collisions (cf. k2 Snowflake / a per-node range).
- Custom-alias namespace sharing the same code space; reserved words; expiry/TTL as an extension.

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
