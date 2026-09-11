// Package drill — C9/04 modulo-negative.
package drill

// Bucket maps a key to one of the buckets by hashing modulo the bucket count.
func Bucket(key int, buckets []string) string {
	idx := key % len(buckets)
	return buckets[idx]
}
