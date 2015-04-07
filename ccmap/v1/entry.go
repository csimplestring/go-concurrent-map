package v1

import (
	"fmt"

	. "github.com/csimplestring/go-concurrent-map/ccmap/key"
)

type Entry interface {
	Key() Key
	Value() interface{}

	SetKey(Key)
	SetValue(interface{})

	String() string
}

// newEntry creates a new entry.
func newEntry(k Key, v interface{}) Entry {
	return &entry{
		k: k,
		v: v,
	}
}

// entry is basic implementation of Entry.
type entry struct {
	k Key
	v interface{}
}

// Key returns the key.
func (e *entry) Key() Key {
	return e.k
}

// SetKey sets the key.
func (e *entry) SetKey(k Key) {
	e.k = k
}

// Value returns the value.
func (e *entry) Value() interface{} {
	return e.v
}

// SetValue sets the value.
func (e *entry) SetValue(v interface{}) {
	e.v = v
}

// String returns a string representation of b.
func (e *entry) String() string {
	return fmt.Sprintf("[%s %v]", e.k.String(), e.v)
}

// linkedEntry inplements Entry and links to next entry.
type linkedEntry struct {
	Entry
	next *linkedEntry
}

// newLinkedEntry new a linkedEntry
func newLinkedEntry(en Entry, next *linkedEntry) *linkedEntry {
	return &linkedEntry{
		Entry: en,
		next:  next,
	}
}
