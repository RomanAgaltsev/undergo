// Package widget has four methods and two ways of calling one of them.
//
// The linker drops what it can prove is unreachable. Reflection makes that
// proof impossible.
package widget

import "reflect"

// Widget has four methods, none of which does anything interesting.
type Widget struct{ N int }

// Alpha returns N plus one.
func (w Widget) Alpha() int { return w.N + 1 }

// Bravo returns N plus two.
func (w Widget) Bravo() int { return w.N + 2 }

// Charlie returns N plus three.
func (w Widget) Charlie() int { return w.N + 3 }

// Delta returns N plus four.
func (w Widget) Delta() int { return w.N + 4 }

// CallDirect names a method at compile time.
func CallDirect(w Widget) int { return w.Alpha() }

// CallByName reaches a method through reflection, so which one is needed is not
// known until the program runs.
func CallByName(w Widget, name string) int {
	m := reflect.ValueOf(w).MethodByName(name)
	if !m.IsValid() {
		return -1
	}
	return int(m.Call(nil)[0].Int())
}
