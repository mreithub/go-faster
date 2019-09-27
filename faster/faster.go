package faster

import (
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/mreithub/go-faster/faster/internal"
)

// TODO tracking execution time might cause performance issues (e.g. in virtualized environments gettimeofday() might be slow)
//   if that turns out to be the case, deactivate Data.TotalNsec

// Faster -- A simple, go-style key-based reference counter that can be used for profiling your application (main class)
type Faster struct {
	tree       internal.RWTree
	data       []data
	histograms []Histogram

	withHistograms bool

	// processed by the run() goroutine
	evChannel chan internal.Event

	// calling TakeSnapshot() triggers an internal.EvSnapshot which in turn causes
	// the run() goroutine to take one and push it here -- while it's not guaranteed
	// that each TakeSnapshot() call will read the response to its own request
	// (multiple TakeSnapshot() calls might block concurrently), that doesn't really
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

	return []string{"src"} // couldn't determine caller -> track using the top level 'src' key
}

func (*Faster) parseCaller(fn string) []string {
	// frame.Function is in the format "golang.org/qualified/package.(*type).function"
	// (with type being optional - and the '(*)' only there for pointer receivers)

	// -> we're only interested in the stuff after the last slash
	var parts = strings.Split(fn, "/")
	fn = parts[len(parts)-1]

	parts = strings.Split(fn, ".")
	var rc = make([]string, 0, len(parts)+1)
	rc = append(rc, "src")

	if len(parts) == 2 {
		// "package.function"
		rc = append(rc, parts[0], parts[1]+"()")
	} else if len(parts) == 3 {
		// "package.type.function"
		var typeName = parts[1]
		fn = parts[2] + "()"
		if strings.HasPrefix(typeName, "(*") && strings.HasSuffix(typeName, ")") {
			// "package.(*type).function()" -> "package.type.*function()"
			// (we don't want two entries - "(*type)" and "type" - for types that have methods with both value and pointer receivers)
			typeName = typeName[2 : len(typeName)-1]
			fn = "*" + fn
		}

		rc = append(rc, parts[0], typeName, fn)
	}

	return rc
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
				f.onTrack(msg.Path)
			case internal.EvDone:
				f.onDone(msg.Path, msg.Took)
			case internal.EvSnapshot:
				var snap = f.takeSnapshot(time.Now())
				f.snapshotChannel <- snap
			case internal.EvReset:
				f.onReset()
			case internal.EvStop:
				return // TODO stop this GoFaster instance safely
			default:
				panic("unsupported go-faster event type")
			}
		case history := <-f.tickChan:
			//log.Print("tick: ", history)
			var snap = f.takeSnapshot(time.Now())
			history.push(snap)
		}
	}
}

// getData -- returns a pointer to the internal.Data object with the given index (extending f.data if necessary)
func (f *Faster) getData(index int) *data {
	if index >= len(f.data) {
		f.data = append(f.data, make([]data, index-len(f.data)+1)...)
	}
	return &f.data[index]
}

// getDataForPath -- returns a pointer to the internal.Data object with the given path (or creates it if necessary)
func (f *Faster) getDataForPath(path ...string) *data {
	var index = f.tree.GetIndex(path...)
	return f.getData(index)
}

func (f *Faster) onDone(path []string, took time.Duration) {
	f.getDataForPath(path...).Done(took)
}

func (f *Faster) onReset() {
	f.data = nil
	f.histograms = nil
	f.tree.Reset()
}

func (f *Faster) onTrack(path []string) {
	f.getDataForPath(path...).active++
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

// TakeSnapshot -- tells the go-faster goroutine to take and return a deep copy of its current state
func (f *Faster) TakeSnapshot() *Snapshot {
	f.do(internal.EvSnapshot, nil, 0)
	return <-f.snapshotChannel
}

// takeSnapshot -- internal (-> thread-unsafe) method taking a deep copy of the current state
//
// should only ever be called from within the run() goroutine
func (f *Faster) takeSnapshot(now time.Time) *Snapshot {
	var rc = Snapshot{
		tree:       f.tree.Clone(),
		data:       make([]data, len(f.data)),
		histograms: make([]Histogram, len(f.histograms)),
		TS:         now,
	}
	copy(rc.data, f.data)
	copy(rc.histograms, f.histograms)

	return &rc
}

// Reset -- Resets this GoFaster instance to its initial state
func (f *Faster) Reset() {
	f.do(internal.EvReset, nil, 0)
}

// New -- Construct a new root-level GoFaster instance
func New(withHistograms bool) *Faster {
	rc := &Faster{
		withHistograms: withHistograms,

		evChannel:       make(chan internal.Event, 100),
		snapshotChannel: make(chan *Snapshot, 5),

		tickChan: make(chan *History),
		history:  make(map[string]*History),
		StartTS:  time.Now(),
	}

	go rc.run()

	return rc
}
