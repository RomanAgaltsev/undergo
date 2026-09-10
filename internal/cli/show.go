package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RomanAgaltsev/undergo/internal/manifest"
)

// find returns the task with this ID.
func find(e Env, id string) (*manifest.Task, error) {
	tasks, err := manifest.Walk(e.TasksDir())
	if err != nil {
		return nil, err
	}
	for _, t := range tasks {
		if t.ID == id {
			return t, nil
		}
	}
	return nil, fmt.Errorf("no task %q — try: undergo list", id)
}

// Show prints a task's manifest summary and its README.
func Show(e Env, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: undergo show <id>")
	}
	t, err := find(e, args[0])
	if err != nil {
		return err
	}

	fmt.Fprintf(e.Out, "%s  [%s, track %s, difficulty %d]\n%s\n\n",
		t.ID, t.Mode, t.Track, t.Difficulty, t.Title)
	if ok, why := manifest.Gradeable(t, manifest.CurrentEnv()); !ok {
		fmt.Fprintf(e.Out, "NOT GRADEABLE HERE: %s\n\n", why)
	}
	readme, err := os.ReadFile(filepath.Join(t.Dir, "README.md"))
	if err != nil {
		return err
	}
	_, err = e.Out.Write(readme)
	return err
}
