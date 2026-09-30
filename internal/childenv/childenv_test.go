package childenv_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/childenv"
)

// value returns the last value env gives name, which is the one a child sees.
func value(env []string, name string) (string, bool) {
	v, found := "", false
	for _, kv := range env {
		k, val, _ := strings.Cut(kv, "=")
		if strings.EqualFold(k, name) {
			v, found = val, true
		}
	}
	return v, found
}

// Environ drops the caller's GOFLAGS and GODEBUG and sets GOFLAGS to a
// default of its own: an empty GOFLAGS is not a reset, the go command then
// reads whatever `go env -w` stored.
func TestEnvironNeutralisesGOFLAGSAndGODEBUG(t *testing.T) {
	t.Setenv("GOFLAGS", "-gcflags=all=-N")
	t.Setenv("GODEBUG", "gctrace=1")
	env := childenv.Environ()
	if v, _ := value(env, "GOFLAGS"); v != "-buildvcs=auto" {
		t.Errorf("GOFLAGS = %q, want -buildvcs=auto", v)
	}
	if v, ok := value(env, "GODEBUG"); ok {
		t.Errorf("GODEBUG survived as %q", v)
	}
	if n := slices.IndexFunc(env, func(kv string) bool { return strings.HasPrefix(kv, "PATH=") || strings.HasPrefix(kv, "Path=") }); n < 0 {
		t.Error("the rest of the environment was dropped")
	}
}

// A task that needs its own setting passes it as extra, and it wins.
func TestEnvironKeepsExtraLast(t *testing.T) {
	t.Setenv("GODEBUG", "gctrace=1")
	env := childenv.Environ("GODEBUG=asyncpreemptoff=1")
	if v, _ := value(env, "GODEBUG"); v != "asyncpreemptoff=1" {
		t.Errorf("GODEBUG = %q, want the task's own asyncpreemptoff=1", v)
	}
}
