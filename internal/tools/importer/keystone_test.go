package main

import (
	"strings"
	"testing"
)

const kataREADME = `# Kata k1 — Rate limiter

## Problem
Design an in-process API rate limiter.

## Design should also address
- Token-bucket vs sliding-window vs fixed-window.
- How you'd extend this to distributed rate limiting.

## Deliverables
- DESIGN.md
`

func TestKataID(t *testing.T) {
	num, slug, err := kataID("k01-rate-limiter")
	if err != nil {
		t.Fatal(err)
	}
	if num != "01" || slug != "rate-limiter" {
		t.Fatalf("got %q, %q", num, slug)
	}
	padded, _, err := kataID("k9-thing")
	if err != nil {
		t.Errorf("a single-digit kata must be zero-padded, not rejected: %v", err)
	}
	if padded != "09" {
		t.Errorf("kataID(k9-thing) = %q, want %q", padded, "09")
	}
	if _, _, err := kataID("_template"); err == nil {
		t.Error("the template directory must be rejected")
	}
}

func TestSectionExtractsTheChecklist(t *testing.T) {
	got, err := section(kataREADME, "Design should also address")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "Token-bucket") || !strings.Contains(got, "distributed rate limiting") {
		t.Errorf("section = %q", got)
	}
	if strings.Contains(got, "Deliverables") {
		t.Error("the section must stop at the next heading")
	}
	if _, err := section(kataREADME, "Nonexistent"); err == nil {
		t.Error("a missing section must be an error, not an empty string")
	}
}

// 33 of keystone's 36 katas suffix the heading with "(in DESIGN.md)". An exact
// match would find the section in only 3 of them.
func TestSectionMatchesTheSuffixedHeadingVariant(t *testing.T) {
	body := strings.Replace(kataREADME,
		"## Design should also address",
		"## Design should also address (in DESIGN.md)", 1)

	got, err := section(body, "Design should also address")
	if err != nil {
		t.Fatalf("section: %v", err)
	}
	if !strings.Contains(got, "Token-bucket") {
		t.Errorf("section = %q", got)
	}
}

func TestDeepenRubricPaths(t *testing.T) {
	in := "> Graded against `../../rubric/design-rubric.md`. Fill each section.\n"
	got := deepenRubricPaths(in)
	if !strings.Contains(got, "../../../rubric/design-rubric.md") {
		t.Errorf("got %q — a kata lands one level deeper in undergo than in keystone", got)
	}
}
