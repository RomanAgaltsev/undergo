package cli

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
)

// relPathRE matches a ../-rooted path to a file, however it is written: as a
// markdown link target, inside backticks, or bare. The imported corpora use
// backticks far more often than links, so matching only links would check
// almost nothing.
var relPathRE = regexp.MustCompile(`(?:\.\./)+[A-Za-z0-9_./-]+\.[A-Za-z0-9]+`)

// brokenRelPaths returns the relative paths in body that do not resolve from dir.
func brokenRelPaths(dir, body string) []string {
	var broken []string
	for _, ref := range relPathRE.FindAllString(body, -1) {
		target := filepath.Join(dir, filepath.FromSlash(ref))
		if _, err := os.Stat(target); err != nil && !slices.Contains(broken, ref) {
			broken = append(broken, ref)
		}
	}
	return broken
}
