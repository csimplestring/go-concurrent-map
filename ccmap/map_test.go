package ccmap

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/csimplestring/go-concurrent-map/algo/random"
)

var (
	benchmarkKeys []Key
)

func init() {
	benchmarkKeys = make([]Key, 10000)
	for i := 0; i < 10000; i++ {
		benchmarkKeys[i] = NewStringKey(random.NewLen(15))
	}
}

func TestMapPut(t *testing.T) {
	m := NewMap()

	for i := 0; i < 10; i++ {
		key := NewStringKey(fmt.Sprintf("%d", i))
		m.Put(key, i)
	}
}

func TestMapGet(t *testing.T) {
	m := NewMap()

	for i := 0; i < 100000; i++ {
		key := NewStringKey(fmt.Sprintf("%s", i))
		m.Put(key, i)
	}

	for i := 0; i < 100000; i++ {
		assert.Equal(t, i, m.Get(NewStringKey(fmt.Sprintf("%s", i))))
	}
}

func BenchmarkMapPut(b *testing.B) {
	m := NewMap()

	for i, k := range benchmarkKeys {
		m.Put(k, i)
	}
}

func BenchmarkMapGet(b *testing.B) {
	m := NewMap()

	size := len(benchmarkKeys)
	for i := 0; i < size/2; i++ {
		m.Put(benchmarkKeys[i], i)
	}
	b.StopTimer()
	b.StartTimer()

	for _, k := range benchmarkKeys {
		m.Get(k)
	}
}

func BenchmarkMapDelete(b *testing.B) {
	m := NewMap()

	size := len(benchmarkKeys)
	for i := 0; i < size/2; i++ {
		m.Put(benchmarkKeys[i], i)
	}
	b.StopTimer()
	b.StartTimer()

	for _, k := range benchmarkKeys {
		m.Delete(k)
	}
}

func BenchmarkStandardMapPut(b *testing.B) {
	m := make(map[string]interface{}, 16)

	for i, k := range benchmarkKeys {
		m[k.String()] = i
	}
}

func BenchmarkStandardMapGet(b *testing.B) {
	m := make(map[string]interface{}, 16)

	size := len(benchmarkKeys)
	for i := 0; i < size/2; i++ {
		m[benchmarkKeys[i].String()] = i
	}

	b.StopTimer()
	b.StartTimer()

	for _, k := range benchmarkKeys {
		_ = m[k.String()]
	}
}

func BenchmarkStandardMapDelete(b *testing.B) {
	m := make(map[string]interface{}, 16)

	size := len(benchmarkKeys)
	for i := 0; i < size/2; i++ {
		m[benchmarkKeys[i].String()] = i
	}

	b.StopTimer()
	b.StartTimer()

	for _, k := range benchmarkKeys {
		delete(m, k.String())
	}
}

func BenchmarkStandardMap(b *testing.B) {
	m := make(map[string]interface{}, 16)

	for i := 0; i < 10000; i++ {
		key := NewStringKey(fmt.Sprintf("%s", i))
		m[key.String()] = i
	}

	for i := 0; i < 10000; i++ {
		key := NewStringKey(fmt.Sprintf("%s", i))
		_ = m[key.String()]
	}

}

func showSimpleMap(m *hashMap) {
	for _, b := range m.buckets {
		fmt.Printf("%s\n", b.String())
	}
	fmt.Printf("----------------------\n")
	for _, b := range m.newBuckets {
		fmt.Printf("%s\n", b.String())
	}
}
