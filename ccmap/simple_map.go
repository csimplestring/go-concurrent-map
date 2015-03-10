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
		buckets: buckets,
	}
}

type SimpleMap struct {
	buckets []Bucket
}

func (s *SimpleMap) Put(key Key, val interface{}) error {
	h := key.Hash() % len(s.buckets)
	s.buckets[h].Put(newEntry(key, val))
	return nil
}

func (s *SimpleMap) Get(key Key) interface{} {
	h := key.Hash() % len(s.buckets)
	return s.buckets[h].Get(key).Value()
}
