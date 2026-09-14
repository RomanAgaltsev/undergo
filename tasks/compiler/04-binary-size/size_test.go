package size

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// snippets are the baseline and the five variations on it.
var snippets = []string{"base", "fmtprint", "reflectuse", "generic", "ifacemethod", "bigtable"}

// sizeOf builds testdata/<name>.go.txt and returns the size of the binary.
func sizeOf(t *testing.T, name string) int64 {
	t.Helper()

	src, err := os.ReadFile(filepath.Join("testdata", name+".go.txt"))
	if err != nil {
		t.Fatalf("reading snippet: %v", err)
	}

	dir := t.TempDir()
	write := func(file, body string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, file), []byte(body), 0o644); err != nil {
			t.Fatalf("writing %s: %v", file, err)
		}
	}
	write("go.mod", "module snippet\n\ngo 1.27\n")
	write("main.go", string(src))

	bin := filepath.Join(dir, "snippet")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("building %s: %v\n%s", name, err, out)
	}

	info, err := os.Stat(bin)
	if err != nil {
		t.Fatalf("stat %s: %v", name, err)
	}
	return info.Size()
}

// TestPredictions builds all six and compares each against the baseline.
//
// Every slot is a comparison rather than a byte count: an absolute size depends
// on the platform, the toolchain and the linker, and would be a different number
// on every machine that ran this.
func TestPredictions(t *testing.T) {
	sizes := make(map[string]int64, len(snippets))
	for _, name := range snippets {
		sizes[name] = sizeOf(t, name)
	}

	base := sizes["base"]
	grew := func(name string) bool { return sizes[name] > base }

	largest := ""
	for _, name := range snippets {
		if name == "base" {
			continue
		}
		if largest == "" || sizes[name] > sizes[largest] {
			largest = name
		}
	}

	predict.Check(t, map[string]any{
		"fmt_grows_binary":          grew("fmtprint"),
		"reflect_grows_binary":      grew("reflectuse"),
		"generic_grows_binary":      grew("generic"),
		"iface_method_grows_binary": grew("ifacemethod"),
		"unused_table_grows_binary": grew("bigtable"),
		"largest_addition":          largest,
	})
}
