package bytes

import (
	"bytes"
	"strconv"
	"strings"

	uc "github.com/PlayerR9/lib_units/common"
)

// filter_equals returns the indices of the other in the data.
//
// Parameters:
//   - indices: The indices.
//   - data: The data.
//   - other: The other value.
//   - offset: The offset to start the search from.
//
// Returns:
//   - []int: The indices.
func filter_equals(indices []int, data []byte, other byte, offset int) []int {
	var top int

	for i := 0; i < len(indices); i++ {
		idx := indices[i]

		if data[idx+offset] == other {
			indices[top] = idx
			top++
		}
	}

	indices = indices[:top]

	return indices
}

// Indices returns the indices of the separator in the data.
//
// Parameters:
//   - data: The data.
//   - sep: The separator.
//   - exclude_sep: Whether the separator is inclusive. If set to true, the indices will point to the character right after the
//     separator. Otherwise, the indices will point to the character right before the separator.
//
// Returns:
//   - []int: The indices.
func IndicesOf(data []byte, sep []byte, exclude_sep bool) []int {
	if len(data) == 0 || len(sep) == 0 {
		return nil
	}

	var indices []int

	for i := 0; i < len(data)-len(sep); i++ {
		if data[i] == sep[0] {
			indices = append(indices, i)
		}
	}

	if len(indices) == 0 {
		return nil
	}

	for i := 1; i < len(sep); i++ {
		other := sep[i]

		indices = filter_equals(indices, data, other, i)

		if len(indices) == 0 {
			return nil
		}
	}

	if exclude_sep {
		for i := 0; i < len(indices); i++ {
			indices[i] += len(sep)
		}
	}

	return indices
}

// FirstIndex returns the first index of the target in the tokens.
//
// Parameters:
//   - tokens: The slice of tokens in which to search for the target.
//   - target: The target to search for.
//
// Returns:
//   - int: The index of the target. -1 if the target is not found.
//
// If either tokens or the target are empty, it returns -1.
func FirstIndex(tokens [][]byte, target []byte) int {
	if len(tokens) == 0 || len(target) == 0 {
		return -1
	}

	for i, token := range tokens {
		if bytes.Equal(token, target) {
			return i
		}
	}

	return -1
}

// FindContentIndexes searches for the positions of opening and closing
// tokens in a slice of strings.
//
// Parameters:
//   - op_token: The string that marks the beginning of the content.
//   - cl_token: The string that marks the end of the content.
//   - tokens: The slice of strings in which to search for the tokens.
//
// Returns:
//   - result: An array of two integers representing the start and end indexes
//     of the content.
//   - err: Any error that occurred while searching for the tokens.
//
// Errors:
//   - *uc.ErrInvalidParameter: If the closingToken is an empty string.
//   - *ErrTokenNotFound: If the opening or closing token is not found in the
//     content.
//   - *ErrNeverOpened: If the closing token is found without any
//     corresponding opening token.
//
// Behaviors:
//   - The first index of the content is inclusive, while the second index is
//     exclusive.
//   - This function returns a partial result when errors occur. ([-1, -1] if
//     errors occur before finding the opening token, [index, 0] if the opening
//     token is found but the closing token is not found.
func FindContentIndexes(op_token, cl_token []byte, tokens [][]byte) (result [2]int, err error) {
	result[0] = -1
	result[1] = -1

	if len(cl_token) == 0 {
		err = uc.NewErrInvalidParameter("cl_token", uc.NewErrEmpty(cl_token))
		return
	}

	op_tok_idx := FirstIndex(tokens, op_token)
	if op_tok_idx == -1 {
		err = NewErrTokenNotFound(op_token, true)
		return
	}

	result[0] = op_tok_idx + 1

	balance := 1
	cl_tok_idx := -1

	for i := result[0]; i < len(tokens) && cl_tok_idx == -1; i++ {
		curr_tok := tokens[i]

		if bytes.Equal(curr_tok, cl_token) {
			balance--

			if balance == 0 {
				cl_tok_idx = i
			}
		} else if bytes.Equal(curr_tok, op_token) {
			balance++
		}
	}

	if cl_tok_idx != -1 {
		result[1] = cl_tok_idx + 1
		return
	}

	if balance < 0 {
		err = NewErrNeverOpened(op_token, cl_token)
		return
	} else if balance != 1 || bytes.Equal(cl_token, []byte("\n")) {
		err = NewErrTokenNotFound(cl_token, false)
		return
	}

	result[1] = len(tokens)
	return
}

