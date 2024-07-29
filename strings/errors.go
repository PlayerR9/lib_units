package strings

import (
	"strconv"
	"strings"
)

// ErrTokenNotFound is a struct that represents an error when a token is not
// found in the content.
type ErrTokenNotFound struct {
	// Token is the token that was not found in the content.
	Token string

	// IsOpening is the type of the token (opening or closing).
	IsOpening bool
}

// Error implements the error interface.
//
// Message: "{Type} token {Token} is not in the content"
func (e *ErrTokenNotFound) Error() string {
	var str_type string

	if e.IsOpening {
		str_type = "opening"
	} else {
		str_type = "closing"
	}

	values := []string{
		str_type,
		"token",
		"(",
		strconv.Quote(e.Token),
		")",
		"is not in the content",
	}

	msg := strings.Join(values, " ")

	return msg
}

// NewErrTokenNotFound is a constructor of ErrTokenNotFound.
//
// Parameters:
//   - token: The token that was not found in the content.
//   - is_opening: The type of the token (opening or closing).
//
// Returns:
//   - *ErrTokenNotFound: A pointer to the newly created error.
func NewErrTokenNotFound(token string, is_opening bool) *ErrTokenNotFound {
	e := &ErrTokenNotFound{
		Token:     token,
		IsOpening: is_opening,
	}
	return e
}

// ErrNeverOpened is a struct that represents an error when a closing
// token is found without a corresponding opening token.
type ErrNeverOpened struct {
	// OpeningToken is the opening token that was never closed.
	OpeningToken string

	// ClosingToken is the closing token that was found without a corresponding
	// opening token.
	ClosingToken string
}

// Error implements the error interface.
//
// Message:
//   - "closing token {ClosingToken} found without a corresponding opening token {OpeningToken}".
func (e *ErrNeverOpened) Error() string {
	values := []string{
		"closing token",
		"(",
		strconv.Quote(e.ClosingToken),
		")",
		"found without a corresponding opening token",
		"(",
		strconv.Quote(e.OpeningToken),
		")",
	}

	msg := strings.Join(values, " ")

	return msg
}

// NewErrNeverOpened is a constructor of ErrNeverOpened.
//
// Parameters:
//   - openingToken: The opening token that was never closed.
//   - closingToken: The closing token that was found without a corresponding opening token.
//
// Returns:
//   - *ErrNeverOpened: A pointer to the newly created error.
func NewErrNeverOpened(openingToken, closingToken string) *ErrNeverOpened {
	e := &ErrNeverOpened{
		OpeningToken: openingToken,
		ClosingToken: closingToken,
	}
	return e
}
