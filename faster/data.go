package faster

import "time"

// Data -- Reference counter Snapshot data
type Data struct {
	// Currently active invocations
	Active int32 `json:"active"`

	// Total number of (finished) invocations
	Count int64 `json:"count"`

	// Total time spent
	Duration time.Duration `json:"duration"`

	// Computed average run time, provided for convenience
	Average time.Duration `json:"average"`
}

// Fills a Data object with the values from an (internal) data object
//
// Copies all the duplicate fields over and calculates the convenience fields.
func newData(src *data) Data {
	var average time.Duration
	if src.count > 0 {
		average = src.totalTime / time.Duration(src.count)
	}

	return Data{
		Active:   src.active,
		Count:    src.count,
		Duration: time.Duration(src.totalTime),
		Average:  average,
	}
}
