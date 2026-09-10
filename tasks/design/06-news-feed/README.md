# Kata k06 — News feed

**Track:** F — Real-time & fan-out · **Tier:** T2

## Problem
Design a social news feed: each user follows others, and their home timeline is the
recent posts of the people they follow, newest first. The core tension is **fan-out on
write** (push each post into followers' timelines) vs **fan-out on read** (assemble the
timeline by pulling from followees at read time) — the heart of Twitter/Instagram feeds.

## Requirements
**Functional**
- Follow / unfollow another user.
- Post a message as a user.
- Read a user's home timeline: recent posts from everyone they follow, newest first.

**Non-functional**
- A read (timeline) is cheap for the common case.
- A post from a normal user reaches followers' timelines promptly.
- A user with a huge follower count (celebrity) does not make posting pathologically slow.

## Given
- Millions of users; a skewed follow graph (a few celebrities, a long tail of normal users).
- Read-heavy: far more timeline reads than posts.

## Your core must
Expose a `Feed` with:
- `Follow(follower, followee string)` — establish a follow edge (and an unfollow counterpart).
- `Post(user string, body string) Post` — publish a post (timestamped).
- `Timeline(user string, limit int) []Post` — the user's home timeline, newest first.

Implement a **hybrid** fan-out (push for normal authors, pull for celebrities). **Tests
must show:** a timeline contains exactly the posts of followed users, newest first;
the write-fanout path and the read-fanout path yield the same timeline for a given user;
unfollowing removes an author's future posts from the timeline.

## Design should also address (in DESIGN.md)
- Fan-out-on-write vs on-read: the write-amplification vs read-cost tradeoff, and the
  hybrid threshold (celebrity cutoff).
- Timeline storage & trimming (bounded per-user timeline cache).
- Ordering & merging posts from multiple authors (event-time vs arrival).
- Backfill when a user follows someone new; eventual consistency of the timeline.

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
