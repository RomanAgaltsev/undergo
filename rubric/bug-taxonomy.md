# Bug taxonomy (C1–C15)

The fifteen defect categories loupe drills are built around. Each drill's `ANSWERS.md` cites one
of these IDs per planted defect. The full taxonomy is the review lens even when a milestone
only ships drills for a subset (M1 = C1, C3, C5). Each category names a `cc-skills-golang`
skill as its depth-source — read that skill when a category's tells aren't obvious yet.

## C1 — Concurrency & races
**Tells:** a shared field/map/slice written from more than one goroutine without a lock or
atomic; `map` accessed concurrently; a goroutine with no exit path (leak); a channel closed by
the wrong party or twice; `WaitGroup.Add` called inside the goroutine.
**Depth:** `golang-concurrency`.

## C2 — Nil & memory safety
**Tells:** dereference of a pointer/interface that can be nil; write to a nil map; append that
aliases a shared backing array; unchecked type assertion (`x.(T)` without the comma-ok).
**Depth:** `golang-safety`.

## C3 — Error handling
**Tells:** an error assigned to `_` or otherwise ignored; wrapping with `%v` instead of `%w`
so `errors.Is/As` break; comparing a wrapped error with `==` instead of `errors.Is`; an error
both logged and returned (handled twice); a swallowed error that lets a zero value flow on.
**Depth:** `golang-error-handling`.

## C4 — Context & cancellation
**Tells:** a function that should take `context.Context` but doesn't; `ctx` accepted but not
passed to the call that does I/O; no deadline/timeout on a blocking call; work that keeps going
after `ctx.Done()`; a `cancel` from `WithCancel` never called.
**Depth:** `golang-context`.

## C5 — Resource leaks
**Tells:** `http.Response.Body`, `sql.Rows`, `*os.File`, or similar opened and never closed on
some path; `defer x.Close()` **inside a loop**; `time.NewTicker`/`NewTimer` never `Stop`ed;
`rows.Err()` never checked after iteration.
**Depth:** `golang-safety`, `golang-database`.

## C6 — API & interface design
**Tells:** an interface defined at the producer instead of the consumer; a value receiver that
mutates (or a needless pointer receiver); an exported signature changed in a breaking way; a
leaky abstraction exposing internal types.
**Depth:** `golang-structs-interfaces`, `golang-naming`.

## C7 — Performance
**Tells:** a query inside a loop (N+1); allocation in a hot path that could be reused/pooled;
an unbounded buffer or slice growth; a `make` without a known capacity; string concatenation in
a loop instead of `strings.Builder`.
**Depth:** `golang-performance`.

## C8 — Security
**Tells:** SQL built by string concatenation (injection); `exec` of a shell string with user
input; weak or misused crypto; unvalidated external input crossing a trust boundary; a secret
or token written to logs.
**Depth:** `golang-security`.

## C9 — Correctness & edge cases
**Tells:** off-by-one or wrong boundary; integer overflow or a narrowing conversion that
truncates; a zero value treated as valid; a missing case in a switch; wrong operator precedence.
**Depth:** `golang-safety`.

## C10 — Testing gaps
**Tells:** a test that asserts nothing (or only `err == nil`); a test that would still pass with
the bug present; the changed branch or its error path untested; a flaky/time- or
ordering-dependent test; a table test whose cases don't actually differ.
**Depth:** `golang-testing`.

## C11 — Generics misuse
**Tells:** an `any` type parameter with a runtime type switch/assertion inside (a generic that
should have a real constraint, or shouldn't be generic at all); a `func(any)` callback on a
generic that should be `func(T)`; `reflect.DeepEqual` where a `comparable` constraint + `==`
would do; a container that stores `[]any` internally so the type parameter is cosmetic; a
generic method/return that widens back to `any`; a lookup returning only `V` (no `ok`) so a
miss is indistinguishable from the zero value.
**Depth:** `golang-generics`, `golang-structs-interfaces`.

## C12 — JSON & (de)serialization
**Tells:** an **unexported** struct field expected in JSON (a tag can't save it); `Unmarshal`
into a non-pointer or with its error ignored; `omitempty` on a field whose zero value is
meaningful (drops a real `false`/`0`); a large `int64` id decoded through `any` and silently
turned into a lossy `float64`; a stray `json:"-"`; a nested pointer left unchecked after decode.
**Depth:** `golang-json`.

## C13 — Time & clocks
**Tells:** comparing `time.Time` with `==` instead of `.Equal()`; a wrong `Format`/`Parse`
layout (Go uses the reference `2006-01-02`, not `YYYY-MM-DD`); advancing a calendar day with
`Add(24*time.Hour)` (DST-unsafe) instead of `AddDate`; parsing without a location and assuming
local; a `time.After`/`NewTimer` in a loop that leaks; two `time.Now()` reads used as one
snapshot.
**Depth:** `golang-time`.

## C14 — HTTP client
**Tells:** `http.Get`/`http.DefaultClient` with no timeout (hangs forever); `resp.StatusCode`
never checked (an error page treated as data); a new `http.Client` built per request (no
connection reuse); a request built without the caller's `context`; a body not closed on some
path or not drained before `Close` (breaks keep-alive).
**Depth:** `golang-http`, `golang-context`.

## C15 — Typed-nil & interface nil
**Tells:** a concrete nil pointer (e.g. `*MyErr`) returned as an `error`/interface so
`err != nil` is wrongly true; a nil `*T` boxed into `any`/an interface field and then compared
`== nil` (always false); a typed-nil stored in a `map[...]error` or `[]error` that passes a
`!= nil` filter; an interface value assumed nil because its underlying pointer is nil.
**Depth:** `golang-safety`, `golang-structs-interfaces`.
