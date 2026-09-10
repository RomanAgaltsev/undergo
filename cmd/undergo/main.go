// Command undergo runs the undergo gym: listing, starting, verifying and
// revealing tasks.
package main

import (
	"os"

	"github.com/RomanAgaltsev/undergo/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
