package ccmap

import "errors"

const (
	BUCKET_SIZE_DEFAULT = 16
	StatusRehashing     = -1
)

var (
	ErrPut = errors.New("Put() Error")
	ErrGet = errors.New("Get() Error")
)

// Map defines the functions that a map should support
// TODO:
// 2. shrink table
// 3. rehash table in background
type Map interface {
	Put(key Key, val interface{}) bool
	Get(key Key) (interface{}, bool)
	Delete(key Key) bool
}

// NewMap new a Map.
func NewMap() Map {
	table := make([]*htable, 2)
	table[0] = newHtable(4)
	table[1] = nil

	return &hashMap{
		entryCnt:  0,
		table:     table,
		rehashIdx: -1,
	}
}

// hashMap is the default implementation of Map.
type hashMap struct {
	// -1: no rehash; otherwise it is rehashing
	rehashIdx int
	entryCnt  int
	table     []*htable
}

// Put puts <key, val> pair in correct slot.
// It returns true if succeed; otherwise false.
func (h *hashMap) Put(key Key, val interface{}) bool {
	entry := newEntry(key, val)

	if h.rehashIdx == -1 {
		if ok := h.putEntry(0, entry); !ok {
			return false
		}

		if h.entryCnt > h.table[0].size {
			h.beginRehash()
			h.rehash()
		}
		return true
	}

	if ok := h.putEntry(1, entry); !ok {
		return false
	}

	h.rehash()
	return true
}

// Get gets the value based on key.
// If value exists, it returns value and TRUE;
// otherwise it returns nil and FALSE.
func (h *hashMap) Get(key Key) (interface{}, bool) {
	if h.rehashIdx != -1 {
		if en, ok := h.table[1].get(key); ok {
			h.rehash()
			return en.Value(), true
		}
	}

	if en, ok := h.table[0].get(key); ok {
		return en.Value(), true
	}
	return nil, false
}

// Delete deletes value based on key.
// It returns TRUE if key exists; otherise FALSE.
func (h *hashMap) Delete(key Key) bool {
	deleted := 0
	_, cnt := h.table[0].delete(key)
	deleted += cnt

	if h.rehashIdx != -1 {
		_, cnt := h.table[1].delete(key)
		deleted += cnt
		h.rehash()
	}

	h.entryCnt -= deleted

	if deleted > 0 {
		return true
	}
	return false
}

// beginRehash sets rehashIdx to be 0, creates new htable for
// table[1].
func (h *hashMap) beginRehash() {
	h.rehashIdx = 0
	newSize := h.table[0].size * 2
	h.table[1] = newHtable(newSize)
}

// stopRehash switches old and new htable internally, resets
// rehashIdx to be -1.
func (h *hashMap) stopRehash() {
	h.table[0] = h.table[1]
	h.table[1] = nil
	h.rehashIdx = -1
}

// rehash moves table[0]'s entry to table[1]
func (h *hashMap) rehash() {
	// find the non-empty bucket
	b := h.table[0].buckets[h.rehashIdx]
	for b.Size() == 0 && h.rehashIdx < h.table[0].size {
		b = h.table[0].buckets[h.rehashIdx]
		h.rehashIdx++
	}

	// move old entries
	for en, ok := b.Pop(); ok; en, ok = b.Pop() {
		h.table[1].push(en)
	}

	// rehash ends
	if h.rehashIdx == h.table[0].size {
		h.stopRehash()
	}
}

// putEntry puts en into table[tableIdx].
// It returns true if succeeds, otherwise false.
func (h *hashMap) putEntry(tableIdx int, en Entry) bool {
	ok := h.table[tableIdx].put(en)
	if !ok {
		return false
	}

	h.entryCnt++
	return true
}

// htable is the underlying hash table. It stores
// <key, value> pairs in buckets.
type htable struct {
	mask    int
	size    int
	buckets []Bucket
}

// newHtable creates a new empty htable with specified size.
// Note that size should always be 2^n.
func newHtable(size int) *htable {
	buckets := make([]Bucket, size)
	for i := 0; i < size; i++ {
		buckets[i] = newBucket()
	}

	return &htable{
		mask:    size - 1,
		size:    size,
		buckets: buckets,
	}
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
func (ht *htable) put(en Entry) bool {
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
