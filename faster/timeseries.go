package faster

import (
	"time"
)

// TimeSeries -- excerpt from a History
type TimeSeries struct {
	// Path -- path of the time series
	Path []string
	// Data -- data over time (index 0 was taken at StartTS)
	Data []DataPoint
	// StartTS -- timestamp of the first Data point
	StartTS time.Time
	// Interval -- interval between Data points (note that )
	Interval time.Duration
}

// GetTimestamp -- returns the time.Time matching the Data point with the given index
func (s *TimeSeries) GetTimestamp(index int) time.Time {
	return s.StartTS.Add(time.Duration(index) * s.Interval)
}

// Relative -- returns a relative version of this time series
//
// This method simply subtracts each value from the one before it, resulting in a TimeSeries one shorter than this
//
// TimeSeries objects created by History.GetData() contain absolute (and monotonically increasing) values.
// This makes it easier to calculate aggregates over arbitrary time ranges.
//
// This function however returns a list of Snapshot where each item's value is the difference between two absolute Snapshot object
// (which is the kind of data you'd want to plot with charts).
//
// The returned list will be one shorter than this one (unless this one is empty or nil in which case nil is returned).
// Calling this method more than once might result in something resembling the second, third, ... derivatives
func (s TimeSeries) Relative() TimeSeries {
	if len(s.Data) == 0 {
		return s
	}

	var rc = TimeSeries{
		Data:     make([]DataPoint, 0, len(s.Data)-1),
		Interval: s.Interval,
		Path:     s.Path,
		StartTS:  s.StartTS, // TODO think about modifying the timestamps (i.e. startTS += interval/2)
	}

	for i := 1; i < len(s.Data); i++ {
		rc.Data = append(rc.Data, s.Data[i].Sub(s.Data[i-1]))
	}

	return rc
}
