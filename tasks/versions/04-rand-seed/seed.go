// Package seed asks whether math/rand still listens to rand.Seed.
package seed

import "github.com/RomanAgaltsev/undergo/internal/goline"

// Program seeds the global source twice with the same value and reports
// whether it got the same number back both times.
//
// It draws Int63 rather than a small Intn: when seeding is a no-op the two
// draws come from one unseeded stream, and two draws from a small range would
// collide by chance often enough to make the answer flaky. Across 2^63 the
// collision is not worth a thought.
const Program = `package main

import (
	"fmt"
	"math/rand"
)

func main() {
	rand.Seed(42)
	first := rand.Int63()
	rand.Seed(42)
	fmt.Print(first == rand.Int63())
}
`

// SeedHonoured runs Program at the given go line with the given extra
// environment and reports whether seeding was reproducible.
func SeedHonoured(line string, env ...string) (bool, error) {
	res, err := goline.Run(line, Program, env...)
	if err != nil {
		return false, err
	}
	return res.Output == "true", nil
}
