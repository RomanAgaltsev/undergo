package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// Run had no test at all, including for the exit codes it is the only place to
// decide. Everything but an unknown command exited 1, so "this solution does not
// pass" and "you typed it wrong" were indistinguishable to a caller.
func TestRunExitCodes(t *testing.T) {
	// Run resolves the repository from the working directory, so the tests have
	// to stand in one. The checkout itself is the closest thing to hand.
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Chdir(wd)

	tests := []struct {
		name string
		args []string
		want int
	}{
		{"no arguments prints usage", nil, 0},
		{"explicit help", []string{"help"}, 0},
		{"unknown command", []string{"nosuchcommand"}, 2},
		{"a command invoked wrongly", []string{"verify"}, 2},
		{"too many arguments", []string{"show", "a", "b"}, 2},
		{"a flag on a command that takes none", []string{"validate", "--race"}, 2},
		{"an unknown task is a runtime failure, not a usage error", []string{"show", "no/99-such"}, 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var out, errb bytes.Buffer
			if got := Run(tc.args, &out, &errb); got != tc.want {
				t.Errorf("exit = %d, want %d\nstdout: %s\nstderr: %s",
					got, tc.want, out.String(), errb.String())
			}
		})
	}
}

// Four commands took their arguments and ignored them, so `undergo ci-stubs
// --race` exited 0 having raced nothing.
func TestCommandsThatTakeNoArgumentsRefuseThem(t *testing.T) {
	for _, name := range []string{"validate", "ci-stubs", "radar-check", "progress"} {
		t.Run(name, func(t *testing.T) {
			err := noArgs(name, []string{"--race"})
			if err == nil {
				t.Fatalf("undergo %s accepted an argument it does not have", name)
			}
			if !strings.Contains(err.Error(), name) {
				t.Errorf("the error should name the command; got: %v", err)
			}
		})
	}
	if err := noArgs("validate", nil); err != nil {
		t.Errorf("no arguments must be fine: %v", err)
	}
}
