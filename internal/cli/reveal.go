package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/RomanAgaltsev/undergo/internal/progress"
	"github.com/RomanAgaltsev/undergo/internal/seal"
)

// Names of the two rungs inside a sealed blob.
const (
	HintFile        = "HINT.md"
	ExplanationFile = "EXPLANATION.md"
)

func readSeal(e Env, id string) (string, error) {
	t, err := find(e, id)
	if err != nil {
		return "", err
	}
	b, err := os.ReadFile(filepath.Join(t.Dir, SealName))
	if err != nil {
		return "", fmt.Errorf("%s has no sealed solution yet", id)
	}
	return string(b), nil
}

// Hint prints the first rung. It is never gated: being stuck is not a failure.
func Hint(e Env, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: undergo hint <id>")
	}
	id := args[0]
	blob, err := readSeal(e, id)
	if err != nil {
		return err
	}
	body, err := seal.File(blob, HintFile)
	if err != nil {
		return err
	}
	if _, err := e.Out.Write(body); err != nil {
		return err
	}

	rec, err := progress.Load(e.ProgressPath())
	if err != nil {
		return err
	}
	rec.RecordHint(id)
	return rec.Save(e.ProgressPath())
}

// Reveal opens the full solution. It refuses while the task is unsolved unless
// --stuck is given, and records that override so the solver's record stays honest.
func Reveal(e Env, args []string) error {
	fs := flag.NewFlagSet("reveal", flag.ContinueOnError)
	fs.SetOutput(e.Err)
	stuck := fs.Bool("stuck", false, "reveal even though the task is not passing")
	id, flags := splitID(args)
	if err := fs.Parse(flags); err != nil {
		return err
	}
	if id == "" || fs.NArg() != 0 {
		return fmt.Errorf("usage: undergo reveal <id> [--stuck]")
	}

	rec, err := progress.Load(e.ProgressPath())
	if err != nil {
		return err
	}
	solved := rec.Get(id).Solved
	if !solved && !*stuck {
		return fmt.Errorf("%s is not passing yet.\n"+
			"Try: undergo hint %s\n"+
			"If you are genuinely stuck: undergo reveal %s --stuck (recorded as a peek)", id, id, id)
	}

	blob, err := readSeal(e, id)
	if err != nil {
		return err
	}
	dest := e.RevealDir(id)
	if err := seal.Extract(blob, dest); err != nil {
		return err
	}
	if body, err := seal.File(blob, ExplanationFile); err == nil {
		if _, err := e.Out.Write(body); err != nil {
			return err
		}
	}
	fmt.Fprintf(e.Out, "\n→ %s\n", dest)

	if !solved {
		rec.RecordPeek(id)
		return rec.Save(e.ProgressPath())
	}
	return nil
}
