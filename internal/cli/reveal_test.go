package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/progress"
	"github.com/RomanAgaltsev/undergo/internal/seal"
)

func sealedRepo(t *testing.T) Env {
	t.Helper()
	e := repo(t)
	src := t.TempDir()
	if err := os.MkdirAll(filepath.Join(src, "solution"), 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(rel, body string) {
		if err := os.WriteFile(filepath.Join(src, filepath.FromSlash(rel)), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("HINT.md", "Count the padding after the last field.\n")
	write("EXPLANATION.md", "Alignment rounds the struct up.\n")
	write("solution/padding.go", "package padding // solved\n")

	blob, err := seal.Seal(src)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(e.TaskDir("layout/01-struct-padding"), SealName),
		[]byte(blob), 0o644); err != nil {
		t.Fatal(err)
	}
	return e
}

func TestHintIsAlwaysAvailableAndRecorded(t *testing.T) {
	e := sealedRepo(t)
	var out bytes.Buffer
	e.Out = &out

	if err := Hint(e, []string{"layout/01-struct-padding"}); err != nil {
		t.Fatalf("Hint: %v", err)
	}
	if !strings.Contains(out.String(), "Count the padding") {
		t.Errorf("hint not printed:\n%s", out.String())
	}
	rec, err := progress.Load(e.ProgressPath())
	if err != nil {
		t.Fatal(err)
	}
	if !rec.Get("layout/01-struct-padding").HintUsed {
		t.Error("HintUsed was not recorded")
	}
}

func TestRevealRefusesWhileUnsolved(t *testing.T) {
	e := sealedRepo(t)
	e.Out, e.Err = &bytes.Buffer{}, &bytes.Buffer{}

	err := Reveal(e, []string{"layout/01-struct-padding"})
	if err == nil {
		t.Fatal("Reveal opened an unsolved task without --stuck")
	}
	if !strings.Contains(err.Error(), "--stuck") {
		t.Errorf("the refusal must name the escape hatch, got: %v", err)
	}
}

func TestRevealWithStuckExtractsAndRecordsThePeek(t *testing.T) {
	e := sealedRepo(t)
	var out bytes.Buffer
	e.Out, e.Err = &out, &bytes.Buffer{}

	if err := Reveal(e, []string{"--stuck", "layout/01-struct-padding"}); err != nil {
		t.Fatalf("Reveal: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(e.RevealDir("layout/01-struct-padding"), "solution", "padding.go"))
	if err != nil {
		t.Fatalf("solution not extracted: %v", err)
	}
	if !strings.Contains(string(got), "solved") {
		t.Errorf("extracted = %q", got)
	}
	rec, err := progress.Load(e.ProgressPath())
	if err != nil {
		t.Fatal(err)
	}
	if !rec.Get("layout/01-struct-padding").Peeked {
		t.Error("the peek was not recorded")
	}
}

func TestRevealIsFreeOnceSolved(t *testing.T) {
	e := sealedRepo(t)
	e.Out, e.Err = &bytes.Buffer{}, &bytes.Buffer{}

	rec, err := progress.Load(e.ProgressPath())
	if err != nil {
		t.Fatal(err)
	}
	rec.RecordRun("layout/01-struct-padding", true)
	if err := rec.Save(e.ProgressPath()); err != nil {
		t.Fatal(err)
	}

	if err := Reveal(e, []string{"layout/01-struct-padding"}); err != nil {
		t.Fatalf("Reveal refused a solved task: %v", err)
	}
	rec, err = progress.Load(e.ProgressPath())
	if err != nil {
		t.Fatal(err)
	}
	if rec.Get("layout/01-struct-padding").Peeked {
		t.Error("reading the solution after solving must not be recorded as a peek")
	}
}

// The usage line and the refusal message both document `reveal <id> --stuck`.
// Go's flag package stops at the first positional, so this order has to be
// handled explicitly or the tool rejects its own advertised command.
func TestRevealAcceptsTheFlagAfterTheID(t *testing.T) {
	e := sealedRepo(t)
	e.Out, e.Err = &bytes.Buffer{}, &bytes.Buffer{}

	if err := Reveal(e, []string{"layout/01-struct-padding", "--stuck"}); err != nil {
		t.Fatalf("Reveal <id> --stuck: %v", err)
	}
	if _, err := os.Stat(e.RevealDir("layout/01-struct-padding")); err != nil {
		t.Errorf("nothing extracted: %v", err)
	}
}

func TestStartAcceptsTheFlagAfterTheID(t *testing.T) {
	e := sealedRepo(t)
	e.Out, e.Err = &bytes.Buffer{}, &bytes.Buffer{}
	if err := os.WriteFile(filepath.Join(e.TaskDir("layout/01-struct-padding"), "README.md"),
		[]byte("# x\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Start(e, []string{"layout/01-struct-padding"}); err != nil {
		t.Fatal(err)
	}
	if err := Start(e, []string{"layout/01-struct-padding", "--force"}); err != nil {
		t.Fatalf("Start <id> --force: %v", err)
	}
}
