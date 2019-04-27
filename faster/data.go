package faster

import (
	"time"

	"github.com/mreithub/go-faster/faster/internal"
)

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
func newData(src *internal.Data) Data {
	var average time.Duration
	if src.Count > 0 {
		average = src.TotalTime / time.Duration(src.Count)
	}

	return Data{
		Active:   src.Active,
		Count:    src.Count,
		Duration: time.Duration(src.TotalTime),
		Average:  average,
	}
}
