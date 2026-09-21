package cli

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/manifest"
	"github.com/RomanAgaltsev/undergo/internal/seal"
)

func TestCIVerifyNeverPrintsSolutionSource(t *testing.T) {
	e := sealedRepo(t)
	var out, errOut bytes.Buffer
	e.Out, e.Err = &out, &errOut

	// The fixture's sealed solution is not a compilable task, so this is
	// expected to report a failure — what matters is that it stays quiet
	// about the contents.
	_ = CIVerify(e, nil)

	combined := out.String() + errOut.String()
	if strings.Contains(combined, "package padding // solved") {
		t.Fatal("ci-verify leaked the reference solution into its output")
	}
	if !strings.Contains(combined, "layout/01-struct-padding") {
		t.Errorf("ci-verify did not report on the task:\n%s", combined)
	}
}

// runnableRepo seals a blob that gate 2 can carry all the way to running the
// tests: a reference prediction so a predict task is not rejected earlier, and
// an overlay that does not compile so the run fails for a readable reason.
func runnableRepo(t *testing.T) Env {
	t.Helper()
	e := repo(t)
	src := t.TempDir()
	if err := os.MkdirAll(filepath.Join(src, "solution"), 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(rel, body string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(src, filepath.FromSlash(rel)), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("HINT.md", "Count the padding after the last field.\n")
	write("EXPLANATION.md", "Alignment rounds the struct up.\n")
	write("prediction.yaml", "sizeof_header: 24\n")
	write("solution/padding.go", "package padding\n\nfunc broken() { totallyUndefinedHelper() }\n")

	blob, err := seal.Seal(src)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(e.TaskDir("layout/01-struct-padding"), SealName),
		[]byte(blob), 0o644); err != nil {
		t.Fatal(err)
	}
	return e
}

// A captured failure is useless if nobody can read it, so gate 2 writes the
// output to a gitignored file and prints only the path. Both halves matter: the
// path must reach the log, and the contents must not.
func TestCIVerifyWritesFailureLogButNotItsContents(t *testing.T) {
	e := runnableRepo(t)
	var out, errOut bytes.Buffer
	e.Out, e.Err = &out, &errOut

	if err := CIVerify(e, nil); err == nil {
		t.Fatal("expected gate 2 to report a failure for an overlay that does not compile")
	}

	body, err := os.ReadFile(e.FailureLogPath("layout/01-struct-padding"))
	if err != nil {
		t.Fatalf("gate 2 did not write a failure log: %v", err)
	}
	if !strings.Contains(string(body), "totallyUndefinedHelper") {
		t.Errorf("the failure log does not say why the run failed:\n%s", body)
	}

	combined := out.String() + errOut.String()
	if !strings.Contains(combined, "ci-failures") {
		t.Errorf("gate 2 did not say where the output went:\n%s", combined)
	}
	if strings.Contains(combined, "totallyUndefinedHelper") {
		t.Errorf("gate 2 printed the captured output instead of only its path:\n%s", combined)
	}
}

// A predict task's slot verdicts are safe to print — the names are public and
// predict.Check never prints a measured value — so a CI log can say which
// question failed. Nothing else from the captured output may be included.
func TestSlotVerdictsOnlyForPredictAndOnlyNames(t *testing.T) {
	const output = `=== RUN   TestPredictions
    sb_test.go:21: slot "plain_00_observed": correct
    sb_test.go:21: slot "atomic_00_observed": incorrect
    sb_test.go:21: slot "third_slot": no prediction
    sb_test.go:21: 1/3 slots correct
package padding // solved
--- FAIL: TestPredictions`

	got := slotVerdicts(&manifest.Task{Mode: manifest.ModePredict}, output)
	if !strings.Contains(got, `slot "atomic_00_observed": incorrect`) {
		t.Errorf("the failing slot was not named: %q", got)
	}
	if !strings.Contains(got, `slot "third_slot": no prediction`) {
		t.Errorf("a missing prediction was not named: %q", got)
	}
	if strings.Contains(got, "plain_00_observed") {
		t.Errorf("a passing slot should not be listed: %q", got)
	}
	if strings.Contains(got, "slots correct") {
		t.Errorf("the summary line should not be listed: %q", got)
	}
	if strings.Contains(got, "package padding") {
		t.Errorf("non-slot output leaked into the printed verdicts: %q", got)
	}

	// A build task's output can contain the reference source, so none of it is
	// printed — not even lines that look like slot verdicts.
	for _, mode := range []manifest.Mode{manifest.ModeBuild, manifest.ModeOptimize, manifest.ModeReview} {
		if v := slotVerdicts(&manifest.Task{Mode: mode}, output); v != "" {
			t.Errorf("mode %s should print nothing, got %q", mode, v)
		}
	}
}

// The log lives under .undergo/, which .gitignore excludes. If that ever stops
// being true, a reference solution could be committed by accident.
func TestFailureLogPathStaysUnderDotUndergo(t *testing.T) {
	e := Env{Root: filepath.FromSlash("/repo")}
	got := e.FailureLogPath("layout/01-struct-padding")
	want := filepath.FromSlash("/repo/.undergo/ci-failures/layout/01-struct-padding.log")
	if got != want {
		t.Errorf("FailureLogPath = %q, want %q", got, want)
	}
}

// toolchainRepo is a repo whose only task declares requires.toolchains.
func toolchainRepo(t *testing.T) Env {
	t.Helper()
	e := repo(t)
	e.Out = io.Discard
	e.Err = io.Discard
	p := filepath.Join(e.TaskDir("layout/01-struct-padding"), "task.yaml")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	y := strings.Replace(string(b), `requires: {go: "1.27"}`,
		"requires:\n  go: \"1.27\"\n  toolchains: [go1.22.12]", 1)
	if err := os.WriteFile(p, []byte(y), 0o644); err != nil {
		t.Fatal(err)
	}
	return e
}

// A skip nobody is told about is indistinguishable from success: if CI lost
// the network for a moment these tasks would stop being proven and the gate
// would stay green. --require-toolchains is what makes that impossible.
func TestCIVerifyRequireToolchainsTurnsAToolchainSkipIntoAFailure(t *testing.T) {
	orig := manifest.ToolchainAvailable
	t.Cleanup(func() { manifest.ToolchainAvailable = orig })
	manifest.ToolchainAvailable = func(string) bool { return false }

	e := toolchainRepo(t)

	if err := CIVerify(e, nil); err != nil {
		t.Fatalf("without the flag an unobtainable toolchain is a skip: %v", err)
	}
	err := CIVerify(e, []string{"--require-toolchains"})
	if err == nil {
		t.Fatal("with the flag it must fail")
	}
	// An undefined flag also returns an error, which would satisfy the check
	// above while proving nothing. The failure has to be about the task.
	if !strings.Contains(err.Error(), "layout/01-struct-padding") {
		t.Fatalf("the failure must name the unproven task; got: %v", err)
	}
}

// archRepo is a repo whose only task declares both a toolchain and an arch it
// cannot run on here.
func archRepo(t *testing.T) Env {
	t.Helper()
	e := repo(t)
	e.Out = io.Discard
	e.Err = io.Discard
	p := filepath.Join(e.TaskDir("layout/01-struct-padding"), "task.yaml")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	y := strings.Replace(string(b), `requires: {go: "1.27"}`,
		"requires:\n  go: \"1.27\"\n  arch: [nosucharch]\n  toolchains: [go1.22.12]", 1)
	if err := os.WriteFile(p, []byte(y), 0o644); err != nil {
		t.Fatal(err)
	}
	return e
}

// --require-toolchains is about toolchains. A task skipped because its answers
// are specific to another architecture has nothing to do with the network, and
// reporting it as an unobtainable toolchain would send whoever reads the log
// looking for a fetch that never failed.
func TestCIVerifyRequireToolchainsStillSkipsAnArchMismatch(t *testing.T) {
	orig := manifest.ToolchainAvailable
	t.Cleanup(func() { manifest.ToolchainAvailable = orig })
	manifest.ToolchainAvailable = func(string) bool { return true }

	e := archRepo(t)

	if err := CIVerify(e, []string{"--require-toolchains"}); err != nil {
		t.Fatalf("an arch mismatch is a skip even with the flag: %v", err)
	}
}
