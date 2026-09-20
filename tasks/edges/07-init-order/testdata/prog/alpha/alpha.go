package alpha

import "fmt"

// A is initialised when this package is initialised, which is before any
// package that imports it.
var A = announce("alpha.A")

func announce(s string) string { fmt.Println(s); return s }

func init() { fmt.Println("alpha.init") }
