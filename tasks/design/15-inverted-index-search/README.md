# Kata k15 — Inverted-index search

**Track:** B — Storage & indexing · **Tier:** T2

## Problem
Design a full-text search index: map each term to the documents that contain it, then
answer ranked queries — the core of Lucene / Elasticsearch.

## Requirements
**Functional**
- Index a document (id + text).
- Query with one or more terms and return ranked hits.
- Support boolean AND / OR over terms.

**Non-functional**
- Query cost proportional to the matching postings, not the corpus size.
- Deterministic ranking with a stable tie-break.

## Given
- Hundreds of thousands of documents; interactive query latency.
- Single process, in-memory.

## Your core must
Expose an `Index` with:
- `Add(docID string, text string)` — tokenize and index a document.
- `Search(query string) []Hit` — ranked hits, where `Hit` carries the docID and a score.

Build a tokenizer + per-term postings lists + TF-IDF scoring. **Tests must show:** boolean
AND returns only documents containing *all* terms and OR returns documents containing
*any*; results are ordered by TF-IDF score descending; ties break deterministically (e.g.
by docID).

## Design should also address (in DESIGN.md)
- Tokenization / normalization (casefolding, stop-words, stemming) and its effect on recall.
- Postings-list storage and intersection (skip pointers / galloping search).
- TF-IDF vs BM25 ranking and why BM25 is usually preferred.
- Phrase / positional queries (storing positions in postings) as an extension.

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
