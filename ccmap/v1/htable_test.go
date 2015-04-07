package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHtableOk(t *testing.T) {
	h, err := newHtable(3)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(h.buckets))

	h, _ = newHtable(15)
	assert.Equal(t, 16, len(h.buckets))

	h, _ = newHtable(24)
	assert.Equal(t, 32, len(h.buckets))
}

func TestNewHtableError(t *testing.T) {
	_, err := newHtable(-1)
	assert.Error(t, err)
}

func TestHtableIndexFor(t *testing.T) {
	ht := &htable{
		mask: 3,
	}
	assert.Equal(t, 1, ht.indexFor(1))
	assert.Equal(t, 0, ht.indexFor(4))
	assert.Equal(t, 3, ht.indexFor(7))
}

func TestHtableGet(t *testing.T) {

}
