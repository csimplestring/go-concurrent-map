package ccmap

import "fmt"

type Entry interface {
	Key() Key
	Value() interface{}

	SetKey(Key)
	SetValue(interface{})

	String() string
}

// newBucket creates a new buckect.
func newEntry(k Key, v interface{}) Entry {
	return &entry{
		k: k,
		v: v,
	}
}

// bucket stores value and key.
type entry struct {
	k Key
	v interface{}
}

func (e *entry) Key() Key {
	return e.k
}

func (e *entry) SetKey(k Key) {
	e.k = k
}

func (e *entry) Value() interface{} {
	return e.v
}

func (e *entry) SetValue(v interface{}) {
	e.v = v
}

// String returns a string representation of b.
func (e *entry) String() string {
	return fmt.Sprintf("[%s %v]", e.k.String(), e.v)
}