// TrimEmpty removes empty bytes from a slice of bytes; including any empty
// spaces at the beginning and end of the bytes.
//
// Parameters:
//   - values: The values to trim.
//
// Returns:
//   - [][]byte: The trimmed values.
func TrimEmpty(values [][]byte) [][]byte {
	if len(values) == 0 {
		return values
	}

	var top int

	for i := 0; i < len(values); i++ {
		current_value := values[i]

		res := bytes.TrimSpace(current_value)
		if len(res) > 0 {
			values[top] = res
			top++
		}
	}

	values = values[:top]

	return values
}

// AndString is a function that returns a string representation of a slice
// of bytes.
//
// Parameters:
//   - values: The values to convert to a string.
//   - quote: Whether to quote the values.
//
// Returns:
//   - string: The string representation of the values.
func AndString(values [][]byte, quote bool) string {
	values = TrimEmpty(values)
	if len(values) == 0 {
		return ""
	}

	if len(values) == 1 {
		if !quote {
			return string(values[0])
		} else {
			return strconv.Quote(string(values[0]))
		}
	}

	elems := make([]string, 0, len(values))

	if quote {
		for i := 0; i < len(values); i++ {
			elems = append(elems, strconv.Quote(string(values[i])))
		}
	} else {
		for i := 0; i < len(values); i++ {
			elems = append(elems, string(values[i]))
		}
	}

	var builder strings.Builder

	builder.WriteString(elems[0])

	if len(elems) > 2 {
		builder.WriteString(strings.Join(elems[1:len(elems)-1], ", "))
		builder.WriteRune(',')
	}

	builder.WriteString(" and ")
	builder.WriteString(elems[len(elems)-1])

	return builder.String()
}

// EitherOrString is a function that returns a string representation of a slice
// of bytes.
//
// Parameters:
//   - values: The values to convert to a string.
//   - quote: True if the values should be quoted, false otherwise.
//
// Returns:
//   - string: The string representation.
//
// Example:
//
//	EitherOrString(bytes.Fields([]byte("a b c"}, false) // "a, b or c"
func EitherOrString(values [][]byte, quote bool) string {
	values = TrimEmpty(values)
	if len(values) == 0 {
		return ""
	}

	if len(values) == 1 {
		if !quote {
			return string(values[0])
		} else {
			return strconv.Quote(string(values[0]))
		}
	}

	elems := make([]string, 0, len(values))

	if quote {
		for _, v := range values {
			elems = append(elems, strconv.Quote(string(v)))
		}
	} else {
		for _, v := range values {
			elems = append(elems, string(v))
		}
	}

	var builder strings.Builder

	builder.WriteString("either ")
	builder.WriteString(elems[0])

	if len(values) > 2 {
		builder.WriteString(strings.Join(elems[1:len(elems)-1], ", "))
		builder.WriteRune(',')
	}

	builder.WriteString(" or ")
	builder.WriteString(elems[len(elems)-1])

	return builder.String()
}

