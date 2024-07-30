package maps

import (
	"cmp"
	"slices"

	luc "github.com/PlayerR9/lib_units/common"
)

// OrderedMap is a map that is ordered by the keys.
type OrderedMap[K cmp.Ordered, V any] struct {
	// values is a map of the values in the map.
	values map[K]V

	// keys is a slice of the keys in the map.
	keys []K
}

// Iterator implements the common.Iterable interface.
//
// Never returns nil.
func (m *OrderedMap[K, V]) Iterator() luc.Iterater[*Entry[K, V]] {
	return &OMIterator[K, V]{
		m:   m,
		pos: 0,
	}
}

// NewOrderedMap creates a new OrderedMap.
//
// Returns:
//   - *OrderedMap: A pointer to the newly created OrderedMap.
//     Never returns nil.
func NewOrderedMap[K cmp.Ordered, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		values: make(map[K]V),
		keys:   make([]K, 0),
	}
}

// Add adds a key-value pair to the map.
//
// Parameters:
//   - key: The key to add.
//   - value: The value to add.
//   - force: If true, the value will be added even if the key already exists. If
//     false, the value will not be added if the key already exists.
//
// Returns:
//   - bool: True if the value was added to the map, false otherwise.
func (m *OrderedMap[K, V]) Add(key K, value V, force bool) bool {
	pos, ok := slices.BinarySearch(m.keys, key)

	if !ok {
		m.keys = slices.Insert(m.keys, pos, key)
	}

	if ok && !force {
		return false
	}

	m.values[key] = value

	return true
}

// KeyIterator is a method that returns an iterator over the keys in the map.
//
// Returns:
//   - common.Iterater[K]: An iterator over the keys in the map. Never returns nil.
func (m *OrderedMap[K, V]) KeyIterator() luc.Iterater[K] {
	return luc.NewSimpleIterator(m.keys)
}

// Size is a method that returns the number of keys in the map.
//
// Returns:
//   - int: The number of keys in the map.
func (m *OrderedMap[K, V]) Size() int {
	return len(m.keys)
}

// GetMap is a method that returns the map of the values in the map.
//
// Returns:
//   - map[K]V: The map of the values in the map. Never returns nil.
func (m *OrderedMap[K, V]) GetMap() map[K]V {
	return m.values
}
