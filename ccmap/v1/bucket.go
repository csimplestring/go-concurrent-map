package v1

import . "github.com/csimplestring/go-concurrent-map/ccmap/key"

const (
	entryErr     = 0
	entryAdd     = 1
	entryReplace = 2
)

type Bucket interface {
	Push(Entry) bool
	Put(Entry) int
	Get(key Key) (Entry, bool)
	Delete(key Key) (Entry, int)

	Entries() []Entry
	Pop() (Entry, bool)
	Size() int

	String() string
}

// newBucket creates a new bucket.
func newBucket() Bucket {
	return &bucket{
		cnt:  0,
		head: newLinkedEntry(nil, nil),
	}
}

// bucket is the default Bucket implementation.
type bucket struct {
	cnt  int
	head *linkedEntry
}

// Push appends en at the end of b.
func (b *bucket) Push(en Entry) bool {
	tail := b.head
	for current := b.head.next; current != nil; current = current.next {
		tail = current
	}
	tail.next = newLinkedEntry(en, nil)
	b.cnt++
	return true
}

// Put appends en at the beginning of b.
func (b *bucket) Put(en Entry) int {
	var p *linkedEntry = nil
	for current := b.head.next; current != nil; current = current.next {
		if current.Key().Equal(en.Key()) {
			p = current
		}
	}

	if p == nil {
		b.head.next = newLinkedEntry(en, b.head.next)
		b.cnt++
		return entryAdd
	} else {
		p.SetValue(en.Value())
		return entryReplace
	}
}

// Get finds entry based on key.
func (b *bucket) Get(key Key) (Entry, bool) {
	for current := b.head.next; current != nil; current = current.next {
		if current.Key().Equal(key) {
			return current.Entry, true
		}
	}
	return nil, false
}

// Delete deletes an entry based on key.
// Returns the deleted entry and the number of deleted entries.
func (b *bucket) Delete(key Key) (Entry, int) {
	var i = 0
	var d Entry

	prev := b.head
	for current := b.head.next; current != nil; current = current.next {
		if current.Key().Equal(key) {
			prev.next = current.next
			if i == 0 {
				d = current.Entry
			}
			i++
		}
		prev = current
	}

	b.cnt = b.cnt - i
	return d, i
}

// Pop pops the first entry. Returns false if no entry in b.
func (b *bucket) Pop() (Entry, bool) {
	if first := b.head.next; first != nil {
		b.head.next = first.next
		b.cnt--
		return first.Entry, true
	}
	return nil, false
}

// Entries returns a slice of all the Entries in b.
func (b *bucket) Entries() []Entry {
	entries := make([]Entry, b.cnt)
	i := 0
	for current := b.head.next; current != nil; current = current.next {
		entries[i] = current.Entry
		i++
	}

	return entries
}

// Size returns the number of entry in b.
func (b *bucket) Size() int {
	return b.cnt
}

// String returns a string representation of b.
func (b *bucket) String() string {
	str := "["
	current := b.head.next
	for current != nil {
		str += current.Entry.String() + ","
		current = current.next
	}
	str += "]"
	return str
}
