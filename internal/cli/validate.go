package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/RomanAgaltsev/undergo/internal/manifest"
)

// Validate is CI gate 3: every manifest parses and validates, IDs are unique,
// and no plaintext solution has been committed.
func Validate(e Env, _ []string) error {
	tasks, err := manifest.Walk(e.TasksDir())
	if err != nil {
		return err
	}
	seen := map[string]string{}
	for _, t := range tasks {
		if prev, dup := seen[t.ID]; dup {
			return fmt.Errorf("duplicate id %q in %s and %s", t.ID, prev, t.Dir)
		}
		seen[t.ID] = t.Dir

		if _, err := os.Stat(filepath.Join(t.Dir, AuthoringDir)); err == nil {
			return fmt.Errorf("%s: %s/ is committed — run: undergo seal %s", t.ID, AuthoringDir, t.ID)
		}
		if _, err := os.Stat(filepath.Join(t.Dir, "README.md")); err != nil {
			return fmt.Errorf("%s: no README.md", t.ID)
		}

		if t.Mode == manifest.ModePredict {
			body, err := os.ReadFile(filepath.Join(t.Dir, "README.md"))
			if err != nil {
				return err
			}
			if cmd := instrumentCommand(string(body)); cmd != "" {
				return fmt.Errorf("%s: README.md prints an instrument command (%q) — it belongs in HINT.md, see CONTRIBUTING.md", t.ID, cmd)
			}
		}

		entries, err := os.ReadDir(t.Dir)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
				continue
			}
			body, err := os.ReadFile(filepath.Join(t.Dir, entry.Name()))
			if err != nil {
				return err
			}
			if broken := brokenRelPaths(t.Dir, string(body)); len(broken) > 0 {
				return fmt.Errorf("%s: %s points at %v, which does not exist",
					t.ID, entry.Name(), broken)
			}
		}
	}
	if err := checkReadmeCounts(e.Root, countCatalogue(tasks)); err != nil {
		return err
	}

	fmt.Fprintf(e.Out, "%d manifests valid\n", len(tasks))
	return nil
}

// QuestionsHeading opens the written-questions section, which every predict
// task carries. Everything above it is the task body.
const QuestionsHeading = "## Questions to answer in writing"

// goSubcommands are the go tool's verbs. A "go " line is only a command when
// what follows is one of these: `go func() {…}` launches a goroutine and
// `go 1.21` is a go.mod directive, and a task about the go line prints both.
var goSubcommands = map[string]bool{
	"build": true, "clean": true, "doc": true, "env": true, "fix": true,
	"fmt": true, "generate": true, "get": true, "install": true, "list": true,
	"mod": true, "run": true, "test": true, "tool": true, "version": true,
	"vet": true, "work": true,
}

// commandFences are the fence info strings that hold commands. A block tagged
// with a language — ```go, ```yaml — holds source, not something to run.
var commandFences = map[string]bool{
	"": true, "sh": true, "bash": true, "shell": true, "console": true, "text": true,
}

// instrumentCommand returns the first line of a fenced block in the task body
// that looks like an instrument being run, or "" if there is none.
//
// The policy: a predict task's README names the instrument that judges the
// answer but does not print the command that runs it — that belongs in HINT.md,
// which is never gated. A command on screen invites measuring before
// predicting, which turns a prediction into a transcription.
//
// Only the body is scanned. The written questions are post-prediction work, and
// several of them deliberately ask the solver to re-run the subject under a
// control; a README with no questions section is scanned end to end.
func instrumentCommand(readme string) string {
	inFence, isCommandFence := false, false
	for line := range strings.SplitSeq(readme, "\n") {
		trimmed := strings.TrimSpace(strings.TrimRight(line, "\r"))
		if trimmed == QuestionsHeading {
			return ""
		}
		if strings.HasPrefix(trimmed, "```") {
			if inFence {
				inFence = false
				continue
			}
			inFence = true
			isCommandFence = commandFences[strings.ToLower(strings.TrimPrefix(trimmed, "```"))]
			continue
		}
		if !inFence || !isCommandFence {
			continue
		}
		if strings.HasPrefix(trimmed, "GODEBUG=") {
			return trimmed
		}
		if rest, ok := strings.CutPrefix(trimmed, "go "); ok {
			verb, _, _ := strings.Cut(strings.TrimSpace(rest), " ")
			if goSubcommands[verb] {
				return trimmed
			}
		}
	}
	return ""
}
