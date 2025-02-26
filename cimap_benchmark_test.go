package cimap_test

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/projectbarks/cimap"
)

const (
	_UPPERCASE = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	_LOWERCASE = "abcdefghijklmnopqrstuvwxyz"
	_UNICODE   = "¡¢£¤¥¦§¨©ª«¬®¯°±²³´µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿ"
)

type keyGroup struct {
	name string
	base string
	keys []string
}

func randomString(n int, base string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, n)
	for i := range b {
		b[i] = base[r.Intn(len(base))]
	}
	return string(b)
}

func generateKeyGroups(num, min, max int) []keyGroup {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	result := []keyGroup{
		{name: "ASCII Mismatch", base: _UPPERCASE},
		{name: "ASCII Match", base: _LOWERCASE},
		{name: "Unicode", base: _UNICODE},
	}

	for i, group := range result {
		for range num {
			key := randomString(r.Intn(max-min)+min, group.base) + strconv.Itoa(i)
			group.keys = append(group.keys, key)
		}
		result[i] = group
	}

	return result
}

// ---------------------------------------------------------------------
// Benchmarks
// ---------------------------------------------------------------------

func BenchmarkAdd(b *testing.B) {
	b.StopTimer()
	groups := generateKeyGroups(b.N*2, 5, 50)

	for _, group := range groups {
		b.Run(group.name, func(b *testing.B) {
			b.Run("Base", func(b *testing.B) {
				m := &InsenstiveStubMap[string]{keys: make(map[string]string, b.N)}

				b.ReportAllocs()
				b.StartTimer()
				defer b.StopTimer()

				for i := 0; i < b.N; i++ {
					m.Add(group.keys[i%len(group.keys)], "some-value")
				}
			})

			b.Run("CIMap", func(b *testing.B) {
				cm := cimap.New[string]()
				b.ReportAllocs()
				b.StartTimer()
				defer b.StopTimer()

				for i := 0; i < b.N; i++ {
					cm.Add(group.keys[i%len(group.keys)], "some-value")
				}
			})
		})
	}
}

func BenchmarkGet(b *testing.B) {
	const numKeys = 100000
	groups := generateKeyGroups(numKeys, 5, 50)

	for _, group := range groups {
		b.Run(group.name, func(b *testing.B) {
			mBase := &InsenstiveStubMap[string]{keys: make(map[string]string, numKeys)}
			cm := cimap.New[string](numKeys)
			for _, k := range group.keys {
				mBase.Add(k, "some-value")
				cm.Add(k, "some-value")
			}

			b.ResetTimer()

			b.Run("Base", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					_, _ = mBase.Get(group.keys[i%numKeys])
				}
			})

			b.Run("CIMap", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					_, _ = cm.Get(group.keys[i%numKeys])
				}
			})
		})
	}
}

// ---------------------------------------------------------------------
// Benchmark: Delete
// ---------------------------------------------------------------------

func BenchmarkDelete(b *testing.B) {
	const numKeys = 100000
	groups := generateKeyGroups(numKeys, 5, 50)

	for _, group := range groups {

		b.Run(group.name, func(b *testing.B) {
			b.Run("Base", func(b *testing.B) {
				m := &InsenstiveStubMap[string]{keys: make(map[string]string, numKeys)}
				for _, k := range group.keys {
					m.Add(k, "some-value")
				}

				b.StartTimer()
				b.ReportAllocs()
				defer b.StopTimer()
				for i := 0; i < b.N; i++ {
					m.Delete(group.keys[i%numKeys])
				}
			})

			b.Run("CIMap", func(b *testing.B) {
				cm := cimap.New[string](numKeys)
				for _, k := range group.keys {
					cm.Add(k, "some-value")
				}

				b.StartTimer()
				b.ReportAllocs()
				defer b.StopTimer()
				for i := 0; i < b.N; i++ {
					cm.Delete(group.keys[i%numKeys])
				}
			})
		})
	}
}
