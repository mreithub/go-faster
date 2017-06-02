package goref

import "time"

// Snapshot -- point-in-time copy of a GoRef instance
type Snapshot struct {
	// Child GoRef instance data
	Children map[string]Snapshot `json:"_children,omitempty"`

	// Snapshot data
	Data map[string]Data `json:"data,omitempty"`

	// Creation timestamp
	Ts time.Time `json:"ts"`
}

// Keys -- List all keys of this read-only instance
func (s *Snapshot) Keys() []string {
	rc := make([]string, 0, len(s.Data))

	for k := range s.Data {
		rc = append(rc, k)
	}

	return rc
}
