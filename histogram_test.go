package goref

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const NS = time.Nanosecond
const US = time.Microsecond

func TestBuckets(t *testing.T) {
	var values = map[time.Duration]int{
		// everything lower than h.resolution should end up in bucket 0
		999 * time.Nanosecond: -1,

		1 * US:   0,
		2 * US:   1,
		3 * US:   1,
		4 * US:   2,
		7 * US:   2,
		8 * US:   3,
		15 * US:  3,
		16 * US:  4,
		31 * US:  4,
		127 * US: 6,
		128 * US: 7,
	}

	var h = newHistogram(US, 8)
	for value, expected := range values {
		var actual = h.getBucket(value)
		assert.Equal(t, expected, actual, fmt.Sprintf("for value %dns", value/time.Nanosecond))
	}
}
