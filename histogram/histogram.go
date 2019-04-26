package histogram

import "time"

// Histogram -- Keeps track of time.Duration values and their distribution
//
// The default Histogram implementation does no locking (but all instances returned
// by goref should be considered immutable)
type Histogram interface {
	Count() int64
	Sum() time.Duration
	Average() time.Duration
	Median() time.Duration

	// GetValues -- returns each bucket's lower bound and its value count
	//
	// will skip empty buckets at the start and end
	GetValues() ([]time.Duration, []int)

	// GetPercentiles -- returns estimates for the percentile values in question
	//
	// Make sure to order 'values' - GetPercentiles() will only traverse the value buckets once
	GetPercentiles(values ...int) []time.Duration

	// Since -- subtracts the values of both histograms and returns the difference as new
	// object (will return nil if the histograms are incompatible (e.g. different resolution))
	Since(other Histogram) Histogram
}

type histogram struct {
	count int64
	sum   time.Duration

	// buckets[0]
	buckets [64]int
}

func (h *histogram) Average() time.Duration { return h.sum / time.Duration(h.count) }
func (h *histogram) Count() int64           { return h.count }
func (h *histogram) Sum() time.Duration     { return h.sum }
func (h *histogram) Median() time.Duration  { return h.GetPercentiles(50)[0] }

// getBucket -- returns the right bucket for the given value (0 if value <= 0, 1..63 for anything else)
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

	// in theory this should never happen (value is a signed 64bit integer so rc should never be greater than 63)
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

func (h *histogram) GetPercentile(value int) time.Duration {
	return h.GetPercentiles(value)[0]
}

func (h *histogram) GetPercentiles(values ...int) []time.Duration {
	var rc = make([]time.Duration, 0, len(values))

	var buckets = h.buckets[:]
	var currentBucket time.Duration = 1

	var total int64
	for _, percentile := range values {
		if percentile < 0 {
			percentile = 0
		}
		if percentile > 100 {
			percentile = 100
		}

		var cutoffValue = h.count * int64(percentile) / 100
		//log.Printf("cutoff for %d%%: %d (count=%d)", percentile, cutoffValue, h.count)

		// the inner loop will only traverse each bucket once
		for total < cutoffValue && len(buckets) > 0 {
			var count = int64(buckets[0])
			buckets = buckets[1:]
			total += count
			currentBucket <<= 1
		}

		//log.Printf("- total: %d, currentBucket: %s", total, currentBucket)

		// TODO do some interpolation
		rc = append(rc, currentBucket>>1)
	}

	return rc
}

func (h *histogram) GetValues() ([]time.Duration, []int) {
	var skipFrom, skipTo = 0, len(h.buckets)

	for i, value := range h.buckets {
		skipFrom = i
		if value > 0 {
			break
		}
	}

	for i := len(h.buckets) - 1; i >= skipFrom; i-- {
		if h.buckets[i] > 0 {
			break
		}
		skipTo = i
	}

	return minValues[skipFrom:skipTo], h.buckets[skipFrom:skipTo]
}

func (h *histogram) Since(other Histogram) Histogram {
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

// minValues -- static list containing each bucket's lower bound
var minValues = func() [64]time.Duration {
	var rc [64]time.Duration
	var v time.Duration = 1
	for i := 1; i < len(rc); i++ {
		rc[i] = v
		v <<= 1
	}
	return rc
}()
