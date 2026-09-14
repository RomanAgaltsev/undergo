package budget

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

var (
	canRE    = regexp.MustCompile(`can inline (\w+) with cost (\d+)`)
	cannotRE = regexp.MustCompile(`cannot inline (\w+): (.*)`)
	budgetRE = regexp.MustCompile(`cost (\d+) exceeds budget (\d+)`)
)

// decisions compiles this package with the inliner's reasoning on and returns
// the cost of every function it would inline, plus the reason for each it would
// not.
func decisions(t *testing.T) (costs map[string]int, refused map[string]string) {
	t.Helper()

	out, err := exec.Command("go", "build", "-gcflags=-m=2", ".").CombinedOutput()
	if err != nil {
		t.Fatalf("go build -gcflags=-m=2: %v\n%s", err, out)
	}

	costs = map[string]int{}
	refused = map[string]string{}
	for line := range strings.SplitSeq(string(out), "\n") {
		if m := canRE.FindStringSubmatch(line); m != nil {
			n, err := strconv.Atoi(m[2])
			if err != nil {
				t.Fatalf("parsing cost: %v", err)
			}
			costs[m[1]] = n
			continue
		}
		if m := cannotRE.FindStringSubmatch(line); m != nil {
			refused[m[1]] = m[2]
		}
	}
	if len(costs) == 0 {
		t.Fatal("the compiler reported no inlining decisions")
	}
	return costs, refused
}

// budgetAndCost pulls both numbers out of a "too complex" refusal.
func budgetAndCost(t *testing.T, reason string) (cost, budget int) {
	t.Helper()
	m := budgetRE.FindStringSubmatch(reason)
	if m == nil {
		t.Fatalf("no cost/budget in refusal: %q", reason)
	}
	c, _ := strconv.Atoi(m[1])
	b, _ := strconv.Atoi(m[2])
	return c, b
}

// TestPredictions reads the compiler's own cost model and compares it with
// yours. It never prints a cost.
func TestPredictions(t *testing.T) {
	costs, refused := decisions(t)

	sprawlingCost, budget := budgetAndCost(t, refused["Sprawling"])
	_, deferRefused := refused["WithDefer"]
	_, goRefused := refused["WithGo"]

	predict.Check(t, map[string]any{
		"inline_budget":        budget,
		"trivial_cost":         costs["Trivial"],
		"loop_cost":            costs["Loop"],
		"calls_inlinable_cost": costs["CallsInlinable"],
		"sprawling_cost":       sprawlingCost,
		"with_defer_inlines":   !deferRefused,
		"with_go_inlines":      !goRefused,
	})
}
