package faster

import (
	"log"
	"time"

	"github.com/mreithub/go-faster/faster/internal"
)

// Instance - Trackable instance
type Instance struct {
	parent    *Faster
	key       string
	startTime time.Time
}

// Deref -- Dereference an instance of 'key'
func (i *Instance) Deref() {
	if i.parent == nil {
		log.Print("GoFaster warning: possible double Deref()")
		return
	}

	var now time.Time
	var took time.Duration
	if !i.startTime.IsZero() {
		// only measure time if startTime was set
		now = time.Now()
		took = now.Sub(i.startTime)
	}

	i.parent.do(internal.EvDeref, i.key, took)
	i.parent = nil // prevent double Deref()
}
