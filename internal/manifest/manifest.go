// Package manifest reads and validates task.yaml files.
package manifest

import (
	"fmt"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v3"
)

// SchemaVersion is the task.yaml schema this build understands.
const SchemaVersion = 1

// Mode is how a task is solved and graded.
type Mode string

// The five task modes.
const (
	ModeBuild    Mode = "build"
	ModePredict  Mode = "predict"
	ModeOptimize Mode = "optimize"
	ModeReview   Mode = "review"
	ModeDesign   Mode = "design"
)

// Requires states the environment a task can be graded in.
type Requires struct {
	// Go is a floor ABOVE the module's own go directive, and is almost always
	// empty.
	//
	// It used to say "1.27" on all 267 tasks, including the 171 prose ones that
	// need no toolchain at all, so it discriminated nothing — and it made
	// `undergo doctor` answer "0 of 267 gradeable" on a Go 1.26 machine, which
	// was both useless and false. It is redundant below the module's floor in
	// any case: go.mod declares a version and GOTOOLCHAIN fetches it, so nobody
	// runs this repo on anything older. Validate rejects a pin at or below that
	// floor, because such a pin cannot discriminate and is indistinguishable
	// from not having thought about the question.
	Go string `yaml:"go"`

	// MaxGo is the counterpart that can discriminate: a ceiling, for a task
	// whose answer stopped holding at some release.
	MaxGo      string   `yaml:"max_go"`
	Arch       []string `yaml:"arch"`
	OS         []string `yaml:"os"`
	Toolchains []string `yaml:"toolchains"`

	// DefaultBuild marks a task whose answers hold only for an unmodified
	// build, so the race gate must leave it alone.
	//
	// The race detector is not a neutral observer: it changes which values
	// the compiler may keep on the stack, and a task measuring an
	// optimisation is then measuring the instrumentation instead. Racing
	// such a task reports a failure that says nothing about the solution.
	//
	// There are three ways it is true, one found per milestone so far.
	//
	// The first is a task measuring a compiler optimisation that -race
	// disables, such as slice stack allocation (types/01-cap-growth).
	//
	// The second is a task whose subject *is* a data race: the memmodel track
	// plants races deliberately, and under -race the detector aborts the
	// program rather than letting it answer the question.
	//
	// The third is a task measuring scheduling order. Instrumenting every
	// memory access inserts work between a goroutine being scheduled and
	// reaching the operation being observed, which reorders goroutines that
	// were not racing at all. sched/01-runnext-ordering measures 4 1 2 3 on an
	// unmodified build and 4 2 1 3 under -race; the runnext slot survives and
	// the queue order behind it does not.
	//
	// Set it only where one of those actually applies, and say why in the
	// task's explanation — a blanket opt-out would hide the concurrency bugs
	// the gate exists to find. A task that merely uses goroutines does not
	// qualify: memmodel/04-seqlock plants no race and is raced like any other.
	DefaultBuild bool `yaml:"default_build"`
}

// Predict configures a predict-mode task.
type Predict struct {
	Slots []string `yaml:"slots"`
}

// Optimize configures an optimize-mode task.
type Optimize struct {
	Baseline     string  `yaml:"baseline"`
	TargetMetric string  `yaml:"target_metric"`
	Target       float64 `yaml:"target"`
}

// Task is one exercise.
type Task struct {
	Schema     int       `yaml:"schema"`
	ID         string    `yaml:"id"`
	Title      string    `yaml:"title"`
	Mode       Mode      `yaml:"mode"`
	Track      string    `yaml:"track"`
	Difficulty int       `yaml:"difficulty"`
	Estimate   string    `yaml:"estimate"`
	Tags       []string  `yaml:"tags"`
	Requires   Requires  `yaml:"requires"`
	Predict    *Predict  `yaml:"predict"`
	Optimize   *Optimize `yaml:"optimize"`
	InspiredBy string    `yaml:"inspired_by"`

	// Deprecated retires a task without renaming or deleting it.
	//
	// Ids are permanent — renaming one invalidates every solver's progress file
	// — so this is the only way out, and §11's SemVer contract rests on it. It
	// is honoured by list (hidden unless --deprecated), by show (which says so
	// before printing the README) and by the README counts (a retired task is
	// not one that ships). Gate 2 still proves it, so a retired task cannot rot
	// into a failing seal unnoticed.
	Deprecated bool `yaml:"deprecated"`

	// Dir is the directory the manifest was loaded from. Not serialised.
	Dir string `yaml:"-"`
}

// Load reads and parses a task.yaml. It does not validate; call Validate.
func Load(path string) (*Task, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var t Task
	if err := yaml.Unmarshal(b, &t); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	t.Dir = filepath.Dir(path)
	return &t, nil
}
