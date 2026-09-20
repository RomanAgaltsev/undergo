# 08 — The suite that ran nothing

`testdata/suite` has three tests. The middle one calls `os.Exit(0)` — the
success exit code — part way through:

```go
func TestAlpha(t *testing.T) { fmt.Println("ALPHA RAN") }

func TestBeta(t *testing.T) {
	fmt.Println("BETA RAN")
	os.Exit(0)
}

func TestGamma(t *testing.T) { fmt.Println("GAMMA RAN") }
```

Nothing here is contrived. A test reaches production code whose shutdown path,
or whose flag parsing, or whose "configuration is invalid, give up cleanly"
branch calls `os.Exit`. The test does not know it is about to.

Predict:

| slot | question |
|---|---|
| `go_test_reports_failure` | does the suite report failure overall? |
| `output_names_os_exit` | does the output mention `os.Exit(0)`? |
| `earlier_test_ran` | did `TestAlpha` run? |
| `later_test_ran` | did `TestGamma` run? |

```
undergo verify edges/08-exit-during-test
```

The judge is the toolchain: the test runs the suite in `testdata/suite` and
grades its output and its exit status.

## Questions to answer in writing

1. `os.Exit(0)` is the success code. Say why the testing package treats it as a
   failure anyway, in terms of what an exit status is supposed to mean here.
2. Say exactly how much of the suite ran, and what the exit status would have
   been if the testing package did *not* intervene.
3. This behaviour was added to Go. Describe what a CI pipeline looked like
   before it — be specific about what the dashboard showed and what was true of
   the code.
4. `os.Exit` skips deferred functions. Name two things a test suite loses
   because of that, beyond the tests that never ran.
5. `TestMain` is *allowed* to call `os.Exit`, and must. Say why that is not a
   contradiction, and what distinguishes the two cases.
