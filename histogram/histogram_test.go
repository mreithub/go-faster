package histogram

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const NS = time.Nanosecond
const US = time.Microsecond

func TestGetBucket(t *testing.T) {
	var h Histogram

	// <= 0 go into bucket 0
	assert.Equal(t, 0, h.getBucket(-1234*NS), "negative values should end up in bucket 0")
	assert.Equal(t, 0, h.getBucket(0))

	// only '1' goes into bucket 1
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

	assert.Equal(t, 63, h.getBucket(math.MaxInt64))
	var maxPlus1 time.Duration = math.MaxInt64
	maxPlus1++
	assert.Equal(t, 0, h.getBucket(maxPlus1))
}

func TestEmptyBuckets(t *testing.T) {
	var h Histogram

	assert.Equal(t, []time.Duration{}, h.GetPercentiles())
	assert.Equal(t, []time.Duration{0, 0, 0, 0, 0, 0, 0}, h.GetPercentiles(-88, 0, 25, 50, 75, 100, 200))

	var buckets, values = h.GetValues()
	assert.Empty(t, buckets)
	assert.Empty(t, values)
}

func TestMinValues(t *testing.T) {
	var h Histogram

	for i, d := range minValues {
		if i > 0 {
			assert.Equal(t, i-1, h.getBucket(d-1), fmt.Sprintf("for value #%d: %d-1", i, d))
		} else {
			assert.Equal(t, 0, h.getBucket(d-1))
		}
		assert.Equal(t, i, h.getBucket(d), fmt.Sprintf("for value #%d: %d", i, d))
	}
}

func TestHistogram(t *testing.T) {
	var h = Histogram{
		buckets: [64]int32{1, 1, 1, 1, 1, 1, 1, 1},
		count:   8,
		sum:     (1 + 2 + 4 + 8 + 16 + 32 + 64 + 128) * NS,
	}

	assert.EqualValues(t, []time.Duration{0, 2 * NS, 8 * NS, 32 * NS, 128 * NS}, h.GetPercentiles(0, 25, 50, 75, 100))

	h = Histogram{
		buckets: [64]int32{0, 0, 0, 3, 2, 11, 8, 16, 24, 8, 2, 0},
		count:   3 + 2 + 11 + 8 + 16 + 24 + 8 + 2,
	}

	var _, values = h.GetValues()
	assert.EqualValues(t, h.buckets[3:11], values)

	//assert.EqualValues(t, []time.Duration{4 * NS, 512 * NS}, h.GetPercentiles(0, 100))
}
