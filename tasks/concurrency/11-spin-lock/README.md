# 11 — A lock made of one atomic

`sync.Mutex` parks a goroutine that cannot get in. Build a lock that never
parks: a goroutine that finds it taken tries again until it gets it.

```go
var l spinlock.SpinLock // the zero value is unlocked
l.Lock()
counter++               // only one goroutine at a time gets here
l.Unlock()
```

## The contract

- `SpinLock` implements `sync.Locker`, and its zero value is unlocked.
- It is built from **one value from `sync/atomic`**, and nothing from `sync`.
- **Mutual exclusion:** at most one goroutine holds it at any moment.
- **It never parks, but it lets others run:** while `Lock` waits, other
  goroutines — including the one that holds the lock — must get to run, even
  with `GOMAXPROCS=1`.

```
undergo verify concurrency/11-spin-lock
```

`TestMutualExclusion` has eight goroutines add to a plain counter under your
lock, thirty rounds over; any round that comes out short fails.
`TestYieldsWhileWaiting` hands the lock between two goroutines fifty times on a
single P, within a quarter of a second. A lock that deadlocks fails after thirty
seconds with a message rather than hanging.

Both tests catch a broken lock without the race detector. To add its view of
your code as well — it needs a C toolchain — run, from the repository root:

```
go test -race ./work/concurrency/11-spin-lock/
```

or, without one:

```
docker run --rm -v "$PWD":/src -w /src golang:1.27 go test -race ./work/concurrency/11-spin-lock/
```

(`task race:docker` is not a substitute: it races the reference solutions, not
yours.)

## Questions to answer in writing

1. If your first `Lock` was a `Load` followed by a `Store`, describe the exact
   interleaving of two goroutines that lets both in, and say why
   `CompareAndSwap` closes it. Why could the race detector not have told you?
2. Remove the yield from your `Lock` and run `TestYieldsWhileWaiting`. Account
   for the time it now takes, round by round, in terms of the P and the
   scheduler's preemption interval.
3. `sync.Mutex` spins too, before it parks. Find `sync_runtime_canSpin` in the
   runtime and say why each of its conditions is there.
4. Under contention every waiter hammers one cache line with `CompareAndSwap`.
   Describe test-and-test-and-set and say what it saves.
5. Your lock has no queue. Describe a workload in which one goroutine never gets
   it, and say what `sync.Mutex` does about that which yours does not.
