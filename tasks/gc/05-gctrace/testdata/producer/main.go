// Command producer allocates enough to trigger a dozen collections, and forces
// two of them, so that a recorded trace contains both kinds of line.
package main

import "runtime"

var sink [][]byte

func main() {
	for i := range 40 {
		sink = make([][]byte, 0, 256)
		for range 256 {
			sink = append(sink, make([]byte, 4096))
		}
		if i == 10 || i == 25 {
			runtime.GC()
		}
	}
}
