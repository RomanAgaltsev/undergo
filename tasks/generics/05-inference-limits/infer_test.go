package infer

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// compiles reports whether the snippet in testdata/<name>.go.txt builds.
//
// Each snippet is a complete package main carrying its own declarations, built
// in its own temporary module, so one snippet's verdict never depends on
// another's. It never reports what the compiler said.
func compiles(t *testing.T, name string) bool {
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

	cmd := exec.Command("go", "build", ".")
	cmd.Dir = dir
	return cmd.Run() == nil
}

// TestPredictions compiles each call site and compares the verdict with yours.
// It never prints a compiler diagnostic.
func TestPredictions(t *testing.T) {
	predict.Check(t, map[string]any{
		"call1_compiles": compiles(t, "call1"),
		"call2_compiles": compiles(t, "call2"),
		"call3_compiles": compiles(t, "call3"),
		"call4_compiles": compiles(t, "call4"),
		"call5_compiles": compiles(t, "call5"),
		"call6_compiles": compiles(t, "call6"),
	})
}
