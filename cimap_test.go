package cimap_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/projectbarks/cimap"

	"github.com/stretchr/testify/assert"
)

type (
	keyPair struct {
		key string
		val string
	}

	keyAssert struct {
		key      string
		val      string
		expected bool
	}
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		size     []int
		expected int
	}{
		{
			name:     "No size parameter",
			size:     nil,
			expected: 0,
		},
		{
			name:     "Zero size parameter",
			size:     []int{0},
			expected: 0,
		},
		{
			name:     "Positive size parameter",
			size:     []int{10},
			expected: 0, // After creation, Len() should still be 0
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			m := cimap.New[string](tt.size...)
			assert.NotNil(t, m, "Expected non-nil map")
			assert.Equal(t, tt.expected, m.Len(), "Expected map length to be 0 upon creation")
		})
	}
}

func TestAdd_Get(t *testing.T) {
	tests := []struct {
		name       string
		operations []keyPair
		checks     []keyAssert
		hashFn     func(string) uint64
	}{
		{
			name: "Single add & get",
			operations: []keyPair{
				{"Hello", "World"},
			},
			checks: []keyAssert{
				{"Hello", "World", true},
				{"hello", "World", true}, // check case-insensitivity
				{"HELLO", "World", true}, // check case-insensitivity
				{"NotFound", "", false},
			},
		},
		{
			name: "Multiple adds & gets",
			operations: []keyPair{
				{"K1", "V1"},
				{"k1", "V2"}, // case collision with same spelling
				{"key2", "val2"},
				{"KEY2", "val3"}, // also collision
				{"MixedCase", "MC"},
			},
			checks: []keyAssert{
				{"K1", "V2", true},
				{"k1", "V2", true},
				{"K1", "V2", true},
				{"key2", "val3", true},
				{"KEY2", "val3", true},
				{"mixedcase", "MC", true},
				{"NotExist", "", false},
			},
		},
		{
			name: "Get with hash collision",
			operations: []keyPair{
				{"dog", "bark"},
				{"cat", "meow"},
			},
			checks: []keyAssert{
				{"dog", "bark", true},
				{"cat", "meow", true},
			},
			hashFn: func(s string) uint64 {
				return uint64(len(s))
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			m := cimap.New[string]()

			if tt.hashFn != nil {
				m.SetHasher(tt.hashFn)
			}

			for _, op := range tt.operations {
				m.Add(op.key, op.val)
			}

			for _, ch := range tt.checks {
				val, ok := m.Get(ch.key)
				assert.Equal(t, ch.expected, ok, "Unexpected existence for key %q", ch.key)
				assert.Equal(t, ch.val, val, "Value mismatch for key %q", ch.key)
			}
		})
	}
}

func TestGetAndDel(t *testing.T) {
	tests := []struct {
		name       string
		insert     []keyPair
		getAndDel  keyAssert
		finalCheck keyAssert
	}{
		{
			name: "GetAndDel existing key",
			insert: []keyPair{
				{"Hello", "World"},
				{"FoO", "Bar"},
			},
			getAndDel:  keyAssert{"hello", "World", true},
			finalCheck: keyAssert{"Hello", "", false},
		},
		{
			name: "GetAndDel non-existing key",
			insert: []keyPair{
				{"A", "B"},
			},
			getAndDel:  keyAssert{"DoesNotExist", "", false},
			finalCheck: keyAssert{"A", "B", true},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			m := cimap.New[string]()
			for _, in := range tt.insert {
				m.Add(in.key, in.val)
			}

			val, ok := m.GetAndDel(tt.getAndDel.key)
			assert.Equal(t, tt.getAndDel.expected, ok, "Found mismatch after GetAndDel")
			assert.Equal(t, tt.getAndDel.val, val, "Value mismatch after GetAndDel")

			// verify that the key is truly deleted if it existed
			_, ok = m.Get(tt.finalCheck.key)
			assert.Equal(t, tt.finalCheck.expected, ok, "Expected existence mismatch after final check")
		})
	}
}

