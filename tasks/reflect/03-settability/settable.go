// Package settable performs eight reflect operations, some of which panic.
//
// You do not implement anything here. Work out which ones panic, write your
// answers into prediction.yaml, then verify.
package settable

import "reflect"

// Config is the subject of most of the operations below.
type Config struct {
	Name    string
	Timeout int
	tag     string //nolint:unused // present so an operation can try to set it
}

// SetThroughValue takes the struct by value and sets a field.
func SetThroughValue() {
	c := Config{Name: "a"}
	reflect.ValueOf(c).FieldByName("Name").SetString("b")
}

// SetThroughPointer takes the struct by pointer and sets a field.
func SetThroughPointer() {
	c := Config{Name: "a"}
	reflect.ValueOf(&c).Elem().FieldByName("Name").SetString("b")
}

// SetUnexported reaches an unexported field through a pointer.
func SetUnexported() {
	c := Config{}
	reflect.ValueOf(&c).Elem().FieldByName("tag").SetString("b")
}

// SetWrongType assigns a string where an int lives.
func SetWrongType() {
	c := Config{}
	reflect.ValueOf(&c).Elem().FieldByName("Timeout").SetString("b")
}

// SetSliceElement sets an element of a slice obtained by value.
//
// The slice header is a copy; the backing array is not.
func SetSliceElement() {
	s := []int{1, 2, 3}
	reflect.ValueOf(s).Index(0).SetInt(99)
}

// SetMapValueInPlace reaches into a map value and sets a field on it.
func SetMapValueInPlace() {
	m := map[string]Config{"a": {Name: "x"}}
	reflect.ValueOf(m).MapIndex(reflect.ValueOf("a")).FieldByName("Name").SetString("y")
}

// ElemOfNilPointer dereferences a nil pointer and sets through it.
func ElemOfNilPointer() {
	var c *Config
	reflect.ValueOf(c).Elem().FieldByName("Name").SetString("b")
}

// AppendAndAssign appends through reflect and assigns the result back.
func AppendAndAssign() {
	s := []int{1, 2, 3}
	v := reflect.ValueOf(&s).Elem()
	v.Set(reflect.Append(v, reflect.ValueOf(4)))
}
