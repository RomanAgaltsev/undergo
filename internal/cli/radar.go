package cli

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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

var radarTags = []string{"candidate", "invalidates", "no action"}

// parseRadarEntries reads a radar file: a "Source:" line, then one "## heading"
// per finding, each followed by a "→ tag | target" line.
func parseRadarEntries(body string) ([]radarEntry, error) {
	var source string
	for _, line := range strings.Split(body, "\n") {
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
		case "invalidates":
			if !known[e.Target] {
				return fmt.Errorf("entry %q invalidates %q, which is not a task in this repo",
					e.Heading, e.Target)
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
func RadarCheck(e Env, _ []string) error {
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

	var total, candidates, invalidates int
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
		total += len(entries)
		for _, entry := range entries {
			switch entry.Tag {
			case "candidate":
				candidates++
			case "invalidates":
				invalidates++
			}
		}
	}
	fmt.Fprintf(e.Out, "radar: %d entries across %d releases — %d candidates, %d invalidations\n",
		total, len(files), candidates, invalidates)
	return nil
}
