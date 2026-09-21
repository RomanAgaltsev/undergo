package cli

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/RomanAgaltsev/undergo/internal/seal"
)

// AuthoringDir holds the plaintext solution while a task is being written. The
// underscore keeps the Go tool out of it, and .gitignore keeps it uncommitted.
const AuthoringDir = "_solution"

// SealCmd packs a task's _solution/ into solution.sealed and removes the plaintext.
func SealCmd(e Env, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: undergo seal <id>: %w", ErrUsage)
	}
	t, err := find(e, args[0])
	if err != nil {
		return err
	}

	src := filepath.Join(t.Dir, AuthoringDir)
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("%s has no %s/ to seal", t.ID, AuthoringDir)
	}
	blob, err := seal.Seal(src)
	if err != nil {
		return err
	}

	// Prove the blob opens, and holds what went in, before anything is
	// destroyed or overwritten. _solution/ is gitignored, so it exists in
	// exactly one place: a blob that does not decode would take the hint, the
	// explanation, the reference solution and the answer sheet with it, and the
	// loss would stay silent until the next gate 2. The verification also runs
	// before the write, so a bad reseal cannot clobber a good blob.
	if err := verifySeal(src, blob); err != nil {
		return fmt.Errorf("%s: not sealing, and keeping %s/: %w", t.ID, AuthoringDir, err)
	}

	if err := os.WriteFile(filepath.Join(t.Dir, SealName), []byte(blob), 0o644); err != nil {
		return err
	}
	if err := os.RemoveAll(src); err != nil {
		return err
	}
	fmt.Fprintf(e.Out, "sealed %s (%d bytes)\n", t.ID, len(blob))
	return nil
}

// verifySeal reports whether blob decodes back to exactly the tree under src.
func verifySeal(src, blob string) error {
	back, err := seal.Entries(blob)
	if err != nil {
		return fmt.Errorf("the blob does not decode: %w", err)
	}
	var count int
	err = filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(src, p)
		if err != nil {
			return err
		}
		name := filepath.ToSlash(rel)
		// G122: src is the author's own _solution/ in this checkout, walked
		// moments after they asked for it to be sealed.
		want, err := os.ReadFile(p) //nolint:gosec // G122: repo-owned path
		if err != nil {
			return err
		}
		got, ok := back[name]
		if !ok {
			return fmt.Errorf("%s is missing from the blob", name)
		}
		if !bytes.Equal(got, want) {
			return fmt.Errorf("%s does not survive the round trip", name)
		}
		count++
		return nil
	})
	if err != nil {
		return err
	}
	if count != len(back) {
		return fmt.Errorf("the blob holds %d files and %s holds %d", len(back), AuthoringDir, count)
	}
	return nil
}
