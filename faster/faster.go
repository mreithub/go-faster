package faster

import (
	"sync"
	"time"

	"github.com/mreithub/go-faster/faster/internal"
)

// TODO tracking execution time might cause performance issues (e.g. in virtualized environments gettimeofday() might be slow)
//   if that turns out to be the case, deactivate Data.TotalNsec

// event types (for internal communication):
const (
	// stop the goroutine handling this GoFaster instance
	evStop = iota
	// resets this GoFaster instance
	evReset = iota
	// Takes a snapshot and sends it to snapshotChannel
	evSnapshot = iota
	// increments a ref counter
	evRef = iota
	// decrements a ref counter (and updates the total count + time)
	evDeref = iota
)

type event struct {
	typ  int
	key  string
	took time.Duration
}

// Faster -- A simple, go-style key-based reference counter that can be used for profiling your application (main class)
type Faster struct {
	name   string
	parent *Faster

	data map[string]*internal.Data

	_children map[string]*Faster
	childLock sync.Mutex

	evChannel       chan event
	snapshotChannel chan Snapshot
}

func (g *Faster) do(evType int, key string, took time.Duration) {
	g.evChannel <- event{
		typ:  evType,
		key:  key,
		took: took,
	}
}

// get -- Get the Data object for the specified key (or create it) - thread safe
func (g *Faster) get(key string) *internal.Data {
	rc, ok := g.data[key]
	if !ok {
		rc = &internal.Data{}
		g.data[key] = rc
	}

	return rc
}

// Ref -- References an instance of 'key'
func (g *Faster) Ref(key string) *Instance {
	g.do(evRef, key, 0)

	return &Instance{
		parent:    g,
		key:       key,
		startTime: time.Now(),
	}
}

func (g *Faster) run() {
	for msg := range g.evChannel {
		//log.Print("~~gofaster: ", msg)
		switch msg.typ {
		case evRef:
			g.get(msg.key).Active++
			break
		case evDeref:
			d := g.get(msg.key)
			d.Active--
			d.Count++
			d.TotalTime += msg.took
			break
		case evSnapshot:
			g.takeSnapshot()
			break
		case evReset:
			g.data = map[string]*internal.Data{}
			break
		case evStop:
			return // TODO stop this GoFaster instance safely
		default:
			panic("unsupported GoFaster event type")
		}
	}
}

// GetChild -- Gets (or creates) a specific child instance (recursively)
func (g *Faster) GetChild(path ...string) *Faster {
	if len(path) == 0 {
		return g
	}

	firstSegment := path[0]

	var child *Faster
	{ // keep the lock as short as possible
		g.childLock.Lock()
		defer g.childLock.Unlock()

		var ok bool
		child, ok = g._children[firstSegment]
		if !ok {
			// create a new child transparently
			child = newFaster(firstSegment, g)
			g._children[firstSegment] = child
		}
	}

	return child.GetChild(path[1:]...)
}

// GetChildren -- Creates a point-in-time copy of this GoFaster instance's children
func (g *Faster) GetChildren() map[string]*Faster {
	g.childLock.Lock()
	defer g.childLock.Unlock()

	// simply copy all entries
	var rc = make(map[string]*Faster, len(g._children))
	for name, child := range g._children {
		rc[name] = child
	}

	return rc
}

// GetParent -- Get the parent of this GoFaster instance (will return nil for root instances)
func (g *Faster) GetParent() *Faster {
	// g.parent is immutable -> no locking necessary
	return g.parent
}

// GetPath -- Get this GoFaster instance's path (i.e. its parents' and its own name)
//
// Root instances have empty names, all the others have the name you give them
// when creating them with GetChild().
//
// To get a single string path, you can use strings.Join()
//
// ```go
// strings.Join(g.GetPath(), "/")
// ```
func (g *Faster) GetPath() []string {
	var rc []string
	// this method needs no thread safety mechanisms.
	// A GoFaster's name and parent are immutable.
	if g.parent != nil {
		rc = append(g.parent.GetPath(), g.name)
	}

	return rc
}

// GetSnapshot -- Creates and returns a deep copy of the current state (including child instance states)
func (g *Faster) GetSnapshot() Snapshot {
	g.do(evSnapshot, "", 0)

	// get child snapshots while we wait
	children := g.GetChildren()
	childData := make(map[string]Snapshot, len(children))

	for name, child := range children {
		childData[name] = child.GetSnapshot()
	}

	rc := <-g.snapshotChannel
	rc.Children = childData
	return rc
}

// takeSnapshot -- internal (-> thread-unsafe) method taking a deep copy of the current state and sending it to snapshotChannel
func (g *Faster) takeSnapshot() {
	// copy entries
	data := make(map[string]Data, len(g.data))
	for key, d := range g.data {
		data[key] = newData(d)
	}

	// send Snapshot
	g.snapshotChannel <- Snapshot{
		Data: data,
		Ts:   time.Now(),
	}
}

// Reset -- Resets this GoFaster instance to its initial state
func (g *Faster) Reset() {
	g.do(evReset, "", 0)
}

// New -- Construct a new root-level GoFaster instance
func New() *Faster {
	return newFaster("", nil)
}

func newFaster(name string, parent *Faster) *Faster {
	rc := &Faster{
		name:            name,
		parent:          parent,
		data:            map[string]*internal.Data{},
		_children:       map[string]*Faster{},
		evChannel:       make(chan event, 100),
		snapshotChannel: make(chan Snapshot, 5),
	}
	go rc.run()

	return rc
}
