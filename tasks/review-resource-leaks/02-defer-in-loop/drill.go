// Package drill — C5/02 defer-in-loop.
package drill

import (
	"io"
	"os"
)

// ConcatFiles reads every path and returns their contents concatenated.
func ConcatFiles(paths []string) ([]byte, error) {
	var out []byte
	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		b, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}
		out = append(out, b...)
	}
	return out, nil
}

// ReadOne reads a single file and closes it promptly.
func ReadOne(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}