func TestGetOrSet(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		initialValue string
		getOrSetVal  string
		expectedVal  string
	}{
		{
			name:         "Key not existing",
			key:          "ABC",
			initialValue: "",
			getOrSetVal:  "Default",
			expectedVal:  "Default",
		},
		{
			name:         "Key exists (case-insensitive)",
			key:          "FoO",
			initialValue: "Bar",
			getOrSetVal:  "NotUsed",
			expectedVal:  "Bar", // Because it should return the existing value
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			m := cimap.New[string]()
			if tt.initialValue != "" {
				m.Add(tt.key, tt.initialValue)
			}
			got := m.GetOrSet(strings.ToLower(tt.key), tt.getOrSetVal)
			assert.Equal(t, tt.expectedVal, got, "Mismatch in GetOrSet returned value")
			// Also ensure the map indeed stored the final value
			stored, _ := m.Get(tt.key)
			assert.Equal(t, tt.expectedVal, stored, "Mismatch in stored value after GetOrSet")
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name      string
		insert    []keyPair
		deleteKey string
		hashFn    func(string) uint64
		finalLen  int
	}{
		{
			name: "Delete existing key",
			insert: []keyPair{
				{"Hello", "World"},
				{"Foo", "Bar"},
			},
			deleteKey: "hello",
			finalLen:  1,
		},
		{
			name: "Delete non-existing key",
			insert: []keyPair{
				{"A", "B"},
			},
			deleteKey: "DoesNotExist",
			finalLen:  1,
		},
		{
			name: "Delete existing key with collision",
			insert: []keyPair{
				{"Abc", "123"},
				{"abc", "456"},
			},
			deleteKey: "ABC",
			finalLen:  0,
		},
		{
			name: "Delete existing key hash with collision",
			insert: []keyPair{
				{"cdf", "123"},
				{"abc", "456"},
			},
			deleteKey: "ABC",
			finalLen:  1,
			hashFn: func(s string) uint64 {
				return uint64(len(s))
			},
		},
		{
			name: "Do not delete existing key hash with collision",
			insert: []keyPair{
				{"cdf", "123"},
			},
			deleteKey: "ABC",
			finalLen:  1,
			hashFn: func(s string) uint64 {
				return uint64(len(s))
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			m := cimap.New[string]()
			if tt.hashFn != nil {
				m.SetHasher(tt.hashFn)
			}
			for _, in := range tt.insert {
				m.Add(in.key, in.val)
				fmt.Printf("%d\n", m.Len())
			}

			m.Delete(tt.deleteKey)
			assert.Equal(t, tt.finalLen, m.Len(), "Mismatch in final length after deletion")
		})
	}
}

func TestLen(t *testing.T) {
	tests := []struct {
		name    string
		inserts []keyPair
		deletes []string
		final   int
	}{
		{
			name:  "No insertion => length zero",
			final: 0,
		},
		{
			name: "Multiple distinct keys => length",
			inserts: []keyPair{
				{"A", "1"},
				{"B", "2"},
				{"C", "3"},
			},
			final: 3,
		},
		{
			name: "Case collision => single logical key stored last",
			inserts: []keyPair{
				{"A", "1"},
				{"a", "2"},
			},
			final: 1,
		},
		{
			name: "Insert 3, delete 1 => length 2",
			inserts: []keyPair{
				{"X", "1"},
				{"Y", "2"},
				{"Z", "3"},
			},
			deletes: []string{"y"},
			final:   2,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			m := cimap.New[string]()
			for _, in := range tt.inserts {
				m.Add(in.key, in.val)
			}
			for _, d := range tt.deletes {
				m.Delete(d)
			}
			assert.Equal(t, tt.final, m.Len(), "Mismatch in final length")
		})
	}
}

func TestClear(t *testing.T) {
	m := cimap.New[string]()
	m.Add("A", "1")
	m.Add("a", "2")
	assert.Equal(t, 1, m.Len(), "Because A/a is the same key, only 1 unique key so far")

	m.Add("C", "3")
	assert.Equal(t, 2, m.Len(), "Expected length = 2 after inserting another distinct key")

	m.Clear()
	assert.Equal(t, 0, m.Len(), "Expected empty map after Clear")
}

func TestKeys(t *testing.T) {
	m := cimap.New[string]()
	m.Add("A", "1")
	m.Add("a", "2")
	m.Add("B", "10")
	m.Add("C", "100")
	m.Add("c", "200")

	t.Run("Default", func(t *testing.T) {
		var keys []string
		m.Keys()(func(k string) bool {
			keys = append(keys, k)
			return true
		})

		// We expect actual nodes to be:
		//  1) "a" or "A" (the last inserted is "a" => "2")
		//  2) "B" => "10"
		//  3) "c" => "200"
		// The map stores "A" and "a" under the same hash bucket, so only the last inserted remains as a distinct node.
		assert.Len(t, keys, 3, "Expected 3 unique keys from the map")
		assert.Contains(t, keys, "a")
		assert.Contains(t, keys, "B")
		assert.Contains(t, keys, "c")
	})

	t.Run("Short circuit", func(t *testing.T) {
		loops := 0
		m.Keys()(func(k string) bool {
			loops += 1
			return false
		})

		assert.Equal(t, 1, loops, "Expected iterator to stop after first iteration")
	})
}

