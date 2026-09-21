// Package progress keeps the solver's local record. It lives under .undergo/,
// which is gitignored: the public repo is the etalon and never carries anyone's
// results.
package progress

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"go.yaml.in/yaml/v3"
)

// Entry is one task's history.
type Entry struct {
	Attempts    int    `yaml:"attempts"`
	Solved      bool   `yaml:"solved"`
	FirstSolved string `yaml:"first_solved,omitempty"`
	HintUsed    bool   `yaml:"hint_used"`
	Peeked      bool   `yaml:"peeked"`
	LastRun     string `yaml:"last_run,omitempty"`
}

// File is the whole record.
type File struct {
	Tasks map[string]*Entry `yaml:"tasks"`
}

// Load reads the record. A missing file is an empty record, not an error.
func Load(path string) (*File, error) {
	b, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return &File{Tasks: map[string]*Entry{}}, nil
	}
	if err != nil {
		return nil, err
	}
	var f File
	if err := yaml.Unmarshal(b, &f); err != nil {
		return nil, err
	}
	if f.Tasks == nil {
		f.Tasks = map[string]*Entry{}
	}
	return &f, nil
}

// Save writes the record, creating its directory if needed.
//
// The write goes to a temporary file in the same directory and is renamed over
// the target, because os.WriteFile truncates before it writes. This is the only
// record of a solver's work, it is deliberately gitignored so there is no
// backup, and Save runs on every verify, hint and reveal — an interrupt during
// a long verify is an ordinary event, not a disaster scenario. A same-directory
// rename is atomic on both Windows and POSIX.
func (f *File) Save(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	b, err := yaml.Marshal(f)
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp(dir, ".progress-*.yaml")
	if err != nil {
		return err
	}
	// A no-op once the rename has succeeded; on any earlier return it is what
	// stops a failed save leaving litter beside the record.
	defer func() { _ = os.Remove(tmp.Name()) }()

	if _, err := tmp.Write(b); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}

// Get returns a task's entry, creating it if absent.
func (f *File) Get(id string) *Entry {
	if f.Tasks == nil {
		f.Tasks = map[string]*Entry{}
	}
	e, ok := f.Tasks[id]
	if !ok {
		e = &Entry{}
		f.Tasks[id] = e
	}
	return e
}

// RecordRun logs one verification attempt.
func (f *File) RecordRun(id string, passed bool) {
	e := f.Get(id)
	e.Attempts++
	e.LastRun = now()
	if passed && !e.Solved {
		e.Solved = true
		e.FirstSolved = e.LastRun
	}
}

// RecordHint logs that the first rung was opened.
func (f *File) RecordHint(id string) { f.Get(id).HintUsed = true }

// RecordPeek logs that the solution was revealed.
func (f *File) RecordPeek(id string) { f.Get(id).Peeked = true }

func now() string { return time.Now().UTC().Format(time.RFC3339) }
