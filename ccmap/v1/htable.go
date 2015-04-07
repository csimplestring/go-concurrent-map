package v1

import (
	"fmt"

	. "github.com/csimplestring/go-concurrent-map/ccmap/key"
)

// htable is the underlying hash tables. It stores
// <key, value> pairs in buckets.
type htable struct {
	mask    int
	buckets []Bucket
}

// newHtable creates a new empty htable with specified size.
// Note that size should always be 2^n.
func newHtable(size int) (*htable, error) {
	if size < 0 {
		return nil,
			fmt.Errorf("Illegal arg: %d, size of tables should be positive.")
	}

	n := 1
	for n < size {
		n = n << 1
	}
	size = n

	buckets := make([]Bucket, size)
	for i := 0; i < size; i++ {
		buckets[i] = newBucket()
	}

	return &htable{
		mask:    size - 1,
		buckets: buckets,
	}, nil
}

// indexFor gives index of bucket for hash. It equals MOD operator.
func (ht *htable) indexFor(hash int) int {
	return hash & ht.mask
}

// get gets Entry based on key.
func (ht *htable) get(key Key) (Entry, bool) {
	index := ht.indexFor(key.Hash())
	return ht.buckets[index].Get(key)
}

// put puts en at the beginning of bucket.
func (ht *htable) put(en Entry) int {
	index := ht.indexFor(en.Key().Hash())
	return ht.buckets[index].Put(en)
}

// delete deletes value based on key.
func (ht *htable) delete(key Key) (Entry, int) {
	index := ht.indexFor(key.Hash())
	return ht.buckets[index].Delete(key)
}

// push inserts en at the end of bucket.
func (ht *htable) push(en Entry) bool {
	index := ht.indexFor(en.Key().Hash())
	return ht.buckets[index].Push(en)
}

// size returns the number of entries in the buckets.
func (ht *htable) size() int {
	cnt := 0
	for _, b := range ht.buckets {
		cnt += b.Size()
	}
	return cnt
}
