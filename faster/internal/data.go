package internal

import "time"

// Data -- internal GoFaster data structure - thread unsafe (only for use in the Faster.run() goroutine)
type Data struct {
	// currently active invocations
	Active int32
	// number of finished invocations
	Count int64
	// time spent in those invocations (in nanoseconds)
	TotalTime time.Duration

	// child instances
	Children map[string]*Data
}

// Average -- returns the average time spent in each invocation
func (d *Data) Average() time.Duration {
	var rc time.Duration
	if d.Count > 0 {
		rc = d.TotalTime / time.Duration(d.Count)
	}
	return rc
}

// GetChild -- Returns the Data instance matching the given path (recursively)
//
// will create new children (also recursively) if not found
// an empty key results in the current object to be returned
//
// This method is thread unsafe! It is assumed that it's only ever accessed from Faster.run() (which runs in its own goroutine)
func (d *Data) GetChild(key ...string) *Data {
	if len(key) == 0 {
		return d
	}

	if d.Children == nil {
		d.Children = make(map[string]*Data)
	}

	var head, tail = key[0], key[1:]
	var child *Data
	var ok bool
	if child, ok = d.Children[head]; !ok {
		child = new(Data)
		d.Children[head] = child
	}
	return child.GetChild(tail...)
}
