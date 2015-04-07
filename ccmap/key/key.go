package key

import "github.com/csimplestring/go-concurrent-map/algo/hash"

// Key defines function set of a key in map.
type Key interface {
	Hash() int
	Equal(Key) bool
	String() string
}

// stringKey implements Key for a string.
type stringKey struct {
	h   int
	str string
}

// NewstringKey return a new string key.
func NewStringKey(str string) Key {
	return &stringKey{
		str: str,
		h:   (int)(hash.BKDRHash(str)),
	}
}

// Hash returns hash code for s.
func (s *stringKey) Hash() int {
	return s.h
}

// Equal checks if s == k.
func (s *stringKey) Equal(k Key) bool {
	other, ok := k.(*stringKey)
	if !ok {
		return false
	}

	return s.str == other.str
}

// String returns string representation of s.
func (s *stringKey) String() string {
	return s.str
}

// nilKey just used as place holder.
type nilKey uint8

// newNilKey new a nilKey.
func NewNilKey() Key {
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
