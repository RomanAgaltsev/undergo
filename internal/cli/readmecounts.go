package cli

import (
	"errors"
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
	var c catalogueCounts
	tracks := map[string]bool{}
	categories := map[string]bool{}

	for _, t := range tasks {
		// A retired task is not one that ships, so it does not inflate the
		// number the README promises. Without this, deprecating a task would
		// make the README's count wrong and gate 3 would then enforce the
		// wrong number.
		if t.Deprecated {
			continue
		}
		c.total++
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

var (
	claimTasks    = readmeClaim{"tasks", regexp.MustCompile(`([0-9]+) tasks ship today`), func(c catalogueCounts) int { return c.total }}
	claimGraded   = readmeClaim{"machine-graded tasks", regexp.MustCompile(`([0-9]+) machine-graded`), func(c catalogueCounts) int { return c.gradeable }}
	claimTracks   = readmeClaim{"internals tracks", regexp.MustCompile(`all ([a-z]+|[0-9]+) internals tracks`), func(c catalogueCounts) int { return c.tracks }}
	claimDrills   = readmeClaim{"review drills", regexp.MustCompile(`([0-9]+) review drills`), func(c catalogueCounts) int { return c.review }}
	claimCategory = readmeClaim{"review categories", regexp.MustCompile(`across ([0-9]+) categories`), func(c catalogueCounts) int { return c.reviewCategories }}
	claimKatas    = readmeClaim{"design katas", regexp.MustCompile(`([0-9]+) system-design katas`), func(c catalogueCounts) int { return c.design }}
)

// countedDoc is a document that states the catalogue's size, and the claims it
// is required to keep making.
//
// Two documents, because the failure has now happened in both. README.md
// drifted for two milestones before M15 gave it a check, and ROADMAP.md drifted
// one milestone after that — the same bookkeeping task updated the README alone,
// and ROADMAP.md carried a paragraph admitting nothing checked it.
//
// The claims are *required*, not merely checked where present. A regex binds to
// exact prose, so rewording "267 tasks ship today" into "267 tasks are
// available" would silently stop that claim being checked and the gate would
// stay green — a check that quietly goes vacuous, which is the failure this
// whole rule exists to prevent, one level up. Retiring a claim means deleting it
// from this list, deliberately.
type countedDoc struct {
	name   string
	claims []readmeClaim
}

var countedDocs = []countedDoc{
	{"README.md", []readmeClaim{claimTasks, claimGraded, claimTracks, claimDrills, claimCategory, claimKatas}},
	{"ROADMAP.md", []readmeClaim{claimTasks, claimGraded, claimTracks}},
}

// checkReadmeCounts compares every counted document's claims with the catalogue.
func checkReadmeCounts(root string, c catalogueCounts) error {
	var problems []string
	for _, doc := range countedDocs {
		if err := checkOneDoc(root, doc, c); err != nil {
			problems = append(problems, err.Error())
		}
	}
	if len(problems) > 0 {
		return errors.New(strings.Join(problems, "\n"))
	}
	return nil
}

// checkOneDoc checks one document's required claims.
//
// A document that is absent entirely is not a fault: the test fixtures are
// repositories with no README, and a repository without one is not lying about
// anything. A document that exists and has stopped making a claim is a fault,
// because that is indistinguishable from a check that has gone vacuous.
func checkOneDoc(root string, doc countedDoc, c catalogueCounts) error {
	body, err := os.ReadFile(filepath.Join(root, doc.name))
	if err != nil {
		return nil //nolint:nilerr // an absent document makes no claims
	}

	var wrong, missing []string
	for _, claim := range doc.claims {
		m := claim.re.FindStringSubmatch(string(body))
		if m == nil {
			missing = append(missing, claim.what)
			continue
		}
		said, ok := readNumber(m[1])
		if !ok {
			wrong = append(wrong, fmt.Sprintf("%s: says %q, which is not a number this check knows",
				claim.what, m[1]))
			continue
		}
		if want := claim.truth(c); said != want {
			wrong = append(wrong, fmt.Sprintf("%s: says %q, the catalogue has %d", claim.what, m[1], want))
		}
	}
	if len(missing) > 0 {
		wrong = append(wrong, fmt.Sprintf("no longer states %s — restore the wording, or remove "+
			"the claim from countedDocs deliberately", strings.Join(missing, ", ")))
	}
	if len(wrong) > 0 {
		return fmt.Errorf("%s is out of date:\n  %s", doc.name, strings.Join(wrong, "\n  "))
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
