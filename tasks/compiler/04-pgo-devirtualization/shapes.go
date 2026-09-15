// Package shapes asks what a profile lets the compiler assume about an
// interface call.
//
// Area is reached through an interface at two call sites. One implementation
// dominates the profile checked in as default.pgo; the others barely appear.
// The compiler can use that to guess, and a guess it can check cheaply is
// worth making.
//
// The Area implementations are deliberately too large to inline. That is not
// decoration: a profile can only record a call edge that actually happened, and
// an inlined callee makes no call.
package shapes

import (
	"os/exec"
	"strings"
)

// Shape is the interface both call sites go through.
type Shape interface{ Area() float64 }

// Rect dominates the profile.
type Rect struct{ W, H float64 }

// Area returns the rectangle's area, the long way round.
func (r Rect) Area() float64 {
	total := 0.0
	for i := range 8 {
		total += r.W * r.H / 8
		total += float64(i) * 0
		if total < 0 {
			total = -total
		}
	}
	return total
}

// Circle appears rarely.
type Circle struct{ R float64 }

// Area returns the circle's area, the long way round.
func (c Circle) Area() float64 {
	total := 0.0
	for i := range 8 {
		total += 3.14159 * c.R * c.R / 8
		total += float64(i) * 0
		if total < 0 {
			total = -total
		}
	}
	return total
}

// Tri appears rarely.
type Tri struct{ B, H float64 }

// Area returns the triangle's area, the long way round.
func (t Tri) Area() float64 {
	total := 0.0
	for i := range 8 {
		total += t.B * t.H / 16
		total += float64(i) * 0
		if total < 0 {
			total = -total
		}
	}
	return total
}

// Hot is the call site the profile is dominated by.
func Hot(s Shape) float64 { return s.Area() }

// Cold is a second call site, barely exercised.
func Cold(s Shape) float64 { return s.Area() + 1 }

// Notes returns the compiler's optimisation notes for this package.
//
// The profile is collected on demand. A default.pgo is picked up automatically
// only for a main package; for a library it has to be named.
func Notes(withPGO bool) (string, error) {
	args := []string{"build", "-gcflags=-m=2"}
	if withPGO {
		p, perr := profile()
		if perr != nil {
			return "", perr
		}
		args = append(args, "-pgo="+p)
	} else {
		args = append(args, "-pgo=off")
	}
	args = append(args, ".")

	out, err := exec.Command("go", args...).CombinedOutput()
	if err != nil && len(out) == 0 {
		return "", err
	}
	return string(out), nil
}

// MentionsDevirtualizing reports whether the notes record a PGO
// devirtualization of a call to the named concrete method.
func MentionsDevirtualizing(notes, method string) bool {
	for _, line := range strings.Split(notes, "\n") {
		if strings.Contains(line, "PGO devirtualizing") && strings.Contains(line, method) {
			return true
		}
	}
	return false
}

// CountDevirtualizations reports how many interface call sites the compiler
// devirtualized in this package.
func CountDevirtualizations(notes string) int {
	n := 0
	for _, line := range strings.Split(notes, "\n") {
		if strings.Contains(line, "PGO devirtualizing") {
			n++
		}
	}
	return n
}
