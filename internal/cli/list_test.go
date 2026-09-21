package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/manifest"
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
requires: {}
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

// A grouped track means --track has to match a prefix: --track review lists
// every category, --track review/concurrency lists one.
func TestListTrackMatchesAGroupPrefix(t *testing.T) {
	tests := map[string]bool{
		"review":             true,
		"review/concurrency": true,
		"review/nil-safety":  false,
		"rev":                false,
		"":                   true,
	}
	for filter, want := range tests {
		t.Run(filter, func(t *testing.T) {
			if got := trackMatches("review/concurrency", filter); got != want {
				t.Errorf("trackMatches(%q, %q) = %v, want %v", "review/concurrency", filter, got, want)
			}
		})
	}
	if trackMatches("layout", "layout") != true {
		t.Error("an ungrouped track must still match itself exactly")
	}
}

// deprecate marks repo(t)'s only task as retired, and adds a second live task
// so the list has something left to show.
func deprecate(t *testing.T, e Env) {
	t.Helper()
	p := filepath.Join(e.TaskDir("layout/01-struct-padding"), "task.yaml")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, append(b, []byte("\ndeprecated: true\n")...), 0o644); err != nil {
		t.Fatal(err)
	}

	dir := filepath.Join(e.TasksDir(), "alloc", "01-live")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := "schema: 1\nid: alloc/01-live\ntitle: \"Still here\"\nmode: build\n" +
		"track: alloc\ndifficulty: 3\nrequires: {}\ntags: [alloc, escape]\n"
	if err := os.WriteFile(filepath.Join(dir, "task.yaml"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// CONTRIBUTING.md says twice that a retired task gets `deprecated: true` rather
// than a new name, and spec 11 rests the whole ID-permanence contract on it.
// Nothing read the field: list printed it, Gradeable graded it, and the README
// counts counted it. Setting it did literally nothing.
func TestListHidesDeprecatedTasks(t *testing.T) {
	e := repo(t)
	deprecate(t, e)
	var out bytes.Buffer
	e.Out = &out

	if err := List(e, nil); err != nil {
		t.Fatalf("List: %v", err)
	}
	if strings.Contains(out.String(), "layout/01-struct-padding") {
		t.Errorf("a retired task is still listed:\n%s", out.String())
	}
	if !strings.Contains(out.String(), "alloc/01-live") {
		t.Errorf("the live task vanished too:\n%s", out.String())
	}
}

func TestListShowsDeprecatedTasksWhenAsked(t *testing.T) {
	e := repo(t)
	deprecate(t, e)
	var out bytes.Buffer
	e.Out = &out

	if err := List(e, []string{"--deprecated"}); err != nil {
		t.Fatalf("List: %v", err)
	}
	if !strings.Contains(out.String(), "layout/01-struct-padding") {
		t.Errorf("--deprecated did not bring the retired task back:\n%s", out.String())
	}
}

// A retired task is not one that ships, so it must not inflate the number the
// README promises — otherwise deprecating a task would make the README wrong
// and gate 3 would then enforce the wrong number.
func TestDeprecatedTasksAreNotCounted(t *testing.T) {
	e := repo(t)
	deprecate(t, e)
	tasks, err := manifest.Walk(e.TasksDir())
	if err != nil {
		t.Fatal(err)
	}
	if got := countCatalogue(tasks).total; got != 1 {
		t.Errorf("total = %d, want 1 (two tasks, one retired)", got)
	}
}

// tags was parsed on all 267 manifests and read by nothing at all. Filtering on
// it is what makes the field worth carrying.
func TestListFiltersByTag(t *testing.T) {
	e := repo(t)
	deprecate(t, e)
	var out bytes.Buffer
	e.Out = &out

	if err := List(e, []string{"--tag", "escape"}); err != nil {
		t.Fatalf("List: %v", err)
	}
	if !strings.Contains(out.String(), "alloc/01-live") {
		t.Errorf("--tag escape did not match the task carrying it:\n%s", out.String())
	}

	out.Reset()
	if err := List(e, []string{"--tag", "nosuchtag"}); err != nil {
		t.Fatalf("List: %v", err)
	}
	if strings.Contains(out.String(), "alloc/01-live") {
		t.Errorf("--tag nosuchtag matched anyway:\n%s", out.String())
	}
}
