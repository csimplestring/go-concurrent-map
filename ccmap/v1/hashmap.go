package v1

import (
	"sync"

	. "github.com/csimplestring/go-concurrent-map/ccmap/key"

	"github.com/csimplestring/go-concurrent-map/ccmap"
)

// hashMap is the default implementation of Map.
type hashMap struct {
	// -1: no rehash; otherwise it is rehashing
	rehashIdx int
	entryCnt  int
	tables    []*htable
	mutex     sync.RWMutex
}

func NewHashMap(size int) (ccmap.Map, error) {
	var err error

	tables := make([]*htable, 2)
	tables[1] = nil
	tables[0], err = newHtable(size)
	if err != nil {
		return nil, err
	}

	return &hashMap{
		entryCnt:  0,
		tables:    tables,
		rehashIdx: -1,
	}, nil
}

func newHashMap(size int) (*hashMap, error) {
	var err error

	tables := make([]*htable, 2)
	tables[1] = nil
	tables[0], err = newHtable(size)
	if err != nil {
		return nil, err
	}

	return &hashMap{
		entryCnt:  0,
		tables:    tables,
		rehashIdx: -1,
	}, nil
}

// Put puts <key, val> pair in correct slot.
// It returns true if succeed; otherwise false.
func (h *hashMap) Put(key Key, val interface{}) bool {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	entry := newEntry(key, val)
	if h.rehashIdx == -1 {
		if ok := h.putEntry(0, entry); !ok {
			return false
		}

		if h.entryCnt > h.tables[0].size {
			h.beginRehash()
		}
		return true
	}

	if ok := h.putEntry(1, entry); !ok {
		return false
	}

	h.rehash()
	return true
}

// Get gets the value based on key.
// If value exists, it returns value and TRUE;
// otherwise it returns nil and FALSE.
func (h *hashMap) Get(key Key) (interface{}, bool) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if h.rehashIdx != -1 {
		if en, ok := h.tables[1].get(key); ok {
			h.rehash()
			return en.Value(), true
		}
	}

	if en, ok := h.tables[0].get(key); ok {
		return en.Value(), true
	}
	return nil, false
}

// Delete deletes value based on key.
// It returns TRUE if key exists; otherise FALSE.
func (h *hashMap) Delete(key Key) bool {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	deleted := 0
	_, cnt := h.tables[0].delete(key)
	deleted += cnt

	if h.rehashIdx != -1 {
		_, cnt := h.tables[1].delete(key)
		deleted += cnt
		h.rehash()
	}

	h.entryCnt -= deleted

	if deleted > 0 {
		return true
	}
	return false
}

// beginRehash sets rehashIdx to be 0, creates new htable for
// tables[1].
func (h *hashMap) beginRehash() {
	h.rehashIdx = 0
	newSize := h.tables[0].size * 2
	h.tables[1], _ = newHtable(newSize)
}

// stopRehash switches old and new htable internally, resets
// rehashIdx to be -1.
func (h *hashMap) stopRehash() {
	h.tables[0] = h.tables[1]
	h.tables[1] = nil
	h.rehashIdx = -1
}

// rehash moves tables[0]'s entry to tables[1]
func (h *hashMap) rehash() {
	// find the non-empty bucket
	b := h.tables[0].buckets[h.rehashIdx]
	for b.Size() == 0 && h.rehashIdx < h.tables[0].size {
		b = h.tables[0].buckets[h.rehashIdx]
		h.rehashIdx++
	}

	// move old entries
	for en, ok := b.Pop(); ok; en, ok = b.Pop() {
		h.tables[1].push(en)
	}

	// rehash ends
	if h.rehashIdx == h.tables[0].size {
		h.stopRehash()
	}
}

// putEntry puts en into tables[tableIdx].
// It returns true if succeeds, otherwise false.
func (h *hashMap) putEntry(tableIdx int, en Entry) bool {
	ok := h.tables[tableIdx].put(en)
	if !ok {
		return false
	}

	h.entryCnt++
	return true
}
