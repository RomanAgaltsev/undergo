// Package childenv is the environment undergo gives every go command it starts.
//
// A task's answer must come from the task, not from the shell it was verified
// in. Two variables would move it: GOFLAGS, which can change how every package
// is compiled (-gcflags=all=-N turns off inlining and escape analysis), and
// GODEBUG, which changes runtime behaviour without changing a line of source.
// Because a task's frozen tests inherit this environment, so does every go
// command a task starts in turn.
package childenv

import (
	"os"
	"strings"
)

// neutralGOFLAGS is one of GOFLAGS's own defaults. GOFLAGS cannot simply be
// emptied: an empty GOFLAGS is treated as unset, and the go command then reads
// whatever `go env -w` stored.
const neutralGOFLAGS = "GOFLAGS=-buildvcs=auto"

// Environ returns the current environment without GOFLAGS and GODEBUG, with
// GOFLAGS set to a neutral default, followed by extra. A task that needs a
// setting of its own passes it in extra, where it wins.
func Environ(extra ...string) []string {
	env := make([]string, 0, len(os.Environ())+1+len(extra))
	for _, kv := range os.Environ() {
		name, _, _ := strings.Cut(kv, "=")
		if strings.EqualFold(name, "GOFLAGS") || strings.EqualFold(name, "GODEBUG") {
			continue
		}
		env = append(env, kv)
	}
	env = append(env, neutralGOFLAGS)
	return append(env, extra...)
}
