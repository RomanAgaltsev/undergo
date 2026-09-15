# The release radar

The catalogue is not a fixed list. Every Go release gets a triage pass that asks
two questions — **what does this create**, and **what does this invalidate** —
and the answers land here as data.

## Format

One file per release, `versions/go1.NN.md`:

    # Go 1.NN

    Source: https://go.dev/doc/go1.NN

    ## A short statement of the finding
    → candidate | <track>
    Two or three sentences: what changed, and what a solver could be asked to
    predict or measure about it.

    ## Another finding
    → invalidates | <track>/<NN>-<slug>
    What this breaks in an existing task, and why.

    ## A third
    → no action
    Why it is not a task, in one line.

Tags are exactly `candidate`, `invalidates` and `no action`. A candidate names
the **track** it would land in; an invalidation names the **task id** it breaks,
and `undergo radar-check` fails if that task does not exist.

Era E4 is the exception to one-file-per-release. Most of Go 1.0–1.9 has been
superseded, so it is harvested into `versions/go1.0-1.9.md` — the durable
surprises only, each naming the release it came from. Ten thin files would
misrepresent the yield. This is also why `radar-check` reports a count of
*radar files* rather than of releases.

## What earns an entry

A finding must be **observable from Go code** — measurable, printable, or
provable by a test. "The linker is faster" is not an entry. "Small allocations
under 80 bytes are up to 30% cheaper" is.

Prefer changes where the *mechanism* is the answer, and be suspicious of
anything where a solver could be vaguely right.

## Eras

The sweep runs newest-first, because the newest material is the content nobody
else has written:

| Era | Versions | Status |
|---|---|---|
| E1 | 1.23 – 1.27 | done — swept release by release |
| E2 | 1.18 – 1.22 | done — swept release by release |
| E3 | 1.10 – 1.17 | done — swept release by release |
| E4 | 1.0 – 1.9 | done — **harvested into one file** |