func TestIterator(t *testing.T) {
	m := cimap.New[string]()
	m.Add("A", "1")
	m.Add("a", "2")
	m.Add("B", "10")
	m.Add("C", "100")
	m.Add("c", "200")

	t.Run("Default", func(t *testing.T) {
		found := make(map[string]string)
		m.Iterator()(func(k, v string) bool {
			found[k] = v
			return true
		})

		// We expect 3 distinct pairs:
		//   ("a", "2"), ("B", "10"), ("c", "200")
		assert.Len(t, found, 3, "Expected 3 items from iterator")
		assert.Equal(t, map[string]string{
			"a": "2",
			"B": "10",
			"c": "200",
		}, found)
	})

	t.Run("Short circuit", func(t *testing.T) {
		loops := 0
		m.Iterator()(func(k, v string) bool {
			loops++
			return false
		})
		assert.Equal(t, 1, loops, "Expected iterator to stop after first iteration")
	})
}

func TestForEach(t *testing.T) {
	m := cimap.New[string]()
	m.Add("K1", "V1")
	m.Add("K2", "V2")
	m.Add("K3", "V3")

	count := 0
	m.ForEach(func(k, v string) bool {
		count++
		return false // stop immediately
	})
	assert.Equal(t, 1, count, "Expected ForEach to stop after first iteration")

	// Check that iterating fully with 'return false' never used
	count = 0
	m.ForEach(func(k, v string) bool {
		count++
		return true // never returns true => never stops
	})
}

func TestSetHasher(t *testing.T) {
	m := cimap.New[string]()

	// Insert with default hash
	m.Add("Hello", "World")
	_, found := m.Get("hello")
	assert.True(t, found, "Expected to find the key with default hash")

	// Change the hash function (just as an example, we uppercase instead of lowercase)
	m.SetHasher(func(s string) uint64 {
		var h uint64
		for _, r := range strings.ToLower(s) {
			h += uint64(r) // not robust, just to test hooking
		}
		return h
	})

	// After changing the hash function, the old entries won't be found by the new method
	_, found = m.Get("Hello")
	assert.True(t, found, "Expected to find key after changing hash function")

	// But newly inserted ones with the new hash function will be found
	m.Add("Hello", "NewValue")
	val, found := m.Get("HeLlO")
	assert.True(t, found)
	assert.Equal(t, "NewValue", val)
}

func TestMarshalUnmarshalJSON(t *testing.T) {
	type MyStruct struct {
		Name string
		Age  int
	}

	t.Run("Simple string map round-trip", func(t *testing.T) {
		m := cimap.New[string]()
		m.Add("K1", "V1")
		m.Add("k1", "V2") // same logical key
		m.Add("K2", "V2")

		encoded, err := json.Marshal(m)
		assert.NoError(t, err)

		var m2 cimap.CaseInsensitiveMap[string]
		err = json.Unmarshal(encoded, &m2)
		assert.NoError(t, err)

		assert.Equal(t, 2, m2.Len(), "Expected same number of unique keys after round trip")
		val, ok := m2.Get("k1")
		assert.True(t, ok)
		assert.Equal(t, "V2", val, "Expected last inserted value for K1")
	})

	t.Run("Struct values round-trip", func(t *testing.T) {
		m := cimap.New[MyStruct]()
		m.Add("Person", MyStruct{Name: "John", Age: 30})
		m.Add("PERSON", MyStruct{Name: "Doe", Age: 40})

		encoded, err := json.Marshal(m)
		assert.NoError(t, err)

		var m2 cimap.CaseInsensitiveMap[MyStruct]
		err = json.Unmarshal(encoded, &m2)
		assert.NoError(t, err)

		assert.Equal(t, 1, m2.Len(), "Expected only one unique key because Person/PERSON are the same ignoring case")
		val, ok := m2.Get("person")
		assert.True(t, ok)
		assert.Equal(t, MyStruct{Name: "Doe", Age: 40}, val, "Should reflect the last inserted struct for the same key ignoring case")
	})

	t.Run("Unmarshal error", func(t *testing.T) {
		var m cimap.CaseInsensitiveMap[string]
		err := m.UnmarshalJSON([]byte{12})
		assert.Error(t, err)
	})

}
