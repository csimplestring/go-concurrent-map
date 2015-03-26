package ccmap

import "fmt"

const (
	BUCKET_SIZE_DEFAULT = 16
)

// Map defines the functions that a map should support
type Map interface {
	Put(key Key, val interface{}) error
	Get(key Key) interface{}
}

// NewMap new a Map.
func NewMap() *hashMap {
	buckets := make([]Bucket, BUCKET_SIZE_DEFAULT)
	for i := 0; i < BUCKET_SIZE_DEFAULT; i++ {
		buckets[i] = newBucket()
	}

	loadFactor := float32(0.75)
	threshold := (int)((float32)(BUCKET_SIZE_DEFAULT) * loadFactor)

	return &hashMap{
		loadFactor: loadFactor,
		threshold:  threshold,
		entrySize:  0,
		bucketSize: BUCKET_SIZE_DEFAULT,
		buckets:    buckets,

		rehashIdx: 0,

		hashFunc: hashFor,
		slotFunc: indexFor,
	}
}

// hashMap is the default implementation of Map.
type hashMap struct {
	loadFactor float32
	threshold  int
	entrySize  int

	bucketSize int
	buckets    []Bucket

	rehashIdx     int
	newBucketSize int
	newBuckets    []Bucket

	hashFunc func(h int) int
	slotFunc func(h, length int) int
}

func (h *hashMap) allocNewBuckets(length int) {
	h.newBuckets = make([]Bucket, length)
	h.newBucketSize = length
	for i := 0; i < length; i++ {
		h.newBuckets[i] = newBucket()
	}
}

// move moves n old entries to new bucket
func (h *hashMap) moveEntry(n int) error {
	for i := 0; i < n && h.rehashIdx < h.bucketSize; {
		b := h.buckets[h.rehashIdx]
		// this bucket is empty
		if b.Size() == 0 {
			h.rehashIdx++
			continue
		}

		// delete old entry
		en, ok := b.Pop()
		if !ok {
			return fmt.Errorf("Bucket is not empty but can not pop entry")
		}

		if _, delta := b.Delete(en.Key()); delta > 1 {
			h.entrySize -= delta - 1
		}

		// move to new bucket
		slot := h.slotFunc(h.hashFunc(en.Key().Hash()), h.newBucketSize)
		h.newBuckets[slot].Put(en)

		i++
	}

	if h.rehashIdx == h.bucketSize {
		h.switchBuckets()
	}
	return nil
}

// switchBuckets remove old buckets
func (h *hashMap) switchBuckets() {
	h.buckets = h.newBuckets

	// update rehashIdx, bucketSize, threshold
	h.rehashIdx = 0
	h.bucketSize = h.newBucketSize
	h.threshold += (int)((float32)(h.newBucketSize) * h.loadFactor)

	// clean up
	h.newBucketSize = 0
	h.newBuckets = nil
}

func (h *hashMap) Put(key Key, val interface{}) error {
	entry := newEntry(key, val)

	if (h.entrySize + 1) < h.threshold {
		slot := h.slotFunc(h.hashFunc(key.Hash()), h.bucketSize)
		if ok := h.buckets[slot].Put(entry); ok {
			h.entrySize++
		}
		return nil
	}

	// allocate new bucket slice if needed
	if h.newBucketSize == 0 {
		h.allocNewBuckets(2 * h.bucketSize)
	}

	// insert entry to new bucket slice
	slot := h.slotFunc(h.hashFunc(key.Hash()), h.newBucketSize)
	if ok := h.newBuckets[slot].Put(entry); ok {
		h.entrySize++
	}

	// move 2 old entries to new bucket
	h.moveEntry(2)
	return nil
}

func (h *hashMap) Get(key Key) interface{} {
	if h.newBuckets != nil {

		slot := h.slotFunc(h.hashFunc(key.Hash()), h.newBucketSize)
		if en, ok := h.newBuckets[slot].Get(key); ok {
			return en.Value()
		}
	}

	slot := h.slotFunc(h.hashFunc(key.Hash()), h.bucketSize)
	en, ok := h.buckets[slot].Get(key)

	if !ok {
		return nil
	}
	return en.Value()
}

func (h *hashMap) Delete(key Key) bool {
	deleted := 0
	if h.newBuckets != nil {
		slot := h.slotFunc(h.hashFunc(key.Hash()), h.newBucketSize)
		_, cnt := h.newBuckets[slot].Delete(key)
		deleted += cnt
	}

	slot := h.slotFunc(h.hashFunc(key.Hash()), h.bucketSize)
	_, cnt := h.buckets[slot].Delete(key)
	deleted += cnt

	if deleted > 0 {
		return true
	}
	return false
}

func hashFor(h int) int {
	h ^= (h >> 20) ^ (h >> 12)
	return h ^ (h >> 7) ^ (h >> 4)
}

func indexFor(h, length int) int {
	return h & (length - 1)
}
