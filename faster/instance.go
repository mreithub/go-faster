package faster

import (
	"log"
	"time"

	"github.com/mreithub/go-faster/faster/internal"
)

// Instance - Trackable instance
type Instance struct {
	parent    *Faster
	path      []string
	startTime time.Time
}

// Done -- Dereference an instance of 'key'
func (i *Instance) Done() {
	if i.parent == nil {
		log.Print("GoFaster warning: possible double Done()")
		return
	}

	var now time.Time
	var took time.Duration
	if !i.startTime.IsZero() {
		// only measure time if startTime was set
		now = time.Now()
		took = now.Sub(i.startTime)
	}

	i.parent.do(internal.EvDone, i.path, took)
	i.parent = nil // prevent double Done()
}
