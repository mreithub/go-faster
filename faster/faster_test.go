package faster

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBasics(t *testing.T) {
	f := New(true)

	f.Track("hello").Done()

	clone1 := f.GetSnapshot()
	ref := f.Track("world")
	time.Sleep(100 * time.Millisecond)
	clone2 := f.GetSnapshot()
	ref.Done()
	ref = f.Track("hello")
	clone3 := f.GetSnapshot()
	ref.Done()

	// all the assertions are done after the fact (to make sure the different clones
	// keep their own copies of the Data)

	f.GetSnapshot() // wait for run() to catch up

	// final (current) state
	assert.Contains(t, f.root.Children, "hello")
	assert.Contains(t, f.root.Children, "world")
	d := f.root.GetChild("hello")
	assert.Equal(t, int32(0), d.Active)
	assert.Equal(t, int64(2), d.Count)
	assert.True(t, d.TotalTime > 0)
	d = f.root.GetChild("world")
	assert.Equal(t, int32(0), d.Active)
	assert.Equal(t, int64(1), d.Count)
	assert.True(t, d.TotalTime >= 100000000)

	// clone1: Track('hello'), Done('hello')
	keys := clone1.Keys()
	assert.Contains(t, keys, "hello")
	assert.NotContains(t, keys, "world")
	d1 := clone1.Children["hello"]
	assert.Equal(t, int32(0), d1.Active)
	assert.Equal(t, int64(1), d1.Count)
	assert.True(t, d1.Duration > 0)
	assert.Equal(t, 1, len(clone1.Children))

	// clone2: clone1 + Track('world'),  sleep(100ms)
	keys = clone2.Keys()
	assert.Contains(t, keys, "hello")
	assert.Contains(t, keys, "world")
	d2 := clone2.Children["world"]
	assert.Equal(t, int32(1), d2.Active)
	assert.Equal(t, int64(0), d2.Count)
	assert.Equal(t, time.Duration(0), d2.Duration)

	// clone3: clone2 + Done('world'), Track('hello')
	keys = clone3.Keys()
	assert.Contains(t, keys, "hello")
	assert.Contains(t, keys, "world")
	d3 := clone3.Children["world"]
	assert.Equal(t, int32(0), d3.Active)
	assert.Equal(t, int64(1), d3.Count)
	assert.True(t, d3.Duration >= 100000000)
	assert.True(t, clone3.Children["hello"].Duration < 100000)
	assert.NotEqual(t, d1.Duration, d3.Duration)
}

func TestTrackFn(t *testing.T) {
	f := New(true)

	var tracker = f.TrackFn()
	tracker.Done()

	assert.Equal(t, []string{"src", "faster", "TestTrackFn()"}, tracker.path)

	clone := f.GetSnapshot()

	assert.Contains(t, clone.Children, "src")
	assert.Contains(t, clone.Get("src").Children, "faster")
	assert.Contains(t, clone.Get("src", "faster").Children, "TestTrackFn()")
}
