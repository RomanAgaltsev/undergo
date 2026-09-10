package cli

import (
	"fmt"
	"text/tabwriter"

	"github.com/RomanAgaltsev/undergo/internal/manifest"
	"github.com/RomanAgaltsev/undergo/internal/progress"
)

// Progress prints the solver's record, per track.
func Progress(e Env, _ []string) error {
	tasks, err := manifest.Walk(e.TasksDir())
	if err != nil {
		return err
	}
	rec, err := progress.Load(e.ProgressPath())
	if err != nil {
		return err
	}

	type row struct{ total, solved, hinted, peeked int }
	byTrack := map[string]*row{}
	var order []string
	for _, t := range tasks {
		r, ok := byTrack[t.Track]
		if !ok {
			r = &row{}
			byTrack[t.Track] = r
			order = append(order, t.Track)
		}
		r.total++
		entry := rec.Get(t.ID)
		if entry.Solved {
			r.solved++
		}
		if entry.HintUsed {
			r.hinted++
		}
		if entry.Peeked {
			r.peeked++
		}
	}

	w := tabwriter.NewWriter(e.Out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TRACK\tSOLVED\tHINTED\tPEEKED")
	for _, track := range order {
		r := byTrack[track]
		fmt.Fprintf(w, "%s\t%d/%d\t%d\t%d\n", track, r.solved, r.total, r.hinted, r.peeked)
	}
	return w.Flush()
}
