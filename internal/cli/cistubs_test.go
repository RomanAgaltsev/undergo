package cli

import (
	"slices"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/manifest"
)

func TestPartitionPackagesVetsOnlyMachineGradedTasks(t *testing.T) {
	tasks := []*manifest.Task{
		{ID: "layout/01-struct-padding", Mode: manifest.ModePredict},
		{ID: "reflect/01-env-decoder", Mode: manifest.ModeBuild},
		{ID: "alloc/01-zero-alloc-join", Mode: manifest.ModeOptimize},
		{ID: "review-concurrency/01-request-counter", Mode: manifest.ModeReview},
		{ID: "design/01-rate-limiter", Mode: manifest.ModeDesign},
	}
	// Every task has Go files except the design kata, which is prose only.
	hasGo := func(pkg string) bool { return pkg != "./tasks/design/01-rate-limiter" }

	build, vet := partitionPackages(tasks, hasGo)

	if slices.Contains(build, "./tasks/design/01-rate-limiter") {
		t.Error("a task with no Go files must not be handed to go build")
	}
	if !slices.Contains(build, "./tasks/review-concurrency/01-request-counter") {
		t.Error("a review drill must still be proven to compile")
	}
	if slices.Contains(vet, "./tasks/review-concurrency/01-request-counter") {
		t.Error("vetting a review drill names its planted defect — that is the answer")
	}
	for _, want := range []string{
		"./tasks/layout/01-struct-padding",
		"./tasks/reflect/01-env-decoder",
		"./tasks/alloc/01-zero-alloc-join",
	} {
		if !slices.Contains(vet, want) {
			t.Errorf("%s must still be vetted", want)
		}
	}
}
