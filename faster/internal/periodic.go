package internal

import (
	"container/list"
	"time"
)

// Periodic -- sets up a ticker that emits periodic events to tickChannel
//
//
type Periodic struct {
	ticker *time.Ticker
	done   chan struct{}

	Name     string
	Capacity int
	entries  *list.List
}

// Push -- Add a new value to the end of Periodic.Entries (removing the first ones if we've exceeded our capacity)
func (p *Periodic) Push(value interface{}) {
	p.entries.PushBack(value)
	for p.entries.Len() > p.Capacity && p.entries.Len() > 0 {
		p.entries.Remove(p.entries.Front())
	}
}

// Stop -- stop the underlying time.Ticker
func (p *Periodic) Stop() {
	p.ticker.Stop()
	p.done <- struct{}{}
}

// NewPeriodic -- creates a Periodic instance and initialize it as requested
func NewPeriodic(name string, interval time.Duration, keep int, tickChannel chan *Periodic) *Periodic {
	var rc = Periodic{
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
