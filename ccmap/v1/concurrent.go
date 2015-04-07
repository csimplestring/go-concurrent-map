package v1

import (
	"fmt"

	"github.com/csimplestring/go-concurrent-map/ccmap"
	. "github.com/csimplestring/go-concurrent-map/ccmap/key"
)

const (
	MAX_SEGMENTS = 65536
)

type concurrentHashMap struct {
	segmentShift uint
	segmentMask  int
	segments     []*hashMap

	stat map[int]int
}

func NewConcurrentMap(concurrencyLevel int) (ccmap.Map, error) {
	if concurrencyLevel > MAX_SEGMENTS {
		concurrencyLevel = MAX_SEGMENTS
	}

	sshift := 0
	ssize := 1
	for ssize < concurrencyLevel {
		sshift = sshift + 1
		ssize = ssize << 1
	}

	segmentShift := 32 - sshift
	segmentMask := ssize - 1

	var err error
	segments := make([]*hashMap, ssize)
	statMap := make(map[int]int)
	for i := 0; i < ssize; i++ {
		statMap[i] = 0
		segments[i], err = newHashMap(16)
		if err != nil {
			return nil, err
		}
	}

	return &concurrentHashMap{
		segmentMask:  segmentMask,
		segmentShift: (uint)(segmentShift),
		segments:     segments,
		stat:         statMap,
	}, nil
}

func (c *concurrentHashMap) hash(h int) int {
	h += (h << 15) ^ 0xffffcd7d
	h ^= (h >> 10)
	h += (h << 3)
	h ^= (h >> 6)
	h += (h << 2) + (h << 14)
	return h ^ (h >> 16)
}

func (c *concurrentHashMap) segmentFor(hash int) int {
	return (hash >> c.segmentShift) & c.segmentMask
}

func (c *concurrentHashMap) Put(key Key, val interface{}) bool {
	slot := c.segmentFor(c.hash(key.Hash()))
	c.stat[slot]++
	return c.segments[slot].Put(key, val)
}

func (c *concurrentHashMap) Get(key Key) (interface{}, bool) {
	return c.segments[c.segmentFor(c.hash(key.Hash()))].Get(key)
}

func (c *concurrentHashMap) Delete(key Key) bool {
	return false
}

func (c *concurrentHashMap) Stat() {
	fmt.Print(c.stat)
}
