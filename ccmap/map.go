package ccmap

import "fmt"

const (
	BUCKET_SIZE_DEFAULT = 16
)

//
type Map interface {
	Put(key Key, val interface{}) error
	Get(key Key) interface{}
}

func NewSimpleMap() *SimpleMap {
	buckets := make([]Bucket, BUCKET_SIZE_DEFAULT)
	for i := 0; i < BUCKET_SIZE_DEFAULT; i++ {
		buckets[i] = newBucket()
	}

	return &SimpleMap{
		loadFactor: 0.75,
		bucketSize: BUCKET_SIZE_DEFAULT,
		entrySize:  0,
		buckets:    buckets,
	}
}

type SimpleMap struct {
	loadFactor float32
	entrySize  int
	bucketSize int
	buckets    []Bucket
}

func (s *SimpleMap) Put(key Key, val interface{}) error {
	h := hashFor(key.Hash())
	h = indexFor(h, len(s.buckets))

	if ok := s.buckets[h].Put(newEntry(key, val)); ok {
		s.entrySize++
	}

	if s.entrySize > s.bucketSize/2 {
		s.resize(2 * s.bucketSize)
	}

	return nil
}

func (s *SimpleMap) Get(key Key) interface{} {
	h := hashFor(key.Hash())
	h = indexFor(h, len(s.buckets))

	en, ok := s.buckets[h].Get(key)
	if !ok {
		return nil
	}

	return en.Value()
}

func (s *SimpleMap) Delete(key Key) bool {
	h := hashFor(key.Hash())
	h = indexFor(h, len(s.buckets))

	_, cnt := s.buckets[h].Delete(key)
	if cnt == 1 {
		s.entrySize--
	}

	return true
}

func (s *SimpleMap) resize(length int) {
	old := s.buckets

	s.bucketSize = length
	s.buckets = make([]Bucket, length)
	for i := 0; i < length; i++ {
		s.buckets[i] = newBucket()
	}

	for _, b := range old {
		oldEntries := b.Entries()

		for _, oldEntry := range oldEntries {
			h := hashFor(oldEntry.Key().Hash())
			h = indexFor(h, length)
			s.buckets[h].Put(oldEntry)
		}
	}
}

///////////////////////////////////////////
///
///
///
func NewLinkedMap() *LinkedMap {
	buckets := make([]Bucket, BUCKET_SIZE_DEFAULT)
	for i := 0; i < BUCKET_SIZE_DEFAULT; i++ {
		buckets[i] = newLinkedBucket()
	}

	loadFactor := float32(0.75)
	threshold := (int)((float32)(BUCKET_SIZE_DEFAULT) * loadFactor)

	return &LinkedMap{
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

type LinkedMap struct {
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

func (l *LinkedMap) allocNewBuckets(length int) {
	l.newBuckets = make([]Bucket, length)
	l.newBucketSize = length
	for i := 0; i < length; i++ {
		l.newBuckets[i] = newLinkedBucket()
	}
}

// move moves n old entries to new bucket
func (l *LinkedMap) moveEntry(n int) error {
	for i := 0; i < n && l.rehashIdx < l.bucketSize; {
		b := l.buckets[l.rehashIdx]
		// this bucket is empty
		if b.Size() == 0 {
			l.rehashIdx++
			continue
		}

		// delete old entry
		en, ok := b.Pop()
		if !ok {
			return fmt.Errorf("Bucket is not empty but can not pop entry")
		}

		if _, delta := b.Delete(en.Key()); delta > 1 {
			l.entrySize -= delta - 1
		}

		// move to new bucket
		h := l.slotFunc(l.hashFunc(en.Key().Hash()), l.newBucketSize)
		l.newBuckets[h].Put(en)

		i++
	}

	if l.rehashIdx == l.bucketSize {
		l.switchBuckets()
	}
	return nil
}

// switchBuckets remove old buckets
func (l *LinkedMap) switchBuckets() {
	l.buckets = l.newBuckets

	// update rehashIdx, bucketSize, threshold
	l.rehashIdx = 0
	l.bucketSize = l.newBucketSize
	l.threshold += (int)((float32)(l.newBucketSize) * l.loadFactor)

	// clean up
	l.newBucketSize = 0
	l.newBuckets = nil
}

func (l *LinkedMap) Put(key Key, val interface{}) error {
	entry := newEntry(key, val)

	if (l.entrySize + 1) < l.threshold {
		h := l.slotFunc(l.hashFunc(key.Hash()), l.bucketSize)
		if ok := l.buckets[h].Put(entry); ok {
			l.entrySize++
		}
		return nil
	}

	// allocate new bucket slice if needed
	if l.newBucketSize == 0 {
		l.allocNewBuckets(2 * l.bucketSize)
	}

	// insert entry to new bucket slice
	h := l.slotFunc(l.hashFunc(key.Hash()), l.newBucketSize)
	if ok := l.newBuckets[h].Put(entry); ok {
		l.entrySize++
	}

	// move 2 old entries to new bucket
	l.moveEntry(2)
	return nil
}

func (l *LinkedMap) Get(key Key) interface{} {
	if l.newBuckets != nil {

		h := l.slotFunc(l.hashFunc(key.Hash()), l.newBucketSize)
		if en, ok := l.newBuckets[h].Get(key); ok {
			return en.Value()
		}
	}

	h := l.slotFunc(l.hashFunc(key.Hash()), l.bucketSize)
	en, ok := l.buckets[h].Get(key)

	if !ok {
		return nil
	}
	return en.Value()
}

func (l *LinkedMap) Delete(key Key) bool {
	deleted := 0
	if l.newBuckets != nil {
		h := l.slotFunc(l.hashFunc(key.Hash()), l.newBucketSize)
		_, cnt := l.newBuckets[h].Delete(key)
		deleted += cnt
	}

	h := l.slotFunc(l.hashFunc(key.Hash()), l.bucketSize)
	_, cnt := l.buckets[h].Delete(key)
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
