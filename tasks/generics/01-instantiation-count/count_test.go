package shapes

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// symbols builds this package into an archive and returns every symbol name.
//
// An archive, not a binary: linking drops instantiations nothing calls, which
// silently changes the answer. Nothing here pins -gcflags either, because the
// counts are the same with and without inlining.
func symbols(t *testing.T) []string {
	t.Helper()

	archive := filepath.Join(t.TempDir(), "shapes.a")
	if out, err := exec.Command("go", "build", "-o", archive, ".").CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
	out, err := exec.Command("go", "tool", "nm", archive).CombinedOutput()
	if err != nil {
		t.Fatalf("go tool nm: %v\n%s", err, out)
	}

	var names []string
	for line := range strings.SplitSeq(string(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		names = append(names, fields[len(fields)-1])
	}
	if len(names) == 0 {
		t.Fatal("nm reported no symbols")
	}
	return names
}

// counts classifies every symbol mentioning name into the three kinds that
// matter: a dictionary, a GC-shape body, and a fully concrete wrapper.
//
// The exclusions are not decoration. nm also reports .arginfo1, .argliveinfo
// and go:info DWARF entries per instantiation, and counting those inflates
// every answer roughly threefold.
func counts(names []string, name string) (shapeBodies, concreteWrappers, dicts int) {
	seen := map[string]bool{}
	for _, sym := range names {
		switch {
		case !strings.Contains(sym, name),
			strings.HasPrefix(sym, "go:info."),
			strings.Contains(sym, ".arginfo"),
			strings.Contains(sym, ".argliveinfo"),
			seen[sym]:
			continue
		}
		seen[sym] = true

		switch {
		case strings.Contains(sym, "..dict."):
			dicts++
		case strings.Contains(sym, "go.shape."):
			shapeBodies++
		default:
			concreteWrappers++
		}
	}
	return shapeBodies, concreteWrappers, dicts
}

// TestPredictions compares your prediction.yaml against what the compiler
// emitted. It never prints a count or a symbol name.
func TestPredictions(t *testing.T) {
	names := symbols(t)

	mBodies, mWrappers, mDicts := counts(names, "MapM")
	fBodies, fWrappers, fDicts := counts(names, "MapF")

	predict.Check(t, map[string]any{
		"method_shape_bodies":      mBodies,
		"method_concrete_wrappers": mWrappers,
		"method_dicts":             mDicts,
		"func_shape_bodies":        fBodies,
		"func_concrete_wrappers":   fWrappers,
		"func_dicts":               fDicts,
	})
}
