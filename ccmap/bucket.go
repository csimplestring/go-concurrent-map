package ccmap

const (
	ENTRY_SIZE  = 8
	BUCKET_SIZE = 8
)

type Bucket interface {
	Put(Entry)
	Get(key Key) Entry
	Delete(key Key) Entry
	Entries() []Entry
	Size() int
}

// newEntry creates a new Entry.
func newBucket() *bucket {
	return &bucket{
	//entries: make([]Entry, ENTRY_SIZE, ENTRY_SIZE),
	}
}

// bucket uses slice to store Entry.
//
//
//
//
//
//
type bucket struct {
	entries []Entry
}

func (b *bucket) Put(e Entry) {
	key := e.Key()
	_, r := b.findEntry(key)

	if r != nil {
		r.SetValue(e.Value())
		return
	}
	b.entries = append(b.entries, e)
}

func (b *bucket) Get(key Key) Entry {
	_, e := b.findEntry(key)
	return e
}

func (b *bucket) Delete(key Key) Entry {
	i, e := b.findEntry(key)
	if i == -1 {
		return nil
	}

	b.deleteEntry(i)
	return e
}

func (b *bucket) Entries() []Entry {
	return b.entries
}

func (b *bucket) Size() int {
	return len(b.entries)
}

func (b *bucket) deleteEntry(i int) {
	b.entries[i] = b.entries[len(b.entries)-1]
	b.entries[len(b.entries)-1] = nil
	b.entries = b.entries[:len(b.entries)-1]
}

func (b *bucket) findEntry(key Key) (int, Entry) {
	for i, e := range b.entries {
		if e.Key().Equal(key) {
			return i, e
		}
	}
	return -1, nil
}

// String returns a string representation of bucket list.
func (b *bucket) String() string {
	str := "["
	for _, e := range b.entries {
		str += e.String() + ","
	}
	str += "]"

	return str
}

// bucketV1 uses array to store Entry.
//
//
//
//
//
type bucketV1 struct {
	entries [ENTRY_SIZE]Entry
}
