package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	categoryDirRE = regexp.MustCompile(`^C[0-9]{1,2}-([a-z0-9-]+)$`)
	taskDirRE     = regexp.MustCompile(`^([0-9]{2})-([a-z0-9-]+)$`)
	defectRE      = regexp.MustCompile(`(?m)^### [0-9]+\.`)
	// The warning wraps differently across the corpus, so the whitespace
	// between the words has to be flexible.
	answerWarningRE = regexp.MustCompile(`\*\*Do not open\s+` + "`" + `ANSWERS\.md` + "`" + `\s+until you submit\.\*\*`)
	answerFileRE    = regexp.MustCompile("`?ANSWERS\\.md`?")
)

// trackForCategory turns a loupe category directory into an undergo track.
func trackForCategory(dirName string) (string, error) {
	m := categoryDirRE.FindStringSubmatch(dirName)
	if m == nil {
		return "", fmt.Errorf("%q is not a loupe category directory", dirName)
	}
	return "review/" + m[1], nil
}

// splitTaskDir splits NN-slug. Template directories start with _ and are rejected.
func splitTaskDir(dirName string) (num, slug string, err error) {
	m := taskDirRE.FindStringSubmatch(dirName)
	if m == nil {
		return "", "", fmt.Errorf("%q is not an NN-slug directory", dirName)
	}
	return m[1], m[2], nil
}

// humanize renders a slug as a title: "request-counter" -> "Request counter".
func humanize(slug string) string {
	words := strings.ReplaceAll(slug, "-", " ")
	if words == "" {
		return ""
	}
	return strings.ToUpper(words[:1]) + words[1:]
}

// fieldValue reads a "- **Label:** value" line, dropping any trailing
// parenthesised aside.
func fieldValue(body, label string) (string, error) {
	prefix := "- **" + label + ":**"
	for line := range strings.SplitSeq(body, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		value := strings.TrimSpace(strings.TrimPrefix(line, prefix))
		if cut := strings.Index(value, "(see"); cut >= 0 {
			value = value[:cut]
		}
		return strings.TrimSpace(value), nil
	}
	return "", fmt.Errorf("no %q line", label)
}

// defectCount counts the numbered defect headings in an answer key.
func defectCount(answers string) int {
	return len(defectRE.FindAllString(answers, -1))
}

// tierProfile maps a loupe tier onto undergo's difficulty and estimate.
func tierProfile(tier string) (int, string, error) {
	switch tier {
	case "obvious":
		return 2, "20m", nil
	case "subtle":
		return 3, "30m", nil
	case "multi-bug":
		return 4, "45m", nil
	default:
		return 0, "", fmt.Errorf("unknown tier %q", tier)
	}
}

// rewriteDrillREADME points the solver at `undergo reveal` instead of the answer
// file, which no longer exists in the imported task, and re-roots the drill's
// relative paths for the depth its task id implies. It fails rather than
// silently leaving a dangling reference.
func rewriteDrillREADME(body, id string) (string, error) {
	if !strings.Contains(body, "ANSWERS.md") {
		return "", fmt.Errorf("README never mentions ANSWERS.md — not a loupe drill")
	}
	out := answerWarningRE.ReplaceAllString(body,
		"**Do not run `undergo reveal "+id+"` until you submit.**")
	out = answerFileRE.ReplaceAllString(out, "`undergo reveal "+id+"`")
	if strings.Contains(out, "ANSWERS.md") {
		return "", fmt.Errorf("rewrite left a reference to ANSWERS.md")
	}
	return rerootRelPaths(out, id), nil
}

// rerootRelPaths fixes the ../-rooted paths in a drill README for the depth of
// its task id. A loupe drill sits three levels below its repo root and writes
// ../../../rubric/…; a task at tasks/<id> sits one level below that per id
// segment, so a grouped id such as review/concurrency/01-x needs four.
func rerootRelPaths(body, id string) string {
	want := strings.Count(id, "/") + 2 // tasks/ + each id segment above the last
	if want == 3 {
		return body
	}
	return strings.ReplaceAll(body, "../../../", strings.Repeat("../", want))
}

// drillHint is rung 1. Knowing how many defects were planted changes how a
// reviewer reads the code without telling them what any of them are.
func drillHint(category, tier string, defects int) string {
	noun := "defects"
	if defects == 1 {
		noun = "defect"
	}
	return fmt.Sprintf(`# Hint

This drill plants **%d %s**.

- Category: %s
- Tier: %s

Knowing the count is the nudge: if you found one and the count says three, keep
reading. The category tells you which section of `+"`rubric/bug-taxonomy.md`"+` to
re-read before you look again.

Nothing here says what the defects are. That is `+"`undergo reveal`"+`.
`, defects, noun, category, tier)
}

// importLoupe converts every drill in from into an undergo task under root.
// It never reads from/log — that is Roman's own review history, and the public
// repo holds reference content only.
func importLoupe(from, root string) (int, error) {
	drills := filepath.Join(from, "drills")
	categories, err := os.ReadDir(drills)
	if err != nil {
		return 0, err
	}

	var n int
	for _, category := range categories {
		if !category.IsDir() {
			continue
		}
		track, err := trackForCategory(category.Name())
		if err != nil {
			continue // _template and _clean-template
		}
		entries, err := os.ReadDir(filepath.Join(drills, category.Name()))
		if err != nil {
			return n, err
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			num, slug, err := splitTaskDir(entry.Name())
			if err != nil {
				continue
			}
			src := filepath.Join(drills, category.Name(), entry.Name())
			if err := importOneDrill(src, root, track, num, slug); err != nil {
				return n, fmt.Errorf("%s/%s: %w", category.Name(), entry.Name(), err)
			}
			n++
		}
	}
	return n, nil
}

func importOneDrill(src, root, track, num, slug string) error {
	id := track + "/" + num + "-" + slug

	readme, err := os.ReadFile(filepath.Join(src, "README.md"))
	if err != nil {
		return err
	}
	answers, err := os.ReadFile(filepath.Join(src, "ANSWERS.md"))
	if err != nil {
		return err
	}
	category, err := fieldValue(string(readme), "Category")
	if err != nil {
		return err
	}
	tier, err := fieldValue(string(readme), "Tier")
	if err != nil {
		return err
	}
	difficulty, estimate, err := tierProfile(tier)
	if err != nil {
		return err
	}
	rewritten, err := rewriteDrillREADME(string(readme), id)
	if err != nil {
		return err
	}

	files := map[string]string{"README.md": rewritten}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".go" {
			continue
		}
		body, err := os.ReadFile(filepath.Join(src, entry.Name()))
		if err != nil {
			return err
		}
		files[entry.Name()] = string(body)
	}

	return writeTask(root, taskFiles{
		ID: id, Title: humanize(slug), Mode: "review", Track: track,
		Difficulty: difficulty, Estimate: estimate,
		Tags:        []string{"review", strings.TrimPrefix(track, "review/"), tier},
		InspiredBy:  "https://github.com/RomanAgaltsev/loupe",
		Files:       files,
		Hint:        drillHint(category, tier, defectCount(string(answers))),
		Explanation: string(answers),
	})
}
