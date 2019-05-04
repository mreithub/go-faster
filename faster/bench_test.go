package faster

import (
	"fmt"
	"testing"
	"time"
)

var nsec time.Duration
var snap *Snapshot

// BenchmarkMeasureTime -- Measures how long measuring the time takes (using time.Now() and Time.Sub())
func BenchmarkMeasureTime(b *testing.B) {
	for n := 0; n < b.N; n++ {
		start := time.Now()
		end := time.Now()
		nsec = end.Sub(start)
	}

}

// BenchmarkTrackDone -- Measures how long an empty Track().Done() call takes
func BenchmarkTrackDone(b *testing.B) {
	g := New(true)

	for n := 0; n < b.N; n++ {
		g.Track("hello").Done()
	}
	//snap := g.Clone()
	//j, _ := json.Marshal(snap.Data)
	//log.Printf("data: %s", j)
}

// BenchmarkTrackDoneDeferred -- Measures how long an empty Track().Done() call takes (doing the Done() in a defer statement)
func BenchmarkTrackDoneDeferred(b *testing.B) {
	g := New(true)

	for n := 0; n < b.N; n++ {
		r := g.Track("hello")
		defer r.Done()
	}
	//snap := g.Clone()
	//j, _ := json.Marshal(snap.Data)
	//log.Printf("data: %s", j)
}

// benchmarkGetSnapshot -- Measure how long it takes to create a deep copy of the snapshot data
func benchmarkGetSnapshot(count int, b *testing.B) {
	// setup
	g := New(true)
	for n := 0; n < count; n++ {
		g.Track(fmt.Sprintf("ref%d", n)).Done()
	}

	for n := 0; n < b.N; n++ {
		snap = g.GetSnapshot()
	}
}

func BenchmarkGetSnapshot100(b *testing.B) {
	benchmarkGetSnapshot(100, b)
}

func BenchmarkGetSnapshot1000(b *testing.B) {
	benchmarkGetSnapshot(1000, b)
}
