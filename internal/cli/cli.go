// Package cli implements the undergo commands. Every command is a function
// taking an Env and its arguments, so all of it is testable without a terminal.
package cli

import (
	"fmt"
	"io"
	"os"
)

const usage = `undergo — Go katas beneath the surface

  list      list tasks            [--track T] [--mode M] [--unsolved]
  show      show one task         <id>
  start     copy a task into work/ <id> [--force]
  verify    grade your work        <id>
  hint      open the first rung    <id>
  reveal    open the solution      <id> [--stuck]
  progress  your record
  doctor    what this machine can grade
  new       scaffold a task        --id track/NN-slug --mode M --title T
  seal      seal a task's _solution/ <id>
  validate  validate every manifest
  ci-verify unseal and prove every reference solution
  ci-stubs  build every task, vet the machine-graded ones
  radar-check validate the release radar against the catalogue
`

// Run dispatches a command. It returns the process exit code.
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stdout, usage)
		return 0
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	root, err := FindRoot(wd)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	env := Env{Root: root, Out: stdout, Err: stderr}

	cmd, rest := args[0], args[1:]
	var runErr error
	switch cmd {
	case "list":
		runErr = List(env, rest)
	case "show":
		runErr = Show(env, rest)
	case "start":
		runErr = Start(env, rest)
	case "verify":
		runErr = Verify(env, rest)
	case "hint":
		runErr = Hint(env, rest)
	case "reveal":
		runErr = Reveal(env, rest)
	case "progress":
		runErr = Progress(env, rest)
	case "new":
		runErr = New(env, rest)
	case "seal":
		runErr = SealCmd(env, rest)
	case "validate":
		runErr = Validate(env, rest)
	case "ci-verify":
		runErr = CIVerify(env, rest)
	case "ci-stubs":
		runErr = CIStubs(env, rest)
	case "radar-check":
		runErr = RadarCheck(env, rest)
	case "doctor":
		runErr = Doctor(env, rest)
	case "help", "-h", "--help":
		fmt.Fprint(stdout, usage)
	default:
		fmt.Fprintf(stderr, "unknown command %q\n\n%s", cmd, usage)
		return 2
	}
	if runErr != nil {
		fmt.Fprintln(stderr, runErr)
		return 1
	}
	return 0
}
