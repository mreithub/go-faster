package internal

import "time"

// Data -- internal GoFaster data structure
type Data struct {
	// currently active invocations
	Active int32
	// number of finished invocations
	Count int64
	// time spent in those invocations (in nanoseconds)
	TotalTime time.Duration
}
