package manifest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTask(t *testing.T, root, rel, body string) {
	t.Helper()
	dir := filepath.Join(root, rel)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "task.yaml"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestWalkFindsTasksSortedByID(t *testing.T) {
	root := t.TempDir()
	writeTask(t, root, "reflect/01-env-decoder", `
schema: 1
id: reflect/01-env-decoder
title: "Env decoder"
mode: build
track: reflect
difficulty: 3
requires: {go: "1.27"}
`)
	writeTask(t, root, "alloc/01-zero-alloc-join", `
schema: 1
id: alloc/01-zero-alloc-join
title: "Zero-alloc join"
mode: optimize
track: alloc
difficulty: 3
requires: {go: "1.27"}
optimize: {baseline: BenchmarkBaseline, target_metric: allocs, target: 1}
`)

	got, err := Walk(root)
	if err != nil {
		t.Fatalf("Walk: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("found %d tasks, want 2", len(got))
	}
	if got[0].ID != "alloc/01-zero-alloc-join" || got[1].ID != "reflect/01-env-decoder" {
		t.Errorf("not sorted by ID: %q, %q", got[0].ID, got[1].ID)
	}
}

func TestWalkRejectsAnInvalidTask(t *testing.T) {
	root := t.TempDir()
	writeTask(t, root, "bad/01-nope", "schema: 99\nid: bad/01-nope\n")
	if _, err := Walk(root); err == nil {
		t.Fatal("Walk accepted an invalid manifest")
	}
}

func TestGradeable(t *testing.T) {
	// The toolchain cases must not reach for the network: stub the resolver.
	origAvailable := ToolchainAvailable
	t.Cleanup(func() { ToolchainAvailable = origAvailable })
	ToolchainAvailable = func(string) bool { return false }

	env := Env{GoVersion: "go1.27.0", GOOS: "linux", GOARCH: "amd64"}

	tests := []struct {
		name string
		req  Requires
		want bool
	}{
		{"no constraints", Requires{}, true},
		{"go floor met", Requires{Go: "1.24"}, true},
		{"go floor unmet", Requires{Go: "1.29"}, false},
		{"go ceiling exceeded", Requires{MaxGo: "1.25"}, false},
		{"arch allowed", Requires{Arch: []string{"amd64", "arm64"}}, true},
		{"arch excluded", Requires{Arch: []string{"arm64"}}, false},
		{"os excluded", Requires{OS: []string{"darwin"}}, false},
		{"toolchain missing", Requires{Toolchains: []string{"go1.21.13"}}, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			task := &Task{Requires: tc.req}
			ok, reason := Gradeable(task, env)
			if ok != tc.want {
				t.Fatalf("Gradeable = %v (%s), want %v", ok, reason, tc.want)
			}
			if !ok && reason == "" {
				t.Error("a refusal must explain itself")
			}
		})
	}
}

// Availability used to mean "on PATH", scanned from a hardcoded list of six
// names that no task ever asked for. It now means "can be obtained", which is
// a question with a network in it — so Gradeable asks through a hook that a
// test can answer instead.
func TestGradeableResolvesToolchainsRatherThanScanningPATH(t *testing.T) {
	orig := ToolchainAvailable
	t.Cleanup(func() { ToolchainAvailable = orig })

	task := &Task{ID: "versions/06-x", Requires: Requires{Toolchains: []string{"go1.22.12"}}}
	env := Env{GoVersion: "go1.27.1", GOOS: "linux", GOARCH: "amd64"}

	ToolchainAvailable = func(string) bool { return true }
	if ok, why := Gradeable(task, env); !ok {
		t.Fatalf("a resolvable toolchain must grade: %s", why)
	}

	ToolchainAvailable = func(string) bool { return false }
	ok, why := Gradeable(task, env)
	if ok {
		t.Fatal("an unobtainable toolchain must skip")
	}
	if !strings.Contains(why, "go1.22.12") {
		t.Errorf("the reason must name the toolchain; got %q", why)
	}
}
