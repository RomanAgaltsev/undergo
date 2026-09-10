package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/manifest"
	"github.com/RomanAgaltsev/undergo/internal/seal"
)

func sampleTask() taskFiles {
	return taskFiles{
		ID: "review-concurrency/01-request-counter", Title: "Request counter",
		Mode: "review", Track: "review-concurrency", Difficulty: 2, Estimate: "20m",
		Tags:       []string{"review", "concurrency", "obvious"},
		InspiredBy: "https://github.com/RomanAgaltsev/loupe",
		Files: map[string]string{
			"README.md": "# Drill\n\nReview it.\n",
			"drill.go":  "package drill\n",
		},
		Hint:        "# Hint\n\nTwo defects.\n",
		Explanation: "# Answer key\n\nThe map write is racy.\n",
	}
}

func TestWriteTaskProducesAValidSealedTask(t *testing.T) {
	root := t.TempDir()
	if err := writeTask(root, sampleTask()); err != nil {
		t.Fatalf("writeTask: %v", err)
	}

	dir := filepath.Join(root, "tasks", "review-concurrency", "01-request-counter")
	task, err := manifest.Load(filepath.Join(dir, "task.yaml"))
	if err != nil {
		t.Fatalf("manifest does not load: %v", err)
	}
	if err := manifest.Validate(task); err != nil {
		t.Fatalf("manifest does not validate: %v", err)
	}
	if task.Mode != manifest.ModeReview || task.Difficulty != 2 {
		t.Errorf("manifest = %+v", task)
	}

	for _, name := range []string{"README.md", "drill.go", "solution.sealed"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Errorf("missing %s: %v", name, err)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, "_solution")); err == nil {
		t.Error("_solution/ survived — the plaintext key must not be left on disk")
	}

	blob, err := os.ReadFile(filepath.Join(dir, "solution.sealed"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(blob), "racy") {
		t.Error("the seal leaks the answer in plaintext")
	}
	hint, err := seal.File(string(blob), "HINT.md")
	if err != nil {
		t.Fatalf("no HINT.md in the blob — undergo hint would fail: %v", err)
	}
	if !strings.Contains(string(hint), "Two defects") {
		t.Errorf("HINT.md = %q", hint)
	}
	if _, err := seal.File(string(blob), "EXPLANATION.md"); err != nil {
		t.Fatalf("no EXPLANATION.md in the blob — undergo reveal would fail: %v", err)
	}
}
