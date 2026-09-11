// Package drill — C8/07 weak-password-hash.
package drill

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashPassword returns a hex-encoded hash of the password for storage.
func HashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}
