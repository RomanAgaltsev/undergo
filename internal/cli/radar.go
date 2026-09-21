package cli

import (
	"errors"
	"fmt"
	"go/version"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/RomanAgaltsev/undergo/internal/manifest"
)

// radarEntry is one finding in a release's radar file.
type radarEntry struct {
	Heading string
	Tag     string // candidate | invalidates | no action
	Target  string // a track for a candidate, a task id for invalidates
	Source  string
}

// radarTags are the verdicts a radar entry may carry.
//
// "resolved" exists because an invalidation used to have no end state. The one
// the radar has ever recorded was addressed at M5 — layout/01's hint corrected,
// the explanation extended, the task re-sealed — and radar-check still reports
// "1 invalidations" and always will. A counter that only goes up is a counter
// nobody reads: with three of them the summary could not tell "three problems,
// all fixed" from "three problems, none fixed", which is the only question the
// line exists to answer.
var radarTags = []string{"candidate", "invalidates", "resolved", "no action"}

// parseRadarEntries reads a radar file: a "Source:" line, then one "## heading"
// per finding, each followed by a "→ tag | target" line.
func parseRadarEntries(body string) ([]radarEntry, error) {
	var source string
	for line := range strings.SplitSeq(body, "\n") {
		if rest, ok := strings.CutPrefix(strings.TrimSpace(line), "Source:"); ok {
			source = strings.TrimSpace(rest)
			break
		}
	}
	if source == "" {
		return nil, fmt.Errorf("no Source: line — a radar entry without provenance is a rumour")
	}

	var entries []radarEntry
	lines := strings.Split(body, "\n")
	for i, line := range lines {
		heading, ok := strings.CutPrefix(strings.TrimSpace(line), "## ")
		if !ok {
			continue
		}
		heading = strings.TrimSpace(heading)

		var tagLine string
		for _, next := range lines[i+1:] {
			if strings.TrimSpace(next) == "" {
				continue
			}
			tagLine = strings.TrimSpace(next)
			break
		}
		rest, ok := strings.CutPrefix(tagLine, "→")
		if !ok {
			return nil, fmt.Errorf("entry %q: no → tag line", heading)
		}
		tag, target, _ := strings.Cut(strings.TrimSpace(rest), "|")
		tag = strings.TrimSpace(tag)
		if !slices.Contains(radarTags, tag) {
			return nil, fmt.Errorf("entry %q: tag %q, want one of %v", heading, tag, radarTags)
		}
		entries = append(entries, radarEntry{
			Heading: heading,
			Tag:     tag,
			Target:  strings.TrimSpace(target),
			Source:  source,
		})
	}
	return entries, nil
}

// checkTargets cross-references entries against the catalogue. An invalidates
// entry must name a task that exists — a dangling id is how a dataset rots.
func checkTargets(entries []radarEntry, known map[string]bool) error {
	for _, e := range entries {
		switch e.Tag {
		case "invalidates", "resolved":
			if !known[e.Target] {
				return fmt.Errorf("entry %q %s %q, which is not a task in this repo",
					e.Heading, e.Tag, e.Target)
			}
		case "candidate":
			if e.Target == "" {
				return fmt.Errorf("entry %q is a candidate with no track", e.Heading)
			}
		}
	}
	return nil
}

// RadarCheck validates every radar file against the catalogue.
func RadarCheck(e Env, args []string) error {
	if err := noArgs("radar-check", args); err != nil {
		return err
	}
	tasks, err := manifest.Walk(e.TasksDir())
	if err != nil {
		return err
	}
	known := make(map[string]bool, len(tasks))
	for _, t := range tasks {
		known[t.ID] = true
	}

	dir := filepath.Join(e.Root, "radar", "versions")
	files, err := os.ReadDir(dir)
	if errors.Is(err, fs.ErrNotExist) {
		// git does not track empty directories, so a checkout made before the
		// first release was triaged has no radar at all. That is not a failure.
		fmt.Fprintln(e.Out, "radar: no releases triaged yet")
		return nil
	}
	if err != nil {
		return err
	}

	var total, candidates, invalidates, resolved, parsed int
	newest := ""
	for _, f := range files {
		if f.IsDir() || filepath.Ext(f.Name()) != ".md" {
			continue
		}
		body, err := os.ReadFile(filepath.Join(dir, f.Name()))
		if err != nil {
			return err
		}
		entries, err := parseRadarEntries(string(body))
		if err != nil {
			return fmt.Errorf("radar/versions/%s: %w", f.Name(), err)
		}
		if err := checkTargets(entries, known); err != nil {
			return fmt.Errorf("radar/versions/%s: %w", f.Name(), err)
		}
		parsed++
		total += len(entries)
		if v := radarFileVersion(f.Name()); v != "" && (newest == "" || version.Compare(v, newest) > 0) {
			newest = v
		}
		for _, entry := range entries {
			switch entry.Tag {
			case "candidate":
				candidates++
			case "invalidates":
				invalidates++
			case "resolved":
				resolved++
			}
		}
	}
	if err := checkRadarReach(newest, runtime.Version()); err != nil {
		return err
	}
	fmt.Fprintf(e.Out, "radar: %d entries across %d radar files — %d candidates, "+
		"%d open invalidations, %d resolved\n",
		total, parsed, candidates, invalidates, resolved)
	return nil
}

// radarFileVersion reads the release a radar filename names, taking the newest
// when the name is a range: go1.27.md is go1.27, and go1.0-1.9.md is go1.9.
func radarFileVersion(name string) string {
	base := strings.TrimSuffix(name, ".md")
	if _, after, ok := strings.Cut(base, "-"); ok {
		base = "go" + after
	}
	if !version.IsValid(base) {
		return ""
	}
	return base
}

// checkRadarReach fails when a release has shipped that the radar has not
// triaged.
//
// The radar's whole job is noticing that Go moved, and it could not notice that
// itself: radar-check validated what was there and said nothing about what was
// missing, so the day Go 1.28 ships it would report the same cheerful summary
// over a catalogue whose answers had started to rot.
//
// The toolchain building this is the signal. GOTOOLCHAIN resolves the go
// directive, so the compiler here is at least as new as the release this repo
// has adopted, and a bump to that directive is exactly the moment the sweep is
// owed. Renovate already raises that bump as a PR; this is what makes ignoring
// it red rather than quiet.
func checkRadarReach(newest, goVersion string) error {
	if newest == "" {
		return nil
	}
	have := version.Lang(goVersion)
	if version.Compare(newest, have) < 0 {
		return fmt.Errorf("radar: the newest triaged release is %s but this toolchain is %s — "+
			"sweep %s into radar/versions/ before the catalogue's answers start rotting against it",
			newest, have, have)
	}
	return nil
}
