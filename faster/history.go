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

	Name     string
	Capacity int

	interval time.Duration
	entries  *list.List
	// guards History.entries (but not its (immutable) data)
	entryLock sync.RWMutex
}

// Duration -- returns .Capacity * .Interval()
//
// to get the actual time span stored in this History object, use time.Now().Sub(h.FirstTS())
func (h *History) Duration() time.Duration {
	return h.interval * time.Duration(h.Capacity)
}

// First -- returns a reference to the oldest Snapshot entry of this History
// instance (or nil if empty)
func (h *History) First() *Snapshot {
	h.entryLock.Lock()
	defer h.entryLock.Unlock()

	if e := h.entries.Front(); e != nil {
		if s, ok := e.Value.(*Snapshot); ok {
			return s
		}
	}
	return nil
}

// FirstTS -- convenience wrapper around First() returning that snapshot's timestamp (or a .IsZero() one)
func (h *History) FirstTS() time.Time {
	var rc time.Time
	if s := h.First(); s != nil {
		rc = s.TS
	}
	return rc
}

// GetData -- returns the TimeSeries for the given key
func (h *History) GetData(path ...string) TimeSeries {
	var snapshots = h.List()
	var rc = TimeSeries{
		Path:     path,
		Data:     make([]DataPoint, 0, len(snapshots)),
		StartTS:  h.FirstTS(),
		Interval: h.Interval(),
	}

	for _, snapshot := range snapshots {
		if d := snapshot.Get(path...); d != nil {
			rc.Data = append(rc.Data, d)
		}
	}

	return rc
}

// Interval -- returns the interval of the internal History Ticker
func (h *History) Interval() time.Duration {
	// we use a getter here to prevent the user from changing this (which wouldn't affect the internal ticker but might lead to odd behaviour)
	return h.interval
}

// Len -- returns the number of entries currently stored in this History struct (thread safe)
func (h *History) Len() int {
	h.entryLock.Lock()
	defer h.entryLock.Unlock()
	return h.entries.Len()
}

// List -- Returns all the Snapshots stored in this History object
//
// This method is thread safe - but the stored Snapshot values are presumed
// immutable (and therefore not guarded by any locks)
func (h *History) List() Snapshots {
	h.entryLock.Lock()
	defer h.entryLock.Unlock()

	var rc = make([]*Snapshot, 0, h.entries.Len())
	for e := h.entries.Front(); e != nil; e = e.Next() {
		if s, ok := e.Value.(*Snapshot); ok {
			rc = append(rc, s)
		}
	}

	return rc
}

// push -- Storing a new Snapshot in this History object - making sure we don't
// exceed our Capacity) (thread safe)
func (h *History) push(snapshot *Snapshot) {
	h.entryLock.Lock()
	defer h.entryLock.Unlock()

	var entries = h.entries
	entries.PushBack(snapshot)

	// Note: h.Capacity isn't guarded and may be changed by another goroutine while doing this
	for entries.Len() > h.Capacity+1 && entries.Len() > 0 {
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
		interval: interval,
		entries:  list.New(),
	}

	// get initial snapshot
	tickChannel <- &rc

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
