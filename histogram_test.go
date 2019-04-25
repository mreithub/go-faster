package goref

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const NS = time.Nanosecond
const US = time.Microsecond

func TestBuckets(t *testing.T) {
	var h = newHistogram()
	assert.Equal(t, 0, h.getBucket(-1234*NS), "negative values should end up in bucket 0")
	assert.Equal(t, 0, h.getBucket(0))

	assert.Equal(t, 1, h.getBucket(1*NS))
	assert.Equal(t, 2, h.getBucket(2*NS))
	assert.Equal(t, 2, h.getBucket(3*NS))
	assert.Equal(t, 3, h.getBucket(4*NS))
	assert.Equal(t, 3, h.getBucket(7*NS))
	assert.Equal(t, 4, h.getBucket(8*NS))
	assert.Equal(t, 4, h.getBucket(15*NS))
	assert.Equal(t, 5, h.getBucket(16*NS))
	assert.Equal(t, 5, h.getBucket(31*NS))
	assert.Equal(t, 7, h.getBucket(127*NS))
	assert.Equal(t, 8, h.getBucket(128*NS))
	assert.Equal(t, 10, h.getBucket(1023*NS))
	assert.Equal(t, 11, h.getBucket(1024*NS))

	assert.Equal(t, 0, h.getBucket(-3*US))
	assert.Equal(t, 10, h.getBucket(1*US))
	assert.Equal(t, 11, h.getBucket(2*US))
	assert.Equal(t, 12, h.getBucket(3*US))
	assert.Equal(t, 12, h.getBucket(4*US))
}
