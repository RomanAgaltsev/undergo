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
