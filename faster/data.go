package faster

import (
	"time"
)

// DataPoint -- read only go-faster data point
type DataPoint interface {
	Active() int32
	Count() int64
	TotalTime() time.Duration
	Average() time.Duration

	Sub(other DataPoint) DataPoint
}

// data -- internal go-faster data structure - thread unsafe (only the go-faster goroutine does writes and creates read-only copies)
type data struct {
	// currently active invocations
	active int32
	// number of finished invocations
	count int64
	// time spent in those invocations (in nanoseconds)
	totalTime time.Duration
}

func (d *data) Active() int32            { return d.active }
func (d *data) Count() int64             { return d.count }
func (d *data) TotalTime() time.Duration { return d.totalTime }

// Average -- returns the average time spent in each invocation
func (d *data) Average() time.Duration {
	var rc time.Duration
	if d.count > 0 {
		rc = d.totalTime / time.Duration(d.count)
	}
	return rc
}

// Done -- caused by Tracker.Done() (and called by Faster.run())
func (d *data) Done(took time.Duration) {
	d.active--
	d.count++
	d.totalTime += took
}

// Sub -- returns the difference between the two given Data objects (assuming 'this' is the newer one)
func (d *data) Sub(other DataPoint) DataPoint {
	if other == nil {
		return d
	}
	return &data{
		active:    0, // doesn't really make sense (maybe we should use max(this, other))
		count:     d.count - other.Count(),
		totalTime: d.totalTime - other.TotalTime(),
	}
}
