package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/manifest"
	"github.com/RomanAgaltsev/undergo/internal/seal"
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

// `undergo new` could not scaffold a review drill, which is 135 of the 267
// tasks. It derived the track by cutting at the FIRST slash, so a grouped id
// got track "review" where Validate wants "review/concurrency" — and the
// scaffolder emitted a manifest its own gate 3 then rejected. It went unnoticed
// because those 135 came from the importer, not from here.
func TestNewScaffoldsAGroupedID(t *testing.T) {
	e := repo(t)
	e.Out = &bytes.Buffer{}

	if err := New(e, []string{
		"--id", "review/concurrency/01-request-counter",
		"--mode", "review",
		"--title", "A counter that races",
	}); err != nil {
		t.Fatalf("New: %v", err)
	}

	dir := e.TaskDir("review/concurrency/01-request-counter")
	task, err := manifest.Load(filepath.Join(dir, "task.yaml"))
	if err != nil {
		t.Fatalf("scaffolded manifest does not load: %v", err)
	}
	if task.Track != "review/concurrency" {
		t.Errorf("track = %q, want review/concurrency", task.Track)
	}
	if err := manifest.Validate(task); err != nil {
		t.Fatalf("scaffolded manifest does not validate: %v", err)
	}
}

// The scaffold is now built as a value and validated before anything is
// written, so bad flags are refused rather than turned into a directory that
// has to be repaired by hand before `task ci` passes.
func TestNewRefusesFlagsThatCannotMakeAValidTask(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"an id with no track", []string{"--id", "01-lonely", "--mode", "build", "--title", "T"}},
		{"a mode that does not exist", []string{"--id", "alloc/01-x", "--mode", "guess", "--title", "T"}},
		{"a difficulty off the scale", []string{"--id", "alloc/01-x", "--mode", "build", "--title", "T", "--difficulty", "9"}},
		{"an id the grammar rejects", []string{"--id", "alloc/1-x", "--mode", "build", "--title", "T"}},
		{"an upper-case id", []string{"--id", "alloc/01-Zero", "--mode", "build", "--title", "T"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := repo(t)
			e.Out = &bytes.Buffer{}
			if err := New(e, tc.args); err == nil {
				t.Fatalf("New accepted %v", tc.args)
			}
			// Nothing may be left behind: a refusal that still created the
			// directory would make the next attempt fail with "already exists".
			if entries, err := os.ReadDir(e.TasksDir()); err == nil {
				for _, entry := range entries {
					if entry.Name() != "layout" {
						t.Errorf("refused scaffold left %s behind", entry.Name())
					}
				}
			}
		})
	}
}

// The template used to hardcode `go: "1.27"` — the module's own floor, so it
// ruled nothing out — and a `verify:` line no code has ever read. That template
// is how all 267 manifests came to carry the same meaningless pin.
func TestNewScaffoldsNoDeadFields(t *testing.T) {
	e := repo(t)
	e.Out = &bytes.Buffer{}
	if err := New(e, []string{
		"--id", "alloc/02-x", "--mode", "predict", "--title", "T",
	}); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(e.TaskDir("alloc/02-x"), "task.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	for _, dead := range []string{"verify:", `go: "1.27"`} {
		if strings.Contains(string(b), dead) {
			t.Errorf("scaffold still writes %q", dead)
		}
	}
}

// SealCmd deletes the only plaintext copy of an answer: _solution/ is
// gitignored, so a blob that does not decode would take the hint, the
// explanation, the reference solution and the answer sheet with it, silently,
// until the next gate 2. It now proves the round trip first.
func TestSealRefusesWhenTheRoundTripWouldLoseSomething(t *testing.T) {
	src := t.TempDir()
	if err := os.WriteFile(filepath.Join(src, "HINT.md"), []byte("a nudge\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	blob, err := seal.Seal(src)
	if err != nil {
		t.Fatal(err)
	}
	if err := verifySeal(src, blob); err != nil {
		t.Fatalf("a good blob must verify: %v", err)
	}

	// The file changes after the blob was made, which is what a lost or
	// corrupted entry looks like from here.
	if err := os.WriteFile(filepath.Join(src, "HINT.md"), []byte("a different nudge\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := verifySeal(src, blob); err == nil {
		t.Error("verifySeal accepted a blob that no longer matches the plaintext")
	}

	// A file the blob has never heard of.
	if err := os.WriteFile(filepath.Join(src, "EXTRA.md"), []byte("new\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := verifySeal(src, blob); err == nil {
		t.Error("verifySeal accepted a blob missing a file that exists")
	}
}

// After `undergo unseal`, the task directory holds the answer in plaintext.
// copyTask used to carry it into work/ beside the stub — the one place the whole
// design exists to keep it out of. Only a maintainer can reach that state, and a
// maintainer checking their own fix is exactly who would.
func TestStartDoesNotCopyPlaintextSolutions(t *testing.T) {
	e := repo(t)
	e.Out = &bytes.Buffer{}
	dir := e.TaskDir("layout/01-struct-padding")

	for _, name := range []string{AuthoringDir, SolutionSubdir} {
		if err := os.MkdirAll(filepath.Join(dir, name), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, name, "answer.go"),
			[]byte("package answer\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	if err := Start(e, []string{"layout/01-struct-padding"}); err != nil {
		t.Fatalf("Start: %v", err)
	}
	work := e.WorkDir("layout/01-struct-padding")
	for _, name := range []string{AuthoringDir, SolutionSubdir} {
		if _, err := os.Stat(filepath.Join(work, name)); err == nil {
			t.Errorf("start copied %s/ into the solver's work directory", name)
		}
	}
}
