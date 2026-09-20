package cli

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// countingRepo is repo(t) plus a second task, so the totals under test are not
// all 1 and a transposed count would still be caught.
func countingRepo(t *testing.T, readme string) Env {
	t.Helper()
	e := repo(t)
	e.Out = io.Discard

	add := func(id, track string, mode string) {
		dir := filepath.Join(e.TasksDir(), filepath.FromSlash(id))
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		y := "schema: 1\nid: " + id + "\ntitle: \"t\"\nmode: " + mode +
			"\ntrack: " + track + "\ndifficulty: 1\nrequires: {go: \"1.27\"}\n"
		if err := os.WriteFile(filepath.Join(dir, "task.yaml"), []byte(y), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# t\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// repo(t) ships one predict task in the layout track.
	if err := os.WriteFile(filepath.Join(e.TaskDir("layout/01-struct-padding"), "README.md"),
		[]byte("# t\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	add("alloc/01-one", "alloc", "build")
	add("review/nil/01-one", "review/nil", "review")
	add("design/01-one", "design", "design")

	if readme != "" {
		if err := os.WriteFile(filepath.Join(e.Root, "README.md"), []byte(readme), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return e
}

// The counts in the top-level README drifted two milestones without anything
// noticing. A claim that is checked nowhere is a claim that rots.
func TestValidateRejectsAStaleTaskCountInTheReadme(t *testing.T) {
	e := countingRepo(t, "# undergo\n\n99 tasks ship today.\n")
	err := Validate(e, nil)
	if err == nil {
		t.Fatal("expected validate to reject a stale task count")
	}
	if !strings.Contains(err.Error(), "99") || !strings.Contains(err.Error(), "4") {
		t.Errorf("the error should name the claim and the truth; got: %v", err)
	}
}

func TestValidateAcceptsAccurateCounts(t *testing.T) {
	e := countingRepo(t, "# undergo\n\n4 tasks ship today: 2 machine-graded tasks across all "+
		"two internals tracks, 1 review drills across 1 categories, and 1 system-design katas.\n")
	if err := Validate(e, nil); err != nil {
		t.Fatalf("accurate counts must pass: %v", err)
	}
}

// The track count is written as an English word, because "all 16 internals
// tracks" reads worse than "all sixteen".
func TestValidateRejectsAStaleTrackCountWrittenAsAWord(t *testing.T) {
	e := countingRepo(t, "# undergo\n\n4 tasks ship today: 2 machine-graded tasks across all "+
		"nine internals tracks.\n")
	err := Validate(e, nil)
	if err == nil {
		t.Fatal("expected validate to reject a stale track count")
	}
	if !strings.Contains(err.Error(), "nine") {
		t.Errorf("the error should quote the word it read; got: %v", err)
	}
}

// A README that claims nothing is not lying. The rule catches wrong numbers,
// not missing prose — and a repo with no README at all is a fixture, not a fault.
func TestValidateIgnoresAReadmeThatMakesNoClaims(t *testing.T) {
	if err := Validate(countingRepo(t, "# undergo\n\nNo numbers here.\n"), nil); err != nil {
		t.Fatalf("a README with no counts must pass: %v", err)
	}
	if err := Validate(countingRepo(t, ""), nil); err != nil {
		t.Fatalf("no README at all must pass: %v", err)
	}
}
