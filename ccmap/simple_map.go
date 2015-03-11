package ccmap

const (
	BUCKET_SIZE_DEFAULT = 16
)

func NewSimpleMap() *SimpleMap {
	buckets := make([]Bucket, BUCKET_SIZE_DEFAULT)
	for i := 0; i < BUCKET_SIZE_DEFAULT; i++ {
		buckets[i] = newBucket()
	}

	return &SimpleMap{
		loadFactor: 0.75,
		bucketSize: BUCKET_SIZE_DEFAULT,
		entrySize:  0,
		threshold:  BUCKET_SIZE_DEFAULT,
		buckets:    buckets,
	}
}

type SimpleMap struct {
	loadFactor float32
	entrySize  int
	bucketSize int
	threshold  int
	buckets    []Bucket
}

func (s *SimpleMap) Put(key Key, val interface{}) error {
	h := hashFor(key.Hash())
	h = indexFor(h, len(s.buckets))
	s.buckets[h].Put(newEntry(key, val))

	s.entrySize++
	if s.entrySize > s.threshold {
		s.resize(2 * s.bucketSize)
	}

	return nil
}

func (s *SimpleMap) Get(key Key) interface{} {
	h := hashFor(key.Hash())
	h = indexFor(h, len(s.buckets))

	en := s.buckets[h].Get(key)
	if en == nil {
		return nil
	}
	return en.Value()
}

func (s *SimpleMap) Delete(key Key) bool {
	h := hashFor(key.Hash())
	h = indexFor(h, len(s.buckets))

	en := s.buckets[h].Delete(key)
	if en == nil {
		return false
	}
	return true
}

func (s *SimpleMap) resize(length int) {
	old := s.buckets

	s.bucketSize = length
	s.threshold = length

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

func hashFor(h int) int {
	h ^= (h >> 20) ^ (h >> 12)
	return h ^ (h >> 7) ^ (h >> 4)
}

func indexFor(h, length int) int {
	return h & (length - 1)
}
