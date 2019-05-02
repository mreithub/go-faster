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

// Get -- Recursively traverse this Snapshot instance returning the entry matching the given key (or nil if not found)
func (s *Snapshot) Get(key ...string) *Snapshot {
	if len(key) == 0 {
		return s
	}
	var head, tail = key[0], key[1:]
	if child, ok := s.Children[head]; ok {
		return child.Get(tail...)
	} else {
		return nil // not found
	}
}

// Keys -- List all keys of this read-only instance
func (s *Snapshot) Keys() []string {
	rc := make([]string, 0, len(s.Children))

	for k := range s.Children {
		rc = append(rc, k)
	}

	return rc
}
