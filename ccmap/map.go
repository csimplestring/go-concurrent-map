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
// TODO: 1. rehash in Get(), Delete()
// 2. shrink table
// 3. rehash table in background
type Map interface {
	Put(key Key, val interface{}) error
	Get(key Key) interface{}
}

// NewMap new a Map.
func NewMap() *hashMap {
	table := make([]*htable, 2)
	table[0] = newHtable(4)
	table[1] = nil

	return &hashMap{
		entryCnt:  0,
		table:     table,
		rehashIdx: -1,
	}
}

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

type htable struct {
	mask    int
	size    int
	buckets []Bucket
}

func (ht *htable) indexFor(hash int) int {
	return hash & ht.mask
}

func (ht *htable) get(key Key) (Entry, bool) {
	index := ht.indexFor(key.Hash())
	return ht.buckets[index].Get(key)
}

func (ht *htable) put(en Entry) bool {
	index := ht.indexFor(en.Key().Hash())
	return ht.buckets[index].Put(en)
}

func (ht *htable) delete(key Key) (Entry, int) {
	index := ht.indexFor(key.Hash())
	return ht.buckets[index].Delete(key)
}

func (ht *htable) push(en Entry) bool {
	index := ht.indexFor(en.Key().Hash())
	return ht.buckets[index].Push(en)
}

// hashMap is the default implementation of Map.
type hashMap struct {
	entryCnt int

	table []*htable

	// -1: no rehash; otherwise it is rehashing
	rehashIdx int
}

// switchBuckets remove old buckets, update rehashIdx
func (h *hashMap) switchTable() {
	h.table[0] = h.table[1]
	h.table[1] = nil
	h.rehashIdx = -1
}

func (h *hashMap) rehash() {
	// find the non-empty bucket
	b := h.table[0].buckets[h.rehashIdx]
	for b.Size() == 0 && h.rehashIdx < h.table[0].size {
		b = h.table[0].buckets[h.rehashIdx]
		h.rehashIdx++
	}

	// move to new table
	for en, ok := b.Pop(); ok; en, ok = b.Pop() {
		h.table[1].push(en)
		if _, delta := b.Delete(en.Key()); delta > 1 {
			h.entryCnt -= delta - 1
		}
	}

	if h.rehashIdx == h.table[0].size {
		h.switchTable()
	}
}

func (h *hashMap) putEntry(tableIdx int, en Entry) error {
	ok := h.table[tableIdx].put(en)
	if !ok {
		return ErrPut
	}

	h.entryCnt++
	return nil
}

func (h *hashMap) Put(key Key, val interface{}) error {
	entry := newEntry(key, val)

	// insert to table[0]
	if h.entryCnt < h.table[0].size && h.rehashIdx == -1 {
		return h.putEntry(0, entry)
	}

	// start to rehash
	if h.table[1] == nil {
		h.rehashIdx = 0
		newSize := h.table[0].size * 2
		h.table[1] = newHtable(newSize)
	}

	// insert to table[1]
	err := h.putEntry(1, entry)
	if err != nil {
		return err
	}

	h.rehash()
	return nil
}

func (h *hashMap) Get(key Key) interface{} {
	// check table[1] firstly
	if h.rehashIdx != -1 {
		if en, ok := h.table[1].get(key); ok {
			return en.Value()
		}
	}

	// then check table[0]
	if en, ok := h.table[0].get(key); ok {
		return en.Value()
	}
	return nil
}

func (h *hashMap) Delete(key Key) bool {
	deleted := 0
	_, cnt := h.table[0].delete(key)
	deleted += cnt

	if h.rehashIdx != -1 {
		_, cnt := h.table[1].delete(key)
		deleted += cnt
	}

	h.entryCnt -= deleted
	if deleted > 0 {
		return true
	}
	return false
}
