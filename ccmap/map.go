package ccmap

import (
	"github.com/csimplestring/go-concurrent-map/algo/hash"
)

//
type ConcurrentMap interface {
	Put(key Key, val interface{}) error
	Get(key Key) interface{}
}

// Key defines function set of a key in map.
type Key interface {
	Hash() int
	Equal(Key) bool
	String() string
}

// NewStringKey return a new string key.
func NewStringKey(str string) Key {
	return &StringKey{
		str: str,
		h:   (int)(hash.BKDRHash(str)),
	}
}

type StringKey struct {
	h   int
	str string
}

func (s *StringKey) Hash() int {
	return s.h
}

func (s *StringKey) Equal(k Key) bool {
	other, ok := k.(*StringKey)
	if !ok {
		return false
	}

	return s.str == other.str
}

func (s *StringKey) String() string {
	return s.str
}

// nilKey just used as place holder.
type nilKey uint8

func newNilKey() Key {
	return new(nilKey)
}

func (n *nilKey) Hash() int {
	return 0
}

func (n *nilKey) Equal(k Key) bool {
	_, ok := k.(*nilKey)
	return ok
}

func (n *nilKey) String() string {
	return ""
}
