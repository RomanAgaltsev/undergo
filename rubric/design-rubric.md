# Design rubric

Every kata `DESIGN.md` is graded against these eight dimensions **before** you build the core.
A strong design isn't long — it's one that makes its choices and their limits explicit.

## 1. Scope clarity
Functional and non-functional requirements are stated, and what's explicitly **out** of scope.
A design that tries to do everything has decided nothing.

## 2. Capacity estimation
Back-of-envelope numbers — QPS, storage, bandwidth, memory — derived from stated assumptions.
The point isn't precision; it's knowing which resource you'll run out of first.

## 3. API & data model
The public surface and the core types fit the requirements — no missing operation, no field
that exists "just in case." Someone could use it from the signatures alone.

## 4. Scaling story
How the system grows past one box: sharding, replication, caching, and where each stops
helping. Names the axis that scales and the axis that doesn't.

## 5. Consistency justification
The consistency/ordering guarantee is **named** (linearizable? per-key order? eventual?) and
**defended** against the requirements — not defaulted to "strong" out of habit.

## 6. Failure-mode coverage
What happens on node loss, overload, and network partition, and how the system **degrades**
rather than collapses. Every remote call is a failure waiting to happen.

## 7. Tradeoff articulation
The alternatives you considered and why you rejected them. A design with no rejected
alternatives hasn't been thought about.

## 8. Bottleneck identification
The single thing that gives out first, and at roughly what scale. If you can't name it, you
don't yet understand where the system's limit lives.

---

## How grading works

Claude scores each dimension **strong / adequate / thin** with one line of feedback, and calls
out the weakest dimension to strengthen. Iterate the design until no dimension is "thin" — then
build the core. The build will prove or break the design; that's the point.
