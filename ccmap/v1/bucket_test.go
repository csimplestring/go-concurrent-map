package v1

import (
	"testing"

	. "github.com/csimplestring/go-concurrent-map/ccmap/key"
	"github.com/stretchr/testify/assert"
)

func TestBucketFindEntry(t *testing.T) {
}

func TestBucketPush(t *testing.T) {
	b := newBucket()

	b.Push(newEntry(NewStringKey("k1"), 1))
	assert.Equal(t, "[[k1 1],]", b.String())

	b.Push(newEntry(NewStringKey("k2"), 2))
	assert.Equal(t, "[[k1 1],[k2 2],]", b.String())

	b.Push(newEntry(NewStringKey("k3"), 3))
	assert.Equal(t, "[[k1 1],[k2 2],[k3 3],]", b.String())
}

func TestBucketPut(t *testing.T) {
	tests := []struct {
		b   Bucket
		str string
	}{
		{
			newBucket(),
			"[[k2 7],[k1 1],]",
		},
	}

	for i, test := range tests {
		t.Logf("test[%d]\n", i)

		ok := test.b.Put(newEntry(NewStringKey("k1"), 1))
		assert.Equal(t, entryAdd, ok)

		ok = test.b.Put(newEntry(NewStringKey("k2"), 2))
		assert.Equal(t, entryAdd, ok)

		ok = test.b.Put(newEntry(NewStringKey("k2"), 7))
		assert.Equal(t, entryReplace, ok)

		assert.Equal(t, test.str, test.b.String())
	}
}

func TestBucketGet(t *testing.T) {
	b2 := newBucket()
	b2.Put(newEntry(NewStringKey("k1"), 1))
	b2.Put(newEntry(NewStringKey("k2"), 2))
	b2.Put(newEntry(NewStringKey("k3"), 3))

	tests := []struct {
		b Bucket
	}{
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
	b2 := newBucket()
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
	b2 := newBucket()
	b2.Put(newEntry(NewStringKey("k1"), 1))
	b2.Put(newEntry(NewStringKey("k2"), 2))
	b2.Put(newEntry(NewStringKey("k3"), 3))

	tests := []struct {
		b   Bucket
		key string
	}{
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
	b2 := newBucket()
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
