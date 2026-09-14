package traps

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// marker is what the runtime prints when checkfinalizers finds something.
const marker = "LIKELY CLEANUP/FINALIZER ISSUES"

// diagnosed builds the snippet in testdata/<name>.go.txt, runs it under
// GODEBUG=checkfinalizers=1, and reports whether the runtime complained.
//
// Each snippet is a whole program in its own module, built to a named output
// rather than run with `go run`, and executed as a subprocess — the runtime
// throws a fatal error when it diagnoses something, which would take this test
// process down with it.
func diagnosed(t *testing.T, name string) bool {
	t.Helper()

	src, err := os.ReadFile(filepath.Join("testdata", name+".go.txt"))
	if err != nil {
		t.Fatalf("reading snippet: %v", err)
	}

	dir := t.TempDir()
	write := func(file, body string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, file), []byte(body), 0o644); err != nil {
			t.Fatalf("writing %s: %v", file, err)
		}
	}
	write("go.mod", "module snippet\n\ngo 1.27\n")
	write("main.go", string(src))

	bin := filepath.Join(dir, "snippet")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = dir
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("building %s: %v\n%s", name, err, out)
	}

	run := exec.Command(bin)
	run.Dir = dir
	run.Env = append(os.Environ(), "GODEBUG=checkfinalizers=1")
	out, _ := run.CombinedOutput() // a diagnosis exits non-zero; that is the point
	return strings.Contains(string(out), marker)
}

// TestPredictions runs all four snippets and compares the verdicts with yours.
// It never prints what the runtime said.
func TestPredictions(t *testing.T) {
	predict.Check(t, map[string]any{
		"trap1_diagnosed": diagnosed(t, "trap1"),
		"trap2_diagnosed": diagnosed(t, "trap2"),
		"trap3_diagnosed": diagnosed(t, "trap3"),
		"trap4_diagnosed": diagnosed(t, "trap4"),
	})
}
