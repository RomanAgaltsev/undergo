// Command byname reaches a Widget method through reflection, using a name that
// is not known until the program runs.
package main

import (
	"fmt"
	"os"

	widget "github.com/RomanAgaltsev/undergo/tasks/reflect/05-methodbyname-linker"
)

func main() {
	name := "Alpha"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	fmt.Println(widget.CallByName(widget.Widget{N: 1}, name))
}
