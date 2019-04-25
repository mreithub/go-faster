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

	GetValues() ([]time.Duration, []int)
	GetPercentiles(values ...int) []time.Duration

	// Since -- subtracts the values of both histograms and returns the difference as new
	// object (will return nil if the histograms are incompatible (e.g. different resolution))
	Since(other Histogram) Histogram
}

type histogram struct {
	count int64
	sum   time.Duration

	buckets [64]int
}

// NewHistogram -- returns a Histogram instance
func newHistogram() *histogram {
	var rc = histogram{}
	return &rc
}

func (h *histogram) Average() time.Duration { return h.sum / time.Duration(h.count) }
func (h *histogram) Count() int64           { return h.count }
func (h *histogram) Sum() time.Duration     { return h.sum }

// getBucket -- returns the right bucket for the given value (0 if value <= 0)
func (h *histogram) getBucket(value time.Duration) int {
	var rc = 0

	for value >= 0x100 {
		value >>= 8
		rc += 8
	}

	for value > 0 {
		value >>= 1
		rc++
	}

	// in theory this should never happen (rc should never be greater than 63)
	if rc >= len(h.buckets) {
		rc = len(h.buckets) - 1
	}
	return rc
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
	var _, otherValues = other.GetValues()

	var rc = histogram{
		count: h.count - other.Count(),
		sum:   h.sum - other.Sum(),
	}

	for i, myVal := range h.buckets {
		rc.buckets[i] = myVal - otherValues[i]
	}

	return &rc
}
