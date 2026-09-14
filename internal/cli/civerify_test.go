package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

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
