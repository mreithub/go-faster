package faster

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSingleton(t *testing.T) {
	Ref("hello").Deref()
	ref := Ref("world")
	snap1 := GetSnapshot()
	ref.Deref()
	snap2 := GetSnapshot()

	// current state
	assert.Contains(t, Singleton.root.Children, "hello")
	assert.Contains(t, Singleton.root.Children, "world")
	d := Singleton.root.GetChild("hello")
	assert.Equal(t, int32(0), d.Active)
	assert.Equal(t, int64(1), d.Count)
	d = Singleton.root.GetChild("world")
	assert.Equal(t, int32(0), d.Active)
	assert.Equal(t, int64(1), d.Count)

	// reset instance
	Reset()
	Ref("bla").Deref()

	GetSnapshot() // synchronize

	assert.NotContains(t, Singleton.root.Children, "hello")
	assert.Contains(t, Singleton.root.Children, "bla")

	//
	// check Snapshot data after the fact
	//

	// snap1: Ref('hello'), Deref('hello'), Ref('world')
	d1 := snap1.Children["hello"]
	assert.Equal(t, int32(0), d1.Active)
	assert.Equal(t, int64(1), d1.Count)
	d2 := snap1.Children["world"]
	assert.Equal(t, int32(1), d2.Active)
	assert.Equal(t, int64(0), d2.Count)
	assert.Equal(t, time.Duration(0), d2.Duration)
	assert.Equal(t, 2, len(snap1.Children))

	// snap2: snap1 + Deref('world')
	d1 = snap2.Children["hello"]
	assert.Equal(t, int32(0), d1.Active)
	assert.Equal(t, int64(1), d1.Count)
	d2 = snap2.Children["world"]
	assert.Equal(t, int32(0), d2.Active)
	assert.Equal(t, int64(1), d2.Count)
	assert.True(t, d2.Duration > 0)
	assert.Equal(t, 2, len(snap2.Children))
}
