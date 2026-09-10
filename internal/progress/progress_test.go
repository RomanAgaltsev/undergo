package progress

import (
	"path/filepath"
	"testing"
)

func TestLoadOfAMissingFileIsEmptyNotAnError(t *testing.T) {
	f, err := Load(filepath.Join(t.TempDir(), "progress.yaml"))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(f.Tasks) != 0 {
		t.Errorf("Tasks = %v, want empty", f.Tasks)
	}
}

func TestRecordRunAndRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.yaml")
	f, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}

	f.RecordRun("layout/01-struct-padding", false)
	f.RecordHint("layout/01-struct-padding")
	f.RecordRun("layout/01-struct-padding", true)
	if err := f.Save(path); err != nil {
		t.Fatal(err)
	}

	got, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	e := got.Get("layout/01-struct-padding")
	if e.Attempts != 2 {
		t.Errorf("Attempts = %d, want 2", e.Attempts)
	}
	if !e.Solved {
		t.Error("Solved = false, want true")
	}
	if !e.HintUsed {
		t.Error("HintUsed = false, want true")
	}
	if e.Peeked {
		t.Error("Peeked = true, want false")
	}
	if e.FirstSolved == "" {
		t.Error("FirstSolved is empty")
	}
}

func TestFirstSolvedIsNotOverwritten(t *testing.T) {
	f := &File{}
	f.RecordRun("a/01-b", true)
	first := f.Get("a/01-b").FirstSolved
	f.RecordRun("a/01-b", true)
	if f.Get("a/01-b").FirstSolved != first {
		t.Error("FirstSolved was overwritten by a later pass")
	}
}
