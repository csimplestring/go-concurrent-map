package ccmap

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/csimplestring/go-concurrent-map/algo/random"
)

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

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

func BenchmarkSimpleMapPut(b *testing.B) {
	m := NewSimpleMap()

	for i := 0; i < 10000; i++ {
		key := NewStringKey(random.NewLen(20))
		m.Put(key, i)
	}
}

func BenchmarkSimpleMapGet(b *testing.B) {
	m := NewSimpleMap()

	for i := 0; i < 10000; i++ {
		key := NewStringKey(random.NewLen(20))
		m.Put(key, i)
	}
	b.StopTimer()
	b.StartTimer()

	for i := 0; i < 10000; i++ {
		key := NewStringKey(random.NewLen(20))
		_ = m.Get(key)
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

	for i := 0; i < 10000; i++ {
		key := NewStringKey(random.NewLen(20))
		m[key.String()] = i
	}
}

func BenchmarkStandardMapGet(b *testing.B) {
	m := make(map[string]interface{}, 16)

	for i := 0; i < 10000; i++ {
		key := NewStringKey(random.NewLen(20))
		m[key.String()] = i
	}
	b.StopTimer()
	b.StartTimer()

	for i := 0; i < 10000; i++ {
		key := NewStringKey(random.NewLen(20))
		_ = m[key.String()]
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
