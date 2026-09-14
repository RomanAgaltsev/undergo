package frames

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

var localsRE = regexp.MustCompile(`\.(\w+) STEXT.*locals=0x([0-9a-f]+)`)

// frames compiles this package and returns each function's local frame size.
//
// Symbols carry the full import path, so the name is taken from the last
// dot-separated component of the header line.
func frames(t *testing.T) map[string]int {
	t.Helper()

	out, err := exec.Command("go", "build", "-gcflags=-S", ".").CombinedOutput()
	if err != nil {
		t.Fatalf("go build -gcflags=-S: %v\n%s", err, out)
	}

	sizes := map[string]int{}
	for line := range strings.SplitSeq(string(out), "\n") {
		m := localsRE.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		n, err := strconv.ParseInt(m[2], 16, 64)
		if err != nil {
			continue
		}
		sizes[m[1]] = int(n)
	}

	for _, name := range []string{"Helper", "Unused", "One", "TwoLive", "TwoDisjoint"} {
		if _, ok := sizes[name]; !ok {
			t.Fatalf("%s did not appear in the assembly listing", name)
		}
	}
	return sizes
}

// TestPredictions compares frames against each other and against your answers.
// It never prints a size.
func TestPredictions(t *testing.T) {
	f := frames(t)

	// "About double" rather than exactly: the frame carries a little
	// bookkeeping beyond the arrays, and that overhead differs by
	// architecture.
	ratio := float64(f["TwoLive"]) / float64(f["One"])

	predict.Check(t, map[string]any{
		"helper_frame_is_zero":         f["Helper"] == 0,
		"unused_array_costs_frame":     f["Unused"] > 0,
		"one_array_costs_frame":        f["One"] > 0,
		"two_live_is_about_double_one": ratio > 1.9 && ratio < 2.1,
		"disjoint_scopes_reduce_frame": f["TwoDisjoint"] < f["TwoLive"],
	})
}
