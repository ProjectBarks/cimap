package cimap

import (
	"encoding/json"
	"iter"
	"strings"
	"unicode"
)

type (
	hash64 = uint64

	node[T any] struct {
		Value T
		Key   string
		Next  *node[T]
	}

	CaseInsensitiveMap[T any] struct {
		size        int
		hashString  func(string) hash64
		internalMap map[hash64]*node[T]
	}
)

const (
	offset64 = hash64(14695981039346656037)
	prime64  = hash64(1099511628211)
)

// New creates a new case-insensitive map.
func New[T any](size ...int) *CaseInsensitiveMap[T] {
	if len(size) > 0 && size[0] > 0 {
		return &CaseInsensitiveMap[T]{
			internalMap: make(map[hash64]*node[T], size[0]),
			hashString:  defaultHashString,
		}
	}
	return &CaseInsensitiveMap[T]{
		internalMap: make(map[hash64]*node[T]),
		hashString:  defaultHashString,
	}
}

// Add adds a key-value pair to the map.
func (c *CaseInsensitiveMap[T]) Add(k string, val T) {
	if n, ok := c.internalMap[c.hashString(k)]; ok {
		if !n.insertOrReplace(k, val) {
			c.size++
		}
	} else {
		newNode := node[T]{Value: val, Key: k}
		c.internalMap[c.hashString(k)] = &newNode
		c.size++
	}
}

// Get returns the value associated with the key.
// If the key is not found, it returns the zero value of T.
func (c *CaseInsensitiveMap[T]) Get(k string) (T, bool) {
	for n := c.internalMap[c.hashString(k)]; n != nil; n = n.Next {
		if !strings.EqualFold(n.Key, k) {
			continue
		}
		return n.Value, true
	}

	var def T
	return def, false
}

// Get returns the value associated with the key.
// If the key is not found, it returns the zero value of T.
func (c *CaseInsensitiveMap[T]) GetAndDel(k string) (T, bool) {
	// TODO: add more performant version
	if v, ok := c.Get(k); ok {
		c.Delete(k)
		return v, true
	}
	var def T
	return def, false
}

// GetOrSet returns the value associated with the key.
// If the key is not found, it sets the value to the zero value of T and returns it.
func (c *CaseInsensitiveMap[T]) GetOrSet(k string, val T) T {
	// TODO: add more performant version
	if v, ok := c.Get(k); ok {
		return v
	}
	c.Add(k, val)
	return val
}

// Delete deletes a key-value pair from the map.
func (c *CaseInsensitiveMap[T]) Delete(k string) {
	if n, ok := c.internalMap[c.hashString(k)]; ok && n.delete(k) {
		delete(c.internalMap, c.hashString(k))
		c.size--
	}
}

// Len returns the number of key-value pairs in the map.
func (c *CaseInsensitiveMap[T]) Len() int {
	return c.size
}

// Clear clears the map.
func (c *CaseInsensitiveMap[T]) Clear() {
	c.internalMap = make(map[hash64]*node[T])
	c.size = 0
}

// Keys returns an iterator over all keys in the map.
func (c *CaseInsensitiveMap[T]) Keys() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, v := range c.internalMap {
			for ; v != nil; v = v.Next {
				if !yield(v.Key) {
					return
				}
			}
		}
	}
}

// Iterator returns an iterator over all key-value pairs in the map.
func (c *CaseInsensitiveMap[T]) Iterator() iter.Seq2[string, T] {
	return func(yield func(string, T) bool) {
		for _, v := range c.internalMap {
			for ; v != nil; v = v.Next {
				if !yield(v.Key, v.Value) {
					return
				}
			}
		}
	}
}

// SetHasher sets the hash function for the map.
func (c *CaseInsensitiveMap[T]) SetHasher(hashString func(string) hash64) {
	c.hashString = hashString
	// we need to rehash the map
	if c.size > 0 {
		newMap := make(map[hash64]*node[T], c.size)
		for _, v := range c.internalMap {
			newMap[hashString(v.Key)] = v
		}
		c.internalMap = newMap
	}
}

// ForEach iterates over all key-value pairs in the map and calls the given function.
//
// If the function returns false, the iteration stops.
func (c *CaseInsensitiveMap[T]) ForEach(fn func(string, T) bool) {
	for _, v := range c.internalMap {
		for ; v != nil; v = v.Next {
			if !fn(v.Key, v.Value) {
				return
			}
		}
	}
}

func (c *CaseInsensitiveMap[T]) UnmarshalJSON(data []byte) error {
	var m map[string]T
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	c.internalMap = make(map[hash64]*node[T], len(m))
	c.size = 0 // it's 0 since we are going to remove elements by cases collision
	if c.hashString == nil {
		c.hashString = defaultHashString
	}
	for k, v := range m {
		c.Add(k, v)
	}
	return nil
}

func (c *CaseInsensitiveMap[T]) MarshalJSON() ([]byte, error) {
	m := make(map[string]T, c.size)
	for _, v := range c.internalMap {
		for ; v != nil; v = v.Next {
			m[v.Key] = v.Value
		}
	}

	return json.Marshal(m)
}

////////////////////////////////////////////////////////////
// HASH METHODS
////////////////////////////////////////////////////////////

// hashString computes the FNV-1a hash for s.
// It manually converts uppercase ASCII letters to lowercase
// on a per-byte basis, avoiding any allocation.
func defaultHashString(key string) hash64 {
	h := offset64
	for _, r := range key {
		h *= prime64
		h ^= uint64(unicode.ToLower(r))
	}
	return h
}

////////////////////////////////////////////////////////////
// NODE METHODS
////////////////////////////////////////////////////////////

func (n *node[T]) delete(key string) bool {
	if strings.EqualFold(n.Key, key) {
		n = n.Next
		return true
	}
	for prev := n; prev.Next != nil; prev = prev.Next {
		if strings.EqualFold(prev.Next.Key, key) {
			prev.Next = prev.Next.Next
			return true
		}
	}
	return false
}

// make a node function called insert or replace which uses key to insert or replace a node
// loop through the linked list and if the key exists, replace the node
// if the key does not exist, insert a new node
//
// return true if the node existed
func (n *node[T]) insertOrReplace(key string, val T) bool {
	var prev *node[T] = nil
	for cur := n; cur != nil; prev, cur = cur, cur.Next {
		if !strings.EqualFold(cur.Key, key) {
			continue
		}
		cur.Key = key
		cur.Value = val
		return true
	}
	prev.Next = &node[T]{Key: key, Value: val}
	return false
}
