package ccmap

import "github.com/csimplestring/go-concurrent-map/ccmap/key"

const (
	BUCKET_SIZE_DEFAULT = 16
	StatusRehashing     = -1
)

// Map defines the functions that a map should support
// TODO:
// 2. shrink table
// 3. rehash table in background
type Map interface {
	Put(k key.Key, val interface{}) bool
	Get(k key.Key) (interface{}, bool)
	Delete(k key.Key) bool
}
