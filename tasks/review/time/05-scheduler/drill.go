// Package drill — C13/05 scheduler.
package drill

import "time"

// Job is a unit of scheduled work.
type Job struct {
	Deadline time.Time
	Run      func()
}

// IsDue reports whether the job's deadline has arrived as of now.
func IsDue(j Job, now time.Time) bool {
	return now == j.Deadline
}

// PollUntil runs due jobs, waking every interval, until ctx-less stop is signaled.
func PollUntil(jobs []Job, interval time.Duration, stop <-chan struct{}) {
	for {
		select {
		case <-stop:
			return
		case <-time.After(interval):
			start := time.Now()
			for _, j := range jobs {
				if IsDue(j, time.Now()) {
					j.Run()
				}
			}
			_ = time.Now().Sub(start) // record how long the sweep took
		}
	}
}
