package cli

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

const paddingID = "layout/01-struct-padding"

// unseal is the inverse of seal: the round trip must return the bytes the
// author sealed, or editing a shipped hint is guesswork.
func TestUnsealRoundTripsWhatSealWrote(t *testing.T) {
	e := sealedRepo(t)
	e.Out = io.Discard

	if err := Unseal(e, []string{paddingID}); err != nil {
		t.Fatalf("Unseal: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(e.TaskDir(paddingID), "_solution", "HINT.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "Count the padding after the last field.\n" {
		t.Fatalf("hint round-tripped as %q", got)
	}
}

// The seal holds more than the hint; unsealing must restore all of it or a
// reseal would silently drop the explanation.
func TestUnsealRestoresEveryEntry(t *testing.T) {
	e := sealedRepo(t)
	e.Out = io.Discard
	if err := Unseal(e, []string{paddingID}); err != nil {
		t.Fatalf("Unseal: %v", err)
	}

	dir := filepath.Join(e.TaskDir(paddingID), "_solution")
	for _, rel := range []string{"HINT.md", "EXPLANATION.md", filepath.Join("solution", "padding.go")} {
		if _, err := os.Stat(filepath.Join(dir, rel)); err != nil {
			t.Errorf("missing from the unsealed tree: %s", rel)
		}
	}
}

// Overwriting an author's in-progress plaintext would destroy unsealed work
// that exists nowhere else, because _solution/ is gitignored.
func TestUnsealRefusesToClobberExistingPlaintext(t *testing.T) {
	e := sealedRepo(t)
	e.Out = io.Discard
	dir := filepath.Join(e.TaskDir(paddingID), "_solution")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "HINT.md"), []byte("in progress\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Unseal(e, []string{paddingID}); err == nil {
		t.Fatal("expected a refusal while _solution/ is present")
	}
	got, err := os.ReadFile(filepath.Join(dir, "HINT.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "in progress\n" {
		t.Fatalf("refusal must not have touched the file; got %q", got)
	}
}

func TestUnsealRequiresAnID(t *testing.T) {
	if err := Unseal(quietSealedRepo(t), nil); err == nil {
		t.Fatal("expected a usage error")
	}
}

// quietSealedRepo is sealedRepo with a writer, because Env's zero Out is nil
// and every command that prints would panic on it.
func quietSealedRepo(t *testing.T) Env {
	t.Helper()
	e := sealedRepo(t)
	e.Out = io.Discard
	return e
}
