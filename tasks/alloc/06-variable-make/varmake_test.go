package varmake

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// reportEnv switches TestReport on. Only measure sets it, in a child process
// whose output the parent reads, so a solver's own run never prints a count.
const reportEnv = "UNDERGO_VARMAKE_REPORT"

// sizes are the calls TestReport measures, by name.
var sizes = []struct {
	name string
	fn   func(int) int
	n    int
}{
	{"bytes_1", Bytes, 1},
	{"bytes_32", Bytes, 32},
	{"bytes_33", Bytes, 33},
	{"words_4", Words, 4},
	{"words_5", Words, 5},
}

// TestReport prints one "name=allocs" line per size. It is not a test of
// anything: it is the instrument measure runs in a child process.
func TestReport(t *testing.T) {
	if os.Getenv(reportEnv) == "" {
		t.Skip("run by this task's own measurement, not directly")
	}
	for _, s := range sizes {
		allocs := testing.AllocsPerRun(100, func() { s.fn(s.n) })
		fmt.Printf("%s=%d\n", s.name, int(allocs))
	}
}

var reportLine = regexp.MustCompile(`^([a-z0-9_]+)=([0-9]+)$`)

// childEnv is the environment for every child process. GOFLAGS is replaced so
// that a solver's own GOFLAGS=-gcflags=... cannot change the answers.
//
// Not with an empty value: the go command treats an empty variable as unset
// and falls back to the go env file, so a `go env -w GOFLAGS=...` would still
// apply. -buildvcs=auto is that flag's own default — a no-op every one of
// build, vet and test accepts — and being non-empty, it wins.
func childEnv(extra ...string) []string {
	return append(append(os.Environ(), "GOFLAGS=-buildvcs=auto"), extra...)
}

// measure runs TestReport in a fresh `go test` of this package, compiled with
// the given extra compiler flags, and returns the allocation counts.
//
// A child process, because the question is what the compiler emits, and a
// compiler flag cannot be changed inside a binary that is already built.
func measure(t *testing.T, gcflags string) map[string]int {
	t.Helper()

	args := []string{"test", "-count=1", "-v", "-run", "^TestReport$"}
	if gcflags != "" {
		args = append(args, "-gcflags="+gcflags)
	}
	cmd := exec.Command("go", append(args, ".")...)
	cmd.Env = childEnv(reportEnv + "=1")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("measuring with gcflags %q: %v", gcflags, err)
	}

	got := map[string]int{}
	for line := range strings.SplitSeq(strings.ReplaceAll(string(out), "\r\n", "\n"), "\n") {
		m := reportLine.FindStringSubmatch(strings.TrimSpace(line))
		if m == nil {
			continue
		}
		n, err := strconv.Atoi(m[2])
		if err != nil {
			t.Fatalf("measuring with gcflags %q: unreadable count for %s", gcflags, m[1])
		}
		got[m[1]] = n
	}
	if len(got) != len(sizes) {
		t.Fatalf("measuring with gcflags %q: read %d of %d sizes", gcflags, len(got), len(sizes))
	}
	return got
}

// saysEscapes compiles this package with -m (plus extra flags) and reports
// whether the compiler says Bytes's make escapes to the heap.
func saysEscapes(t *testing.T, extra string) bool {
	t.Helper()

	flags := "-m"
	if extra != "" {
		flags += " " + extra
	}
	cmd := exec.Command("go", "build", "-gcflags="+flags, ".")
	cmd.Env = childEnv()
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("compiling with %q: %v", flags, err)
	}

	found := false
	for line := range strings.SplitSeq(string(out), "\n") {
		if !strings.Contains(line, "make([]byte, n)") {
			continue
		}
		found = true
		if strings.Contains(line, "escapes to heap") {
			return true
		}
	}
	if !found {
		t.Fatalf("compiling with %q: the compiler said nothing about make([]byte, n)", flags)
	}
	return false
}

// TestPredictions grades seven slots: five allocation counts and two verdicts
// from the compiler's escape diagnostics, with the stack buffer on and off.
func TestPredictions(t *testing.T) {
	const hashOff = "-d=variablemakehash=n"

	on := measure(t, "")
	off := measure(t, hashOff)

	predict.Check(t, map[string]any{
		"m_says_escapes":          saysEscapes(t, ""),
		"bytes_32_allocs":         on["bytes_32"],
		"bytes_33_allocs":         on["bytes_33"],
		"words_4_allocs":          on["words_4"],
		"words_5_allocs":          on["words_5"],
		"bytes_1_allocs_hash_off": off["bytes_1"],
		"m_says_escapes_hash_off": saysEscapes(t, hashOff),
	})
}
