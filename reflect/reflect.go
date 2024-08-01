package reflect

import (
	"fmt"
	"reflect"
	"strings"
	// dbg "github.com/PlayerR9/lib_units/debug"
)

// AssertIfZero panics if the element is zero.
//
// Parameters:
//   - elem: The element to check.
//   - msg: The message to show if the element is zero.
//
// The panic message is the string msg.
func AssertIfZero(elem any, msg string) {
	value := reflect.ValueOf(elem)
	ok := value.IsZero()
	if ok {
		panic(msg)
	}
}

// TypeOf returns the type of the value as a string.
//
// Parameters:
//   - value: The value to get the type of.
//
// Returns:
//   - string: The type of the value.
func TypeOf(value any) string {
	if value == nil {
		return "nil"
	}

	return reflect.TypeOf(value).String()
}

// IsEmpty returns true if the element is empty.
//
// Parameters:
//   - elem: The element to check.
//
// Returns:
//   - bool: True if the element is empty, false otherwise.
func IsEmpty(elem any) bool {
	if elem == nil {
		return true
	}

	value := reflect.ValueOf(elem)
	return value.IsZero()
}

// FunctionCall represents a function call.
type FunctionCall struct {
	// Call is the string representation of the function call.
	Call string

	// Dependencies are the dependencies of the function call.
	Dependencies []string
}

func NewFunctionCall(call string, dependencies []string) FunctionCall {
	return FunctionCall{
		Call:         call,
		Dependencies: dependencies,
	}
}

// GetStringOf returns the string function call for the given element.
//
// Parameters:
//   - type_name: The name of the type.
//   - elem: The element to get the string of.
//   - custom: The custom strings to use. Empty values are ignored.
//
// Returns:
//   - FunctionCall: The function call.
func GetStringOf(type_name string, elem any, custom map[string][]string) FunctionCall {
	if elem == nil {
		return NewFunctionCall("\"nil\"", nil)
	}

	to := reflect.TypeOf(elem)
	// dbg.Assert(to != nil, "value must be non-nil")

	if custom != nil {
		values, ok := custom[type_name]
		if ok && len(values) > 0 {
			return NewFunctionCall(values[0], values[1:])
		}
	}

	var builder strings.Builder
	var dependencies []string

	switch to.String() {
	case "bool":
		builder.WriteString("strconv.FormatBool(")
		builder.WriteString(type_name)
		builder.WriteString(")")

		dependencies = append(dependencies, "strconv")
	case "byte":
		builder.WriteString("string(")
		builder.WriteString(type_name)
		builder.WriteString(")")
	case "complex64":
		builder.WriteString("strconv.FormatComplex(")
		builder.WriteString(type_name)
		builder.WriteString(", 'f', -1, 64)")

		dependencies = append(dependencies, "strconv")
	case "complex128":
		builder.WriteString("strconv.FormatComplex(")
		builder.WriteString(type_name)
		builder.WriteString(", 'f', -1, 128)")

		dependencies = append(dependencies, "strconv")
	case "float32":
		builder.WriteString("strconv.FormatFloat(")
		builder.WriteString(type_name)
		builder.WriteString(", 'f', -1, 32)")

		dependencies = append(dependencies, "strconv")
	case "float64":
		builder.WriteString("strconv.FormatFloat(")
		builder.WriteString(type_name)
		builder.WriteString(", 'f', -1, 64)")

		dependencies = append(dependencies, "strconv")
	case "int", "int8", "int16", "int32":
		builder.WriteString("strconv.FormatInt(int64(")
		builder.WriteString(type_name)
		builder.WriteString("), 10)")

		dependencies = append(dependencies, "strconv")
	case "int64":
		builder.WriteString("strconv.FormatInt(")
		builder.WriteString(type_name)
		builder.WriteString(", 10)")

		dependencies = append(dependencies, "strconv")
	case "rune":
		builder.WriteString("string(")
		builder.WriteString(type_name)
		builder.WriteString(")")
	case "string":
		builder.WriteString(type_name)
	case "uint", "uint8", "uint16", "uint32", "uintptr":
		builder.WriteString("strconv.FormatUint(uint64(")
		builder.WriteString(type_name)
		builder.WriteString("), 10)")

		dependencies = append(dependencies, "strconv")
	case "uint64":
		builder.WriteString("strconv.FormatUint(")
		builder.WriteString(type_name)
		builder.WriteString(", 10)")

		dependencies = append(dependencies, "strconv")
	default:
		ok := to.Implements(reflect.TypeOf((*fmt.Stringer)(nil)).Elem())
		if ok {
			builder.WriteString(type_name)
			builder.WriteString(".String()")

			return NewFunctionCall(builder.String(), dependencies)
		}

		ok = to.Implements(reflect.TypeOf((*fmt.GoStringer)(nil)).Elem())
		if ok {
			builder.WriteString(type_name)
			builder.WriteString(".GoString()")

			return NewFunctionCall(builder.String(), dependencies)
		}

		ok = to.Implements(reflect.TypeOf((*error)(nil)).Elem())
		if ok {
			builder.WriteString(type_name)
			builder.WriteString(".Error()")
		} else {
			builder.WriteString("fmt.Sprintf(\"%v\", ")
			builder.WriteString(type_name)

			dependencies = append(dependencies, "fmt")
		}
	}

	return NewFunctionCall(builder.String(), dependencies)
}
