package faster

import (
	"log"
	"time"

	"github.com/mreithub/go-faster/faster/internal"
)

// Tracker - Trackable instance
type Tracker struct {
	parent    *Faster
	path      []string
	startTime time.Time
}

// Done -- Dereference an instance of 'key'
func (t *Tracker) Done() {
	if t.parent == nil {
		log.Print("go-faster warning: possible double Done()")
		return
	}

	var now time.Time
	var took time.Duration
	if !t.startTime.IsZero() {
		// only measure time if startTime was set
		now = time.Now()
		took = now.Sub(t.startTime)
	}

	t.parent.do(internal.EvDone, t.path, took)
	t.parent = nil // prevent double Done()
}
