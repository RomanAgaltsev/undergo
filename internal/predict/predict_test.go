package predict

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// The interface exists to make Check testable; this is the assertion that it
// still describes the thing it stands in for.
var _ TB = (*testing.T)(nil)

// recorder stands in for *testing.T so that a grading can be inspected instead
// of being obeyed. Its Fatalf ends the calling goroutine exactly as the real one
// does, because a Check that stops at the first fatal and a Check that carries on
// are different functions, and the one under test has to be the one that ships.
type recorder struct {
	errs  []string
	logs  []string
	fatal string
}

func (r *recorder) Helper() {}

func (r *recorder) Errorf(format string, args ...any) {
	r.errs = append(r.errs, fmt.Sprintf(format, args...))
}

func (r *recorder) Logf(format string, args ...any) {
	r.logs = append(r.logs, fmt.Sprintf(format, args...))
}

func (r *recorder) Fatalf(format string, args ...any) {
	r.fatal = fmt.Sprintf(format, args...)
	runtime.Goexit()
}

// failed reports what a real *testing.T would have reported: any Errorf, or a
// Fatalf, means the task did not pass.
func (r *recorder) failed() bool { return len(r.errs) > 0 || r.fatal != "" }

func (r *recorder) all() string {
	return strings.Join(append(append([]string{}, r.errs...), r.logs...), "\n") + "\n" + r.fatal
}

// grade writes a prediction file and runs Check against it as a SOLVER sees it,
// on its own goroutine so that a Fatalf can end that goroutine without ending
// the test.
func grade(t *testing.T, file string, measured map[string]any) *recorder {
	t.Helper()
	return gradeIn(t, file, measured, false)
}

// gradeAsCI is the same grading as ci-verify sees it, where the reader is a
// maintainer and the failing slots may be named.
func gradeAsCI(t *testing.T, file string, measured map[string]any) *recorder {
	t.Helper()
	return gradeIn(t, file, measured, true)
}

