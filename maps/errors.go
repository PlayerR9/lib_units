package maps

import (
	"cmp"
	"strconv"
	"strings"

	lustr "github.com/PlayerR9/lib_units/strings"
)

// ErrKeyNotFound is an error that is returned when a key is not found in the
// map.
type ErrKeyNotFound[K cmp.Ordered] struct {
	// Key is the key that was not found.
	Key K
}

// Error implements the error interface.
//
// Message: "key ("{{ .Key }}") not found"
func (e *ErrKeyNotFound[K]) Error() string {
	var builder strings.Builder

	builder.WriteString("key (")
	builder.WriteString(strconv.Quote(lustr.GoStringOf(e.Key)))
	builder.WriteString(") not found")

	return builder.String()
}

// NewErrKeyNotFound creates a new ErrKeyNotFound error.
//
// Parameters:
//   - key: The key that was not found.
//
// Returns:
//   - *ErrKeyNotFound: A pointer to the newly created ErrKeyNotFound.
//     Never returns nil.
func NewErrKeyNotFound[K cmp.Ordered](key K) *ErrKeyNotFound[K] {
	return &ErrKeyNotFound[K]{
		Key: key,
	}
}
