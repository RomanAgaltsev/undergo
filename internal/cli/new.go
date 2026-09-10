package cli

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/RomanAgaltsev/undergo/internal/manifest"
)

// New scaffolds a task directory plus its authoring _solution/ skeleton.
func New(e Env, args []string) error {
	fs := flag.NewFlagSet("new", flag.ContinueOnError)
	fs.SetOutput(e.Err)
	id := fs.String("id", "", "task id, e.g. alloc/01-zero-alloc-join")
	mode := fs.String("mode", "", "build|predict|optimize|review|design")
	title := fs.String("title", "", "one-line title")
	difficulty := fs.Int("difficulty", 3, "1..5")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *id == "" || *mode == "" || *title == "" {
		return fmt.Errorf("usage: undergo new --id track/NN-slug --mode M --title T [--difficulty N]")
	}

	track, _, ok := strings.Cut(*id, "/")
	if !ok {
		return fmt.Errorf("id %q must be track/NN-slug", *id)
	}
	dir := e.TaskDir(*id)
	if _, err := os.Stat(dir); err == nil {
		return fmt.Errorf("%s already exists", dir)
	}
	if err := os.MkdirAll(filepath.Join(dir, "_solution"), 0o755); err != nil {
		return err
	}

	var extra string
	switch manifest.Mode(*mode) {
	case manifest.ModePredict:
		extra = "predict:\n  slots: [replace_me]\n"
	case manifest.ModeOptimize:
		extra = "optimize:\n  baseline: BenchmarkBaseline\n  target_metric: allocs\n  target: 1\n"
	}

	files := map[string]string{
		"task.yaml": fmt.Sprintf(`schema: 1
id: %s
title: %q
mode: %s
track: %s
difficulty: %d
estimate: 30m
tags: []
requires:
  go: "1.27"
verify: "go test ./..."
%sinspired_by: ""
deprecated: false
`, *id, *title, *mode, track, *difficulty, extra),

		"README.md": fmt.Sprintf(`# %s

%s

## The question

Describe what the solver must work out. State it so that guessing is possible
and being vaguely right is not.

## Questions to answer in writing

1.
2.
3.
`, path.Base(*id), *title),

		"_solution/HINT.md":        "One nudge. Name the mechanism to look at, not the answer.\n",
		"_solution/EXPLANATION.md": "Why the answer is what it is, in terms of the mechanism.\n",
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(dir, filepath.FromSlash(name)), []byte(body), 0o644); err != nil {
			return err
		}
	}

	fmt.Fprintf(e.Out, "→ %s\nWrite the stub, the frozen tests and _solution/, then: undergo seal %s\n", dir, *id)
	return nil
}
