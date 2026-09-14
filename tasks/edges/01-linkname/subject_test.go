package subject

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// builds reports whether testdata/<name>.go.txt links, with the given extra
// build arguments.
func builds(t *testing.T, name string, args ...string) bool {
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
	full := append([]string{"build"}, args...)
	full = append(full, "-o", bin, ".")

	cmd := exec.Command("go", full...)
	cmd.Dir = dir
	return cmd.Run() == nil
}

// TestPredictions links each snippet and compares the outcome with yours. It
// never prints the linker's message.
func TestPredictions(t *testing.T) {
	predict.Check(t, map[string]any{
		"gcstart_links":             builds(t, "gcstart"),
		"nanotime_links":            builds(t, "nanotime"),
		"mallocgc_links":            builds(t, "mallocgc"),
		"own_package_links":         builds(t, "ownpackage"),
		"gcstart_links_with_optout": builds(t, "gcstart", "-ldflags=-checklinkname=0"),
	})
}
