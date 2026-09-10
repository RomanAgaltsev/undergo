// Command importer converts a source repository's exercises into undergo tasks.
//
// It is a one-off: each source is imported once, in a single `import:` commit,
// and the tool is kept so the conversion stays reproducible and reviewable.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	source := flag.String("source", "", "loupe | keystone")
	from := flag.String("from", "", "path to the source repository")
	to := flag.String("to", ".", "path to the undergo checkout")
	flag.Parse()

	if *source == "" || *from == "" {
		fmt.Fprintln(os.Stderr, "usage: importer -source loupe|keystone -from <path> [-to <path>]")
		os.Exit(2)
	}
	root, err := filepath.Abs(*to)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var n int
	switch *source {
	case "loupe":
		n, err = importLoupe(*from, root)
	default:
		err = fmt.Errorf("unknown source %q", *source)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("imported %d tasks from %s\n", n, *source)
}
