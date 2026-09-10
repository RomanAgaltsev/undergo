package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/manifest"
)

func TestNewScaffoldsAValidTask(t *testing.T) {
	e := repo(t)
	e.Out = &bytes.Buffer{}

	err := New(e, []string{
		"--id", "alloc/01-zero-alloc-join",
		"--mode", "optimize",
		"--title", "Join without allocating",
		"--difficulty", "3",
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	dir := e.TaskDir("alloc/01-zero-alloc-join")
	task, err := manifest.Load(filepath.Join(dir, "task.yaml"))
	if err != nil {
		t.Fatalf("scaffolded manifest does not load: %v", err)
	}
	if err := manifest.Validate(task); err != nil {
		t.Fatalf("scaffolded manifest does not validate: %v", err)
	}
	for _, want := range []string{"README.md", "_solution/HINT.md", "_solution/EXPLANATION.md"} {
		if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(want))); err != nil {
			t.Errorf("missing %s: %v", want, err)
		}
	}
}

func TestSealReplacesPlaintextWithABlob(t *testing.T) {
	e := repo(t)
	e.Out = &bytes.Buffer{}
	if err := New(e, []string{
		"--id", "alloc/01-zero-alloc-join", "--mode", "build",
		"--title", "Join", "--difficulty", "3",
	}); err != nil {
		t.Fatal(err)
	}

	if err := SealCmd(e, []string{"alloc/01-zero-alloc-join"}); err != nil {
		t.Fatalf("SealCmd: %v", err)
	}
	dir := e.TaskDir("alloc/01-zero-alloc-join")
	blob, err := os.ReadFile(filepath.Join(dir, SealName))
	if err != nil {
		t.Fatalf("no seal produced: %v", err)
	}
	if strings.Contains(string(blob), "HINT") {
		t.Error("the seal contains plaintext")
	}
	if _, err := os.Stat(filepath.Join(dir, "_solution")); err == nil {
		t.Error("_solution/ still exists after sealing")
	}
}
