package cli

import (
	"flag"
	"fmt"
	"slices"
	"strings"
	"text/tabwriter"

	"github.com/RomanAgaltsev/undergo/internal/manifest"
	"github.com/RomanAgaltsev/undergo/internal/progress"
)

// List prints the catalogue, optionally filtered.
func List(e Env, args []string) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	fs.SetOutput(e.Err)
	track := fs.String("track", "", "only this track")
	mode := fs.String("mode", "", "only this mode")
	tag := fs.String("tag", "", "only tasks carrying this tag")
	unsolved := fs.Bool("unsolved", false, "hide solved tasks")
	deprecated := fs.Bool("deprecated", false, "include retired tasks")
	if err := fs.Parse(args); err != nil {
		return err
	}

	tasks, err := manifest.Walk(e.TasksDir())
	if err != nil {
		return err
	}
	rec, err := progress.Load(e.ProgressPath())
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(e.Out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "STATUS\tID\tMODE\tDIFF\tTITLE")
	for _, t := range tasks {
		if !trackMatches(t.Track, *track) {
			continue
		}
		if *mode != "" && string(t.Mode) != *mode {
			continue
		}
		if *tag != "" && !slices.Contains(t.Tags, *tag) {
			continue
		}
		// A retired task is kept only so that an existing progress file still
		// resolves; it is not something to pick up today.
		if t.Deprecated && !*deprecated {
			continue
		}
		entry := rec.Get(t.ID)
		if *unsolved && entry.Solved {
			continue
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n", statusOf(entry), t.ID, t.Mode, t.Difficulty, t.Title)
	}
	return w.Flush()
}

func statusOf(e *progress.Entry) string {
	switch {
	case e.Solved && e.Peeked:
		return "peeked"
	case e.Solved && e.HintUsed:
		return "hinted"
	case e.Solved:
		return "solved"
	case e.Attempts > 0:
		return "started"
	default:
		return "-"
	}
}

// trackMatches reports whether a task's track satisfies a --track filter.
// Tracks may be grouped, so the filter matches the whole track or a leading
// group of it: --track review selects every review/<category>, and
// --track review/concurrency selects one of them.
func trackMatches(track, filter string) bool {
	switch {
	case filter == "":
		return true
	case track == filter:
		return true
	default:
		return strings.HasPrefix(track, filter+"/")
	}
}
