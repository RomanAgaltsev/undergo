package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const radarFile = `# Go 1.27

Source: https://go.dev/doc/go1.27

## Size-specialized allocation routines
→ candidate | alloc
Small (<80 byte) allocations are up to 30% cheaper. A predict task can ask for
the per-op cost before and after.

## Timer channels are always unbuffered
→ invalidates | sched/04-timer-wakeup
The asynctimerchan GODEBUG is gone, so any answer that depended on a buffered
timer channel is now wrong.

## crypto/mldsa
→ no action
Post-quantum signatures are not what this repo is about.
`

func TestParseRadarEntriesReadsTagAndTarget(t *testing.T) {
	got, err := parseRadarEntries(radarFile)
	if err != nil {
		t.Fatalf("parseRadarEntries: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("found %d entries, want 3", len(got))
	}
	if got[0].Tag != "candidate" || got[0].Target != "alloc" {
		t.Errorf("entry 0 = %+v", got[0])
	}
	if got[1].Tag != "invalidates" || got[1].Target != "sched/04-timer-wakeup" {
		t.Errorf("entry 1 = %+v", got[1])
	}
	if got[2].Tag != "no action" || got[2].Target != "" {
		t.Errorf("entry 2 = %+v", got[2])
	}
	if got[0].Source != "https://go.dev/doc/go1.27" {
		t.Errorf("source = %q — every entry carries the file's source", got[0].Source)
	}
}

func TestParseRadarEntriesRejectsAnUntaggedHeading(t *testing.T) {
	body := "# Go 1.27\n\nSource: https://go.dev/doc/go1.27\n\n## Something happened\nNo tag line.\n"
	_, err := parseRadarEntries(body)
	if err == nil {
		t.Fatal("an entry with no tag line must be rejected")
	}
	if !strings.Contains(err.Error(), "Something happened") {
		t.Errorf("the error must name the offending entry, got: %v", err)
	}
}

func TestParseRadarEntriesRejectsAnUnknownTag(t *testing.T) {
	body := "# Go 1.27\n\nSource: https://go.dev/doc/go1.27\n\n## Thing\n→ maybe | alloc\nx\n"
	if _, err := parseRadarEntries(body); err == nil {
		t.Fatal("only candidate, invalidates and no action are tags")
	}
}

func TestParseRadarEntriesRequiresASource(t *testing.T) {
	body := "# Go 1.27\n\n## Thing\n→ no action\nx\n"
	if _, err := parseRadarEntries(body); err == nil {
		t.Fatal("a radar file with no Source: line must be rejected — provenance is the point")
	}
}

func TestInvalidatesMustNameARealTask(t *testing.T) {
	known := map[string]bool{"alloc/01-zero-alloc-join": true}

	err := checkTargets([]radarEntry{
		{Heading: "A", Tag: "invalidates", Target: "alloc/01-zero-alloc-join"},
	}, known)
	if err != nil {
		t.Fatalf("a real task id must pass: %v", err)
	}

	err = checkTargets([]radarEntry{
		{Heading: "B", Tag: "invalidates", Target: "sched/04-timer-wakeup"},
	}, known)
	if err == nil {
		t.Fatal("an invalidates line naming a task that does not exist must fail")
	}
	if !strings.Contains(err.Error(), "sched/04-timer-wakeup") {
		t.Errorf("the error must name the dangling id, got: %v", err)
	}
}

func TestCandidateTargetIsATrackNotATask(t *testing.T) {
	// A candidate names the track it would land in; that track need not have
	// tasks yet, so it is not cross-checked against the catalogue.
	if err := checkTargets([]radarEntry{
		{Heading: "C", Tag: "candidate", Target: "gc"},
	}, map[string]bool{}); err != nil {
		t.Errorf("a candidate naming an empty track must pass: %v", err)
	}
	if err := checkTargets([]radarEntry{
		{Heading: "D", Tag: "candidate", Target: ""},
	}, map[string]bool{}); err == nil {
		t.Error("a candidate with no track is not actionable and must fail")
	}
}

// radarRepo extends repo with a radar/versions directory holding the given
// files, keyed by name.
func radarRepo(t *testing.T, files map[string]string) Env {
	t.Helper()
	e := repo(t)
	dir := filepath.Join(e.Root, "radar", "versions")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return e
}

// The summary counts files, not releases: era E4 is one file covering ten of
// them, so calling the number "releases" would misreport the dataset.
func TestRadarCheckSummaryCountsFiles(t *testing.T) {
	e := radarRepo(t, map[string]string{
		"go1.27.md":    "# Go 1.27\n\nSource: https://go.dev/doc/go1.27\n\n## A finding\n→ candidate | alloc\nWhy.\n",
		"go1.0-1.9.md": "# Go 1.0 – 1.9\n\nSource: https://go.dev/doc/devel/release\n\n## Another\n→ no action\nWhy not.\n",
	})
	var out bytes.Buffer
	e.Out = &out

	if err := RadarCheck(e, nil); err != nil {
		t.Fatalf("RadarCheck: %v", err)
	}
	got := out.String()
	if !strings.Contains(got, "2 radar files") {
		t.Errorf("the summary should count files:\n%s", got)
	}
	if strings.Contains(got, "releases") {
		t.Errorf("the summary should not call files releases:\n%s", got)
	}
}

// A directory or a non-Markdown file in radar/versions is skipped by the loop,
// so it must not be counted either.
func TestRadarCheckCountsOnlyParsedFiles(t *testing.T) {
	e := radarRepo(t, map[string]string{
		"go1.27.md": "# Go 1.27\n\nSource: https://go.dev/doc/go1.27\n\n## A finding\n→ candidate | alloc\nWhy.\n",
		"notes.txt": "not markdown",
	})
	if err := os.MkdirAll(filepath.Join(e.Root, "radar", "versions", "drafts"), 0o755); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	e.Out = &out

	if err := RadarCheck(e, nil); err != nil {
		t.Fatalf("RadarCheck: %v", err)
	}
	if got := out.String(); !strings.Contains(got, "1 radar file") {
		t.Errorf("only the one parsed file should be counted:\n%s", got)
	}
}
