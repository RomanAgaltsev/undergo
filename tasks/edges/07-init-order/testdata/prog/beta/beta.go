package beta

import (
	"fmt"

	"prog/alpha"
)

// B depends on alpha.A, so alpha must be fully initialised first.
var B = announce("beta.B (sees " + alpha.A + ")")

func announce(s string) string { fmt.Println("beta.B"); return s }

func init() { fmt.Println("beta.init") }
