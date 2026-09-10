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
func (f *File) Save(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := yaml.Marshal(f)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
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
