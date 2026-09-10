package main

import (
	"strings"
	"testing"
)

const drillREADME = "# Drill: 01-request-counter\n" + `
- **Category:** C1 — Concurrency & races (see ` + "`../../../rubric/bug-taxonomy.md`" + `)
- **Tier:** obvious

## PR description

> Adds a tiny ` + "`Server`" + ` that tracks how many times each URL path is hit.

## Files to review

- ` + "`drill.go`" + `

## Your task

Review against ` + "`../../../rubric/review-rubric.md`" + `; write findings in the
` + "`../../../rubric/submission-format.md`" + ` format. **Do not open
` + "`ANSWERS.md`" + ` until you submit.**
`

const drillANSWERS = `# Answer key: 01-request-counter

> SPOILERS

## Planted defects

### 1. Concurrent map write with no synchronization
- **Severity:** blocker

### 2. Unsynchronised read in Counts
- **Severity:** major
`

func TestTrackForCategory(t *testing.T) {
	tests := map[string]string{
		"C1-concurrency": "review-concurrency",
		"C15-typed-nil":  "review-typed-nil",
		"C10-testing":    "review-testing",
		"C6-api-design":  "review-api-design",
	}
	for in, want := range tests {
		got, err := trackForCategory(in)
		if err != nil {
			t.Fatalf("%s: %v", in, err)
		}
		if got != want {
			t.Errorf("trackForCategory(%q) = %q, want %q", in, got, want)
		}
	}
	if _, err := trackForCategory("drills"); err == nil {
		t.Error("a non-category directory must be rejected")
	}
}

func TestSplitTaskDirAndHumanize(t *testing.T) {
	num, slug, err := splitTaskDir("01-request-counter")
	if err != nil {
		t.Fatal(err)
	}
	if num != "01" || slug != "request-counter" {
		t.Fatalf("got %q, %q", num, slug)
	}
	if got := humanize(slug); got != "Request counter" {
		t.Errorf("humanize = %q", got)
	}
	if _, _, err := splitTaskDir("_template"); err == nil {
		t.Error("a template directory must be rejected")
	}
}

func TestFieldValueReadsTierAndCategory(t *testing.T) {
	tier, err := fieldValue(drillREADME, "Tier")
	if err != nil {
		t.Fatal(err)
	}
	if tier != "obvious" {
		t.Errorf("tier = %q", tier)
	}
	cat, err := fieldValue(drillREADME, "Category")
	if err != nil {
		t.Fatal(err)
	}
	if cat != "C1 — Concurrency & races" {
		t.Errorf("category = %q — the trailing (see ...) must be stripped", cat)
	}
}

func TestDefectCount(t *testing.T) {
	if got := defectCount(drillANSWERS); got != 2 {
		t.Errorf("defectCount = %d, want 2", got)
	}
}

func TestTierProfile(t *testing.T) {
	tests := map[string]struct {
		difficulty int
		estimate   string
	}{
		"obvious":   {2, "20m"},
		"subtle":    {3, "30m"},
		"multi-bug": {4, "45m"},
	}
	for tier, want := range tests {
		d, e, err := tierProfile(tier)
		if err != nil {
			t.Fatalf("%s: %v", tier, err)
		}
		if d != want.difficulty || e != want.estimate {
			t.Errorf("%s = %d/%s, want %d/%s", tier, d, e, want.difficulty, want.estimate)
		}
	}
	if _, _, err := tierProfile("impossible"); err == nil {
		t.Error("an unknown tier must be rejected, not silently defaulted")
	}
}

func TestRewriteDrillREADMELeavesNoMentionOfTheAnswerFile(t *testing.T) {
	got, err := rewriteDrillREADME(drillREADME, "review-concurrency/01-request-counter")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(got, "ANSWERS.md") {
		t.Errorf("the rewritten README still names the answer file:\n%s", got)
	}
	if !strings.Contains(got, "undergo reveal review-concurrency/01-request-counter") {
		t.Errorf("the rewritten README does not tell the solver how to open the key:\n%s", got)
	}
	if !strings.Contains(got, "../../../rubric/review-rubric.md") {
		t.Error("the rubric path must survive untouched — the new depth is the same")
	}
}

func TestRewriteDrillREADMERejectsAnUnknownPhrasing(t *testing.T) {
	if _, err := rewriteDrillREADME("# Drill: x\n\nNothing to rewrite.\n", "a/01-b"); err == nil {
		t.Error("a README that never mentions the answer file is not a loupe drill")
	}
}

func TestDrillHintNamesTheShapeWithoutTheAnswer(t *testing.T) {
	got := drillHint("C1 — Concurrency & races", "subtle", 3)
	for _, want := range []string{"3 defects", "C1 — Concurrency & races", "subtle"} {
		if !strings.Contains(got, want) {
			t.Errorf("hint missing %q:\n%s", want, got)
		}
	}
}
