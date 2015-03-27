package ccmap

import (
	"fmt"
	"testing"

	"github.com/csimplestring/go-concurrent-map/algo/random"
	"github.com/stretchr/testify/assert"
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

	for i := 0; i < 30; i++ {
		key := NewStringKey(fmt.Sprintf("%s", i))
		m.Put(key, i)
	}

	for i := 0; i < 30; i++ {
		key := NewStringKey(fmt.Sprintf("%s", i))
		m.Put(key, i*2)
	}

	for i := 0; i < 30; i++ {
		assert.Equal(t, i*2, m.Get(NewStringKey(fmt.Sprintf("%s", i))))
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

func BenchmarkNativePut(b *testing.B) {
	m := make(map[string]interface{}, 16)

	for i, k := range benchmarkKeys {
		m[k.String()] = i
	}
}

func BenchmarkNativeGet(b *testing.B) {
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

func BenchmarkNativeDelete(b *testing.B) {
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

func BenchmarkNative(b *testing.B) {
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
	for _, b := range m.table[0].buckets {
		fmt.Printf("%s\n", b.String())
	}
	fmt.Printf("----------------------\n")
	if m.table[1] != nil {
		for _, b := range m.table[1].buckets {
			fmt.Printf("%s\n", b.String())
		}
	}
}
