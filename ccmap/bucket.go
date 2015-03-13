package ccmap

const (
	ENTRY_SIZE  = 8
	BUCKET_SIZE = 8
)

type Bucket interface {
	Put(Entry) bool
	Get(key Key) (Entry, bool)
	Delete(key Key) (Entry, int)
	Entries() []Entry
	Pop() (Entry, bool)
	Size() int

	String() string
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

func (b *bucket) Put(e Entry) bool {
	key := e.Key()
	_, r := b.findEntry(key)

	if r != nil {
		r.SetValue(e.Value())
		return false
	}

	b.entries = append(b.entries, e)
	return true
}

func (b *bucket) Get(key Key) (Entry, bool) {
	i, e := b.findEntry(key)
	if i == -1 {
		return nil, false
	}

	return e, true
}

func (b *bucket) Delete(key Key) (Entry, int) {
	i, e := b.findEntry(key)
	if i == -1 {
		return nil, 0
	}

	b.deleteEntry(i)
	return e, 1
}

func (b *bucket) Pop() (Entry, bool) {
	return nil, true
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

// bucketV1 uses linked list to store Entry.
// optimised for writes.
//
//
//
//
//

func newLinkedBucket() Bucket {
	return &linkedBucket{
		head: newLinkedEntry(nil, nil),
	}
}

func newLinkedEntry(en Entry, next *linkedEntry) *linkedEntry {
	return &linkedEntry{
		Entry: en,
		next:  next,
	}
}

type linkedBucket struct {
	totalEntryCnt int
	head          *linkedEntry
}

// Put appends en at the beginning of linkedEntry list.
func (l *linkedBucket) Put(en Entry) bool {
	l.head.next = newLinkedEntry(en, l.head.next)
	l.totalEntryCnt++
	return true
}

func (l *linkedBucket) Get(key Key) (Entry, bool) {
	for current := l.head.next; current != nil; current = current.next {
		if current.Key().Equal(key) {
			return current.Entry, true
		}
	}
	return nil, false
}

func (l *linkedBucket) Delete(key Key) (Entry, int) {
	var (
		i = 0
		d Entry
	)

	prev := l.head
	for current := l.head.next; current != nil; current = current.next {
		if current.Key().Equal(key) {
			prev.next = current.next
			if i == 0 {
				d = current.Entry
			}
			i++
		}
		prev = current
	}

	l.totalEntryCnt = l.totalEntryCnt - i

	return d, i
}

func (l *linkedBucket) Pop() (Entry, bool) {
	if first := l.head.next; first != nil {
		l.head.next = first.next
		l.totalEntryCnt--
		return first.Entry, true
	}
	return nil, false
}

func (l *linkedBucket) Entries() []Entry {
	entries := make([]Entry, l.totalEntryCnt)
	i := 0
	for current := l.head.next; current != nil; current = current.next {
		entries[i] = current.Entry
		i++
	}

	return entries
}

func (l *linkedBucket) Size() int {
	return l.totalEntryCnt
}

func (l *linkedBucket) String() string {
	str := "["
	current := l.head.next
	for current != nil {
		str += current.Entry.String() + ","
		current = current.next
	}
	str += "]"
	return str
}

type linkedEntry struct {
	Entry
	next *linkedEntry
}
