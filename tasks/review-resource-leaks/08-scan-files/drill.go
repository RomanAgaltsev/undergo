// Package drill — C5/08 scan-files.
package drill

import (
	"bufio"
	"fmt"
	"os"
)

// CountLines returns the total number of lines across all paths.
func CountLines(paths []string) (int, error) {
	total := 0
	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			return 0, err
		}
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			if len(sc.Bytes()) > 1<<20 {
				return 0, fmt.Errorf("%s: line too long", p)
			}
			total++
		}
		defer f.Close()
	}
	return total, nil
}
