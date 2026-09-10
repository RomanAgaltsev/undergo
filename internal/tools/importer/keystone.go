package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var kataDirRE = regexp.MustCompile(`^k([0-9]{1,2})-([a-z0-9-]+)$`)

// kataID splits a keystone kata directory, zero-padding the number so it fits
// undergo's two-digit id format.
func kataID(dirName string) (num, slug string, err error) {
	m := kataDirRE.FindStringSubmatch(dirName)
	if m == nil {
		return "", "", fmt.Errorf("%q is not a keystone kata directory", dirName)
	}
	n, err := strconv.Atoi(m[1])
	if err != nil {
		return "", "", err
	}
	return fmt.Sprintf("%02d", n), m[2], nil
}

// section returns the body under a "## heading", up to the next heading.
//
// The heading is matched by prefix: 33 of keystone's 36 katas write it as
// "## Design should also address (in DESIGN.md)" and only 3 use the bare form,
// so an exact match would silently find almost nothing.
func section(body, heading string) (string, error) {
	want := "## " + heading
	lines := strings.Split(body, "\n")
	for i, line := range lines {
		if !strings.HasPrefix(strings.TrimSpace(line), want) {
			continue
		}
		var out []string
		for _, next := range lines[i+1:] {
			if strings.HasPrefix(next, "## ") {
				break
			}
			out = append(out, next)
		}
		return strings.TrimSpace(strings.Join(out, "\n")), nil
	}
	return "", fmt.Errorf("no %q section", heading)
}

// deepenRubricPaths fixes the one path change the move causes: a keystone kata
// sits two levels below the repo root, an undergo task sits three.
func deepenRubricPaths(body string) string {
	return strings.ReplaceAll(body, "../../rubric/", "../../../rubric/")
}

// kataExplanation is rung 2 for a design kata. keystone ships no reference
// designs, so this is honest about being a grading checklist rather than
// pretending to be a model answer.
func kataExplanation(title, checklist string) string {
	return fmt.Sprintf(`# How to grade your design: %s

There is no model answer here, and pretending otherwise would be worse than
saying so. A design is graded against a rubric and against the questions the
kata asked, not against one right shape.

## Grade against the rubric

Score your DESIGN.md on every dimension in `+"`rubric/design-rubric.md`"+`. Be
harder on yourself than a reviewer would be: the dimension you skipped is the
one you did not want to think about.

## The kata's own checklist

Your design should have addressed each of these explicitly. Any you did not
name at all is a gap, not a preference:

%s

## Then build the core

The design is graded before the code exists. Once you have scored it, build the
core the kata asks for and write REFLECTION.md on what the building taught you
that the designing did not.
`, title, checklist)
}

// importKeystone converts every kata in from into an undergo design task.
func importKeystone(from, root string) (int, error) {
	katas := filepath.Join(from, "katas")
	entries, err := os.ReadDir(katas)
	if err != nil {
		return 0, err
	}

	var n int
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		num, slug, err := kataID(entry.Name())
		if err != nil {
			continue // _template
		}
		if err := importOneKata(filepath.Join(katas, entry.Name()), root, num, slug); err != nil {
			return n, fmt.Errorf("%s: %w", entry.Name(), err)
		}
		n++
	}
	return n, nil
}

func importOneKata(src, root, num, slug string) error {
	id := "design/" + num + "-" + slug
	title := humanize(slug)

	files := map[string]string{}
	for _, name := range []string{"README.md", "DESIGN.md", "REFLECTION.md"} {
		body, err := os.ReadFile(filepath.Join(src, name))
		if err != nil {
			return err
		}
		files[name] = deepenRubricPaths(string(body))
	}
	checklist, err := section(files["README.md"], "Design should also address")
	if err != nil {
		return err
	}

	return writeTask(root, taskFiles{
		ID: id, Title: title, Mode: "design", Track: "design",
		Difficulty: 4, Estimate: "90m",
		Tags:       []string{"design", "system-design"},
		InspiredBy: "https://github.com/RomanAgaltsev/keystone",
		Files:      files,
		Hint: "# Hint\n\nBefore you look at anything else, re-read " +
			"`rubric/design-rubric.md` and check your DESIGN.md against the " +
			"kata's own \"Design should also address\" list. Most gaps are there.\n",
		Explanation: kataExplanation(title, checklist),
	})
}
