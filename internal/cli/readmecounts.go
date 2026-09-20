package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/RomanAgaltsev/undergo/internal/manifest"
)

// catalogueCounts is what the catalogue actually contains. The top-level
// README states these numbers in prose, where nothing used to check them: they
// were two milestones stale before anybody noticed, because a sentence cannot
// fail a build.
type catalogueCounts struct {
	total            int
	gradeable        int // neither review nor design — the machine grades these
	review           int
	design           int
	tracks           int // internals tracks, so review/ and design/ do not count
	reviewCategories int
}

func countCatalogue(tasks []*manifest.Task) catalogueCounts {
	c := catalogueCounts{total: len(tasks)}
	tracks := map[string]bool{}
	categories := map[string]bool{}

	for _, t := range tasks {
		top, _, _ := strings.Cut(t.Track, "/")
		switch t.Mode {
		case manifest.ModeReview:
			c.review++
			categories[t.Track] = true
		case manifest.ModeDesign:
			c.design++
		default:
			c.gradeable++
		}
		if top != "review" && top != "design" {
			tracks[top] = true
		}
	}
	c.tracks = len(tracks)
	c.reviewCategories = len(categories)
	return c
}

// numberWords covers the range a track count can plausibly reach. A README
// reads better saying "all sixteen internals tracks" than "all 16", so the
// claim is allowed to be a word — but then it has to be a word we can check.
var numberWords = map[string]int{
	"one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6,
	"seven": 7, "eight": 8, "nine": 9, "ten": 10, "eleven": 11, "twelve": 12,
	"thirteen": 13, "fourteen": 14, "fifteen": 15, "sixteen": 16,
	"seventeen": 17, "eighteen": 18, "nineteen": 19, "twenty": 20,
}

// readmeClaim is one number the README asserts about the catalogue.
type readmeClaim struct {
	what  string
	re    *regexp.Regexp
	truth func(catalogueCounts) int
}

var readmeClaims = []readmeClaim{
	{"tasks", regexp.MustCompile(`([0-9]+) tasks ship today`), func(c catalogueCounts) int { return c.total }},
	{"machine-graded tasks", regexp.MustCompile(`([0-9]+) machine-graded`), func(c catalogueCounts) int { return c.gradeable }},
	{"internals tracks", regexp.MustCompile(`all ([a-z]+|[0-9]+) internals tracks`), func(c catalogueCounts) int { return c.tracks }},
	{"review drills", regexp.MustCompile(`([0-9]+) review drills`), func(c catalogueCounts) int { return c.review }},
	{"review categories", regexp.MustCompile(`across ([0-9]+) categories`), func(c catalogueCounts) int { return c.reviewCategories }},
	{"design katas", regexp.MustCompile(`([0-9]+) system-design katas`), func(c catalogueCounts) int { return c.design }},
}

// checkReadmeCounts compares the top-level README's claims with the catalogue.
//
// A claim the README does not make is not checked: a README that says no
// numbers is not lying, and a repository without one at all is a fixture rather
// than a fault. The rule catches a number that has stopped being true, which is
// the only failure this has ever actually had.
func checkReadmeCounts(root string, c catalogueCounts) error {
	body, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil {
		return nil //nolint:nilerr // no README is not a claim
	}

	var wrong []string
	for _, claim := range readmeClaims {
		m := claim.re.FindStringSubmatch(string(body))
		if m == nil {
			continue
		}
		said, ok := readNumber(m[1])
		if !ok {
			wrong = append(wrong, fmt.Sprintf("%s: README says %q, which is not a number this check knows",
				claim.what, m[1]))
			continue
		}
		if want := claim.truth(c); said != want {
			wrong = append(wrong, fmt.Sprintf("%s: README says %q, the catalogue has %d", claim.what, m[1], want))
		}
	}
	if len(wrong) > 0 {
		return fmt.Errorf("README.md is out of date:\n  %s", strings.Join(wrong, "\n  "))
	}
	return nil
}

func readNumber(s string) (int, bool) {
	if n, err := strconv.Atoi(s); err == nil {
		return n, true
	}
	n, ok := numberWords[s]
	return n, ok
}
