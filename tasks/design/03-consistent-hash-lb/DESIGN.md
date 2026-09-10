# Design: consistent-hash load balancer

> Graded against `../../../rubric/design-rubric.md`. Fill each section; keep it tight.

## Problem & scope
<What you're designing, and what's explicitly out of scope.>

## Requirements (functional / non-functional)
<The must-dos and the -ilities.>

## Capacity estimation
<Ring size / virtual nodes per member / memory; lookup cost.>

## API
<The public operations and their signatures.>

## Data model
<The ring representation and how members map onto it.>

## Architecture
<How a lookup resolves a key to a member; how add/remove mutate the ring.>

## Scaling
<How it grows to more members/keys; the axis that scales and the one that doesn't.>

## Consistency
<What "minimal remapping" guarantees; behavior during a membership change.>

## Failure modes
<Member loss / skew / a hot key; how it degrades.>

## Tradeoffs & alternatives
<Consistent hashing vs rendezvous; virtual-node count; bounded loads.>

## Bottleneck
<The first thing to give out, and at what scale.>
