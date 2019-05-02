package faster

import (
	"sync"
	"time"

	"github.com/mreithub/go-faster/faster/internal"
)

// TODO tracking execution time might cause performance issues (e.g. in virtualized environments gettimeofday() might be slow)
//   if that turns out to be the case, deactivate Data.TotalNsec

// Faster -- A simple, go-style key-based reference counter that can be used for profiling your application (main class)
type Faster struct {
	// internal data structure -- only ever accessed from within the run() goroutine // TODO check if this is true
	root *internal.Data

	// processed by the run() goroutine
	evChannel chan internal.Event

	// calling GetSnapshot() triggers an internal.EvSnapshot which in turn causes
	// the run() goroutine to take one and push it here -- while it's not guaranteed
	// that each GetSnapshot() call will read the response to its own request
	// (multiple GetSnapshot() calls might block concurrently), that doesn't really
	// change the results of each call
	snapshotChannel chan *Snapshot

	// periodic snapshots
	history map[string]*History
	// guards the history map
	historyLock sync.Mutex
	// indicates a History ticker expired (and expects to be sent a new Snapshot)
	tickChan chan *History
}

func (g *Faster) do(evType internal.EventType, path []string, took time.Duration) {
	g.evChannel <- internal.Event{
		Type: evType,
		Path: path,
		Took: took,
	}
}

// SetTicker -- Sets up periodic snapshots
//
// - name is the unique name of the given History ticker
// - interval specifies how often these ticks should happen (if 0, the ticker will be deleted)
// - keep is the number of past Snapshots stored for the given History object
func (g *Faster) SetTicker(name string, interval time.Duration, keep int) {
	g.historyLock.Lock()
	defer g.historyLock.Unlock()

	if ticker, ok := g.history[name]; ok {
		// replacing/removing an existing ticker -> stop the old one
		ticker.Stop()
	}

	if interval == 0 {
		return
	}

	g.history[name] = NewHistory(name, interval, keep, g.tickChan)
}

// Track -- Tracks an instance of 'key'
func (g *Faster) Track(key ...string) *Instance {
	g.do(internal.EvTrack, key, 0)

	return &Instance{
		parent:    g,
		path:      key,
		startTime: time.Now(),
	}
}

func (g *Faster) run() {
	for {
		select {
		case msg := <-g.evChannel:
			//log.Print("~~gofaster: ", msg)
			switch msg.Type {
			case internal.EvTrack:
				g.root.GetChild(msg.Path...).Active++
			case internal.EvDone:
				d := g.root.GetChild(msg.Path...)
				d.Active--
				d.Count++
				d.TotalTime += msg.Took
			case internal.EvSnapshot:
				var snap = g.takeSnapshotRec(g.root, time.Now())
				g.snapshotChannel <- snap
			case internal.EvReset:
				g.root = new(internal.Data)
			case internal.EvStop:
				return // TODO stop this GoFaster instance safely
			default:
				panic("unsupported go-faster event type")
			}
		case ticker := <-g.tickChan:
			//log.Print("tick: ", ticker)
			var snap = g.takeSnapshotRec(g.root, time.Now())
			ticker.Push(snap)
		}
	}
}

// GetSnapshot -- Creates and returns a deep copy of the current state (including child instance states)
func (g *Faster) GetSnapshot() *Snapshot {
	g.do(internal.EvSnapshot, nil, 0)
	return <-g.snapshotChannel
}

// takeSnapshot -- internal (-> thread-unsafe) method taking a deep copy of the current state
//
// should only ever be called from within the run() goroutine
// 'now' is passed all the way down
func (g *Faster) takeSnapshotRec(data *internal.Data, now time.Time) *Snapshot {
	var children = make(map[string]*Snapshot, len(data.Children))

	for key, child := range data.Children {
		children[key] = g.takeSnapshotRec(child, now)
	}

	var rc = Snapshot{
		Active:   data.Active,
		Average:  data.Average(),
		Count:    data.Count,
		Duration: data.TotalTime,
		Ts:       now,
		Children: children,
	}

	return &rc
}

// Reset -- Resets this GoFaster instance to its initial state
func (g *Faster) Reset() {
	g.do(internal.EvReset, nil, 0)
}

// New -- Construct a new root-level GoFaster instance
func New() *Faster {
	rc := &Faster{
		root:            new(internal.Data),
		evChannel:       make(chan internal.Event, 100),
		snapshotChannel: make(chan *Snapshot, 5),

		tickChan: make(chan *History),
		history:  make(map[string]*History),
	}
	go rc.run()

	return rc
}
