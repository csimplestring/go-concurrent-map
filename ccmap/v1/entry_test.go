package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/csimplestring/go-concurrent-map/ccmap/key"
)

func TestNewEntry(t *testing.T) {
	e := newEntry(nil, nil)
	assert.Nil(t, e.Key())
	assert.Nil(t, e.Value())
}

func TestEntryString(t *testing.T) {
	e := &entry{
		k: NewStringKey("k1"),
		v: 1,
	}
	assert.Equal(t, "[k1 1]", e.String())
}
