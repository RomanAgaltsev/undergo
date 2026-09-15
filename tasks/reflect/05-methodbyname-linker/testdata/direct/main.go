// Command direct calls exactly one Widget method, named at compile time.
package main

import (
	"fmt"

	widget "github.com/RomanAgaltsev/undergo/tasks/reflect/05-methodbyname-linker"
)

func main() {
	fmt.Println(widget.CallDirect(widget.Widget{N: 1}))
}
