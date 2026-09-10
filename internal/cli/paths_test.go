package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindRootWalksUp(t *testing.T) {
	root := t.TempDir()
	mod := "module github.com/RomanAgaltsev/undergo\n\ngo 1.27\n"
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte(mod), 0o644); err != nil {
		t.Fatal(err)
	}
	deep := filepath.Join(root, "tasks", "layout", "01-struct-padding")
	if err := os.MkdirAll(deep, 0o755); err != nil {
		t.Fatal(err)
	}

	got, err := FindRoot(deep)
	if err != nil {
		t.Fatalf("FindRoot: %v", err)
	}
	if got != root {
		t.Errorf("FindRoot = %q, want %q", got, root)
	}
}

func TestFindRootFailsOutsideTheRepo(t *testing.T) {
	if _, err := FindRoot(t.TempDir()); err == nil {
		t.Fatal("FindRoot succeeded outside an undergo checkout")
	}
}

func TestPathsHangOffRoot(t *testing.T) {
	e := Env{Root: filepath.FromSlash("/repo")}
	if got, want := e.TasksDir(), filepath.FromSlash("/repo/tasks"); got != want {
		t.Errorf("TasksDir = %q, want %q", got, want)
	}
	if got, want := e.WorkDir("layout/01-x"), filepath.FromSlash("/repo/work/layout/01-x"); got != want {
		t.Errorf("WorkDir = %q, want %q", got, want)
	}
	if got, want := e.ProgressPath(), filepath.FromSlash("/repo/.undergo/progress.yaml"); got != want {
		t.Errorf("ProgressPath = %q, want %q", got, want)
	}
}
