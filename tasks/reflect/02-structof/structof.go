// Package structof builds a struct type at run time and copies a subset of
// another struct's fields into it.
//
// This is what a column projection does when the columns are not known until
// the query arrives, and it is the one job reflect.StructOf exists for.
package structof

// Project returns a new struct value holding only the named fields of v.
//
// The new type keeps each field's name, type and tags, so the result marshals
// to JSON with the same keys the original would have used for those fields.
//
// The contract the frozen tests enforce:
//
//   - the result holds exactly the named fields, in the order they were asked
//     for, with their original types and struct tags;
//   - naming no fields returns an empty struct value and no error;
//   - an unknown field name is an error that names the field;
//   - an unexported field name is an error, not a panic;
//   - a non-struct argument is an error.
//
// A pointer to a struct counts as a struct: dereference it.
func Project(v any, fields ...string) (any, error) {
	panic("undergo: implement Project")
}
