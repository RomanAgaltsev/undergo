package cli

import (
	"fmt"
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
		return fmt.Errorf("usage: undergo seal <id>")
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
	if err := os.WriteFile(filepath.Join(t.Dir, SealName), []byte(blob), 0o644); err != nil {
		return err
	}
	if err := os.RemoveAll(src); err != nil {
		return err
	}
	fmt.Fprintf(e.Out, "sealed %s (%d bytes)\n", t.ID, len(blob))
	return nil
}
