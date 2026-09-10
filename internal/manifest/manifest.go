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
	Go         string   `yaml:"go"`
	MaxGo      string   `yaml:"max_go"`
	Arch       []string `yaml:"arch"`
	OS         []string `yaml:"os"`
	Toolchains []string `yaml:"toolchains"`
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
	Verify     string    `yaml:"verify"`
	Predict    *Predict  `yaml:"predict"`
	Optimize   *Optimize `yaml:"optimize"`
	InspiredBy string    `yaml:"inspired_by"`
	Deprecated bool      `yaml:"deprecated"`

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
