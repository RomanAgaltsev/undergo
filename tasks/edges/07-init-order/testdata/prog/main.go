package main

import (
	"fmt"

	"prog/beta"
)

// Declared first, initialised second: Second depends on First.
var Second = announce("main.Second") + First

var First = announce("main.First")

func announce(s string) string { fmt.Println(s); return s }

func init() { fmt.Println("main.init") }

func main() { fmt.Println("main.main", beta.B != "" && Second != "") }
