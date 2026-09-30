package childenv_test

import (
	"io/fs"
	"os"
	"regexp"
	"strings"
	"testing"
)

// goLaunch matches a line that starts the go command.
var goLaunch = regexp.MustCompile(`exec\.Command(Context)?\([^"]*"go"`)

// lookahead is how many lines after a launch may set its environment, room
// for a comment explaining why.
const lookahead = 10

// Every go command the harness starts must get its environment from Environ;
// one that inherits os.Environ would let a caller's GOFLAGS or GODEBUG back
// in. This walks the machinery's own source, not the tasks: tasks inherit a
// clean environment from RunTests.
func TestEveryGoLaunchUsesEnviron(t *testing.T) {
	root, err := os.OpenRoot("..")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = root.Close() }()
	src := root.FS()
	var missing []string
	err = fs.WalkDir(src, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		data, err := fs.ReadFile(src, path)
		if err != nil {
			return err
		}
		lines := strings.Split(string(data), "\n")
		for i, line := range lines {
			if !goLaunch.MatchString(line) {
				continue
			}
			window := strings.Join(lines[i:min(len(lines), i+lookahead)], "\n")
			if !strings.Contains(window, "childenv.Environ(") {
				missing = append(missing, path+":"+strings.TrimSpace(line))
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range missing {
		t.Errorf("starts go without childenv.Environ: %s", m)
	}
}
