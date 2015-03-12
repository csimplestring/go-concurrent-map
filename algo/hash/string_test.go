package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBKDRHash(t *testing.T) {
	assert.Equal(t, 211780, BKDRHash("cat"))
}