// OrString is a function that returns a string representation of a slice of
// bytes.
//
// Parameters:
//   - values: The values to convert to a string.
//   - quote: True if the values should be quoted, false otherwise.
//   - is_negative: True if the string should use "nor" instead of "or", false
//     otherwise.
//
// Returns:
//   - string: The string representation.
//
// Example:
//
//	OrString(bytes.Fields([]byte("a b c"), false, true) // "a, b, nor c"
func OrString(values [][]byte, quote, is_negative bool) string {
	values = TrimEmpty(values)
	if len(values) == 0 {
		return ""
	}

	if len(values) == 1 {
		if !quote {
			return string(values[0])
		} else {
			return strconv.Quote(string(values[0]))
		}
	}

	var sep string

	if is_negative {
		sep = " nor "
	} else {
		sep = " or "
	}

	elems := make([]string, 0, len(values))

	if quote {
		for _, v := range values {
			elems = append(elems, strconv.Quote(string(v)))
		}
	} else {
		for _, v := range values {
			elems = append(elems, string(v))
		}
	}

	var builder strings.Builder

	builder.WriteString(elems[0])

	if len(values) > 2 {
		builder.WriteString(strings.Join(elems[1:len(elems)-1], ", "))
		builder.WriteRune(',')
	}

	builder.WriteString(sep)
	builder.WriteString(elems[len(elems)-1])

	return builder.String()
}

// FindByte searches for the first occurrence of a byte in a byte slice starting from a given index.
//
// Parameters:
//   - data: the byte slice to search in.
//   - from: the index to start the search from. If negative, it is treated as 0.
//   - sep: the byte to search for.
//
// Returns:
//   - int: the index of the first occurrence of the byte in the byte slice, or -1 if not found.
func FindByte(data []byte, from int, sep byte) int {
	if len(data) == 0 || from >= len(data) {
		return -1
	}

	len_data := len(data)

	if from < 0 {
		from = 0
	}

	for i := from; i < len_data; i++ {
		if data[i] == sep {
			return i
		}
	}

	return -1
}

// FindByteReversed searches for the first occurrence of a byte in a byte slice starting from a given index in reverse order.
//
// Parameters:
//   - data: the byte slice to search in.
//   - from: the index to start the search from. If greater than or equal to the length of the byte slice,
//     it is treated as the length of the byte slice minus 1.
//   - sep: the byte to search for.
//
// Returns:
//   - int: the index of the first occurrence of the byte in the byte slice in reverse order, or -1 if not found.
func FindByteReversed(data []byte, from int, sep byte) int {
	if len(data) == 0 || from < 0 {
		return -1
	}

	len_data := len(data)

	if from >= len_data {
		from = len_data - 1
	}

	for i := from; i >= 0; i-- {
		if data[i] == sep {
			return i
		}
	}

	return -1
}

// ReverseSearch searches for the last occurrence of a byte in a byte slice.
//
// Parameters:
//   - data: the byte slice to search in.
//   - from: the index to start the search from. If greater than or equal to the length of the byte slice,
//     it is treated as the length of the byte slice minus 1.
//   - sep: the byte to search for.
//
// Returns:
//   - int: the index of the last occurrence of the byte in the byte slice, or -1 if not found.
func ReverseSearch(data []byte, from int, sep []byte) int {
	if from < 0 || len(sep) == 0 || len(data) == 0 {
		return -1
	}

	sep_len := len(sep)

	if from+sep_len >= len(data) {
		from = len(data) - sep_len
	}

	if sep_len == 1 {
		return FindByteReversed(data, from, sep[0])
	}

	for {
		idx := FindByteReversed(data, from, sep[0])
		if idx == -1 {
			return -1
		}

		if bytes.Equal(data[idx:idx+sep_len], sep) {
			return idx
		}

		from = idx
	}
}

// ForwardSearch searches for the first occurrence of a byte in a byte slice.
//
// Parameters:
//   - data: the byte slice to search in.
//   - from: the index to start the search from. If negative, it is treated as 0.
//   - sep: the byte to search for.
//
// Returns:
//   - int: the index of the first occurrence of the byte in the byte slice, or -1 if not found.
func ForwardSearch(data []byte, from int, sep []byte) int {
	if len(sep) == 0 || len(data) == 0 || from+len(sep) >= len(data) {
		return -1
	}

	sep_len := len(sep)

	if from < 0 {
		from = 0
	}

	if sep_len == 1 {
		return FindByte(data, from, sep[0])
	}

	for {
		idx := FindByte(data, from, sep[0])
		if idx == -1 {
			return -1
		}

		if bytes.Equal(data[idx:idx+sep_len], sep) {
			return idx
		}

		from = idx
	}
}
