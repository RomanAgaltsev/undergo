package manifest

import (
	"os"
	"path/filepath"
	"testing"
)

const sample = `
schema: 1
id: layout/01-struct-padding
title: "Padding and alignment"
mode: predict
track: layout
difficulty: 2
estimate: 25m
tags: [unsafe, layout]
requires:
  go: "1.27"
  arch: [amd64, arm64]
verify: "go test ./..."
predict:
  slots: [sizeof_header, offset_flags]
inspired_by: "https://go.dev/ref/spec"
`

func TestLoadReadsEveryField(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.yaml")
	if err := os.WriteFile(path, []byte(sample), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.ID != "layout/01-struct-padding" {
		t.Errorf("ID = %q", got.ID)
	}
	if got.Mode != ModePredict {
		t.Errorf("Mode = %q", got.Mode)
	}
	if got.Difficulty != 2 {
		t.Errorf("Difficulty = %d", got.Difficulty)
	}
	if got.Requires.Go != "1.27" {
		t.Errorf("Requires.Go = %q", got.Requires.Go)
	}
	if got.Predict == nil || len(got.Predict.Slots) != 2 {
		t.Fatalf("Predict = %+v", got.Predict)
	}
	if got.Dir != dir {
		t.Errorf("Dir = %q, want %q", got.Dir, dir)
	}
}
