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

// knownArches and knownOSes are the values Buildable and Gradeable compare a
// task's pins against.
//
// A typo here is not a syntax error and used to pass validation: `arch: [amd46]`
// produced a task that no machine in the world could build or grade, skipped
// everywhere and reported to nobody. The lists are the platforms Go supports
// that a task here could plausibly want; extend them rather than working around
// them.
var (
	knownArches = []string{
		"386", "amd64", "arm", "arm64", "loong64", "mips", "mips64",
		"mips64le", "mipsle", "ppc64", "ppc64le", "riscv64", "s390x", "wasm",
	}
	knownOSes = []string{
		"aix", "android", "darwin", "dragonfly", "freebsd", "illumos", "ios",
		"js", "linux", "netbsd", "openbsd", "plan9", "solaris", "wasip1", "windows",
	}
)

// minToolchain is the first release GOTOOLCHAIN can select, because GOTOOLCHAIN
// is itself a Go 1.21 feature. internal/toolchain enforces this at run time;
// here it is enforced where a task author meets it.
const minToolchain = "go1.21"

// validateRequires checks the values of a requires block, not merely its shape.
func validateRequires(t *Task) error {
	for _, a := range t.Requires.Arch {
		if !slices.Contains(knownArches, a) {
			return fmt.Errorf("%s: requires.arch %q is not a GOARCH", t.ID, a)
		}
	}
	for _, o := range t.Requires.OS {
		if !slices.Contains(knownOSes, o) {
			return fmt.Errorf("%s: requires.os %q is not a GOOS", t.ID, o)
		}
	}
	for _, name := range t.Requires.Toolchains {
		// A toolchain is named in full — go1.22.12, not go1.22. Getting this
		// wrong makes the toolchain unobtainable rather than the manifest
		// invalid, and Gradeable's message would then blame the network.
		if !version.IsValid(name) || version.Lang(name) == name {
			return fmt.Errorf("%s: requires.toolchains %q is not a full toolchain name "+
				"(want go1.22.12, not go1.22)", t.ID, name)
		}
		if version.Compare(version.Lang(name), minToolchain) < 0 {
			return fmt.Errorf("%s: requires.toolchains %q predates GOTOOLCHAIN, which arrived in %s",
				t.ID, name, minToolchain)
		}
	}
	if t.Requires.Go != "" && t.Requires.MaxGo != "" &&
		version.Compare("go"+t.Requires.Go, "go"+t.Requires.MaxGo) > 0 {
		return fmt.Errorf("%s: requires.go %s is above requires.max_go %s, so nothing satisfies it",
			t.ID, t.Requires.Go, t.Requires.MaxGo)
	}
	return nil
}

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
	if err := validateRequires(t); err != nil {
		return err
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