func gradeIn(t *testing.T, file string, measured map[string]any, ci bool) *recorder {
	t.Helper()
	path := filepath.Join(t.TempDir(), DefaultFile)
	if err := os.WriteFile(path, []byte(file), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("UNDERGO_PREDICTION", path)
	if ci {
		t.Setenv(CIEnv, "1")
	} else {
		t.Setenv(CIEnv, "")
	}

	r := &recorder{}
	done := make(chan struct{})
	go func() {
		defer close(done)
		Check(r, measured)
	}()
	<-done
	return r
}

// measured is one slot of each kind Check has to compare: an integer, a boolean
// and a string. Between them they cover every branch of canonical.
func measured() map[string]any {
	return map[string]any{
		"sizeof_header": 24,
		"escapes":       true,
		"caps":          "1 2 4 8",
	}
}

// This is the test the repository did not have. Check grades 77 of the 96
// machine-graded tasks, and gate 2 cannot prove it works: gate 2 feeds it only
// reference answers, which are correct by construction, so a Check that accepted
// everything would leave every gate green and every task ungraded.
func TestCheckRejectsWrongAnswers(t *testing.T) {
	tests := []struct {
		name       string
		file       string
		wantFailed bool
		mentions   string
	}{
		{
			name: "every slot correct",
			file: "sizeof_header: 24\nescapes: true\ncaps: 1 2 4 8\n",
		},
		{
			name:       "one integer wrong",
			file:       "sizeof_header: 25\nescapes: true\ncaps: 1 2 4 8\n",
			wantFailed: true,
			mentions:   `slot "sizeof_header": incorrect`,
		},
		{
			name:       "boolean inverted",
			file:       "sizeof_header: 24\nescapes: false\ncaps: 1 2 4 8\n",
			wantFailed: true,
			mentions:   `slot "escapes": incorrect`,
		},
		{
			name:       "sequence off by one element",
			file:       "sizeof_header: 24\nescapes: true\ncaps: 1 2 4 16\n",
			wantFailed: true,
			mentions:   `slot "caps": incorrect`,
		},
		{
			name:       "a slot left out",
			file:       "sizeof_header: 24\nescapes: true\n",
			wantFailed: true,
			mentions:   `slot "caps": no prediction`,
		},
		{
			name:       "a slot left blank",
			file:       "sizeof_header: 24\nescapes: true\ncaps:\n",
			wantFailed: true,
			mentions:   `slot "caps": no prediction`,
		},
		{
			name:       "a slot this task does not have",
			file:       "sizeof_header: 24\nescapes: true\ncaps: 1 2 4 8\nbogus: 1\n",
			wantFailed: true,
			mentions:   `slot "bogus": not a slot in this task`,
		},
		{
			name:       "every slot wrong",
			file:       "sizeof_header: 0\nescapes: false\ncaps: nothing\n",
			wantFailed: true,
		},
		{
			// canonical exists so that a solver is graded on the value and not
			// on how they wrote it down.
			name: "hex, quotes and padding all accepted",
			file: "sizeof_header:   0x18   # a comment\nescapes: \"true\"\ncaps: '1 2 4 8'\n",
		},
		{
			name:       "a file that is not key: value",
			file:       "this line has no colon\n",
			wantFailed: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// The verdict must be the same for both readers. Only the detail
			// differs, which is what the next two tests are about.
			for _, as := range []struct {
				name  string
				grade func(*testing.T, string, map[string]any) *recorder
			}{{"solver", grade}, {"ci", gradeAsCI}} {
				t.Run(as.name, func(t *testing.T) {
					r := as.grade(t, tc.file, measured())
					if got := r.failed(); got != tc.wantFailed {
						t.Errorf("failed = %v, want %v\n%s", got, tc.wantFailed, r.all())
					}
					if as.name == "ci" && tc.mentions != "" && !strings.Contains(r.all(), tc.mentions) {
						t.Errorf("CI output does not mention %q:\n%s", tc.mentions, r.all())
					}
				})
			}
		})
	}
}

// `undergo verify` runs `go test -v`, so a solver reads every line of this. 34
// of the 77 predict tasks grade booleans only, several of them nine or ten at a
// time — and naming the wrong ones turns that into a perfect oracle: answer all
// true, read back which came out wrong, flip exactly those, done in two runs
// with no understanding at all. That is not what the seal trades away; opening
// a seal is deliberate and is recorded as a peek, while this was free.
func TestCheckDoesNotTellASolverWhichSlotsAreWrong(t *testing.T) {
	r := grade(t, "sizeof_header: 999\nescapes: false\ncaps: wrong\n", measured())
	if !r.failed() {
		t.Fatal("expected a failure to inspect")
	}
	for _, slot := range []string{"sizeof_header", "escapes", "caps"} {
		if strings.Contains(r.all(), `slot "`+slot+`": incorrect`) {
			t.Errorf("a solver was told %q is the wrong one:\n%s", slot, r.all())
		}
	}
	// The count is the signal the solver does need.
	if !strings.Contains(r.all(), "3 of 3 slots incorrect") {
		t.Errorf("a solver should still learn how many are wrong:\n%s", r.all())
	}
}

// ci-verify is the other reader: the answers are already known to be correct,
// the whole output is captured, and naming the failing slots is what turned two
// blind gate-2 failures into a one-run diagnosis in M10.
func TestCheckTellsCIWhichSlotsAreWrong(t *testing.T) {
	r := gradeAsCI(t, "sizeof_header: 999\nescapes: true\ncaps: 1 2 4 8\n", measured())
	if !r.failed() {
		t.Fatal("expected a failure to inspect")
	}
	if !strings.Contains(r.all(), `slot "sizeof_header": incorrect`) {
		t.Errorf("CI was not told which slot failed:\n%s", r.all())
	}
}

