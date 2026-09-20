// Package seeds is about what hash/maphash actually promises.
//
// A map's iteration order is deliberately randomised, and the hash behind it is
// seeded per process. hash/maphash exposes that hash, seed and all — which
// makes it excellent for hash tables and wrong for anything that has to agree
// with another process, or with itself tomorrow.
package seeds

import (
	"hash/maphash"
	"os/exec"
	"path/filepath"
	"strings"
)

// Key is the one input every slot hashes.
const Key = "undergo"

// InProcess hashes Key with two independently made seeds and reports the three
// comparisons that can be made without leaving this process.
func InProcess() (sameSeedAgrees, differentSeedsAgree, stringAndBytesAgree bool) {
	a, b := maphash.MakeSeed(), maphash.MakeSeed()

	first := maphash.String(a, Key)
	return first == maphash.String(a, Key),
		first == maphash.String(b, Key),
		first == maphash.Bytes(a, []byte(Key))
}

// AcrossProcesses runs testdata/twice twice and reports whether the two
// processes printed the same hash for the same key.
func AcrossProcesses() (bool, error) {
	run := func() (string, error) {
		cmd := exec.Command("go", "run", ".")
		cmd.Dir = filepath.Join("testdata", "twice")
		out, err := cmd.CombinedOutput()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(out)), nil
	}
	first, err := run()
	if err != nil {
		return false, err
	}
	second, err := run()
	if err != nil {
		return false, err
	}
	return first == second, nil
}
