package manifest

import "testing"

func valid() *Task {
	return &Task{
		Schema: 1, ID: "layout/01-struct-padding", Title: "Padding",
		Mode: ModePredict, Track: "layout", Difficulty: 2,
		Requires: Requires{Go: "1.27"},
		Predict:  &Predict{Slots: []string{"a"}},
		Dir:      "tasks/layout/01-struct-padding",
	}
}

func TestValidateAcceptsAGoodTask(t *testing.T) {
	if err := Validate(valid()); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func TestValidateRejects(t *testing.T) {
	tests := map[string]func(*Task){
		"wrong schema":       func(x *Task) { x.Schema = 2 },
		"bad id":             func(x *Task) { x.ID = "Layout/1_padding" },
		"id disagrees":       func(x *Task) { x.Track = "alloc" },
		"unknown mode":       func(x *Task) { x.Mode = "guess" },
		"difficulty range":   func(x *Task) { x.Difficulty = 9 },
		"predict no slots":   func(x *Task) { x.Predict = &Predict{} },
		"predict block gone": func(x *Task) { x.Predict = nil },
		"stray optimize":     func(x *Task) { x.Optimize = &Optimize{Baseline: "B"} },
		"bad go version":     func(x *Task) { x.Requires.Go = "tomorrow" },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			task := valid()
			mutate(task)
			if err := Validate(task); err == nil {
				t.Fatal("Validate accepted an invalid task")
			}
		})
	}
}

func TestValidateOptimizeMetric(t *testing.T) {
	task := valid()
	task.Mode = ModeOptimize
	task.Predict = nil
	task.Optimize = &Optimize{Baseline: "BenchmarkBaseline", TargetMetric: "binary_size"}
	if err := Validate(task); err == nil {
		t.Fatal("binary_size is not supported in v1 and must be rejected")
	}
	task.Optimize.TargetMetric = "allocs"
	if err := Validate(task); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

// A review drill lives at tasks/review/<category>/<NN>-slug, so a task id may
// carry one grouping segment between the track root and the task itself. The
// track is then the whole path minus the task, which keeps the track/id
// agreement rule below unchanged.
func TestValidateAcceptsAGroupedTrack(t *testing.T) {
	task := valid()
	task.ID = "review/concurrency/01-request-counter"
	task.Track = "review/concurrency"
	task.Mode = ModeReview
	task.Predict = nil

	if err := Validate(task); err != nil {
		t.Fatalf("Validate rejected a grouped id: %v", err)
	}
}

func TestValidateRejectsIDShapes(t *testing.T) {
	tests := map[string]string{
		"three grouping segments": "review/go/concurrency/01-request-counter",
		"no task segment":         "review/concurrency",
		"upper case":              "review/Concurrency/01-request-counter",
		"one-digit number":        "review/concurrency/1-request-counter",
		"trailing slash":          "review/concurrency/01-request-counter/",
	}
	for name, id := range tests {
		t.Run(name, func(t *testing.T) {
			task := valid()
			task.ID = id
			task.Track = "review/concurrency"
			task.Mode = ModeReview
			task.Predict = nil
			if err := Validate(task); err == nil {
				t.Fatalf("Validate accepted %q", id)
			}
		})
	}
}

// A requires block used to be checked for shape and not for values, so
// `arch: [amd46]` validated and produced a task no machine in the world could
// build or grade — skipped everywhere and reported to nobody.
func TestValidateChecksRequiresValues(t *testing.T) {
	tests := []struct {
		name       string
		requires   Requires
		wantReject bool
	}{
		{name: "nothing pinned"},
		{name: "real arch and os", requires: Requires{Arch: []string{"amd64"}, OS: []string{"linux"}}},
		{
			name:       "a plausible GOARCH typo",
			requires:   Requires{Arch: []string{"amd46"}},
			wantReject: true,
		},
		{
			name:       "a GOOS that does not exist",
			requires:   Requires{OS: []string{"linus"}},
			wantReject: true,
		},
		{
			name:     "a full toolchain name",
			requires: Requires{Toolchains: []string{"go1.22.12"}},
		},
		{
			// A language version is not a toolchain. Getting this wrong makes
			// the toolchain unobtainable rather than the manifest invalid, and
			// the error would then blame the network.
			name:       "a language version where a toolchain belongs",
			requires:   Requires{Toolchains: []string{"go1.22"}},
			wantReject: true,
		},
		{
			// GOTOOLCHAIN is itself a Go 1.21 feature, so nothing older can be
			// selected this way.
			name:       "a toolchain older than GOTOOLCHAIN itself",
			requires:   Requires{Toolchains: []string{"go1.19.13"}},
			wantReject: true,
		},
		{
			name:       "a floor above its own ceiling",
			requires:   Requires{Go: "1.28", MaxGo: "1.27"},
			wantReject: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			task := valid()
			task.Requires = tc.requires
			err := Validate(task)
			if tc.wantReject && err == nil {
				t.Fatalf("Validate accepted %+v", tc.requires)
			}
			if !tc.wantReject && err != nil {
				t.Fatalf("Validate rejected %+v: %v", tc.requires, err)
			}
		})
	}
}
