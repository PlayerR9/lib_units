package maps

import (
	"cmp"
	"strconv"

	luc "github.com/PlayerR9/lib_units/common"
	lustr "github.com/PlayerR9/lib_units/strings"
)

// Entry is a key-value pair in an OrderedMap.
type Entry[K cmp.Ordered, V any] struct {
	// Key is the key of the entry.
	Key K

	// Value is the value of the entry.
	Value V
}

// OMIterator is an iterator for an OrderedMap.
type OMIterator[K cmp.Ordered, V any] struct {
	// m is the map to iterate over.
	m *OrderedMap[K, V]

	// pos is the current position in the iterator.
	pos int
}

// Consume implements the common.Iterater interface.
func (i *OMIterator[K, V]) Consume() (*Entry[K, V], error) {
	luc.AssertNil(i.m, "i.m")

	if i.pos >= len(i.m.keys) {
		return nil, luc.NewErrExhaustedIter()
	}

	key := i.m.keys[i.pos]
	i.pos++

	val, ok := i.m.values[key]
	luc.AssertOk(ok, "i.m.values[%s]", strconv.Quote(lustr.GoStringOf(key)))

	return &Entry[K, V]{
		Key:   key,
		Value: val,
	}, nil
}

// Restart implements the common.Iterater interface.
func (i *OMIterator[K, V]) Restart() {
	i.pos = 0
}

// NewOMIterator creates a new OMIterator.
//
// Parameters:
//   - m: The map to iterate over.
//
// Returns:
//   - *OMIterator: A pointer to the newly created OMIterator. Nil if m is nil.
func NewOMIterator[K cmp.Ordered, V any](m *OrderedMap[K, V]) *OMIterator[K, V] {
	if m == nil {
		return nil
	}

	return &OMIterator[K, V]{
		m:   m,
		pos: 0,
	}
}
