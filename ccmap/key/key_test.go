package key

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringKeyHash(t *testing.T) {
	k := NewStringKey("cat")
	assert.Equal(t, 211780, k.Hash())
}

func TestStringKeyEqual(t *testing.T) {
	k := NewStringKey("cat")

	k2 := NewStringKey("cat")
	assert.Equal(t, true, k.Equal(k2))

	k3 := NewStringKey("cat1")
	assert.Equal(t, false, k.Equal(k3))
}

func TestStringKeyString(t *testing.T) {
	k := NewStringKey("s1")
	assert.Equal(t, "s1", k.String())
}
