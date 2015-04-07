package v1

import (
	"fmt"
	"testing"

	"github.com/csimplestring/go-concurrent-map/algo/random"
	. "github.com/csimplestring/go-concurrent-map/ccmap/key"
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

func TestHashMapSize(t *testing.T) {
	m, _ := newHashMap(6)

	m.Put(NewStringKey("1"), 1)
	m.Put(NewStringKey("2"), 1)
	m.Put(NewStringKey("3"), 1)
	assert.Equal(t, 3, m.Size())

	m.Put(NewStringKey("3"), 3)
	assert.Equal(t, 3, m.Size())

	m.Delete(NewStringKey("2"))
	assert.Equal(t, 2, m.Size())

	m.Delete(NewStringKey("2"))
	assert.Equal(t, 2, m.Size())
}

func TestHashMapPut(t *testing.T) {
	m, _ := NewHashMap(100)

	for i := 0; i < 30; i++ {
		key := NewStringKey(fmt.Sprintf("%d", i))
		m.Put(key, i)
	}
}

func TestHashMapGetOK(t *testing.T) {
	m, _ := NewHashMap(100)

	for i := 0; i < 30; i++ {
		key := NewStringKey(fmt.Sprintf("%d", i))
		m.Put(key, i)
	}

	for i := 0; i < 30; i++ {
		actual, ok := m.Get(NewStringKey(fmt.Sprintf("%d", i)))
		assert.True(t, ok)
		assert.Equal(t, i, actual)
	}

	for i := 0; i < 30; i++ {
		key := NewStringKey(fmt.Sprintf("%d", i))
		m.Put(key, i*2)
	}

	for i := 0; i < 30; i++ {
		actual, ok := m.Get(NewStringKey(fmt.Sprintf("%d", i)))
		assert.True(t, ok)
		assert.Equal(t, i*2, actual)
	}
}

func TestHashMapGetFail(t *testing.T) {
	m, _ := NewHashMap(100)

	for i := 0; i < 30; i++ {
		key := NewStringKey(fmt.Sprintf("%d", i))
		m.Put(key, i)
	}

	for i := 31; i < 60; i++ {
		actual, ok := m.Get(NewStringKey(fmt.Sprintf("%d", i)))
		assert.False(t, ok)
		assert.Nil(t, actual)
	}
}

func TestHashMapDeleteOK(t *testing.T) {
	m, _ := NewHashMap(100)

	for i := 0; i < 30; i++ {
		key := NewStringKey(fmt.Sprintf("%d", i))
		m.Put(key, i)
	}

	for i := 0; i < 30; i++ {
		key := NewStringKey(fmt.Sprintf("%d", i))
		ok := m.Delete(key)
		assert.True(t, ok, "%d", i)
	}

	for i := 0; i < 30; i++ {
		actual, ok := m.Get(NewStringKey(fmt.Sprintf("%d", i)))
		assert.False(t, ok)
		assert.Nil(t, actual)
	}
}

func TestHashMapDeleteFail(t *testing.T) {
	m, _ := NewHashMap(100)

	for i := 0; i < 30; i++ {
		key := NewStringKey(fmt.Sprintf("%d", i))
		m.Put(key, i)
	}

	for i := 31; i < 60; i++ {
		key := NewStringKey(fmt.Sprintf("%d", i))
		ok := m.Delete(key)
		assert.False(t, ok)
	}
}

func BenchmarkHashMapPut(b *testing.B) {
	m, _ := NewHashMap(100)

	for i, k := range benchmarkKeys {
		m.Put(k, i)
	}
}

func BenchmarkHashMapGet(b *testing.B) {
	m, _ := NewHashMap(100)

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

func BenchmarkHashMapDelete(b *testing.B) {
	m, _ := NewHashMap(100)

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
	for _, b := range m.tables[0].buckets {
		fmt.Printf("%s\n", b.String())
	}
	fmt.Printf("----------------------\n")
	if m.tables[1] != nil {
		for _, b := range m.tables[1].buckets {
			fmt.Printf("%s\n", b.String())
		}
	}
}
