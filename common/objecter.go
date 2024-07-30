package common

import (
	"fmt"
)

// Objecter is an interface that defines the behavior of an object that can be
// copied, compared, and converted to a string.
type Objecter interface {
	fmt.Stringer
	Copier
	Equaler
}

// CopyOf creates a copy of the element by either calling the Copy method if the
// element implements the Copier interface or returning the element as is.
//
// Parameters:
//   - elem: The element to copy.
//
// Returns:
//   - any: A copy of the element.
func CopyOf(elem any) any {
	if elem == nil {
		return nil
	}

	switch elem := elem.(type) {
	case Objecter:
		c := elem.Copy()
		return c
	default:
		x := elem

		return x
	}
}

// EqualOf compares two objects of the same type. If any of the objects implements
// the Equaler interface, the Equals method is called. Otherwise, the objects are
// compared using the == operator. However, a is always checked first.
//
// Parameters:
//   - a: The first object to compare.
//   - b: The second object to compare.
//
// Returns:
//   - bool: True if the objects are equal, false otherwise.
//
// Behaviors:
//   - Nil objects are always considered different.
func EqualOf(a, b any) bool {
	if a == nil || b == nil {
		return false
	}

	switch a := a.(type) {
	case Objecter:
		otherB, ok := b.(Objecter)
		if !ok {
			return false
		}

		ok = a.Equals(otherB)
		return ok
	default:
		return a == b
	}
}
