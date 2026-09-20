// Package loopvar runs the most-taught bug in Go under two declared language
// versions, with and without the workaround that used to be necessary.
package loopvar

import (
	"fmt"
	"strings"

	"github.com/RomanAgaltsev/undergo/internal/goline"
)

// Program launches three goroutines that close over a loop variable, once
// plainly and once with the classic `i := i` shadow, and prints both results
// sorted so the output does not depend on scheduling.
//
// The goroutines wait on a start gate that is closed after the loop has
// finished. Without it the answer would be a race with the scheduler rather
// than a fact about the language: a goroutine that happens to run during the
// loop reads whatever the shared variable holds at that moment. The gate is
// what makes "the loop is over when they read it" true instead of likely.
//
// The two slices are separated by a vertical bar: a slice prints with spaces
// inside it, so any whitespace-delimited parse of two slices would be wrong.
const Program = `package main

import (
	"fmt"
	"sort"
	"sync"
)

func collect(shadow bool) []int {
	var mu sync.Mutex
	var got []int
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 3; i++ {
		wg.Add(1)
		if shadow {
			i := i
			go func() {
				defer wg.Done()
				<-start
				mu.Lock()
				got = append(got, i)
				mu.Unlock()
			}()
			continue
		}
		go func() {
			defer wg.Done()
			<-start
			mu.Lock()
			got = append(got, i)
			mu.Unlock()
		}()
	}
	close(start)
	wg.Wait()
	sort.Ints(got)
	return got
}

func main() { fmt.Printf("%v|%v", collect(false), collect(true)) }
`

// Results runs Program at the given go line and returns the plain and shadowed
// outputs, in that order.
func Results(line string) (plain, shadowed string, err error) {
	res, err := goline.Run(line, Program)
	if err != nil {
		return "", "", err
	}
	a, b, ok := strings.Cut(res.Output, "|")
	if !ok {
		return "", "", fmt.Errorf("unexpected output %q: no separator", res.Output)
	}
	return a, b, nil
}
