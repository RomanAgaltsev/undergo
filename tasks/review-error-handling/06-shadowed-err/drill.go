// Package drill — C3/06 shadowed-err.
package drill

import "os"

// Save writes data to path, reporting any error.
func Save(path string, data []byte) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		// write failed — recorded below
	}
	return err
}
