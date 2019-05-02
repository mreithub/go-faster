package faster

import (
	"container/list"
	"sync"
	"time"
)

// History -- periodically takes snapshots of GoFaster instances
//
// All methods are thread safe
type History struct {
	ticker *time.Ticker
	done   chan struct{}

	Name      string
	Capacity  int
	entries   *list.List
	entryLock sync.RWMutex
}

// Push -- Add a new value to the end of Periodic.Entries (removing the first ones if we've exceeded our capacity)
func (h *History) Push(value interface{}) {
	h.entryLock.Lock()
	defer h.entryLock.Unlock()

	var entries = h.entries
	entries.PushBack(value)
	for h.entries.Len() > h.Capacity && entries.Len() > 0 {
		entries.Remove(h.entries.Front())
	}
}

// Stop -- stop the underlying time.Ticker
func (h *History) Stop() {
	h.ticker.Stop()
	h.done <- struct{}{}
}

// NewHistory -- creates a History instance and initialize it as requested
//
// after each interval, it'll push a reference to itself to tickChannel
// (indicating to the underlying Faster instance that it should take another
// snapshot and Push() it to this History instance)
func NewHistory(name string, interval time.Duration, keep int, tickChannel chan *History) *History {
	var rc = History{
		ticker:   time.NewTicker(interval),
		Name:     name,
		Capacity: keep,
		entries:  list.New(),
	}

	go func() {
		for {
			select {
			case <-rc.ticker.C:
				tickChannel <- &rc
			case <-rc.done:
				break
			}
		}
	}()

	return &rc
}
