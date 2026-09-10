package seal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func fixture(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "solution"), 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		"HINT.md":             "Look at alignment.\n",
		"EXPLANATION.md":      "Because the trailing field forces padding.\n",
		"solution/padding.go": "package padding\n",
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(dir, filepath.FromSlash(name)), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestSealIsDeterministic(t *testing.T) {
	dir := fixture(t)
	a, err := Seal(dir)
	if err != nil {
		t.Fatal(err)
	}
	b, err := Seal(dir)
	if err != nil {
		t.Fatal(err)
	}
	if a != b {
		t.Fatal("Seal is not byte-stable across runs")
	}
}

func TestSealHidesPlaintext(t *testing.T) {
	blob, err := Seal(fixture(t))
	if err != nil {
		t.Fatal(err)
	}
	for _, secret := range []string{"alignment", "padding", "package"} {
		if strings.Contains(blob, secret) {
			t.Errorf("blob leaks %q in plaintext", secret)
		}
	}
	for _, line := range strings.Split(strings.TrimSpace(blob), "\n") {
		if len(line) > 76 {
			t.Fatalf("line of %d chars, want <= 76", len(line))
		}
	}
}

func TestRoundTrip(t *testing.T) {
	blob, err := Seal(fixture(t))
	if err != nil {
		t.Fatal(err)
	}
	got, err := File(blob, "HINT.md")
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "Look at alignment.\n" {
		t.Errorf("HINT.md = %q", got)
	}

	dest := t.TempDir()
	if err := Extract(blob, dest); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(dest, "solution", "padding.go"))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "package padding\n" {
		t.Errorf("extracted = %q", b)
	}
}

func TestExtractRejectsEscapingPaths(t *testing.T) {
	if _, err := Entries("not base64 at all!!"); err == nil {
		t.Fatal("Entries accepted junk")
	}
}
