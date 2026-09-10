// Package envdecode decodes environment variables into a struct using field tags.
package envdecode

// Decode fills dst — which must be a non-nil pointer to a struct — from lookup.
//
// For each exported field carrying an `env:"NAME"` tag, Decode looks up NAME and
// assigns the value. Supported field kinds: string, int, int64, bool.
//
// A tag of `env:"NAME,required"` makes a missing variable an error. A field with
// no env tag is left alone. Unexported fields are skipped.
//
// Decode returns an error when dst is not a pointer to a struct, when a required
// variable is absent, or when a value cannot be parsed into the field's kind.
func Decode(dst any, lookup func(string) (string, bool)) error {
	panic("undergo: implement Decode")
}
