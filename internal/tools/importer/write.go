package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/RomanAgaltsev/undergo/internal/seal"
)

// taskFiles is one task about to be written: its manifest fields, the files
// that go in the task directory, and the two rungs that go in the seal.
type taskFiles struct {
	ID         string
	Title      string
	Mode       string
	Track      string
	Difficulty int
	Estimate   string
	Tags       []string
	InspiredBy string

	// Files are written into the task directory, keyed by name.
	Files map[string]string

	Hint        string
	Explanation string
}

const manifestTemplate = `schema: 1
id: %s
title: %q
mode: %s
track: %s
difficulty: %d
estimate: %s
tags: [%s]
requires:
  go: "1.27"
verify: "go test ./..."
inspired_by: %q
deprecated: false
`

// writeTask writes a task directory and its sealed blob. The plaintext key is
// staged in a temporary directory and never touches the repository.
func writeTask(root string, tf taskFiles) error {
	dir := filepath.Join(root, "tasks", filepath.FromSlash(tf.ID))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	body := fmt.Sprintf(manifestTemplate, tf.ID, tf.Title, tf.Mode, tf.Track,
		tf.Difficulty, tf.Estimate, strings.Join(tf.Tags, ", "), tf.InspiredBy)
	if err := os.WriteFile(filepath.Join(dir, "task.yaml"), []byte(body), 0o644); err != nil {
		return err
	}
	for name, content := range tf.Files {
		target := filepath.Join(dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(target, []byte(content), 0o644); err != nil {
			return err
		}
	}

	staging, err := os.MkdirTemp("", "undergo-import-")
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(staging) }()
	if err := os.WriteFile(filepath.Join(staging, "HINT.md"), []byte(tf.Hint), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(staging, "EXPLANATION.md"), []byte(tf.Explanation), 0o644); err != nil {
		return err
	}
	blob, err := seal.Seal(staging)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "solution.sealed"), []byte(blob), 0o644)
}
