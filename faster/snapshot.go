package faster

import "time"

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
}

// Keys -- List all keys of this read-only instance
func (s *Snapshot) Keys() []string {
	rc := make([]string, 0, len(s.Children))

	for k := range s.Children {
		rc = append(rc, k)
	}

	return rc
}
