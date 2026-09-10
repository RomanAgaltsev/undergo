# Design: unique ID generator

> Graded against `../../../rubric/design-rubric.md`. Fill each section; keep it tight.

## Problem & scope
<What you're designing, and what's explicitly out of scope.>

## Requirements (functional / non-functional)
<The must-dos and the -ilities.>

## Capacity estimation
<Bit budget: time bits (range), node bits (count), sequence bits (IDs/ms). Show the math.>

## API
<The public operations and their signatures.>

## Data model
<The 64-bit layout and the generator's state.>

## Architecture
<How an ID is assembled; where the state lives.>

## Scaling
<How it grows to more nodes / higher rate; the axis that scales and the one that doesn't.>

## Consistency
<Ordering guarantee (k-sorted?) and why strict global order isn't provided.>

## Failure modes
<Clock regression / node-id collision / sequence exhaustion; how you handle each.>

## Tradeoffs & alternatives
<Snowflake vs UUID vs ticket server, and the bit-allocation choice.>

## Bottleneck
<The first thing to give out, and at what scale.>
