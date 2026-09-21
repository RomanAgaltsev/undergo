// Package seal packs a directory into a single deterministic, opaque blob.
//
// This is obfuscation, not secrecy: anyone can run `base64 -d | tar xz`. Its
// only job is to keep a solution out of casual sight, out of `grep`, and out of
// GitHub code search, so that reading it is a deliberate act.
package seal

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strings"
)

const lineWidth = 76

// Seal packs every file under dir into a line-wrapped base64 blob. Output is
// byte-stable for identical input: entries are sorted and all timestamps,
// ownership and permission variation are dropped.
func Seal(dir string) (string, error) {
	var names []string
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(dir, p)
		if err != nil {
			return err
		}
		names = append(names, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		return "", err
	}
	if len(names) == 0 {
		return "", fmt.Errorf("seal: %s is empty", dir)
	}
	sort.Strings(names)

	var buf bytes.Buffer
	zw, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return "", err
	}
	tw := tar.NewWriter(zw)
	for _, name := range names {
		body, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(name)))
		if err != nil {
			return "", err
		}
		hdr := &tar.Header{
			Name:     name,
			Mode:     0o644,
			Size:     int64(len(body)),
			Typeflag: tar.TypeReg,
			Format:   tar.FormatUSTAR,
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return "", err
		}
		if _, err := tw.Write(body); err != nil {
			return "", err
		}
	}
	if err := tw.Close(); err != nil {
		return "", err
	}
	if err := zw.Close(); err != nil {
		return "", err
	}

	enc := base64.StdEncoding.EncodeToString(buf.Bytes())
	var out strings.Builder
	for i := 0; i < len(enc); i += lineWidth {
		out.WriteString(enc[i:min(i+lineWidth, len(enc))])
		out.WriteByte('\n')
	}
	return out.String(), nil
}

// Entries decodes a blob into its files.
func Entries(blob string) (map[string][]byte, error) {
	raw, err := base64.StdEncoding.DecodeString(strings.Join(strings.Fields(blob), ""))
	if err != nil {
		return nil, fmt.Errorf("seal: not a sealed blob: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("seal: %w", err)
	}
	defer func() { _ = zr.Close() }()

	out := map[string][]byte{}
	tr := tar.NewReader(zr)
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("seal: %w", err)
		}
		body, err := io.ReadAll(tr)
		if err != nil {
			return nil, fmt.Errorf("seal: %w", err)
		}
		out[hdr.Name] = body
	}
	return out, nil
}

// File returns one entry from a blob.
func File(blob, name string) ([]byte, error) {
	entries, err := Entries(blob)
	if err != nil {
		return nil, err
	}
	body, ok := entries[name]
	if !ok {
		return nil, fmt.Errorf("seal: %s is not in this blob", name)
	}
	return body, nil
}

// Extract writes every entry under dest, refusing any path that escapes it.
//
// A blob is untrusted input. Spec §13 answers "sealed blobs are unreviewable in
// a PR diff" with "maintainers undergo reveal locally on the branch" — which is
// to say the documented way to review a stranger's contribution is to run this
// over an archive that stranger wrote. Containment has to hold against a hostile
// entry name, not only a well-formed one.
//
// os.Root does the refusing, and the hand-rolled predecessor is why. It cleaned
// the name with path.Clean("/"+name), which is the slash-only cleaner: it reads a
// backslash as an ordinary filename character, so `..\..\x` survived it intact
// and filepath.Join — which does treat a backslash as a separator on Windows —
// then resolved the escape. A blob wrote two directories above dest on the
// maintainer's own platform. os.Root applies the operating system's real path
// semantics instead, and refuses a symlink out of the tree as well.
func Extract(blob, dest string) error {
	entries, err := Entries(blob)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return err
	}
	root, err := os.OpenRoot(dest)
	if err != nil {
		return err
	}
	defer func() { _ = root.Close() }()

	// Sorted so that a failure names the same entry on every platform, for the
	// same reason Seal sorts: a map's order is not a fact about the blob.
	for _, name := range slices.Sorted(maps.Keys(entries)) {
		if err := writeEntry(root, name, entries[name]); err != nil {
			return fmt.Errorf("seal: %s: %w", name, err)
		}
	}
	return nil
}

// writeEntry creates one file beneath root.
//
// Seal writes every name with filepath.ToSlash, so a backslash can only reach
// here from an archive this package did not produce. Refusing it keeps the
// meaning of a name identical on every platform — on Windows `a\b` is a nested
// path and on Linux it is a filename — rather than leaving os.Root to decide.
func writeEntry(root *os.Root, name string, body []byte) error {
	if strings.ContainsRune(name, '\\') {
		return errors.New("a blob path is slash-separated; refusing a backslash")
	}
	if dir := path.Dir(name); dir != "." && dir != "/" {
		if err := root.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	f, err := root.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	if _, err := f.Write(body); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}
