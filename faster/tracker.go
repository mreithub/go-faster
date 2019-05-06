package faster

import (
	"time"

	"github.com/mreithub/go-faster/faster/internal"
)

// Tracker - Trackable instance
//
// Note that this struct will only work as expected when it has a backing Faster instance
// (i.e. is acquired by calling Faster.Track() or NewChild())
type Tracker struct {
	parent  *Faster
	path    []string
	startTS time.Time
	took    time.Duration
}

// Done -- Dereference an instance of 'key'
//
// it is safe to call Done() more than once (from the same goroutine - this struct is NOT thread safe)
func (t *Tracker) Done() {
	if t.parent == nil {
		//log.Print("go-faster warning: possible double Done()")
		return
	}

	var now time.Time
	var took time.Duration
	if !t.startTS.IsZero() {
		// only measure time if startTime was set
		now = time.Now()
		took = now.Sub(t.startTS)
	}
	t.took = took

	t.parent.do(internal.EvDone, t.path, took)
	t.parent = nil // prevent double Done()
}

// NewChild -- creates a child with the same startTS and backing Faster instance but different path
//
// won't work after Done() was called on this object (will return nil)
//
// Make sure to call Done() on each child you create here
func (t *Tracker) NewChild(path ...string) *Tracker {
	if t.parent == nil {
		return nil
	}

	var rc = Tracker{
		parent:  t.parent,
		path:    append(t.path, path...),
		startTS: t.startTS,
	}
	t.parent.do(internal.EvTrack, rc.path, 0)

	return &rc
}

// Path -- returns the Faster path this Tracker object is bound to
func (t *Tracker) Path() []string {
	return t.path
}

// StartTS -- Returns the timestamp of this Tracker object's creation
func (t *Tracker) StartTS() time.Time {
	return t.startTS
}

// Took -- returns the time between StartTS() and the (first) call to Done()
//
// before Done() is called, this getter will return 0
func (t *Tracker) Took() time.Duration {
	return t.took
}
