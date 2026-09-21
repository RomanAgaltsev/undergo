// Package pauses is about a metric that split in two, twice.
//
// A stop-the-world pause has two parts: getting every goroutine to stop, and
// the time the world spent stopped. runtime/metrics reported one number for
// years. The task is which names exist under which toolchain — a question with
// an exact answer, unlike the durations behind them.
package pauses

import (
	"fmt"
	"sort"
	"strings"

	"github.com/RomanAgaltsev/undergo/internal/toolchain"
)

// Program prints every runtime/metrics name containing "pauses", sorted and
// separated by "|". It reports names, never values: a duration cannot be
// graded.
const Program = `package main

import (
	"fmt"
	"runtime/metrics"
	"sort"
	"strings"
)

func main() {
	var names []string
	for _, d := range metrics.All() {
		if strings.Contains(d.Name, "pauses") {
			names = append(names, d.Name)
		}
	}
	sort.Strings(names)
	fmt.Print(strings.Join(names, "|"))
}
`

// Local is the toolchain running this test.
const Local = toolchain.Local

// Names runs Program under the named toolchain and returns the sorted metric
// names it printed. The go line is 1.21 in every run and does not vary.
func Names(name string) ([]string, error) {
	res, err := toolchain.Run(name, "1.21", Program)
	if err != nil {
		return nil, err
	}
	if !res.Exited0 {
		return nil, fmt.Errorf("%s: program did not run: %s", name, res.Output)
	}
	if res.Output == "" {
		return nil, fmt.Errorf("%s: no pause metrics at all", name)
	}
	names := strings.Split(res.Output, "|")
	sort.Strings(names)
	return names, nil
}

// Has reports whether names contains want.
func Has(names []string, want string) bool {
	for _, n := range names {
		if n == want {
			return true
		}
	}
	return false
}
