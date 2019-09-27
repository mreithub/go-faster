package faster

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBasics(t *testing.T) {
	f := New(true)

	f.Track("hello").Done()

	snap1 := f.TakeSnapshot()
	ref := f.Track("world")
	time.Sleep(100 * time.Millisecond)
	snap2 := f.TakeSnapshot()
	ref.Done()
	ref = f.Track("hello")
	snap3 := f.TakeSnapshot()
	ref.Done()

	// all the assertions are done after the fact (to make sure the different snapshots
	// keep their own copies of the Data)

	f.TakeSnapshot() // wait for run() to catch up

	// final (current) state
	assert.True(t, f.tree.Exists("hello"))
	assert.True(t, f.tree.Exists("world"))
	d := f.data[f.tree.GetIndex("hello")]
	assert.Equal(t, int32(0), d.Active)
	assert.Equal(t, int64(2), d.Count)
	assert.True(t, d.TotalTime > 0)
	d = f.data[f.tree.GetIndex("world")]
	assert.Equal(t, int32(0), d.Active)
	assert.Equal(t, int64(1), d.Count)
	assert.True(t, d.TotalTime >= 100000000)

	// snap1: Track('hello'), Done('hello')
	assert.True(t, snap1.tree.Exists("hello"))
	assert.False(t, snap1.tree.Exists("world"))
	d1 := snap1.Get("hello")
	assert.Equal(t, int32(0), d1.Active)
	assert.Equal(t, int64(1), d1.Count)
	assert.True(t, d1.TotalTime > 0)
	assert.Equal(t, 0, snap1.tree.GetIndex())
	assert.Equal(t, 1, snap1.tree.GetIndex("hello"))
	assert.Equal(t, -1, snap1.tree.GetIndex("world"))
	assert.Equal(t, 2, len(snap1.data))

	// snap2: snap1 + Track('world'),  sleep(100ms)
	assert.True(t, snap2.tree.Exists("hello"))
	assert.True(t, snap2.tree.Exists("world"))
	d2 := snap2.Get("world")
	assert.Equal(t, int32(1), d2.Active)
	assert.Equal(t, int64(0), d2.Count)
	assert.Equal(t, time.Duration(0), d2.TotalTime)

	// snap3: snap2 + Done('world'), Track('hello')
	assert.True(t, snap3.tree.Exists("hello"))
	assert.True(t, snap3.tree.Exists("world"))
	d3 := snap3.Get("world")
	assert.Equal(t, int32(0), d3.Active)
	assert.Equal(t, int64(1), d3.Count)
	assert.True(t, d3.TotalTime >= 100*time.Millisecond)
	assert.True(t, snap3.Get("hello").TotalTime < 100*time.Microsecond) // arbitrary value, but should be longer than we need
	assert.NotEqual(t, d1.TotalTime, d3.TotalTime)
}

func TestReflection(t *testing.T) {
	f := New(true)

	var tracker = f.TrackFn()
	tracker.Done()

	assert.Equal(t, []string{"src", "faster", "TestReflection()"}, tracker.path)

	snap := f.TakeSnapshot()

	assert.True(t, snap.tree.Exists("src"))
	assert.True(t, snap.tree.Exists("src", "faster"))
	assert.True(t, snap.tree.Exists("src", "faster", "TestReflection()"))

	assert.Equal(t, []string{"src", "foo", "Bar", "*Func()"}, f.parseCaller("github.com/mreithub/foo.(*Bar).Func"))
	assert.Equal(t, []string{"src", "foo", "Func()"}, f.parseCaller("github.com/mreithub/foo.Func"))
}
