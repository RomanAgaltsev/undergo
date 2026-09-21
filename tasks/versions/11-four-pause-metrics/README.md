# 11 — One number became four

Stopping the world is not instantaneous. The runtime decides to stop, waits for
every P to reach a safe point — during which some threads are still running —
and only then is the world actually stopped.

For years `runtime/metrics` reported one distribution for all of that. This
task is about what it reports now, and it grades **names**, not durations: a
duration is a fact about the machine that ran it, a name is a fact about the
toolchain.

## The program

```go
var names []string
for _, d := range metrics.All() {
	if strings.Contains(d.Name, "pauses") {
		names = append(names, d.Name)
	}
}
sort.Strings(names)
fmt.Print(strings.Join(names, "|"))
```

It runs under two toolchains. The go line is `go 1.21` in every run and does
not vary — only the compiler does.

| slot | answer |
|---|---|
| `pause_metric_count_under_go1_21` | how many such names exist on the old toolchain |
| `pause_metric_count_on_local` | how many exist on the one you have |
| `sched_pauses_total_gc_present_under_go1_21` | is `/sched/pauses/total/gc:seconds` there? |
| `sched_pauses_total_gc_present_on_local` | is it there now? |
| `gc_pauses_seconds_present_on_local` | is `/gc/pauses:seconds` **still** there? |

The two counts are integers; the other three are `true` or `false`.

```
undergo verify versions/11-four-pause-metrics
```

**This task fetches a Go toolchain the first time it runs.** If it cannot be
fetched the task skips and says so.

## Questions to answer in writing

1. Name the two questions hiding inside a single stop-the-world number. Which
   one is a subset of the other, and what does the difference between them
   mean?
2. Not every stop-the-world pause is the garbage collector. Name two other
   things that stop it. What did that do to the old metric's name?
3. Predict the count on the new toolchain *before* you count the new names.
   What is Go allowed to do to a metric name that someone's dashboard already
   queries? Did your prediction follow from the split, or from the compatibility
   promise — and which one turned out to decide it?
4. You are paged because pause latency doubled. Which of these metrics do you
   look at first, and what does each possible answer tell you to do next?
