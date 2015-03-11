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

func TestSimpleMapPut(t *testing.T) {
	m := NewSimpleMap()

	m.Put(NewStringKey("k1"), 1)
	m.Put(NewStringKey("k2"), 2)
	m.Put(NewStringKey("k2"), 3)
}

func TestSimpleMapGet(t *testing.T) {
	m := NewSimpleMap()

	for i := 0; i < 10000; i++ {
		key := NewStringKey(fmt.Sprintf("%s", i))
		m.Put(key, i)
	}

	for i := 0; i < 10000; i++ {
		assert.Equal(t, i, m.Get(NewStringKey(fmt.Sprintf("%s", i))))
	}
}

func TestSimpleMapDelete(t *testing.T) {
	m := NewSimpleMap()

	m.Put(NewStringKey("k1"), 1)
	m.Put(NewStringKey("k2"), 2)
	m.Put(NewStringKey("k2"), 3)

	assert.True(t, m.Delete(NewStringKey("k1")))
	assert.Nil(t, m.Get(NewStringKey("k1")))
}

func BenchmarkSimpleMapPut(b *testing.B) {
	m := NewSimpleMap()

	for i, k := range benchmarkKeys {
		m.Put(k, i)
	}
}

func BenchmarkSimpleMapGet(b *testing.B) {
	m := NewSimpleMap()

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

func BenchmarkSimpleMapDelete(b *testing.B) {
	m := NewSimpleMap()

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

func BenchmarkSimpleMap(b *testing.B) {
	m := NewSimpleMap()

	for i := 0; i < 10000; i++ {
		key := NewStringKey(fmt.Sprintf("%s", i))
		m.Put(key, i)
	}

	for i := 0; i < 10000; i++ {
		key := NewStringKey(fmt.Sprintf("%s", i))
		_ = m.Get(key)
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

func showSimpleMap(m *SimpleMap) {
	for _, b := range m.buckets {
		fmt.Printf("%s\n", b.(*bucket).String())
	}
}
