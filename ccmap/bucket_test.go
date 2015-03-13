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
	tests := []struct {
		b   Bucket
		str string
	}{
		{
			newBucket(),
			"[[k1 1],[k2 7],]",
		},
		{
			newLinkedBucket(),
			"[[k2 7],[k2 2],[k1 1],]",
		},
	}

	for i, test := range tests {
		t.Logf("test[%d]\n", i)

		ok := test.b.Put(newEntry(NewStringKey("k1"), 1))
		assert.True(t, ok)

		ok = test.b.Put(newEntry(NewStringKey("k2"), 2))
		assert.True(t, ok)

		ok = test.b.Put(newEntry(NewStringKey("k2"), 7))

		assert.Equal(t, test.str, test.b.String())
	}

}

func TestBucketGet(t *testing.T) {
	b1 := &bucket{
		entries: []Entry{
			newEntry(NewStringKey("k1"), 1),
			newEntry(NewStringKey("k2"), 2),
			newEntry(NewStringKey("k3"), 3),
		},
	}

	b2 := newLinkedBucket()
	b2.Put(newEntry(NewStringKey("k1"), 1))
	b2.Put(newEntry(NewStringKey("k2"), 2))
	b2.Put(newEntry(NewStringKey("k3"), 3))

	tests := []struct {
		b Bucket
	}{
		{
			b1,
		},
		{
			b2,
		},
	}

	for i, test := range tests {
		t.Logf("tests[%d]", i)

		en, ok := test.b.Get(NewStringKey("k2"))
		assert.True(t, ok)
		assert.Equal(t, 2, en.Value())

		en, ok = test.b.Get(NewStringKey("k4"))
		assert.False(t, ok)
		assert.Nil(t, en)
	}
}

func TestBucketDeleteOK(t *testing.T) {
	b1 := &bucket{
		entries: []Entry{
			newEntry(NewStringKey("k1"), 1),
			newEntry(NewStringKey("k2"), 2),
			newEntry(NewStringKey("k3"), 3),
		},
	}

	b2 := newLinkedBucket()
	b2.Put(newEntry(NewStringKey("k1"), 1))
	b2.Put(newEntry(NewStringKey("k2"), 2))
	b2.Put(newEntry(NewStringKey("k3"), 3))

	tests := []struct {
		b   Bucket
		key string
		bs  string
		es  string
		cnt int
	}{
		{
			b1,
			"k2",
			"[[k1 1],[k3 3],]",
			"[k2 2]",
			1,
		},
		{
			b2,
			"k2",
			"[[k3 3],[k1 1],]",
			"[k2 2]",
			1,
		},
	}

	for i, test := range tests {
		t.Logf("tests[%d]", i)

		e, cnt := test.b.Delete(NewStringKey(test.key))
		assert.Equal(t, test.cnt, cnt)
		assert.Equal(t, test.bs, test.b.String())
		assert.Equal(t, test.es, e.String())
	}
}

func TestBucketDeleteFailed(t *testing.T) {
	b1 := &bucket{
		entries: []Entry{
			newEntry(NewStringKey("k1"), 1),
			newEntry(NewStringKey("k2"), 2),
			newEntry(NewStringKey("k3"), 3),
		},
	}

	b2 := newLinkedBucket()
	b2.Put(newEntry(NewStringKey("k1"), 1))
	b2.Put(newEntry(NewStringKey("k2"), 2))
	b2.Put(newEntry(NewStringKey("k3"), 3))

	tests := []struct {
		b   Bucket
		key string
	}{
		{
			b1,
			"k4",
		},
		{
			b2,
			"k4",
		},
	}

	for i, test := range tests {
		t.Logf("tests[%d]", i)

		e, cnt := test.b.Delete(NewStringKey(test.key))
		assert.Equal(t, 0, cnt)
		assert.Nil(t, e)
	}
}

func TestBucketPopOK(t *testing.T) {
	b2 := newLinkedBucket()
	b2.Put(newEntry(NewStringKey("k1"), 1))

	tests := []struct {
		b  Bucket
		bs string
		es string
	}{
		{
			b2,
			"[]",
			"[k1 1]",
		},
	}

	for i, test := range tests {
		t.Logf("tests[%d]", i)

		e, ok := test.b.Pop()
		assert.True(t, ok)
		assert.Equal(t, test.bs, test.b.String())
		assert.Equal(t, test.es, e.String())
	}
}
