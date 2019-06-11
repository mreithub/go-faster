package faster

import (
	"regexp"
	"runtime"
	"sync"
	"time"

	"github.com/mreithub/go-faster/faster/internal"
	"github.com/mreithub/go-faster/histogram"
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

	// StartTS -- timestamp of this Faster object's creation
	StartTS time.Time
}

func (f *Faster) do(evType internal.EventType, path []string, took time.Duration) {
	f.evChannel <- internal.Event{
		Type: evType,
		Path: path,
		Took: took,
	}
}

// getCaller -- returns the given stack trace entry in the format we want it
func (f *Faster) getCaller(skip int) []string {
	pc := make([]uintptr, 5)
	n := runtime.Callers(skip+2, pc) // also skip getCaller() and runtime.Callers()
	if n > 0 {
		frames := runtime.CallersFrames(pc[:n])
		for { // iterate over frames
			frame, more := frames.Next()

			if fn := frame.Function; fn != "" {
				return f.parseCaller(fn)
			}
			if !more {
				break
			}
		}
	}

	return []string{"src"} // couldn't determine caller -> track using the top level '/src' key
}

// frame.Function is in the format "qualified/package/name.(*type).function"
// (with type being optional - and the '(*)' only there for pointer receivers)
var callerPattern = regexp.MustCompile("([^/.]+)(\\.([^./]+))?\\.([^./]+)$")

func (*Faster) parseCaller(fn string) []string {
	var groups = callerPattern.FindStringSubmatch(fn)
	if len(groups) == 0 {
		return []string{"src"}
	}

	var pkg, typeName, function = groups[1], groups[3], groups[4]

	if typeName == "" {
		return []string{"src", pkg, function + "()"}
	}
	return []string{"src", pkg, typeName, function + "()"}
}

// SetTicker -- Sets up periodic snapshots
//
// - name is the unique name of the given History ticker
// - interval specifies how often these ticks should happen (if 0, the ticker will be deleted)
// - keep is the number of past Snapshots stored for the given History object
func (f *Faster) SetTicker(name string, interval time.Duration, keep int) {
	f.historyLock.Lock()
	defer f.historyLock.Unlock()

	if ticker, ok := f.history[name]; ok {
		// replacing/removing an existing ticker -> stop the old one
		ticker.Stop()
	}

	if interval == 0 {
		return
	}

	f.history[name] = NewHistory(name, interval, keep, f.tickChan)
}

// Track -- Tracks an instance of 'key'
func (f *Faster) Track(key ...string) *Tracker {
	f.do(internal.EvTrack, key, 0)

	return &Tracker{
		parent:  f,
		path:    key,
		startTS: time.Now(),
	}
}

// TrackFn -- Tracks the calling function (using ["src", "pkgName", "typeName", "fn()"] as key - omitting typeName if empty)
func (f *Faster) TrackFn() *Tracker {
	var key = f.getCaller(1)
	return f.Track(key...)
}

func (f *Faster) run() {
	for {
		select {
		case msg := <-f.evChannel:
			//log.Print("~~gofaster: ", msg)
			switch msg.Type {
			case internal.EvTrack:
				f.root.GetChild(msg.Path...).Active++
			case internal.EvDone:
				f.root.GetChild(msg.Path...).Done(msg.Took)
			case internal.EvSnapshot:
				var snap = f.takeSnapshotRec(f.root, time.Now())
				f.snapshotChannel <- snap
			case internal.EvReset:
				f.root = new(internal.Data)
			case internal.EvStop:
				return // TODO stop this GoFaster instance safely
			default:
				panic("unsupported go-faster event type")
			}
		case history := <-f.tickChan:
			//log.Print("tick: ", history)
			var snap = f.takeSnapshotRec(f.root, time.Now())
			history.push(snap)
		}
	}
}

// ListTickers -- returns the (currently registered) History tickers (taking periodic snapshots)
func (f *Faster) ListTickers() map[string]*History {
	f.historyLock.Lock()
	defer f.historyLock.Unlock()

	var rc = make(map[string]*History, len(f.history))
	for k, v := range f.history {
		rc[k] = v
	}

	return rc
}

// GetSnapshot -- Creates and returns a deep copy of the current state (including child instance states)
func (f *Faster) GetSnapshot() *Snapshot {
	f.do(internal.EvSnapshot, nil, 0)
	return <-f.snapshotChannel
}

// takeSnapshot -- internal (-> thread-unsafe) method taking a deep copy of the current state
//
// should only ever be called from within the run() goroutine
// 'now' is passed all the way down
func (f *Faster) takeSnapshotRec(data *internal.Data, now time.Time) *Snapshot {
	var children = make(map[string]*Snapshot, len(data.Children))

	for key, child := range data.Children {
		children[key] = f.takeSnapshotRec(child, now)
	}

	var rc = Snapshot{
		Active:   data.Active,
		Average:  data.Average(),
		Count:    data.Count,
		Duration: data.TotalTime,
		Ts:       now,
		Children: children,
	}

	if data.Histogram != nil {
		rc.Histogram = data.Histogram.Copy()
	}

	return &rc
}

// Reset -- Resets this GoFaster instance to its initial state
func (f *Faster) Reset() {
	f.do(internal.EvReset, nil, 0)
}

// New -- Construct a new root-level GoFaster instance
func New(withHistograms bool) *Faster {
	rc := &Faster{
		root:            new(internal.Data),
		evChannel:       make(chan internal.Event, 100),
		snapshotChannel: make(chan *Snapshot, 5),

		tickChan: make(chan *History),
		history:  make(map[string]*History),
		StartTS:  time.Now(),
	}

	if withHistograms {
		rc.root.Histogram = new(histogram.Histogram)
	}
	go rc.run()

	return rc
}
