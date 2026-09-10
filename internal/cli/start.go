package cli

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/RomanAgaltsev/undergo/internal/manifest"
	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// SealName is the sealed solution file inside a task directory.
const SealName = "solution.sealed"

// Start copies a task into the solver's gitignored work directory.
func Start(e Env, args []string) error {
	fs2 := flag.NewFlagSet("start", flag.ContinueOnError)
	fs2.SetOutput(e.Err)
	force := fs2.Bool("force", false, "overwrite existing work")
	id, flags := splitID(args)
	if err := fs2.Parse(flags); err != nil {
		return err
	}
	if id == "" || fs2.NArg() != 0 {
		return fmt.Errorf("usage: undergo start <id> [--force]")
	}

	t, err := find(e, id)
	if err != nil {
		return err
	}
	work := e.WorkDir(t.ID)
	if _, err := os.Stat(work); err == nil && !*force {
		return fmt.Errorf("%s already exists — pass --force to overwrite your work", work)
	}
	if err := copyTask(t.Dir, work); err != nil {
		return err
	}
	if t.Mode == manifest.ModePredict {
		if err := writePredictionStub(work, t); err != nil {
			return err
		}
	}
	fmt.Fprintf(e.Out, "→ %s\nEdit there, then: undergo verify %s\n", work, t.ID)
	return nil
}

// copyTask copies everything except the sealed solution.
func copyTask(src, dst string) error {
	return filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, p)
		if err != nil {
			return err
		}
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(dst, rel), 0o755)
		}
		if d.Name() == SealName {
			return nil
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(dst, rel), b, 0o644)
	})
}

// writePredictionStub creates an empty answer sheet with one line per slot.
func writePredictionStub(work string, t *manifest.Task) error {
	var b strings.Builder
	b.WriteString("# Your predictions. Fill in every slot, then: undergo verify " + t.ID + "\n")
	b.WriteString("# Guess first. Running the test before you commit an answer is the one\n")
	b.WriteString("# way to waste this task.\n\n")
	for _, slot := range t.Predict.Slots {
		b.WriteString(slot + ":\n")
	}
	return os.WriteFile(filepath.Join(work, predict.DefaultFile), []byte(b.String()), 0o644)
}
