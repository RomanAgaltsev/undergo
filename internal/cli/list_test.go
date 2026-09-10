package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func repo(t *testing.T) Env {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"),
		[]byte("module github.com/RomanAgaltsev/undergo\n\ngo 1.27\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(root, "tasks", "layout", "01-struct-padding")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := `
schema: 1
id: layout/01-struct-padding
title: "Padding and the trailing zero-size field"
mode: predict
track: layout
difficulty: 2
requires: {go: "1.27"}
predict: {slots: [sizeof_header]}
`
	if err := os.WriteFile(filepath.Join(dir, "task.yaml"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return Env{Root: root}
}

func TestListShowsIDTitleAndMode(t *testing.T) {
	e := repo(t)
	var out bytes.Buffer
	e.Out = &out

	if err := List(e, nil); err != nil {
		t.Fatalf("List: %v", err)
	}
	got := out.String()
	for _, want := range []string{"layout/01-struct-padding", "predict", "Padding"} {
		if !strings.Contains(got, want) {
			t.Errorf("list output missing %q:\n%s", want, got)
		}
	}
}

func TestListFiltersByTrack(t *testing.T) {
	e := repo(t)
	var out bytes.Buffer
	e.Out = &out

	if err := List(e, []string{"--track", "alloc"}); err != nil {
		t.Fatalf("List: %v", err)
	}
	if strings.Contains(out.String(), "layout/01-struct-padding") {
		t.Errorf("--track alloc listed a layout task:\n%s", out.String())
	}
}
