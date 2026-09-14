# 03 — Which one sets the goal

GOGC asks for a goal proportional to the live set. `GOMEMLIMIT` asks for a goal
that keeps total memory under a ceiling. Both can be set at once, and then only
one of them is actually deciding.

For each configuration, predict whether the **limit** is what held the goal down.

| slot | configuration |
|---|---|
| `clamped_gogc800_limit16` | GOGC=800, limit 16 MiB, ~8 MiB live |
| `clamped_gogc100_limit512` | GOGC=100, limit 512 MiB, ~8 MiB live |
| `clamped_gogc_off_limit32` | GOGC **off**, limit 32 MiB |
| `clamped_gogc100_limit_huge` | GOGC=100, limit 1 TiB |
| `limit_can_be_exceeded` | can a program hold more than the limit allows? |

```
undergo verify gc/03-memory-limit
```

The last one is not about which knob wins. It is about what kind of promise
`GOMEMLIMIT` is, and the naive answer is the wrong one.

## Why this task asks yes-or-no

The clamped goal itself is not reproducible: measuring it eight times under
identical settings gave four different values, because it depends on non-heap
memory that moves between runs. **Whether** the limit is doing the clamping is
stable in all eight. So that is what you are asked for.

## Questions to answer in writing

1. State the rule for which of GOGC and `GOMEMLIMIT` decides the goal, in one
   sentence that predicts all four of the first slots.
2. `GOMEMLIMIT` is documented as a *soft* limit. Say what the runtime will do as
   a program approaches it, and what it will not do — and say what happens if the
   program keeps allocating anyway.
3. Setting `GOMEMLIMIT` and leaving `GOGC=100` is a common pairing and often the
   wrong one. Describe the deployment where it misbehaves, what you would set
   instead, and what you give up by doing so.
4. The clamped goal is not reproducible between runs, but *whether it is clamped*
   is. Say which of those two you would alert on in production, and what the
   alert would mean.
5. Every function here restores GOGC and the memory limit with `defer`. Say what
   would break without that, and why it would be hard to diagnose.
