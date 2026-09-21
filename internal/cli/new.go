package cli

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"

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
		return fmt.Errorf("usage: undergo new --id track/NN-slug --mode M --title T [--difficulty N]: %w", ErrUsage)
	}

	// The track is the id minus its FINAL segment, which is what Validate
	// checks. Cutting at the first slash instead gave a grouped id the track
	// "review" where "review/concurrency" was wanted, so `undergo new --id
	// review/concurrency/NN-slug` — the documented way to add any of the 135
	// review drills — scaffolded a manifest that gate 3 then rejected. Nobody
	// hit it because those 135 came from the importer, not from here.
	track := path.Dir(*id)
	if track == "." || track == "/" {
		return fmt.Errorf("id %q must be track/NN-slug or track/group/NN-slug", *id)
	}

	// Build the manifest as a value and validate it before touching the disk. A
	// scaffold that has to be hand-repaired before `task ci` passes is not a
	// scaffold, and this also catches a bad --mode, an out-of-range
	// --difficulty and a slug that does not match the id grammar — all of which
	// used to be written out happily and rejected later by gate 3.
	t := &manifest.Task{
		Schema: manifest.SchemaVersion, ID: *id, Title: *title,
		Mode: manifest.Mode(*mode), Track: track, Difficulty: *difficulty,
	}
	switch t.Mode {
	case manifest.ModePredict:
		t.Predict = &manifest.Predict{Slots: []string{"replace_me"}}
	case manifest.ModeOptimize:
		t.Optimize = &manifest.Optimize{
			Baseline: "BenchmarkBaseline", TargetMetric: "allocs", Target: 1,
		}
	}
	if err := manifest.Validate(t); err != nil {
		return fmt.Errorf("those flags do not describe a valid task: %w", err)
	}

	dir := e.TaskDir(*id)
	if _, err := os.Stat(dir); err == nil {
		return fmt.Errorf("%s already exists", dir)
	}
	if err := os.MkdirAll(filepath.Join(dir, "_solution"), 0o755); err != nil {
		return err
	}

	var extra string
	switch t.Mode {
	case manifest.ModePredict:
		extra = "predict:\n  slots: [replace_me]\n"
	case manifest.ModeOptimize:
		extra = "optimize:\n  baseline: BenchmarkBaseline\n  target_metric: allocs\n  target: 1\n"
	}

	// No requires block and no verify line. The old template hardcoded
	// `go: "1.27"`, which is the module's own floor and so ruled nothing out —
	// that is how all 267 manifests came to carry the same meaningless pin — and
	// a `verify:` that no code has ever read.
	files := map[string]string{
		"task.yaml": fmt.Sprintf(`schema: 1
id: %s
title: %q
mode: %s
track: %s
difficulty: %d
estimate: 30m
tags: []
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
