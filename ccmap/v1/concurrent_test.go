package v1

import (
	"fmt"
	"testing"

	. "github.com/csimplestring/go-concurrent-map/ccmap/key"
	"github.com/stretchr/testify/assert"
)

func TestNewConcurrentMap(t *testing.T) {
	_, err := NewConcurrentMap(8)
	assert.Nil(t, err)
}

func TestCCHashMapPut(t *testing.T) {
	m, _ := NewConcurrentMap(4)

	for i := 0; i < 4; i++ {
		key := NewStringKey(fmt.Sprintf("%d", i))
		m.Put(key, i)
	}
	m.(*concurrentHashMap).Stat()
}

func TestCCHashMapPut2(t *testing.T) {
	m, _ := NewConcurrentMap(8)

	for i := 0; i < 100000; i++ {
		key := NewStringKey(fmt.Sprintf("%d", i))
		m.Put(key, i)
	}

}

func TestCCHashMapGetOK(t *testing.T) {
	m, _ := NewConcurrentMap(8)

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

func TestCCHashMapGetFail(t *testing.T) {
	m, _ := NewConcurrentMap(8)

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

func BenchmarkCCHashMapPut(b *testing.B) {
	m, _ := NewConcurrentMap(8)

	for i, k := range benchmarkKeys {
		m.Put(k, i)
	}
}
