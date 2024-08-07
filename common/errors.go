// Package errors provides a custom error type for out-of-bound errors.
package common

import (
	"fmt"
	"strings"
)

// ErrPanic represents an error when a panic occurs.
type ErrPanic struct {
	// Value is the value that caused the panic.
	Value any
}

// Error implements the error interface.
//
// Message: "panic: {value}"
func (e *ErrPanic) Error() string {
	var builder strings.Builder

	builder.WriteString("panic: ")
	fmt.Fprintf(&builder, "%v", e.Value)

	str := builder.String()

	return str
}

// NewErrPanic creates a new ErrPanic error.
//
// Parameters:
//   - value: The value that caused the panic.
//
// Returns:
//   - *ErrPanic: A pointer to the newly created ErrPanic.
func NewErrPanic(value any) *ErrPanic {
	e := &ErrPanic{
		Value: value,
	}

	return e
}

// ErrUnexpectedType represents an error when a value has an invalid type.
type ErrUnexpectedType[T any] struct {
	// Elem is the element that caused the error.
	Elem T

	// Kind is the category of the type that was expected.
	Kind string
}

// Error implements the error interface.
//
// Message: "type <type> is not a valid <kind> type"
func (e *ErrUnexpectedType[T]) Error() string {
	values := []string{
		"type",
		fmt.Sprintf("%T", e.Elem),
		"is not a valid",
		e.Kind,
		"type",
	}

	str := strings.Join(values, " ")

	return str
}

// NewErrUnexpectedType creates a new ErrUnexpectedType error.
//
// Parameters:
//   - typeName: The name of the type that was expected.
//   - elem: The element that caused the error.
//
// Returns:
//   - *ErrUnexpectedType: A pointer to the newly created ErrUnexpectedType.
func NewErrUnexpectedType[T any](kind string, elem T) *ErrUnexpectedType[T] {
	e := &ErrUnexpectedType[T]{
		Elem: elem,
		Kind: kind,
	}
	return e
}

// ErrExhaustedIter is an error type that is returned when an iterator
// is exhausted (i.e., there are no more elements to consume).
type ErrExhaustedIter struct{}

// Error implements the error interface.
//
// Message: "iterator is exhausted"
func (e *ErrExhaustedIter) Error() string {
	return "iterator is exhausted"
}

// NewErrExhaustedIter creates a new ErrExhaustedIter error.
//
// Returns:
//   - *ErrExhaustedIter: A pointer to the new error.
func NewErrExhaustedIter() *ErrExhaustedIter {
	e := &ErrExhaustedIter{}
	return e
}
