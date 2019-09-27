package faster

import (
	"time"

	"github.com/mreithub/go-faster/faster/internal"
	"github.com/mreithub/go-faster/histogram"
)

// Snapshot -- point-in-time copy of go-faster's state
type Snapshot struct {
	tree       *internal.Tree
	data       []internal.Data
	histograms []histogram.Histogram

	// Creation timestamp
	TS time.Time `json:"ts"`
}

// Snapshots -- list of Snapshot (adding some useful helper functions)
type Snapshots []*Snapshot

// Children -- returns the names of the direct children of the given path
func (s *Snapshot) Children(path ...string) []string {
	return s.tree.Children(path...)
}

// Get -- Return the entry matching the given key (or nil if not found)
func (s *Snapshot) Get(path ...string) *internal.Data {
	var rc *internal.Data
	if s.tree == nil {
		return nil
	}

	var index = s.tree.GetIndex(path...)
	if index >= 0 && index < len(s.data) {
		rc = &s.data[index]
	}
	return rc
}
