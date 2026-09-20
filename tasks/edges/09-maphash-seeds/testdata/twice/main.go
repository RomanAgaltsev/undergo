// Command twice prints one maphash of a fixed key, using a seed made in this
// process. Run it twice and compare the two lines.
package main

import (
	"fmt"
	"hash/maphash"
)

func main() {
	seed := maphash.MakeSeed()
	fmt.Print(maphash.String(seed, "undergo"))
}
