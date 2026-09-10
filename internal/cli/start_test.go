package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStartCopiesTheTaskWithoutTheSeal(t *testing.T) {
	e := repo(t)
	e.Out = &bytes.Buffer{}
	dir := e.TaskDir("layout/01-struct-padding")
	for name, body := range map[string]string{
		"README.md":       "# Padding\n",
		"padding.go":      "package padding\n",
		"solution.sealed": "H4sIAAAA\n",
	} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	if err := Start(e, []string{"layout/01-struct-padding"}); err != nil {
		t.Fatalf("Start: %v", err)
	}

	work := e.WorkDir("layout/01-struct-padding")
	if _, err := os.Stat(filepath.Join(work, "padding.go")); err != nil {
		t.Errorf("stub not copied: %v", err)
	}
	if _, err := os.Stat(filepath.Join(work, "solution.sealed")); err == nil {
		t.Error("the seal was copied into the work directory")
	}
}

func TestStartWritesABlankPredictionFile(t *testing.T) {
	e := repo(t)
	e.Out = &bytes.Buffer{}
	dir := e.TaskDir("layout/01-struct-padding")
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# x\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Start(e, []string{"layout/01-struct-padding"}); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(e.WorkDir("layout/01-struct-padding"), "prediction.yaml"))
	if err != nil {
		t.Fatalf("prediction.yaml not written: %v", err)
	}
	if !strings.Contains(string(b), "sizeof_header:") {
		t.Errorf("prediction.yaml has no slot for sizeof_header:\n%s", b)
	}
}

func TestStartRefusesToClobberWork(t *testing.T) {
	e := repo(t)
	e.Out = &bytes.Buffer{}
	if err := os.WriteFile(filepath.Join(e.TaskDir("layout/01-struct-padding"), "README.md"),
		[]byte("# x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Start(e, []string{"layout/01-struct-padding"}); err != nil {
		t.Fatal(err)
	}
	if err := Start(e, []string{"layout/01-struct-padding"}); err == nil {
		t.Fatal("Start overwrote existing work without --force")
	}
}
