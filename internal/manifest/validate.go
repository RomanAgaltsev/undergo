package manifest

import (
	"fmt"
	"go/version"
	"path"
	"regexp"
	"slices"
)

// A task id is track/NN-slug, where the track may carry one grouping segment:
// review/concurrency/01-request-counter. The track is always the id minus its
// final segment, which is what Validate checks below.
var idRE = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*(?:/[a-z0-9]+(?:-[a-z0-9]+)*)?/[0-9]{2}-[a-z0-9]+(?:-[a-z0-9]+)*$`)

// Metrics supported by internal/optimize in v1. binary_size is specified but
// not yet implemented, so it is rejected rather than silently ignored.
var okMetrics = []string{"ns", "allocs", "bytes"}

var okModes = []Mode{ModeBuild, ModePredict, ModeOptimize, ModeReview, ModeDesign}

// Validate reports whether a task is well formed.
func Validate(t *Task) error {
	if t.Schema != SchemaVersion {
		return fmt.Errorf("schema = %d, want %d", t.Schema, SchemaVersion)
	}
	if !idRE.MatchString(t.ID) {
		return fmt.Errorf("id %q must look like track/NN-slug or track/group/NN-slug, lower-case", t.ID)
	}
	if t.Title == "" {
		return fmt.Errorf("%s: title is empty", t.ID)
	}
	if !slices.Contains(okModes, t.Mode) {
		return fmt.Errorf("%s: unknown mode %q", t.ID, t.Mode)
	}
	if want := t.Track + "/" + path.Base(t.ID); want != t.ID {
		return fmt.Errorf("id %q disagrees with track %q", t.ID, t.Track)
	}
	if t.Difficulty < 1 || t.Difficulty > 5 {
		return fmt.Errorf("%s: difficulty %d outside 1..5", t.ID, t.Difficulty)
	}
	if t.Requires.Go != "" && !version.IsValid("go"+t.Requires.Go) {
		return fmt.Errorf("%s: requires.go %q is not a Go version", t.ID, t.Requires.Go)
	}
	if t.Requires.MaxGo != "" && !version.IsValid("go"+t.Requires.MaxGo) {
		return fmt.Errorf("%s: requires.max_go %q is not a Go version", t.ID, t.Requires.MaxGo)
	}

	switch t.Mode {
	case ModePredict:
		if t.Predict == nil || len(t.Predict.Slots) == 0 {
			return fmt.Errorf("%s: predict mode needs predict.slots", t.ID)
		}
	case ModeOptimize:
		if t.Optimize == nil || t.Optimize.Baseline == "" {
			return fmt.Errorf("%s: optimize mode needs optimize.baseline", t.ID)
		}
		if !slices.Contains(okMetrics, t.Optimize.TargetMetric) {
			return fmt.Errorf("%s: target_metric %q, want one of %v",
				t.ID, t.Optimize.TargetMetric, okMetrics)
		}
	}
	if t.Mode != ModePredict && t.Predict != nil {
		return fmt.Errorf("%s: predict block on a %s task", t.ID, t.Mode)
	}
	if t.Mode != ModeOptimize && t.Optimize != nil {
		return fmt.Errorf("%s: optimize block on a %s task", t.ID, t.Mode)
	}
	return nil
}
