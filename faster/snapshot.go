package faster

import (
	"time"

	"github.com/mreithub/go-faster/histogram"
)

// Snapshot -- point-in-time copy of a GoFaster instance
type Snapshot struct {
	// Currently active invocations
	Active int32 `json:"active,omitempty"`

	// Total number of (finished) invocations
	Count int64 `json:"count,omitempty"`

	// Total time spent
	Duration time.Duration `json:"duration,omitempty"`

	// Computed average run time, provided for convenience
	Average time.Duration `json:"average,omitempty"`

	// Child GoFaster instance data
	Children map[string]*Snapshot `json:"_children,omitempty"`

	// Creation timestamp
	Ts time.Time `json:"ts"`

	Histogram *histogram.Histogram
}

// Snapshots -- list of Snapshot (adding some useful helper functions)
type Snapshots []*Snapshot

// Get -- Recursively traverse this Snapshot instance returning the entry matching the given key (or nil if not found)
func (s *Snapshot) Get(key ...string) *Snapshot {
	if len(key) == 0 {
		return s
	}
	var head, tail = key[0], key[1:]
	if child, ok := s.Children[head]; ok {
		return child.Get(tail...)
	}

	return nil // not found
}

// Keys -- List all keys of this read-only instance
func (s *Snapshot) Keys() []string {
	rc := make([]string, 0, len(s.Children))

	for k := range s.Children {
		rc = append(rc, k)
	}

	return rc
}

// Sub -- returns the difference between the two given snapshots (assuming 'other' is the newer one)
func (s *Snapshot) Sub(other *Snapshot) *Snapshot {
	if other == nil {
		return s
	}
	var rc = Snapshot{
		Active:   0,   // doesn't really make sense (maybe we should use max(this, other))
		Children: nil, // TODO should we do this recursively?
		Count:    other.Count - s.Count,
		Duration: other.Duration - s.Duration,
		Ts:       s.Ts,
	}
	if rc.Count > 0 {
		rc.Average = rc.Duration / time.Duration(rc.Count)
	}

	return &rc
}

// Relative -- Returns a new list of Snapshots - subtracting each entry's values from the last one's
//
// go-faster Snapshot objects usually store absolute (and monotonically increasing) counts. This makes it easier to calculate aggregates
// over arbitrary time ranges.
//
// This function however returns a list of Snapshot where each item's value is the difference between two absolute Snapshot object
// (which is the kind of data you'd want to plot with charts).
//
// The returned list will be one shorter than this one (unless this one is empty or nil in which case nil is returned).
// Calling this method more than once might result in something resembling the second, third, ... derivatives
func (s Snapshots) Relative() Snapshots {
	if len(s) == 0 {
		return nil
	}
	var rc = make([]*Snapshot, 0, len(s)-1)

	var last *Snapshot
	for _, snap := range s {
		if last != nil {
			rc = append(rc, last.Sub(snap))
		}
		last = snap
	}

	return rc
}
