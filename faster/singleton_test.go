package faster

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSingleton(t *testing.T) {
	Track("hello").Done()
	ref := Track("world")
	snap1 := TakeSnapshot()
	ref.Done()
	snap2 := TakeSnapshot()

	// current state
	assert.True(t, Singleton.tree.Exists("hello"))
	assert.True(t, Singleton.tree.Exists("world"))
	d := Singleton.getDataForPath("hello")
	assert.Equal(t, int32(0), d.Active)
	assert.Equal(t, int64(1), d.Count)
	d = Singleton.getDataForPath("world")
	assert.Equal(t, int32(0), d.Active)
	assert.Equal(t, int64(1), d.Count)

	// reset instance
	Reset()
	Track("bla").Done()

	TakeSnapshot() // synchronize

	assert.False(t, Singleton.tree.Exists("hello"))
	assert.True(t, Singleton.tree.Exists("bla"))

	//
	// check old Snapshot data after the fact
	//

	// snap1: Track('hello'), Done('hello'), Track('world')
	d1 := snap1.Get("hello")
	assert.Equal(t, int32(0), d1.Active)
	assert.Equal(t, int64(1), d1.Count)
	d2 := snap1.Get("world")
	assert.Equal(t, int32(1), d2.Active)
	assert.Equal(t, int64(0), d2.Count)
	assert.Equal(t, time.Duration(0), d2.TotalTime)
	assert.Equal(t, 3, len(snap1.data))
	assert.Equal(t, 2, len(snap1.tree.Children()))

	// snap2: snap1 + Done('world')
	d1 = snap2.Get("hello")
	assert.Equal(t, int32(0), d1.Active)
	assert.Equal(t, int64(1), d1.Count)
	d2 = snap2.Get("world")
	assert.Equal(t, int32(0), d2.Active)
	assert.Equal(t, int64(1), d2.Count)
	assert.True(t, d2.TotalTime > 0)
	assert.Equal(t, 3, len(snap1.data))
	assert.Equal(t, 2, len(snap2.tree.Children()))

	// test reflection
	var tracker = TrackFn()
	tracker.Done()
	assert.Equal(t, []string{"src", "faster", "TestSingleton()"}, tracker.path)
	assert.Equal(t, int64(1), TakeSnapshot().Get("src", "faster", "TestSingleton()").Count)
}
