package cimap_test

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/projectbarks/cimap"
)

// randomString generates a random lowercase string of length n.
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, n)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

func generateRandomKeys(num, min, max int) []string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	out := make([]string, num)

	for i := 0; i < num; i++ {
		out[i] = randomString(r.Intn(max-min)+min) + strconv.Itoa(i)
	}
	return out
}

// ---------------------------------------------------------------------
// Benchmarks
// ---------------------------------------------------------------------

func BenchmarkAdd(b *testing.B) {
	b.StopTimer()
	keys := generateRandomKeys(b.N*2, 5, 20)

	b.Run("Base", func(b *testing.B) {
		m := &InsenstiveStubMap[string]{keys: make(map[string]string, b.N)}

		b.ReportAllocs()
		b.StartTimer()
		defer b.StopTimer()

		for i := 0; i < b.N; i++ {
			m.Add(keys[i%len(keys)], "some-value")
		}
	})

	b.Run("CIMap", func(b *testing.B) {
		cm := cimap.New[string]()
		b.ReportAllocs()
		b.StartTimer()
		defer b.StopTimer()

		for i := 0; i < b.N; i++ {
			cm.Add(keys[i%len(keys)], "some-value")
		}
	})
}

func BenchmarkGet(b *testing.B) {
	const numKeys = 100000
	keys := generateRandomKeys(numKeys, 5, 20)

	// Pre-fill both maps with all keys
	mBase := &InsenstiveStubMap[string]{keys: make(map[string]string, numKeys)}
	cm := cimap.New[string](numKeys)
	for _, k := range keys {
		mBase.Add(k, "some-value")
		cm.Add(k, "some-value")
	}

	b.ResetTimer()

	b.Run("Base", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = mBase.Get(keys[i%numKeys])
		}
	})

	b.Run("CIMap", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = cm.Get(keys[i%numKeys])
		}
	})
}

// ---------------------------------------------------------------------
// Benchmark: Delete
// ---------------------------------------------------------------------

func BenchmarkDelete(b *testing.B) {
	const numKeys = 100000
	keys := generateRandomKeys(numKeys, 5, 10)

	b.Run("Base", func(b *testing.B) {
		m := &InsenstiveStubMap[string]{keys: make(map[string]string, numKeys)}
		for _, k := range keys {
			m.Add(k, "some-value")
		}

		b.StartTimer()
		b.ReportAllocs()
		defer b.StopTimer()
		for i := 0; i < b.N; i++ {
			m.Delete(keys[i%numKeys])
		}
	})

	b.Run("CIMap", func(b *testing.B) {
		cm := cimap.New[string](numKeys)
		for _, k := range keys {
			cm.Add(k, "some-value")
		}

		b.StartTimer()
		b.ReportAllocs()
		defer b.StopTimer()
		for i := 0; i < b.N; i++ {
			cm.Delete(keys[i%numKeys])
		}
	})
}
