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

	// [CaseInsensitiveMap] is a generic map that performs case-insensitive key comparisons.
	//
	// It uses a customizable hash function to store keys in an internal map,
	// handling collisions via separate chaining.
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

// New creates and returns a new [CaseInsensitiveMap] instance.
//
// An optional positive integer can be provided to preallocate the internal map with the given capacity.
//
//	m := cimap.New[int](10)
//	fmt.Println(m.Len()) // Output: 0
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

// Add inserts or updates the key-value pair in the map.
//
// The key comparison is case-insensitive, so if a key differing only by case exists,
// its value will be replaced with the new one.
//
//	m := cimap.New[string]()
//	m.Add("Hello", "World")
//	m.Add("hello", "Gophers")
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

// Get retrieves the value associated with the specified key using a case-insensitive comparison.
//
// It returns the value and a boolean indicating whether the key was found.
//
// If the key is not present, the zero value of T and false are returned.
//
//	m := cimap.New[int]()
//	m.Add("Key", 42)
//	value, ok := m.Get("key") // Output: 42 true
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

// GetAndDel retrieves the value associated with the specified key and then removes the key-value pair from the map.
//
// It returns the value and a boolean indicating whether the key was found.
//
// If the key does not exist, the zero value of T and false are returned.
//
//	m := cimap.New[string]()
//	m.Add("temp", "data")
//	value, ok := m.GetAndDel("temp") // Output: "data" true
//	value, ok = m.Get("temp") // Output: false
func (c *CaseInsensitiveMap[T]) GetAndDel(k string) (T, bool) {
	// TODO: add more performant version
	if v, ok := c.Get(k); ok {
		c.Delete(k)
		return v, true
	}
	var def T
	return def, false
}

// GetOrSet retrieves the value associated with the specified key.
//
// If the key is not present, it sets the value to the provided value and returns it.
// This ensures that the key exists in the map after the call.
//
//	m := cimap.New[int]()
//	v1 := m.GetOrSet("count", 100) // Output: 100
//	v2 := m.GetOrSet("COUNT", 200) // Output: 100
func (c *CaseInsensitiveMap[T]) GetOrSet(k string, val T) T {
	// TODO: add more performant version
	if v, ok := c.Get(k); ok {
		return v
	}
	c.Add(k, val)
	return val
}

// Delete removes the key-value pair associated with the specified key from the map.
//
// The key comparison is performed in a case-insensitive manner.
//
//	m := cimap.New[int]()
//	m.Add("delete", 123)
//	m.Delete("DELETE")
//	m.Get("delete") // Output: false
func (c *CaseInsensitiveMap[T]) Delete(k string) {
	if n, ok := c.internalMap[c.hashString(k)]; ok && n.delete(k) {
		delete(c.internalMap, c.hashString(k))
		c.size--
	}
}

// Len returns the number of key-value pairs currently stored in the map.
//
//	m := cimap.New[int]()
//	m.Add("a", 1)
//	m.Add("A", 2)
//	m.Len() // Output: 1
func (c *CaseInsensitiveMap[T]) Len() int {
	return c.size
}

// Clear removes all key-value pairs from the map, resetting it to an empty state.
//
//	m := cimap.New[string]()
//	m.Add("x", "y")
//	m.Clear()
//	m.Len() // Output: 0
func (c *CaseInsensitiveMap[T]) Clear() {
	c.internalMap = make(map[hash64]*node[T])
	c.size = 0
}

// Keys returns an iterator over all keys stored in the map.
// The iteration order is unspecified.
//
//	m := cimap.New[int]()
//	m.Add("One", 1)
//	m.Add("Two", 2)
//	m.Keys()(func(key string) bool {
//	    fmt.Println(key) // Output: One Two
//	    return true
//	})
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
// The order of iteration is not guaranteed.
//
//	m := cimap.New[string]()
//	m.Add("first", "a")
//	m.Add("second", "b")
//	m.Iterator()(func(key string, value string) bool {
//	    fmt.Printf("%s: %s\n", key, value) // Output: first: a second: b
//	    return true
//	})
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

// SetHasher sets a custom hash function for computing keys in the map.
//
// The provided hash function is used for all subsequent operations,
// and the map is rehashed immediately to reflect the new hashing strategy.
//
// WARNING(a1): Don't use this unless you know what you are doing. This function
// can destroy the performance of this module if not used correctly.
//
//	customHasher := func(s string) uint64 {
//	    return uint64(len(s))
//	}
//	m.SetHasher(customHasher)
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

// ForEach executes the provided function for each key-value pair in the map.
//
// Iteration stops early if the function returns false.
// The order of iteration is undefined.
//
//	m.ForEach(func(key string, value int) bool {
//	    fmt.Printf("%s: %d\n", key, value)
//	    return true
//	})
func (c *CaseInsensitiveMap[T]) ForEach(fn func(string, T) bool) {
	for _, v := range c.internalMap {
		for ; v != nil; v = v.Next {
			if !fn(v.Key, v.Value) {
				return
			}
		}
	}
}

// UnmarshalJSON implements the [json.Unmarshaler] interface.
//
// It decodes JSON data into the map using case-insensitive key handling.
// Any existing data in the map is cleared before unmarshalling.
//
//	data := []byte(`{"Foo": 10, "bar": 20}`)
//	var m cimap.CaseInsensitiveMap[int]
//	if err := json.Unmarshal(data, &m); err != nil {
//	    log.Fatal(err)
//	}
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

// MarshalJSON implements the [json.Marshaler] interface.
//
// It encodes the map into JSON format, preserving the original casing of keys.
//
//   data, err := json.Marshal(m) // Output: {"Key":123}

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
// It manually converts to lowercase
// avoiding any allocation.
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
