package ccmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBucketFindEntry(t *testing.T) {
	b := newBucket()
	pos, result := b.findEntry(NewStringKey("k2"))
	assert.Nil(t, result)
	assert.Equal(t, -1, pos)

	b = &bucket{
		entries: []Entry{
			newEntry(NewStringKey("k1"), 1),
			newEntry(NewStringKey("k2"), 2),
			newEntry(NewStringKey("k3"), 3),
		},
	}

	pos, result = b.findEntry(NewStringKey("k2"))
	assert.Equal(t, "k2", result.Key().String())
	assert.Equal(t, 2, result.Value())
	assert.Equal(t, 1, pos)

	pos, result = b.findEntry(NewStringKey("k4"))
	assert.Nil(t, result)
	assert.Equal(t, -1, pos)
}

func TestBucketPut(t *testing.T) {
	b := newBucket()
	b.Put(newEntry(NewStringKey("k1"), 1))
	b.Put(newEntry(NewStringKey("k2"), 2))
	assert.Equal(t, "[[k1 1],[k2 2],]", b.String())

	b.Put(newEntry(NewStringKey("k2"), 7))
	assert.Equal(t, "[[k1 1],[k2 7],]", b.String())
}

func TestBucketGet(t *testing.T) {
	b := &bucket{
		entries: []Entry{
			newEntry(NewStringKey("k1"), 1),
			newEntry(NewStringKey("k2"), 2),
			newEntry(NewStringKey("k3"), 3),
		},
	}

	assert.Equal(t, 2, b.Get(NewStringKey("k2")).Value())
	assert.Equal(t, "k2", b.Get(NewStringKey("k2")).Key().String())
	assert.Nil(t, b.Get(NewStringKey("k4")))
}

func TestBucketDelete(t *testing.T) {
	b := &bucket{
		entries: []Entry{
			newEntry(NewStringKey("k1"), 1),
			newEntry(NewStringKey("k2"), 2),
			newEntry(NewStringKey("k3"), 3),
		},
	}

	e := b.Delete(NewStringKey("k2"))
	assert.Equal(t, "[[k1 1],[k3 3],]", b.String())
	assert.Equal(t, "[k2 2]", e.String())

	e = b.Delete(NewStringKey("k4"))
	assert.Equal(t, "[[k1 1],[k3 3],]", b.String())
	assert.Nil(t, e)
}
