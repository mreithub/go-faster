package goref

import "time"

// Histogram -- Keeps track of time.Duration values and their distribution
//
// The default Histogram implementation does no locking (but all instances returned
// by goref should be considered immutable)
type Histogram interface {
	Count() int64
	Sum() time.Duration
	Average() time.Duration
	Resolution() time.Duration

	GetBuckets() ([]time.Duration, []int)
	GetPercentiles(values ...int) []time.Duration

	// Since -- subtracts the values of both histograms and returns the difference as new
	// object (will return nil if the histograms are incompatible (e.g. different resolution))
	Since(other Histogram) Histogram
}

type histogram struct {
	resolution time.Duration

	count int64
	sum   time.Duration

	buckets []int
}

// NewHistogram -- returns a Histogram instance
func newHistogram(resolution time.Duration, buckets int) *histogram {
	var rc = histogram{
		resolution: resolution,
		buckets:    make([]int, buckets), // all initialized with 0
	}
	return &rc
}

func (h *histogram) Average() time.Duration { return h.sum / time.Duration(h.count) }
func (h *histogram) Count() int64           { return h.count }
func (h *histogram) Sum() time.Duration     { return h.sum }

// getBucket -- returns the right bucket for the given value
// (-1 if value < h.resolution, >=len(h.buckets) if out of bounds)
func (h *histogram) getBucket(value time.Duration) int {
	value /= h.resolution
	var rc = 0

	for value >= 0x100 {
		value >>= 8
		rc += 8
	}

	for value > 0 {
		value >>= 1
		rc++
	}

	return rc - 1
}

func (h *histogram) Add(value time.Duration) {
	h.sum += value
	h.count++

	var bucket = h.getBucket(value)
	if bucket < 0 || bucket >= len(h.buckets) {
		return // out of bounds -> won't be added to the histogram
	}

	h.buckets[bucket]++
}

func (h *histogram) Since(other Histogram) *histogram {
	if other.Resolution() != h.resolution {
		return nil
	}

	var _, otherValues = other.GetBuckets()
	if len(otherValues) != len(h.buckets) {
		return nil
	}

	var rc = histogram{
		buckets:    make([]int, 0, len(otherValues)),
		count:      h.count - other.Count(),
		resolution: h.resolution,
		sum:        h.sum - other.Sum(),
	}

	for i, myVal := range h.buckets {
		rc.buckets = append(rc.buckets, myVal-otherValues[i])
	}

	return &rc
}
