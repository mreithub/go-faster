package faster

import (
	"time"

	"github.com/mreithub/go-faster/faster/internal"
)

// Snapshot -- point-in-time copy of go-faster's state
type Snapshot struct {
	tree       *internal.Tree
	data       []data
	histograms []Histogram

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
func (s *Snapshot) Get(path ...string) DataPoint {
	var rc DataPoint
	if s.tree == nil {
		return nil
	}

	var index = s.tree.GetIndex(path...)
	if index >= 0 && index < len(s.data) {
		rc = &s.data[index]
	}
	return rc
}

// GetHistogram -- returns the histogram for the given key (or nil if not found/disabled)
func (s *Snapshot) GetHistogram(path ...string) *Histogram {
	var rc *Histogram
	if s.tree == nil || len(s.histograms) == 0 {
		return nil
	}

	var index = s.tree.GetIndex(path...)
	if index >= 0 && index < len(s.histograms) {
		rc = &s.histograms[index]
	}
	return rc
}
