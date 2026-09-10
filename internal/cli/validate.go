package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RomanAgaltsev/undergo/internal/manifest"
)

// Validate is CI gate 3: every manifest parses and validates, IDs are unique,
// and no plaintext solution has been committed.
func Validate(e Env, _ []string) error {
	tasks, err := manifest.Walk(e.TasksDir())
	if err != nil {
		return err
	}
	seen := map[string]string{}
	for _, t := range tasks {
		if prev, dup := seen[t.ID]; dup {
			return fmt.Errorf("duplicate id %q in %s and %s", t.ID, prev, t.Dir)
		}
		seen[t.ID] = t.Dir

		if _, err := os.Stat(filepath.Join(t.Dir, AuthoringDir)); err == nil {
			return fmt.Errorf("%s: %s/ is committed — run: undergo seal %s", t.ID, AuthoringDir, t.ID)
		}
		if _, err := os.Stat(filepath.Join(t.Dir, "README.md")); err != nil {
			return fmt.Errorf("%s: no README.md", t.ID)
		}

		entries, err := os.ReadDir(t.Dir)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
				continue
			}
			body, err := os.ReadFile(filepath.Join(t.Dir, entry.Name()))
			if err != nil {
				return err
			}
			if broken := brokenRelPaths(t.Dir, string(body)); len(broken) > 0 {
				return fmt.Errorf("%s: %s points at %v, which does not exist",
					t.ID, entry.Name(), broken)
			}
		}
	}
	fmt.Fprintf(e.Out, "%d manifests valid\n", len(tasks))
	return nil
}
