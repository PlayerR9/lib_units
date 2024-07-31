package reflect

import (
	"reflect"
	"strings"
)

// ErrInvalidCall represents an error that occurs when a function
// is not called correctly.
type ErrInvalidCall struct {
	// FnName is the name of the function.
	FnName string

	// Signature is the Signature of the function.
	Signature reflect.Type

	// Reason is the Reason for the failure.
	Reason error
}

// Error implements the Unwrapper interface.
//
// Message: "call to {function}({signature}) failed: {reason}".
//
// However, if the reason is nil, the message is "call to {function}({signature})
// failed" instead.
func (e *ErrInvalidCall) Error() string {
	var builder strings.Builder

	builder.WriteString("call to ")
	builder.WriteString(e.FnName)
	builder.WriteString(e.Signature.String())
	builder.WriteString(" failed")

	if e.Reason != nil {
		builder.WriteString(": ")
		builder.WriteString(e.Reason.Error())
	}

	return builder.String()
}

// Unwrap implements the Unwrapper interface.
func (e *ErrInvalidCall) Unwrap() error {
	return e.Reason
}

// ChangeReason implements the Unwrapper interface.
func (e *ErrInvalidCall) ChangeReason(reason error) {
	e.Reason = reason
}

// NewErrInvalidCall creates a new ErrInvalidCall.
//
// Parameters:
//   - functionName: The name of the function.
//   - function: The function that failed.
//   - reason: The reason for the failure.
//
// Returns:
//   - *ErrInvalidCall: A pointer to the new ErrInvalidCall.
func NewErrInvalidCall(functionName string, function any, reason error) *ErrInvalidCall {
	return &ErrInvalidCall{
		FnName:    functionName,
		Signature: reflect.ValueOf(function).Type(),
		Reason:    reason,
	}
}

// ErrEmpty represents an error when a value is empty.
type ErrEmpty[T any] struct {
	// Value is the value that caused the error.
	Value T
}

// Error implements the error interface.
//
// Message: "<type> must not be empty"
func (e *ErrEmpty[T]) Error() string {
	type_str := TypeOf(e.Value)

	var builder strings.Builder

	builder.WriteString(type_str)
	builder.WriteString(" must not be empty")

	return builder.String()
}

// NewErrEmpty creates a new ErrEmpty error.
//
// Parameters:
//   - value: The value that caused the error.
//
// Returns:
//   - *ErrEmpty: A pointer to the newly created ErrEmpty.
func NewErrEmpty[T any](value T) *ErrEmpty[T] {
	e := &ErrEmpty[T]{
		Value: value,
	}
	return e
}
