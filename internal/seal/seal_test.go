package seal

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/fs"
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

func TestEntriesRejectsJunk(t *testing.T) {
	if _, err := Entries("not base64 at all!!"); err == nil {
		t.Fatal("Entries accepted junk")
	}
}

// hostileBlob seals one entry under a name Seal would never produce, because
// Seal writes every name through filepath.ToSlash. Extract's threat model is an
// archive this package did not write: spec §13 answers "sealed blobs are
// unreviewable in a PR diff" with "maintainers undergo reveal locally on the
// branch", so the documented review path runs Extract over a contributor's blob.
func hostileBlob(t *testing.T, name string) string {
	t.Helper()
	var buf bytes.Buffer
	zw, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		t.Fatal(err)
	}
	tw := tar.NewWriter(zw)
	body := []byte("pwned\n")
	hdr := &tar.Header{
		Name: name, Mode: 0o644, Size: int64(len(body)),
		Typeflag: tar.TypeReg, Format: tar.FormatPAX,
	}
	if err := tw.WriteHeader(hdr); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write(body); err != nil {
		t.Fatal(err)
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

// Extract must refuse a path that leaves its destination, on every platform.
//
// The backslash form is the one that actually escaped: path.Clean is slash-only,
// so it read `..\..\x` as a single ordinary filename and filepath.Join then
// resolved it, writing two directories above dest on Windows. The slash form was
// always refused — it is here so the two are held to one rule.
func TestExtractRefusesPathsThatEscapeTheDestination(t *testing.T) {
	for _, name := range []string{
		"../../escaped.txt",
		`..\..\escaped.txt`,
		"/absolute.txt",
		"solution/../../../escaped.txt",
	} {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			dest := filepath.Join(root, "a", "b", "dest")
			if err := os.MkdirAll(dest, 0o755); err != nil {
				t.Fatal(err)
			}

			if err := Extract(hostileBlob(t, name), dest); err == nil {
				t.Errorf("Extract accepted %q", name)
			}

			// An error is not the property under test; containment is. Walk the
			// whole tree above dest and assert nothing landed outside it, so a
			// future Extract that reports an error after writing still fails.
			var strays []string
			if err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
				if err != nil || d.IsDir() {
					return err
				}
				rel, relErr := filepath.Rel(dest, p)
				if relErr != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
					strays = append(strays, p)
				}
				return nil
			}); err != nil {
				t.Fatal(err)
			}
			if len(strays) > 0 {
				t.Errorf("%q wrote outside the destination: %v", name, strays)
			}
		})
	}
}

// The refusal must not cost the ordinary case: a seal's real contents are a
// couple of files at the root and an overlay directory beneath it.
func TestExtractWritesNestedEntries(t *testing.T) {
	blob, err := Seal(fixture(t))
	if err != nil {
		t.Fatal(err)
	}
	dest := filepath.Join(t.TempDir(), "does", "not", "exist", "yet")
	if err := Extract(blob, dest); err != nil {
		t.Fatalf("Extract: %v", err)
	}
	for _, name := range []string{"HINT.md", "EXPLANATION.md", "solution/padding.go"} {
		if _, err := os.Stat(filepath.Join(dest, filepath.FromSlash(name))); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
}
