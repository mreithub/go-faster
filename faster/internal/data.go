package internal

import (
	"time"
)

// Data -- internal GoFaster data structure - thread unsafe (only for use in the Faster.run() goroutine)
type Data struct {
	// currently active invocations
	Active int32
	// number of finished invocations
	Count int64
	// time spent in those invocations (in nanoseconds)
	TotalTime time.Duration
}

// Average -- returns the average time spent in each invocation
func (d *Data) Average() time.Duration {
	var rc time.Duration
	if d.Count > 0 {
		rc = d.TotalTime / time.Duration(d.Count)
	}
	return rc
}

// Done -- caused by Tracker.Done() (and called by Faster.run())
func (d *Data) Done(took time.Duration) {
	d.Active--
	d.Count++
	d.TotalTime += took
}

// Sub -- returns the difference between the two given Data objects (assuming 'this' is the newer one)
func (d *Data) Sub(other *Data) Data {
	if other == nil {
		return *d
	}
	return Data{
		Active:    0, // doesn't really make sense (maybe we should use max(this, other))
		Count:     d.Count - other.Count,
		TotalTime: d.TotalTime - other.TotalTime,
	}
}
