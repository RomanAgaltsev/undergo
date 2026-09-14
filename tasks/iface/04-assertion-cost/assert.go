// Package assert is about what a type assertion compiles to.
//
// You do not implement anything here. Work out which assertions need help from
// the runtime, write your answers into prediction.yaml, then verify.
package assert

// Speaker is the interface under test.
type Speaker interface{ Speak() string }

// Named is a second interface, for interface-to-interface assertions.
type Named interface{ Name() string }

// Dog satisfies both.
type Dog struct{ N string }

// Speak satisfies Speaker.
func (d Dog) Speak() string { return "woof" }

// Name satisfies Named.
func (d Dog) Name() string { return d.N }

// ConcreteCommaOk asserts to a concrete type, two-value form.
func ConcreteCommaOk(v any) (Dog, bool) {
	d, ok := v.(Dog)
	return d, ok
}

// ConcretePanicking asserts to a concrete type, one-value form.
func ConcretePanicking(v any) Dog {
	return v.(Dog)
}

// EmptyToInterface asserts an empty interface to a one-method interface.
func EmptyToInterface(v any) (Named, bool) {
	n, ok := v.(Named)
	return n, ok
}

// InterfaceToInterface asserts one non-empty interface to another.
func InterfaceToInterface(s Speaker) (Named, bool) {
	n, ok := s.(Named)
	return n, ok
}

// TypeSwitchConcrete switches over three concrete types.
func TypeSwitchConcrete(v any) int {
	switch v.(type) {
	case int:
		return 1
	case string:
		return 2
	case Dog:
		return 3
	}
	return 0
}
