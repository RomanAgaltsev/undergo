package subject

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// outcome is what one snippet did when it ran.
type outcome struct {
	recovered bool
	fatal     bool
}

// run builds testdata/<name>.go.txt and runs it as its own process.
//
// Every snippet has the same shape: a deferred recover that prints RECOVERED if
// it caught something. So "recovered" is what the program reports about itself,
// and "fatal" is the runtime's own word for what happened instead.
func run(t *testing.T, name string) outcome {
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

	cmd := exec.Command(bin)
	cmd.Dir = dir
	out, _ := cmd.CombinedOutput() // a crash is a result, not a test failure
	text := string(out)

	return outcome{
		recovered: strings.Contains(text, "RECOVERED"),
		fatal:     strings.Contains(text, "fatal error"),
	}
}

// TestPredictions runs all five and compares the outcomes with yours. It never
// prints what the runtime said.
func TestPredictions(t *testing.T) {
	ordinary := run(t, "ordinary")
	child := run(t, "childgoroutine")
	concurrent := run(t, "concurrentmap")
	nilMap := run(t, "nilmapwrite")
	deadlock := run(t, "deadlock")

	predict.Check(t, map[string]any{
		"ordinary_recovered":        ordinary.recovered,
		"child_goroutine_recovered": child.recovered,
		"concurrent_map_recovered":  concurrent.recovered,
		"nil_map_write_recovered":   nilMap.recovered,
		"deadlock_recovered":        deadlock.recovered,
		"concurrent_map_is_fatal":   concurrent.fatal,
	})
}
