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
	"os"
	"path"
	"path/filepath"
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
func Extract(blob, dest string) error {
	entries, err := Entries(blob)
	if err != nil {
		return err
	}
	for name, body := range entries {
		clean := path.Clean("/" + name)
		target := filepath.Join(dest, filepath.FromSlash(strings.TrimPrefix(clean, "/")))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(target, body, 0o644); err != nil {
			return err
		}
	}
	return nil
}
