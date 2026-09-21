package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RomanAgaltsev/undergo/internal/seal"
)

// Unseal is the authoring inverse of SealCmd: it unpacks solution.sealed back
// into _solution/ so a shipped hint or explanation can be edited and resealed.
//
// It is not the solver's path — that is `undergo reveal <id> --stuck`, which is
// gated and records a peek. Editing a task is not being stuck, so this one
// neither gates nor records; it only refuses to overwrite plaintext that exists
// nowhere else, _solution/ being gitignored.
func Unseal(e Env, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: undergo unseal <id>: %w", ErrUsage)
	}
	t, err := find(e, args[0])
	if err != nil {
		return err
	}

	dest := filepath.Join(t.Dir, AuthoringDir)
	entries, err := os.ReadDir(dest)
	if err == nil && len(entries) > 0 {
		return fmt.Errorf("%s already has a %s/ — move it aside before unsealing over it", t.ID, AuthoringDir)
	}

	blob, err := readSeal(e, t.ID)
	if err != nil {
		return err
	}
	if err := seal.Extract(blob, dest); err != nil {
		return err
	}
	fmt.Fprintf(e.Out, "unsealed %s into %s/\n", t.ID, AuthoringDir)
	return nil
}
