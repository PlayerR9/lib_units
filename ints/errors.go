package ints

import "strings"

// ErrWhileAt represents an error that occurs while performing an operation at a specific index.
type ErrWhileAt struct {
	// Index is the index where the error occurred.
	Index int

	// Element is the element where the index is pointing to.
	Element string

	// Operation is the operation that was being performed.
	Operation string

	// Reason is the reason for the error.
	Reason error
}

// Error implements the Unwrapper interface.
//
// Message: "an error occurred while {operation} at index {index} {element}: {reason}"
//
// However, if the reason is nil, the message is "an error occurred while {operation}
// at index {index} {element}" instead.
func (e *ErrWhileAt) Error() string {
	var builder strings.Builder

	if e.Reason == nil {
		builder.WriteString("an error occurred ")
	}

	builder.WriteString("while ")
	builder.WriteString(e.Operation)
	builder.WriteRune(' ')
	builder.WriteString(GetOrdinalSuffix(e.Index))
	builder.WriteRune(' ')
	builder.WriteString(e.Element)

	if e.Reason != nil {
		builder.WriteString(": ")
		builder.WriteString(e.Reason.Error())
	}

	return builder.String()
}

// Unwrap implements the Unwrapper interface.
func (e *ErrWhileAt) Unwrap() error {
	return e.Reason
}

// ChangeReason implements the Unwrapper interface.
func (e *ErrWhileAt) ChangeReason(reason error) {
	e.Reason = reason
}

// NewErrWhileAt creates a new ErrWhileAt error.
//
// Parameters:
//   - operation: The operation that was being performed.
//   - index: The index where the error occurred.
//   - elem: The element where the index is pointing to.
//   - reason: The reason for the error.
//
// Returns:
//   - *ErrWhileAt: A pointer to the newly created ErrWhileAt.
func NewErrWhileAt(operation string, index int, elem string, reason error) *ErrWhileAt {
	e := &ErrWhileAt{
		Index:     index,
		Operation: operation,
		Element:   elem,
		Reason:    reason,
	}
	return e
}

// ErrAt represents an error that occurs at a specific index.
type ErrAt struct {
	// Index is the index where the error occurred.
	Index int

	// Name is the name of the index.
	Name string

	// Reason is the reason for the error.
	Reason error
}

// Error implements the Unwrapper interface.
//
// Message: "something went wrong at the {index} {name}: {reason}".
//
// However, if the reason is nil, the message is "something went wrong at the {index}
// {name}" instead.
func (e *ErrAt) Error() string {
	var builder strings.Builder

	if e.Reason == nil {
		builder.WriteString("something went wrong at the ")
	}

	var name string

	if e.Name != "" {
		name = e.Name
	} else {
		name = "index"
	}

	builder.WriteString(GetOrdinalSuffix(e.Index))
	builder.WriteRune(' ')
	builder.WriteString(name)

	if e.Reason != nil {
		builder.WriteString(" is invalid: ")
		builder.WriteString(e.Reason.Error())
	}

	return builder.String()
}

// Unwrap implements the Unwrapper interface.
func (e *ErrAt) Unwrap() error {
	return e.Reason
}

// ChangeReason implements the Unwrapper interface.
func (e *ErrAt) ChangeReason(reason error) {
	e.Reason = reason
}

// NewErrAt creates a new ErrAt error.
//
// Parameters:
//   - index: The index where the error occurred.
//   - name: The name of the index.
//   - reason: The reason for the error.
//
// Returns:
//   - *ErrAt: A pointer to the newly created ErrAt.
func NewErrAt(index int, name string, reason error) *ErrAt {
	return &ErrAt{
		Index:  index,
		Name:   name,
		Reason: reason,
	}
}
