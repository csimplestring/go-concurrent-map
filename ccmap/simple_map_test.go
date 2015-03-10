package ccmap

import (
	"fmt"
	"testing"
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

func BenchmarkSimpleMap(b *testing.B) {
	m := NewSimpleMap()

	for i := 0; i < 10000; i++ {
		key := NewStringKey(fmt.Sprintf("%s", i))
		m.Put(key, i)
	}

	// for i := 0; i < 10000; i++ {
	// 	m.Get(NewStringKey(fmt.Sprintf("%s", i)))
	// }

}

func BenchmarkStandardMap(b *testing.B) {
	m := make(map[string]interface{}, 16)

	for i := 0; i < 10000; i++ {
		key := NewStringKey(fmt.Sprintf("%s", i))
		m[key.String()] = i
	}

	// for i := 0; i < 10000; i++ {
	// 	_ = m[fmt.Sprintf("%s", i)]
	// }

}

func showSimpleMap(m *SimpleMap) {
	for _, b := range m.buckets {
		fmt.Printf("%s\n", b.(*bucket).String())
	}
}
