package faster

import (
	"time"

	"github.com/mreithub/go-faster/faster/internal"
)

// TODO tracking execution time might cause performance issues (e.g. in virtualized environments gettimeofday() might be slow)
//   if that turns out to be the case, deactivate Data.TotalNsec

// Faster -- A simple, go-style key-based reference counter that can be used for profiling your application (main class)
type Faster struct {
	// internal data structure -- only ever accessed from within the run() goroutine // TODO check if this is true
	root *internal.Data

	evChannel       chan internal.Event
	snapshotChannel chan *Snapshot
}

func (g *Faster) do(evType internal.EventType, path []string, took time.Duration) {
	g.evChannel <- internal.Event{
		Type: evType,
		Path: path,
		Took: took,
	}
}

// Ref -- References an instance of 'key'
func (g *Faster) Ref(key ...string) *Instance {
	g.do(internal.EvRef, key, 0)

	return &Instance{
		parent:    g,
		path:      key,
		startTime: time.Now(),
	}
}

func (g *Faster) run() {
	for msg := range g.evChannel {
		//log.Print("~~gofaster: ", msg)
		switch msg.Type {
		case internal.EvRef:
			g.root.GetChild(msg.Path...).Active++
			break
		case internal.EvDeref:
			d := g.root.GetChild(msg.Path...)
			d.Active--
			d.Count++
			d.TotalTime += msg.Took
			break
		case internal.EvSnapshot:
			var snap = g.takeSnapshotRec(g.root, time.Now())
			g.snapshotChannel <- snap
			break
		case internal.EvReset:
			g.root = new(internal.Data)
			break
		case internal.EvStop:
			return // TODO stop this GoFaster instance safely
		default:
			panic("unsupported GoFaster event type")
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
	}
	go rc.run()

	return rc
}