// An unanswered slot is named to both readers: which question you left blank is
// not a hint about any answer.
func TestCheckAlwaysNamesAnUnansweredSlot(t *testing.T) {
	r := grade(t, "sizeof_header: 24\nescapes: true\n", measured())
	if !strings.Contains(r.all(), `slot "caps": no prediction`) {
		t.Errorf("a missing answer should name its slot:\n%s", r.all())
	}
}

// Flipping any single slot must fail the task. A grader that catches some slots
// and not others is the same defect as one that catches none, only harder to
// notice — so this asserts the property over every slot rather than sampling it.
func TestCheckFailsOnEverySingleWrongSlot(t *testing.T) {
	correct := map[string]string{
		"sizeof_header": "24",
		"escapes":       "true",
		"caps":          "1 2 4 8",
	}
	wrong := map[string]string{
		"sizeof_header": "999",
		"escapes":       "false",
		"caps":          "nothing like it",
	}

	for slot := range correct {
		t.Run(slot, func(t *testing.T) {
			var b strings.Builder
			for name, value := range correct {
				if name == slot {
					value = wrong[name]
				}
				fmt.Fprintf(&b, "%s: %s\n", name, value)
			}
			r := grade(t, b.String(), measured())
			if !r.failed() {
				t.Errorf("a wrong %q was graded as correct:\n%s", slot, r.all())
			}
		})
	}
}

// Check must never print a measured value: a solver who got a slot wrong would
// otherwise be handed the answer on their first attempt, which is the whole
// reason a predict task stores no answer.
func TestCheckNeverPrintsTheMeasuredValue(t *testing.T) {
	m := map[string]any{"sizeof_header": 24, "escapes": true, "caps": "1 2 4 8"}
	r := grade(t, "sizeof_header: 999\nescapes: false\ncaps: wrong\n", m)
	if !r.failed() {
		t.Fatal("expected a failure to inspect")
	}
	for _, secret := range []string{"24", "true", "1 2 4 8"} {
		if strings.Contains(r.all(), secret) {
			t.Errorf("grading leaked the measured value %q:\n%s", secret, r.all())
		}
	}
}

func TestCheckFatalsWhenThePredictionFileIsMissing(t *testing.T) {
	t.Setenv("UNDERGO_PREDICTION", filepath.Join(t.TempDir(), "absent.yaml"))
	r := &recorder{}
	done := make(chan struct{})
	go func() {
		defer close(done)
		Check(r, measured())
	}()
	<-done

	if r.fatal == "" {
		t.Fatal("a missing prediction file must be fatal, not a silent pass")
	}
	if len(r.errs) > 0 {
		t.Errorf("Fatalf must stop the grading, but it continued: %v", r.errs)
	}
}

func TestParseFlatFile(t *testing.T) {
	got, err := parse([]byte("# a comment\nsizeof_header: 24\nname: \"padded\"\n\nflag: true   # trailing\n"))
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"sizeof_header": "24", "name": "padded", "flag": "true"}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("%s = %q, want %q", k, got[k], v)
		}
	}
}

func TestParseRejectsGarbage(t *testing.T) {
	if _, err := parse([]byte("this line has no colon\n")); err == nil {
		t.Fatal("parse accepted a line with no key")
	}
}

func TestCanonicalNormalisesAcrossTypes(t *testing.T) {
	tests := []struct{ a, b any }{
		{"24", uintptr(24)},
		{"24", 24},
		{"0x18", 24},
		{"true", true},
		{"1.5", 1.5},
	}
	for _, tc := range tests {
		if canonical(tc.a) != canonical(tc.b) {
			t.Errorf("canonical(%v)=%q != canonical(%v)=%q",
				tc.a, canonical(tc.a), tc.b, canonical(tc.b))
		}
	}
	if canonical("25") == canonical(24) {
		t.Error("canonical collapsed two different values")
	}
}
