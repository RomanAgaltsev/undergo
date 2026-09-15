package cli

import (
	"slices"
	"strings"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/manifest"
)

func TestPartitionPackagesVetsOnlyMachineGradedTasks(t *testing.T) {
	tasks := []*manifest.Task{
		{ID: "layout/01-struct-padding", Mode: manifest.ModePredict},
		{ID: "reflect/01-env-decoder", Mode: manifest.ModeBuild},
		{ID: "alloc/01-zero-alloc-join", Mode: manifest.ModeOptimize},
		{ID: "review/concurrency/01-request-counter", Mode: manifest.ModeReview},
		{ID: "design/01-rate-limiter", Mode: manifest.ModeDesign},
	}
	// Every task has Go files except the design kata, which is prose only.
	hasGo := func(pkg string) bool { return pkg != "./tasks/design/01-rate-limiter" }

	build, vet, _ := partitionPackages(tasks, amd64Linux, hasGo)

	if slices.Contains(build, "./tasks/design/01-rate-limiter") {
		t.Error("a task with no Go files must not be handed to go build")
	}
	if !slices.Contains(build, "./tasks/review/concurrency/01-request-counter") {
		t.Error("a review drill must still be proven to compile")
	}
	if slices.Contains(vet, "./tasks/review/concurrency/01-request-counter") {
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

// amd64Linux is a host that can build anything not pinned elsewhere.
var amd64Linux = manifest.Env{GOARCH: "amd64", GOOS: "linux", GoVersion: "go1.27.1"}

// requires.arch gated gate 2 but not gate 1, so an amd64-only .s failed to
// build on the arm64 runner. Gate 1 now skips what the host cannot build.
func TestPartitionSkipsForeignArch(t *testing.T) {
	tasks := []*manifest.Task{
		{ID: "layout/01-x", Mode: manifest.ModePredict},
		{
			ID: "asm/03-y", Mode: manifest.ModePredict,
			Requires: manifest.Requires{Arch: []string{"amd64"}},
		},
	}
	env := manifest.Env{GOARCH: "arm64", GOOS: "linux", GoVersion: "go1.27.1"}

	build, vet, skipped := partitionPackages(tasks, env, func(string) bool { return true })

	if slices.Contains(build, "./tasks/asm/03-y") {
		t.Errorf("an amd64-only task must not be built on arm64: %v", build)
	}
	if !slices.Contains(build, "./tasks/layout/01-x") {
		t.Errorf("a portable task must still be built: %v", build)
	}
	if slices.Contains(vet, "./tasks/asm/03-y") {
		t.Errorf("an unbuildable task must not be vetted either: %v", vet)
	}
	if len(skipped) != 1 || !strings.Contains(skipped[0], "asm/03-y") {
		t.Errorf("the skip must be reported, got %v", skipped)
	}
}

func TestPartitionBuildsMatchingArchAndOS(t *testing.T) {
	tasks := []*manifest.Task{
		{
			ID: "asm/03-y", Mode: manifest.ModePredict,
			Requires: manifest.Requires{Arch: []string{"amd64"}},
		},
		{
			ID: "edges/07-w", Mode: manifest.ModePredict,
			Requires: manifest.Requires{OS: []string{"windows"}},
		},
	}

	build, _, skipped := partitionPackages(tasks, amd64Linux, func(string) bool { return true })
	if !slices.Contains(build, "./tasks/asm/03-y") {
		t.Errorf("an amd64 task must be built on amd64: %v", build)
	}
	if slices.Contains(build, "./tasks/edges/07-w") {
		t.Errorf("a windows-only task must not be built on linux: %v", build)
	}
	if len(skipped) != 1 {
		t.Errorf("exactly the windows task should be skipped, got %v", skipped)
	}
}

// Only arch and OS may gate building. A task needing an extra toolchain still
// builds with the default one, and a task needing a newer Go must fail loudly
// rather than quietly vanish from the gate.
func TestPartitionIgnoresToolchainAndGoVersion(t *testing.T) {
	tasks := []*manifest.Task{
		{
			ID: "edges/09-z", Mode: manifest.ModePredict,
			Requires: manifest.Requires{Go: "1.99", Toolchains: []string{"go1.21.13"}},
		},
	}

	build, _, skipped := partitionPackages(tasks, amd64Linux, func(string) bool { return true })
	if !slices.Contains(build, "./tasks/edges/09-z") {
		t.Errorf("only arch and OS may gate building, got build=%v skipped=%v", build, skipped)
	}
	if len(skipped) != 0 {
		t.Errorf("nothing should be skipped, got %v", skipped)
	}
}
