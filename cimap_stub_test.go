package cimap_test

import "strings"

// Implement a stub for the cimap package using strings.ToLower to hash keys as baseline
// Implement get, add, delete

type InsenstiveStubMap[T any] struct {
	keys map[string]T
}

func (m *InsenstiveStubMap[T]) Get(key string) (T, bool) {
	lowerKey := strings.ToLower(key)
	v, ok := m.keys[lowerKey]
	return v, ok
}

func (m *InsenstiveStubMap[T]) Add(key string, value T) {
	m.keys[strings.ToLower(key)] = value
}

func (m *InsenstiveStubMap[T]) Delete(key string) {
	delete(m.keys, strings.ToLower(key))
}
